package test

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type Test struct {
	description string

	// Test input
	route   string
	body    []byte
	reqType string
	// Expected output
	expectedError bool
	expectedCode  int
	expectedBody  string
}

func testIndex() []Test {
	tests := []Test{

		{
			description:   "index",
			route:         "/",
			body:          []byte{},
			reqType:       "GET",
			expectedCode:  fiber.StatusForbidden,
			expectedBody:  "This is not a valid route",
			expectedError: false,
		},
	}

	return tests
}

func TestRoute(t *testing.T) {
	app, _ := Setup()

	tests := testIndex()

	for _, test := range tests {
		req, _ := http.NewRequest(test.reqType, "/v1/"+test.route, bytes.NewBuffer(test.body))
		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}
