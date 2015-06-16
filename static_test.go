package station

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testFile = "testdata/bar"

func TestStatic(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", testFile, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := Static(".", StaticOptions{})(http.HandlerFunc(testHandler))
	handler.ServeHTTP(rr, r)

	err = testBody(t, rr, r, testFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStaticListDir(t *testing.T) {
	t.Parallel()

	opts := StaticOptions{
		ListDir: true,
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/testdata", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := Static(".", opts)(http.HandlerFunc(testHandler))
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Fatalf("was not able to list dir: %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "bar") {
		t.Fatalf("directory listing did not contain %s, got %q", testFile,
			rr.Body.String())
	}
}

func TestStaticListFalse(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/testdata", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := Static(".", StaticOptions{})(http.HandlerFunc(testHandler))
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Fatalf("wrapped handler was not called: got status %v", rr.Code)
	}
}

func TestStaticNoFile(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", testFile+".txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := Static(".", StaticOptions{})(http.HandlerFunc(testHandler))
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Fatalf("wrapped handler was not called: got status %v", rr.Code)
	}
}

func TestServe(t *testing.T) {
	t.Parallel()

	opts := StaticOptions{NotFoundHandler: http.NotFoundHandler()}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", testFile, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := Serve(".", opts)
	handler.ServeHTTP(rr, r)

	err = testBody(t, rr, r, testFile)
	if err != nil {
		t.Fatal(err)
	}
}

func testBody(t *testing.T, rr *httptest.ResponseRecorder, r *http.Request,
	fname string) error {
	if rr.Code != http.StatusOK {
		fmt.Errorf("was not able to serve static file at %v", r.URL)
	}

	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("test file '%v' does not exist", fname)
	}

	if !bytes.Equal(rr.Body.Bytes(), file) {
		return fmt.Errorf("body mismatch: got %q, expected %q", rr.Body.Bytes(), file)
	}

	return nil
}

func testHandler(w http.ResponseWriter, r *http.Request) {
}
