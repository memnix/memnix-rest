package components_test

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/memnix/memnix-rest/app/v2/views/components"
)

func TestRegisterComponent(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		_ = components.RegisterComponent().Render(context.Background(), w)
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

	// Assert that the username input exists
	if doc.Find("input[name='username']").Length() != 1 {
		t.Errorf("Expected to find a username input")
	}

	// Assert that the register button exists
	if doc.Find("button").Length() != 1 {
		t.Errorf("Expected to find a register button")
	}

	// Assert that the register error div exists
	if doc.Find("#register-error").Length() != 1 {
		t.Errorf("Expected to find a register error div")
	}

	// Assert that the login error div is hidden
	if !doc.Find(`[data-testid="registerError"]`).HasClass("hidden") {
		val, ok := doc.Find(`[data-testid="registerError"]`).Attr("class")
		if !ok {
			t.Errorf("Expected the login error div to have a class, but got none")
		}

		t.Errorf("Expected the login error div to be hidden, but got class: %s", val)
	}
}

func TestRegisterError(t *testing.T) {
	r, w := io.Pipe()
	errorMessage := "Invalid email or password"
	go func() {
		_ = components.RegisterError(errorMessage).Render(context.Background(), w)
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
