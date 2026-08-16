// Command agentmem-site serves the public docs site embedded in the binary.
//
//	go run . [-addr :5555]
//	go build -o agentmem-site . && ./agentmem-site -addr :5555
//
// Features:
//   - site/ embedded via go:embed (fully portable binary)
//   - i18n from site/lang/<code>.yaml (default en; ?lang= or Accept-Language)
//   - html/template layout (shared head, nav, footer)
//   - extensionless routes: /docs/install
//   - .md routes: /docs/install.md → text/markdown
//   - static files under /css, /js, /assets
package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
)

//go:embed all:site
var embeddedSite embed.FS

// version is set at link time: make build VERSION=x.y.z
var version = "dev"

// page is route metadata (not user-facing copy — copy lives in lang/*.yaml).
type page struct {
	// Content is the template name and i18n page key (e.g. "index", "install").
	Content string
	// Nav is the active primary nav key.
	Nav  string
	Home bool
	// Markdown is the path for the "Markdown" footer link.
	Markdown string
	// Prev / Next paths for footer (labels come from i18n).
	Prev, Next string
	// PrevKey / NextKey are i18n keys under common.prev / common.next.
	PrevKey, NextKey string
}

// view is the template data bag.
type view struct {
	page
	Lang    string
	LangTag string // html lang attribute
	pageKey string
	cat     catalog
	URLBase string // absolute base URL (e.g. https://agentmem.thamtech.co) for canonical + OG tags
	Path    string // request path (no query, no fragment)
}

// T looks up page-local then common keys (or absolute dotted paths).
// Method form so templates can call {{.T "key"}}.
func (v view) T(key string) string {
	return v.cat.t(v.pageKey, key)
}

// L builds an internal path that preserves the active locale for htmx
// and full-page navigations. Accepts paths like "/docs/install" or
// "/docs/integrations#claude-code". External URLs and bare fragments
// are returned unchanged.
func (v view) L(p string) string {
	if p == "" || p == "#" || strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "//") || strings.HasPrefix(p, "mailto:") {
		return p
	}
	frag := ""
	if i := strings.IndexByte(p, '#'); i >= 0 {
		frag = p[i:]
		p = p[:i]
	}
	if p == "" {
		return frag
	}
	// Replace an existing lang= query if present; otherwise append.
	if i := strings.IndexByte(p, '?'); i >= 0 {
		q := p[i+1:]
		p = p[:i]
		parts := strings.Split(q, "&")
		out := make([]string, 0, len(parts)+1)
		for _, part := range parts {
			if part == "" || strings.HasPrefix(part, "lang=") {
				continue
			}
			out = append(out, part)
		}
		out = append(out, "lang="+v.Lang)
		return p + "?" + strings.Join(out, "&") + frag
	}
	return p + "?lang=" + v.Lang + frag
}

var routes = map[string]page{
	"/": {
		Content:  "index",
		Home:     true,
		Markdown: "index.md",
	},
	"/docs/install": {
		Content:  "install",
		Nav:      "install",
		Markdown: "docs/install.md",
		Prev:     "/",
		PrevKey:  "home",
		Next:     "/docs/commands",
		NextKey:  "commands",
	},
	"/docs/commands": {
		Content:  "commands",
		Nav:      "commands",
		Markdown: "docs/commands.md",
		Prev:     "/docs/install",
		PrevKey:  "install",
		Next:     "/docs/integrations",
		NextKey:  "integrations",
	},
	"/docs/integrations": {
		Content:  "integrations",
		Nav:      "integrations",
		Markdown: "docs/integrations.md",
		Prev:     "/docs/commands",
		PrevKey:  "commands",
		Next:     "/docs/how-it-works",
		NextKey:  "how_it_works",
	},
	"/docs/how-it-works": {
		Content:  "how-it-works",
		Nav:      "how-it-works",
		Markdown: "docs/how-it-works.md",
		Prev:     "/docs/integrations",
		PrevKey:  "integrations",
		Next:     "/",
		NextKey:  "home",
	},
}

// baseURL is the absolute origin used for canonical + Open Graph URLs.
// Overridden via -url-base; defaults to the production origin.
var baseURL = "https://agentmem.thamtech.co"

