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
