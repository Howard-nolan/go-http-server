package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ilog "github.com/joeynolan/go-http-server/internal/platform/log"
)

const testBaseURL = "https://jnshorter.com"

func newTestHandler(t *testing.T) (*Handler, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	return NewHandler(db, ilog.New(), testBaseURL), mock, func() { _ = db.Close() }
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type = %q, want application/json; charset=utf-8", ct)
	}
	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil || body["status"] != "ok" {
		t.Fatalf("body = %+v, want status ok", body)
	}
}

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name       string
		mock       func(sqlmock.Sqlmock)
		wantStatus int
		wantMsg    string
	}{
		{
			name: "ready",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectPing().WillReturnError(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "db unavailable",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectPing().WillReturnError(errors.New("down"))
			},
			wantStatus: http.StatusServiceUnavailable,
			wantMsg:    "Database unavailable",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, mock, cleanup := newTestHandler(t)
			defer cleanup()
			if tc.mock != nil {
				tc.mock(mock)
			}

			req := httptest.NewRequest(http.MethodGet, "/ready", nil)
			w := httptest.NewRecorder()
			h.ReadyHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantStatus {
				t.Fatalf("status = %d, want %d", res.StatusCode, tc.wantStatus)
			}
			if ct := res.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
				t.Fatalf("content-type = %q, want application/json; charset=utf-8", ct)
			}

			if tc.wantMsg != "" {
				var body ErrorResponse
				_ = json.NewDecoder(res.Body).Decode(&body)
				if body.Message != tc.wantMsg {
					t.Fatalf("msg = %q, want %q", body.Message, tc.wantMsg)
				}
			} else {
				var body map[string]string
				if err := json.NewDecoder(res.Body).Decode(&body); err != nil || body["status"] != "ok" {
					t.Fatalf("body = %+v, want status ok", body)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("sql expectations: %v", err)
			}
		})
	}
}

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mock       func(sqlmock.Sqlmock)
		modifyReq  func(*http.Request) *http.Request
		wantStatus int
		wantMsg    string
		wantShort  bool
		wantCode   string
	}{
		{
			name: "happy path",
			body: `{"url":"https://example.com"}`,
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(`INSERT INTO links`).WithArgs(sqlmock.AnyArg(), "https://example.com", "").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantStatus: http.StatusCreated,
			wantShort:  true,
		},
		{name: "bad json", body: `{"url":`, wantStatus: http.StatusBadRequest, wantMsg: "invalid JSON"},
		{name: "missing url", body: `{}`, wantStatus: http.StatusBadRequest, wantMsg: "url is required"},
		{
			name: "db keeps failing",
			body: `{"url":"https://example.com"}`,
			mock: func(m sqlmock.Sqlmock) {
				for i := 0; i < 3; i++ {
					m.ExpectExec(`INSERT INTO links`).WithArgs(sqlmock.AnyArg(), "https://example.com", "").
						WillReturnError(errors.New("boom"))
				}
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "failed to generate unique code",
		},
		{
			name: "db timeout",
			body: `{"url":"https://example.com"}`,
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(`INSERT INTO links`).WithArgs(sqlmock.AnyArg(), "https://example.com", "").
					WillReturnError(context.DeadlineExceeded)
			},
			wantStatus: http.StatusRequestTimeout,
			wantMsg:    "timeout",
		},
		{
			name: "idempotent duplicate",
			body: `{"url":"https://example.com"}`,
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(`INSERT INTO links`).WithArgs(sqlmock.AnyArg(), "https://example.com", "req-123").
					WillReturnResult(sqlmock.NewResult(1, 0))
				m.ExpectQuery(`SELECT code FROM links WHERE idempotency_key = \?`).
					WithArgs("req-123").
					WillReturnRows(sqlmock.NewRows([]string{"code"}).AddRow("existing"))
			},
			modifyReq: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), middleware.RequestIDKey, "req-123")
				return r.WithContext(ctx)
			},
			wantStatus: http.StatusCreated,
			wantShort:  true,
			wantCode:   "existing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, mock, cleanup := newTestHandler(t)
			defer cleanup()
			if tc.mock != nil {
				tc.mock(mock)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/shorten", strings.NewReader(tc.body))
			if tc.modifyReq != nil {
				req = tc.modifyReq(req)
			}
			w := httptest.NewRecorder()
			h.ShortenHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantStatus {
				t.Fatalf("status = %d, want %d", res.StatusCode, tc.wantStatus)
			}

			if tc.wantShort {
				var body map[string]string
				_ = json.NewDecoder(res.Body).Decode(&body)
				short := body["short"]
				if tc.wantCode != "" {
					if short != fmt.Sprintf("%s/%s", testBaseURL, tc.wantCode) {
						t.Fatalf("short = %q, want %s/%s", short, testBaseURL, tc.wantCode)
					}
				} else if !strings.HasPrefix(short, testBaseURL+"/") {
					t.Fatalf("short = %q, want %s/...", short, testBaseURL)
				}
			}
			if tc.wantMsg != "" {
				var body ErrorResponse
				_ = json.NewDecoder(res.Body).Decode(&body)
				if body.Message != tc.wantMsg {
					t.Fatalf("msg = %q, want %q", body.Message, tc.wantMsg)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("sql expectations: %v", err)
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		seedCache  map[string]string
		mock       func(sqlmock.Sqlmock)
		wantStatus int
		wantLoc    string
		wantMsg    string
	}{
		{name: "missing code", code: "", wantStatus: http.StatusBadRequest, wantMsg: "missing code"},
		{
			name:       "cache hit",
			code:       "abc123",
			seedCache:  map[string]string{"abc123": "https://example.com"},
			wantStatus: http.StatusFound,
			wantLoc:    "https://example.com",
		},
		{
			name: "db hit adds https",
			code: "abc123",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT url FROM links WHERE code = \?`).
					WithArgs("abc123").
					WillReturnRows(sqlmock.NewRows([]string{"url"}).AddRow("example.com"))
			},
			wantStatus: http.StatusFound,
			wantLoc:    "https://example.com",
		},
		{
			name: "not found",
			code: "missing",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT url FROM links WHERE code = \?`).
					WithArgs("missing").
					WillReturnError(sql.ErrNoRows)
			},
			wantStatus: http.StatusNotFound,
			wantMsg:    "code not found",
		},
		{
			name: "db error",
			code: "oops",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT url FROM links WHERE code = \?`).
					WithArgs("oops").
					WillReturnError(errors.New("db down"))
			},
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "lookup failed",
		},
		{
			name: "db timeout",
			code: "slow",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT url FROM links WHERE code = \?`).
					WithArgs("slow").
					WillReturnError(context.DeadlineExceeded)
			},
			wantStatus: http.StatusRequestTimeout,
			wantMsg:    "timeout",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, mock, cleanup := newTestHandler(t)
			defer cleanup()
			if tc.mock != nil {
				tc.mock(mock)
			}
			for k, v := range tc.seedCache {
				h.cache.Add(k, v)
			}

			req := httptest.NewRequest(http.MethodGet, "/v1/r/"+tc.code, nil)
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("code", tc.code)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			w := httptest.NewRecorder()
			h.RedirectHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantStatus {
				t.Fatalf("status = %d, want %d", res.StatusCode, tc.wantStatus)
			}
			if tc.wantLoc != "" {
				if loc := res.Header.Get("Location"); loc != tc.wantLoc {
					t.Fatalf("location = %q, want %q", loc, tc.wantLoc)
				}
			}
			if tc.wantMsg != "" {
				var body ErrorResponse
				_ = json.NewDecoder(res.Body).Decode(&body)
				if body.Message != tc.wantMsg {
					t.Fatalf("msg = %q, want %q", body.Message, tc.wantMsg)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("sql expectations: %v", err)
			}
		})
	}
}
