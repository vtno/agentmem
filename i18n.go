package main

import (
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

// catalog is a nested map loaded from lang/<code>.yaml.
// Root keys are page names (plus "common").
type catalog map[string]any

// loadCatalog reads site/lang/<code>.yaml (fallback en).
func loadCatalog(site fs.FS, code string) (catalog, error) {
	code = strings.TrimSpace(strings.ToLower(code))
	if code == "" {
		code = "en"
	}
	// simple language tag → file (en-US → en)
	if i := strings.IndexByte(code, '-'); i > 0 {
		code = code[:i]
	}
	if i := strings.IndexByte(code, '_'); i > 0 {
		code = code[:i]
	}

	data, err := fs.ReadFile(site, "lang/"+code+".yaml")
	if err != nil {
		if code != "en" {
			return loadCatalog(site, "en")
		}
		return nil, fmt.Errorf("load lang/%s.yaml: %w", code, err)
	}
	var c catalog
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse lang/%s.yaml: %w", code, err)
	}
	return c, nil
}

// asMap coerces nested YAML maps (catalog or map[string]any) for walking.
func asMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case catalog:
		return map[string]any(m), true
	default:
		return nil, false
	}
}

// get walks a dotted path (e.g. "common.nav.install" or "why.label" under page).
func (c catalog) get(path string) (any, bool) {
	if c == nil || path == "" {
		return nil, false
	}
	parts := strings.Split(path, ".")
	var cur any = catalog(c)
	for _, p := range parts {
		m, ok := asMap(cur)
		if !ok {
			return nil, false
		}
		next, ok := m[p]
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}

// str returns a string at path, or "".
func (c catalog) str(path string) string {
	v, ok := c.get(path)
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case fmt.Stringer:
		return t.String()
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

// t looks up key under page first (page.key), then as absolute path, then common.key.
func (c catalog) t(page, key string) string {
	if key == "" {
		return ""
	}
	// Absolute: common.x or page.x
	if strings.Contains(key, ".") {
		if s := c.str(key); s != "" {
			return s
		}
	}
	// Page-local
	if page != "" {
		if s := c.str(page + "." + key); s != "" {
			return s
		}
	}
	// Common fallback
	if s := c.str("common." + key); s != "" {
		return s
	}
	// Missing key marker (helps spot untranslated strings in dev)
	return "⟨" + key + "⟩"
}

// normalizePage maps content template name to yaml root key.
func normalizePage(content string) string {
	switch content {
	case "how-it-works":
		return "how_it_works"
	default:
		return content
	}
}
