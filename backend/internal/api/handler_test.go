package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"calculator/internal/calculator"
)

func doRequest(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, req)
	return rec
}

func TestCalculate_Success(t *testing.T) {
	tests := []struct {
		name string
		body string
		want float64
	}{
		{"add", `{"operation":"add","a":2,"b":3}`, 5},
		{"subtract", `{"operation":"subtract","a":2,"b":3}`, -1},
		{"multiply", `{"operation":"multiply","a":2.5,"b":4}`, 10},
		{"divide", `{"operation":"divide","a":7,"b":2}`, 3.5},
		{"power", `{"operation":"power","a":2,"b":8}`, 256},
		{"sqrt", `{"operation":"sqrt","a":81}`, 9},
		{"percentage", `{"operation":"percentage","a":20,"b":50}`, 10},
		{"zero is a valid operand", `{"operation":"add","a":0,"b":0}`, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := doRequest(t, http.MethodPost, "/api/v1/calculate", tc.body)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
			var resp CalculateResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if resp.Result != tc.want {
				t.Errorf("result = %v, want %v", resp.Result, tc.want)
			}
		})
	}
}

func TestCalculate_Errors(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{"division by zero", `{"operation":"divide","a":1,"b":0}`, 422, "DIVISION_BY_ZERO"},
		{"negative sqrt", `{"operation":"sqrt","a":-1}`, 422, "NEGATIVE_SQRT"},
		{"overflow", `{"operation":"power","a":10,"b":400}`, 422, "RESULT_NOT_FINITE"},
		{"unknown operation", `{"operation":"modulo","a":1,"b":2}`, 400, "UNKNOWN_OPERATION"},
		{"missing b for binary op", `{"operation":"add","a":1}`, 400, "INVALID_OPERAND_COUNT"},
		{"extra b for unary op", `{"operation":"sqrt","a":4,"b":2}`, 400, "INVALID_OPERAND_COUNT"},
		{"missing a", `{"operation":"add","b":1}`, 400, "MISSING_OPERAND"},
		{"missing operation", `{"a":1,"b":2}`, 400, "MISSING_OPERATION"},
		{"string operand", `{"operation":"add","a":"1","b":2}`, 400, "INVALID_JSON"},
		{"malformed json", `{"operation":`, 400, "INVALID_JSON"},
		{"empty body", ``, 400, "INVALID_JSON"},
		{"unknown field", `{"operation":"add","a":1,"b":2,"c":3}`, 400, "INVALID_JSON"},
		{"body too large", `{"operation":"add","a":1,"b":2,"x":"` + strings.Repeat("a", 2048) + `"}`, 400, "INVALID_JSON"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := doRequest(t, http.MethodPost, "/api/v1/calculate", tc.body)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body)
			}
			var resp ErrorBody
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if resp.Error.Code != tc.wantCode {
				t.Errorf("code = %q, want %q", resp.Error.Code, tc.wantCode)
			}
			if resp.Error.Message == "" {
				t.Error("expected a non-empty error message")
			}
		})
	}
}

func TestCalculate_MethodNotAllowed(t *testing.T) {
	rec := doRequest(t, http.MethodGet, "/api/v1/calculate", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}

func TestOperations(t *testing.T) {
	rec := doRequest(t, http.MethodGet, "/api/v1/operations", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var ops []OperationInfo
	if err := json.NewDecoder(rec.Body).Decode(&ops); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(ops) != 7 {
		t.Errorf("got %d operations, want 7", len(ops))
	}
}

func TestHealth(t *testing.T) {
	rec := doRequest(t, http.MethodGet, "/healthz", "")
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestMapError(t *testing.T) {
	tests := []struct {
		err        error
		wantStatus int
		wantCode   string
	}{
		{calculator.ErrUnknownOperation, 400, "UNKNOWN_OPERATION"},
		{calculator.ErrOperandCount, 400, "INVALID_OPERAND_COUNT"},
		{calculator.ErrInvalidOperand, 400, "INVALID_OPERAND"},
		{calculator.ErrDivisionByZero, 422, "DIVISION_BY_ZERO"},
		{calculator.ErrNegativeSqrt, 422, "NEGATIVE_SQRT"},
		{calculator.ErrNotFinite, 422, "RESULT_NOT_FINITE"},
		{errors.New("unexpected"), 500, "INTERNAL_ERROR"},
	}
	for _, tc := range tests {
		t.Run(tc.wantCode, func(t *testing.T) {
			status, code := mapError(fmt.Errorf("wrapped: %w", tc.err))
			if status != tc.wantStatus || code != tc.wantCode {
				t.Errorf("got (%d, %s), want (%d, %s)", status, code, tc.wantStatus, tc.wantCode)
			}
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	h := Logging(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want 418", rec.Code)
	}
}
