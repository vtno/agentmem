package main

import (
	"io/fs"
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
	if !strings.Contains(got, "agentmem") {
		t.Fatalf("title = %q", got)
	}
	if got := c.t("index", "why.label"); got != "ทำไม" {
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
