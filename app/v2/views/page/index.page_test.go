package page_test

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/memnix/memnix-rest/app/v2/views/page"
)

func TestHero(t *testing.T) {
	r, w := io.Pipe()
	const name = "John"
	go func() {
		_ = page.Hero(name).Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the hero exists
	if doc.Find(`[data-testid="hero"]`).Length() != 1 {
		t.Errorf("Expected to find a hero")
	}
}
