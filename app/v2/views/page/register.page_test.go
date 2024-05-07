package page_test

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/memnix/memnix-rest/app/v2/views/page"
	"github.com/memnix/memnix-rest/domain"
)

func TestRegisterContent(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		_ = page.RegisterContent().Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the register content exists
	if doc.Find(`[data-testid="registerContent"]`).Length() != 1 {
		t.Errorf("Expected to find a register content")
	}

	// Assert that the register content has a register component
	if doc.Find(`[data-testid="registerComponent"]`).Length() != 1 {
		t.Errorf("Expected to find a register component")
	}

	// Assert that the register with github button exists
	if doc.Find(`[data-testid="githubRegister"]`).Length() != 1 {
		t.Errorf("Expected to find a register with github button")
	}

	// Assert that the register with discord button exists
	if doc.Find(`[data-testid="discordRegister"]`).Length() != 1 {
		t.Errorf("Expected to find a register with discord button")
	}
}

func TestRegisterPage(t *testing.T) {
	const title = "Register"
	r, w := io.Pipe()
	go func() {
		_ = page.RegisterPage(title, domain.Nonce{}, page.RegisterContent()).Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	if doc.Find(`[data-testid="registerContent"]`).Length() != 1 {
		t.Errorf("Expected to find a register content")
	}

	if doc.Find(`[data-testid="githubRegister"]`).Length() != 1 {
		t.Errorf("Expected to find a register with github button")
	}

	if doc.Find(`[data-testid="discordRegister"]`).Length() != 1 {
		t.Errorf("Expected to find a register with discord button")
	}

	if doc.Find("title").Text() != "Register" {
		t.Errorf("Expected to find a title with the text 'Register'")
	}

	// Assert that the header exists
	if doc.Find("header").Length() != 1 {
		t.Errorf("Expected to find a header")
	}

	// Assert that the footer exists
	if doc.Find("footer").Length() != 1 {
		t.Errorf("Expected to find a footer")
	}

	// Assert that the main exists
	if doc.Find("main").Length() != 1 {
		t.Errorf("Expected to find a main")
	}
}
