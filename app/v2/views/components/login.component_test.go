package components_test

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/memnix/memnix-rest/app/v2/views/components"
)

func TestLoginComponent(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		_ = components.LoginComponent().Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the form exists
	if doc.Find("form").Length() != 1 {
		t.Errorf("Expected to find a form")
	}

	// Assert that the email input exists
	if doc.Find("input[name='email']").Length() != 1 {
		t.Errorf("Expected to find an email input")
	}

	// Assert that the password input exists
	if doc.Find("input[name='password']").Length() != 1 {
		t.Errorf("Expected to find a password input")
	}

	// Assert that the login button exists
	if doc.Find("button").Length() != 1 {
		t.Errorf("Expected to find a login button")
	}

	// Assert that the login error div exists
	if doc.Find("#login-error").Length() != 1 {
		t.Errorf("Expected to find a login error div")
	}

	// Assert that the login error div is hidden
	if !doc.Find(`[data-testid="LoginError"]`).HasClass("hidden") {
		val, ok := doc.Find(`[data-testid="LoginError"]`).Attr("class")
		if !ok {
			t.Errorf("Expected the login error div to have a class, but got none")
		}

		t.Errorf("Expected the login error div to be hidden, but got class: %s", val)
	}
}

func TestLoginError(t *testing.T) {
	r, w := io.Pipe()
	errorMessage := "Invalid email or password"
	go func() {
		_ = components.LoginError(errorMessage).Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the login test is correct
	if doc.Text() != errorMessage {
		t.Errorf("Expected the error message to be %s, got %s", errorMessage, doc.Text())
	}
}
