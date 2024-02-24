package page_test

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/memnix/memnix-rest/app/v2/views/page"
	"github.com/memnix/memnix-rest/pkg/i18n"
)

var localizer = i18n.MockLocalizer("../../../../locales/en.json")

func TestHero(t *testing.T) {
	r, w := io.Pipe()
	const name = "John"
	go func() {
		_ = page.Hero(name, localizer).Render(context.Background(), w)
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

	// Assert that the hero has a hello component
	if doc.Find(`[data-testid="helloComponent"]`).Length() != 1 {
		t.Errorf("Expected to find a hello component")
	}

	// Assert that the h1 is the name
	if doc.Find(`[data-testid="helloH1"]`).Text() != "Welcome back on Memnix , "+name+" !" {
		t.Errorf("Expected to find a h1 with the name: %s, but got %s", name, doc.Find(`[data-testid="helloH1"]`).Text())
	}

	// Assert that there is only one h1
	if doc.Find(`[data-testid="helloH1"]`).Length() != 1 {
		t.Errorf("Expected to find only one h1")
	}
}
