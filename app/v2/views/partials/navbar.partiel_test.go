package partials_test

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/memnix/memnix-rest/app/v2/views/partials"
)

func TestNavbar(t *testing.T) {
	r, w := io.Pipe()
	username, fromProtected := "", false
	go func() {
		_ = partials.Navbar(username, fromProtected).Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the navbar exists
	if doc.Find(`[data-testid="navbar"]`).Length() != 1 {
		t.Errorf("Expected to find a navbar")
	}

	// Assert that the navbar has a login button
	if doc.Find(`[data-testid="loginButton"]`).Length() != 1 {
		t.Errorf("Expected to find a login button")
	}

	// Assert that the navbar has a navbar start
	if doc.Find(`[data-testid="navbarStart"]`).Length() != 1 {
		t.Errorf("Expected to find a navbar start")
	}

	// Assert that the navbar has a navbar end
	if doc.Find(`[data-testid="navbarEnd"]`).Length() != 1 {
		t.Errorf("Expected to find a navbar end")
	}

	// Assert that the navbar has a navbar center
	if doc.Find(`[data-testid="navbarCenter"]`).Length() != 1 {
		t.Errorf("Expected to find a navbar center")
	}

	// Assert that the navbar does not have a username
	if doc.Find(`[data-testid="username"]`).Length() != 0 {
		t.Errorf("Expected to not find a username")
	}

	// Assert that the navbar does not have a logout button
	if doc.Find(`[data-testid="logoutButton"]`).Length() != 0 {
		t.Errorf("Expected to not find a logout button")
	}
}

func TestNavbarProtected(t *testing.T) {
	r, w := io.Pipe()
	username, fromProtected := "Test", true
	go func() {
		_ = partials.Navbar(username, fromProtected).Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the navbar exists
	if doc.Find(`[data-testid="navbar"]`).Length() != 1 {
		t.Errorf("Expected to find a navbar")
	}

	// Assert that the navbar has a logout button
	if doc.Find(`[data-testid="logoutButton"]`).Length() != 1 {
		t.Errorf("Expected to find a logout button")
	}

	// Assert that the navbar has a username
	if doc.Find(`[data-testid="username"]`).Length() != 1 {
		t.Errorf("Expected to find a username")
	}
}
