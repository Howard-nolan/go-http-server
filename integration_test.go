//go:build integration

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	db "github.com/joeynolan/go-http-server/internal/db"
	apphttp "github.com/joeynolan/go-http-server/internal/http"
	"github.com/joeynolan/go-http-server/internal/http/handlers"
	ilog "github.com/joeynolan/go-http-server/internal/platform/log"
)

func startIntegrationServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	sqlDB, err := db.OpenAndMigrate(":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}

	logger := ilog.New()

	r := chi.NewRouter()
	r.Use(apphttp.MetricsMiddleware)
	r.Use(apphttp.RequestLogger(logger.Desugar()))
	h := handlers.NewHandler(sqlDB, logger)
	apphttp.Register(r, h)

	srv := httptest.NewServer(r)

	cleanup := func() {
		srv.Close()
		_ = sqlDB.Close()
		logger.Sync()
	}

	return srv, cleanup
}

func TestIntegration_ShortenAndRedirect(t *testing.T) {
	srv, cleanup := startIntegrationServer(t)
	defer cleanup()

	shortenReq := `{"url":"example.com"}`
	res, err := http.Post(srv.URL+"/v1/shorten", "application/json", strings.NewReader(shortenReq))
	if err != nil {
		t.Fatalf("POST /v1/shorten: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusCreated)
	}

	var shortenResp map[string]string
	if err := json.NewDecoder(res.Body).Decode(&shortenResp); err != nil {
		t.Fatalf("decode shorten response: %v", err)
	}
	shortURL := shortenResp["short"]
	if !strings.HasPrefix(shortURL, "https://short.example/") {
		t.Fatalf("short url = %q, want https://short.example/...", shortURL)
	}
	code := strings.TrimPrefix(shortURL, "https://short.example/")
	if code == "" {
		t.Fatalf("no generated code")
	}

	client := &http.Client{
		// do not follow redirects so we can assert Location
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	redirectRes, err := client.Get(srv.URL + "/v1/r/" + code)
	if err != nil {
		t.Fatalf("GET /v1/r/%s: %v", code, err)
	}
	defer redirectRes.Body.Close()

	if redirectRes.StatusCode != http.StatusFound {
		t.Fatalf("status = %d, want %d", redirectRes.StatusCode, http.StatusFound)
	}
	if loc := redirectRes.Header.Get("Location"); loc != "https://example.com" {
		t.Fatalf("location = %q, want %q", loc, "https://example.com")
	}
}
