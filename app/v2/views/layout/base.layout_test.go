package layout_test

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/memnix/memnix-rest/app/v2/views/layout"
)

type BaseArgs struct {
	title         string
	username      string
	errMsgs       []string
	sucMsgs       []string
	fromProtected bool
	isError       bool
}

func TestBase(t *testing.T) {
	r, w := io.Pipe()
	args := BaseArgs{
		title:         "Memnix",
		username:      "",
		fromProtected: false,
		isError:       false,
		errMsgs:       nil,
		sucMsgs:       nil,
	}
	go func() {
		_ = layout.Base(args.title, args.username, args.fromProtected, args.isError, args.errMsgs, args.sucMsgs).Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("Error reading document: %s", err)
	}

	// Assert that the title exists
	if doc.Find("title").Text() != args.title {
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
}
