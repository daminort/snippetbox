package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"snippetbox.demien.net/internal/models/mocks"
)

type TestServer struct {
	*httptest.Server
}

type TestResponse struct {
	status  int
	headers http.Header
	cookies []*http.Cookie
	body    string
}

func newTestApplication(t *testing.T) *Application {

	// template cache
	templateCache, err := newTemplatesCache()
	if err != nil {
		t.Fatal(err)
	}

	// form decoder
	formDecoder := form.NewDecoder()

	// session manager
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &Application{
		logger:         slog.New(slog.DiscardHandler),
		snippets:       &mocks.SnippetsService{},
		users:          &mocks.UsersService{},
		templateCache:  templateCache,
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
	}
}

func newTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{ts}
}

func extractCSRFToken(t *testing.T, body string) string {
	rx := regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)
	matches := rx.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("CSRF token not found")
	}

	return html.UnescapeString(matches[1])
}

func (ts *TestServer) Get(t *testing.T, urlPath string) *TestResponse {
	req, err := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return &TestResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}

func (ts *TestServer) ResetCookies(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar
}

func (ts *TestServer) PostForm(t *testing.T, urlPath string, form url.Values) *TestResponse {
	req, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return &TestResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}
