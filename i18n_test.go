package main

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoadCatalogEN(t *testing.T) {
	site, err := fs.Sub(embeddedSite, "site")
	if err != nil {
		t.Fatal(err)
	}
	c, err := loadCatalog(site, "en")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.str("common.nav.install"); got != "Install" {
		t.Fatalf("common.nav.install = %q", got)
	}
	got := c.t("index", "title")
	if strings.Contains(got, "title") && strings.HasPrefix(got, "\u27e8") {
		t.Fatalf("missing title: %q", got)
	}
	if !strings.Contains(got, "vectorless") {
		t.Fatalf("title = %q", got)
	}
	if got := c.t("index", "why.label"); got != "Why" {
		t.Fatalf("why.label = %q", got)
	}
}

func TestLoadCatalogTH(t *testing.T) {
	site, err := fs.Sub(embeddedSite, "site")
	if err != nil {
		t.Fatal(err)
	}
	c, err := loadCatalog(site, "th")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.str("common.nav.install"); got != "ติดตั้ง" {
		t.Fatalf("common.nav.install = %q", got)
	}
	got := c.t("index", "title")
	if strings.HasPrefix(got, "\u27e8") {
		t.Fatalf("missing title: %q", got)
	}
	if !strings.Contains(strings.ToLower(got), "agentmem") {
		t.Fatalf("title = %q", got)
	}
	if got := c.t("index", "why.label"); got != "ทำไมถึงต้องใช้" {
		t.Fatalf("why.label = %q", got)
	}
}

func TestNormalizeLangCode(t *testing.T) {
	cases := map[string]string{
		"en":    "en",
		"EN":    "en",
		"en-US": "en",
		"th":    "th",
		"th-TH": "th",
		"fr":    "en",
		"":      "en",
	}
	for in, want := range cases {
		if got := normalizeLangCode(in); got != want {
			t.Fatalf("normalizeLangCode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestViewL(t *testing.T) {
	v := view{Lang: "th"}
	cases := map[string]string{
		"/":                              "/?lang=th",
		"/docs/install":                  "/docs/install?lang=th",
		"/docs/integrations#claude-code": "/docs/integrations?lang=th#claude-code",
		"/docs/install?lang=en":          "/docs/install?lang=th",
		"https://github.com/vtno/agentmem": "https://github.com/vtno/agentmem",
		"#section":                       "#section",
	}
	for in, want := range cases {
		if got := v.L(in); got != want {
			t.Fatalf("L(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPickLang(t *testing.T) {
	req := func(modify func(*http.Request)) *http.Request {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		if modify != nil {
			modify(r)
		}
		return r
	}
	if got := pickLang(req(nil)); got != "en" {
		t.Fatalf("default = %q", got)
	}
	if got := pickLang(req(func(r *http.Request) {
		r.URL.RawQuery = "lang=th"
	})); got != "th" {
		t.Fatalf("query = %q", got)
	}
	if got := pickLang(req(func(r *http.Request) {
		r.AddCookie(&http.Cookie{Name: "lang", Value: "th"})
	})); got != "th" {
		t.Fatalf("cookie = %q", got)
	}
	if got := pickLang(req(func(r *http.Request) {
		r.Header.Set("Accept-Language", "fr-FR,th;q=0.8,en;q=0.5")
	})); got != "th" {
		t.Fatalf("accept-language = %q", got)
	}
	// Query wins over cookie
	if got := pickLang(req(func(r *http.Request) {
		r.URL.RawQuery = "lang=en"
		r.AddCookie(&http.Cookie{Name: "lang", Value: "th"})
	})); got != "en" {
		t.Fatalf("query over cookie = %q", got)
	}
}