func main() {
	addr := flag.String("addr", ":5555", "listen address")
	urlBase := flag.String("url-base", baseURL, "absolute base URL for canonical + OG tags")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	baseURL = strings.TrimRight(*urlBase, "/")

	if *showVersion {
		fmt.Println(version)
		return
	}

	site, err := fs.Sub(embeddedSite, "site")
	if err != nil {
		log.Fatal(err)
	}

	// Preload default English (fail fast if missing)
	if _, err := loadCatalog(site, "en"); err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFS(site, "templates/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	h := &server{fsys: site, tmpl: tmpl}
	log.Printf("agentmem site %s  http://127.0.0.1%s  (embedded, i18n)", version, *addr)
	if err := http.ListenAndServe(*addr, h); err != nil {
		log.Fatal(err)
	}
}

type server struct {
	fsys fs.FS
	tmpl *template.Template
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	p := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/"))
	if p != "/" {
		p = strings.TrimSuffix(p, "/")
	}

	if strings.HasPrefix(p, "/css/") || strings.HasPrefix(p, "/js/") || strings.HasPrefix(p, "/assets/") {
		s.serveStatic(w, r, strings.TrimPrefix(p, "/"))
		return
	}
	if p == "/robots.txt" {
		s.serveStatic(w, r, "static/robots.txt")
		return
	}
	if p == "/sitemap.xml" {
		s.serveStatic(w, r, "static/sitemap.xml")
		return
	}
	if p == "/install" {
		s.serveEmbedded(w, r, "install", "text/plain; charset=utf-8")
		return
	}

	if strings.HasSuffix(p, ".md") {
		if !s.tryMarkdown(w, r, p) {
			http.NotFound(w, r)
		}
		return
	}

	if wantsMarkdown(r) {
		md := p
		if md == "/" {
			md = "/index.md"
		} else {
			md = md + ".md"
		}
		if s.tryMarkdown(w, r, md) {
			return
		}
	}

	pg, ok := routes[p]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Locale for HTML is only ?lang= (default en). Cookie / Accept-Language
	// must not affect the body: Cloudflare would BYPASS private/varied HTML
	// and every first hit to / would go to origin. Client JS redirects TH
	// users to ?lang=th (see layout.tmpl + site/js/index.js).
	lang := "en"
	if strings.TrimSpace(r.URL.Query().Get("lang")) != "" {
		lang = pickLang(r)
	}
	cat, err := loadCatalog(s.fsys, lang)
	if err != nil {
		log.Printf("i18n: %v", err)
		http.Error(w, "language pack error", http.StatusInternalServerError)
		return
	}
	v := view{
		page:    pg,
		Lang:    lang,
		LangTag: lang,
		pageKey: normalizePage(pg.Content),
		cat:     cat,
		URLBase: baseURL,
		Path:    p,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Language", lang)
	// Distinct URLs for en/th; no Set-Cookie (CF will not cache those).
	w.Header().Set("Cache-Control", "public, max-age=60, s-maxage=86400")
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}
	if err := s.tmpl.ExecuteTemplate(w, "layout", v); err != nil {
		log.Printf("template %s: %v", p, err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// supportedLangs is the set of site/lang/<code>.yaml packs we ship.
var supportedLangs = map[string]bool{
	"en": true,
	"th": true,
}

func pickLang(r *http.Request) string {
	if q := strings.TrimSpace(r.URL.Query().Get("lang")); q != "" {
		return normalizeLangCode(q)
	}
	if c, err := r.Cookie("lang"); err == nil && c.Value != "" {
		return normalizeLangCode(c.Value)
	}
	// Accept-Language: first supported tag (not merely first tag)
	if al := r.Header.Get("Accept-Language"); al != "" {
		for _, raw := range strings.Split(al, ",") {
			part := strings.TrimSpace(strings.Split(raw, ";")[0])
			if part == "" {
				continue
			}
			code := normalizeLangCode(part)
			// normalizeLangCode maps unknown → en; only accept if the raw prefix is supported
			base := strings.ToLower(part)
			if i := strings.IndexAny(base, "-_"); i > 0 {
				base = base[:i]
			}
			if supportedLangs[base] {
				return code
			}
		}
	}
	return "en"
}

func normalizeLangCode(code string) string {
	code = strings.TrimSpace(strings.ToLower(code))
	if i := strings.IndexAny(code, "-_"); i > 0 {
		code = code[:i]
	}
	if !supportedLangs[code] {
		return "en"
	}
	return code
}

func wantsMarkdown(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if accept == "" || accept == "*/*" {
		return false
	}
	md := strings.Index(accept, "text/markdown")
	html := strings.Index(accept, "text/html")
	if md < 0 {
		return false
	}
	if html < 0 {
		return true
	}
	return md < html
}

func (s *server) tryMarkdown(w http.ResponseWriter, r *http.Request, urlPath string) bool {
	rel := strings.TrimPrefix(urlPath, "/")
	if !embeddedFile(s.fsys, rel) {
		return false
	}
	s.serveEmbedded(w, r, rel, "text/markdown; charset=utf-8")
	return true
}

func (s *server) serveStatic(w http.ResponseWriter, r *http.Request, rel string) {
	if !embeddedFile(s.fsys, rel) {
		http.NotFound(w, r)
		return
	}
	ct := mime.TypeByExtension(path.Ext(rel))
	if ct == "" {
		ct = "application/octet-stream"
	}
	if strings.HasPrefix(ct, "text/") && !strings.Contains(ct, "charset") {
		ct += "; charset=utf-8"
	}
	s.serveEmbedded(w, r, rel, ct)
}

func (s *server) serveEmbedded(w http.ResponseWriter, r *http.Request, rel, contentType string) {
	data, err := fs.ReadFile(s.fsys, rel)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	if strings.HasPrefix(rel, "css/") || strings.HasPrefix(rel, "js/") || strings.HasPrefix(rel, "assets/") {
		w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=86400")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=60, s-maxage=86400")
	}
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write(data)
}

func embeddedFile(fsys fs.FS, rel string) bool {
	rel = path.Clean("/" + rel)
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" || strings.HasPrefix(rel, "..") {
		return false
	}
	st, err := fs.Stat(fsys, rel)
	return err == nil && st.Mode().IsRegular()
}
