package page_test

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/memnix/memnix-rest/app/v2/views/page"
)

func TestLoginContent(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		_ = page.LoginContent().Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the login content exists
	if doc.Find(`[data-testid="loginContent"]`).Length() != 1 {
		t.Errorf("Expected to find a login content")
	}

	// Assert that the login content has a login component
	if doc.Find(`[data-testid="loginComponent"]`).Length() != 1 {
		t.Errorf("Expected to find a login component")
	}

	// Assert that the login with github button exists
	if doc.Find(`[data-testid="githubLogin"]`).Length() != 1 {
		t.Errorf("Expected to find a login with github button")
	}

	// Assert that the login with discord button exists
	if doc.Find(`[data-testid="discordLogin"]`).Length() != 1 {
		t.Errorf("Expected to find a login with discord button")
	}
}

func TestLoginPage(t *testing.T) {
	r, w := io.Pipe()
	const title = "Login"
	const isError = false
	errMsgs := make([]string, 0)
	sucMsgs := make([]string, 0)
	go func() {
		_ = page.LoginPage(title, isError, errMsgs, sucMsgs, page.LoginContent()).Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the title exists
	if doc.Find("title").Text() != title {
		t.Errorf("Expected to find a title")
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

	// Assert that the login content exists
	if doc.Find(`[data-testid="loginContent"]`).Length() != 1 {
		t.Errorf("Expected to find a login content")
	}

	// Assert that the login content has a login component
	if doc.Find(`[data-testid="loginComponent"]`).Length() != 1 {
		t.Errorf("Expected to find a login component")
	}

	// Assert that the login with github button exists
	if doc.Find(`[data-testid="githubLogin"]`).Length() != 1 {
		t.Errorf("Expected to find a login with github button")
	}

	// Assert that the login with discord button exists
	if doc.Find(`[data-testid="discordLogin"]`).Length() != 1 {
		t.Errorf("Expected to find a login with discord button")
	}
}
