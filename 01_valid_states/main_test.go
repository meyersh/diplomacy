package main // import github.com/meyersh/diplomacy/01_valid_states

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDuplicateInRange(t *testing.T) {
	want := true
	if got := duplicateInRange("11"); got != want {
		t.Errorf("duplicateInRange(\"11\") = %t, want %t", got, want)
	}
}

// func validateRow(board string, row int) bool
func TestValidateRow(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}

	tests := []test{
		{input: "123456789........................................................................", expected: true},
		{input: "113456789........................................................................", expected: false},
	}

	for _, tc := range tests {
		if got := validateRow(tc.input, 0); got != tc.expected {
			t.Errorf("validateRow(\"%s\", 0) = %t, want %t", tc.input, got, tc.expected)
		}
	}
}

func TestValidateCol(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}

	tests := []test{
		{input: "123456789........................................................................", expected: true},
		{input: "113456789........................................................................", expected: true},
		{input: "1234567891.......................................................................", expected: false},
		{input: "123456789.1......................................................................", expected: true},
	}

	for _, tc := range tests {
		if got := validateCol(tc.input, 0); got != tc.expected {
			t.Errorf("validateCol(\"%s\", 0) = %t, want %t", tc.input, got, tc.expected)
		}
	}
}

func TestValidateSquare(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}

	tests := []test{
		{input: "123456789........................................................................", expected: true},
		{input: "113456789........................................................................", expected: false},
		{input: "1234567891.......................................................................", expected: false},
		{input: "123456789.1......................................................................", expected: false},
	}

	for _, tc := range tests {
		if got := validateSquare(tc.input, 0); got != tc.expected {
			t.Errorf("validateSquare(\"%s\", 0) = %t, want %t", tc.input, got, tc.expected)
		}
	}
}

// func hello(w http.ResponseWriter, req *http.Request) {
func TestHello(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/hello", nil)
		response := httptest.NewRecorder()

		hello(response, request)

		got := response.Body.String()
		want := "hello\n"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestSudokuValidator(t *testing.T) {
	type test struct {
		input    string
		method   string
		expected int
	}
	tests := []test{
		{input: "board=123456789........................................................................", expected: 200, method: http.MethodPost},
		{input: "board=113456789........................................................................", expected: 405, method: http.MethodPost},
		{input: "board=1234567891.......................................................................", expected: 405, method: http.MethodPost},
		{input: "board=123456789.1......................................................................", expected: 405, method: http.MethodPost},
		{input: "board=123456789.1.....................................................................", expected: 422, method: http.MethodPost},
		{input: "", expected: 405, method: http.MethodGet},
		{input: "brrrd=123456789.1.....................................................................", expected: 422, method: http.MethodPost},
	}

	t.Run("Sudoku Validator", func(t *testing.T) {
		for _, tc := range tests {
			request, _ := http.NewRequest(tc.method, "/sudokuValidator", strings.NewReader(tc.input))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
			response := httptest.NewRecorder()

			sudokuValidator(response, request)

			got := response.Code
			want := tc.expected

			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		}
	})
}
