package docs

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func newTestEmbedder() *Embedder {
	return NewEmbedder(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestSearchReturnsMostRelevantDocument(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"economia.md":  "A inflação é o aumento da oferta monetária e a queda do poder de compra da moeda.",
		"culinaria.md": "Receita de bolo de cenoura com cobertura de chocolate e farinha de trigo.",
		"ignored.json": "inflação inflação inflação",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	e := newTestEmbedder()
	if err := e.Index(dir); err != nil {
		t.Fatalf("Index: %v", err)
	}
	if len(e.docs) != 2 {
		t.Fatalf("indexed %d chunks, want 2 (non .md/.txt/.pdf files must be skipped)", len(e.docs))
	}

	results, err := e.Search("oferta monetária e inflação", 1)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 || results[0].Filename != "economia.md" {
		t.Fatalf("got %+v, want economia.md first", results)
	}
}

func TestIndexMissingDirectory(t *testing.T) {
	e := newTestEmbedder()
	if err := e.Index(filepath.Join(t.TempDir(), "missing")); err != nil {
		t.Fatalf("Index on missing dir: %v", err)
	}
	results, err := e.Search("anything", 3)
	if err != nil || len(results) != 0 {
		t.Fatalf("got %v, %v; want no results", results, err)
	}
}
