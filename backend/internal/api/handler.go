// Package api exposes the calculator domain over HTTP/JSON.
package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"calculator/internal/calculator"
)

const maxBodyBytes = 1 << 10 // 1 KiB is plenty for a calculation request.

// CalculateRequest is the JSON body for POST /api/v1/calculate.
// Pointers let us tell "missing" apart from a legitimate 0.
type CalculateRequest struct {
	Operation string   `json:"operation"`
	A         *float64 `json:"a"`
	B         *float64 `json:"b,omitempty"`
}

// CalculateResponse is returned on success.
type CalculateResponse struct {
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
	Result    float64   `json:"result"`
}

// OperationInfo describes a supported operation for GET /api/v1/operations.
type OperationInfo struct {
	Name  string `json:"name"`
	Arity int    `json:"arity"`
}

// ErrorBody is the envelope for every error response.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail carries a stable machine-readable code plus a human message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewRouter wires all API routes. Uses the Go 1.22+ method-aware ServeMux,
// so no third-party router is needed.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/calculate", handleCalculate)
	mux.HandleFunc("GET /api/v1/operations", handleOperations)
	mux.HandleFunc("GET /healthz", handleHealth)
	return mux
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req CalculateRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must be valid JSON with numeric operands")
		return
	}
	if req.Operation == "" {
		writeError(w, http.StatusBadRequest, "MISSING_OPERATION", "field 'operation' is required")
		return
	}
	if req.A == nil {
		writeError(w, http.StatusBadRequest, "MISSING_OPERAND", "field 'a' is required")
		return
	}

	operands := []float64{*req.A}
	if req.B != nil {
		operands = append(operands, *req.B)
	}

	op := calculator.Operation(req.Operation)
	result, err := calculator.Calculate(op, operands...)
	if err != nil {
		status, code := mapError(err)
		writeError(w, status, code, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, CalculateResponse{
		Operation: req.Operation,
		Operands:  operands,
		Result:    result,
	})
}

func handleOperations(w http.ResponseWriter, _ *http.Request) {
	ops := calculator.Operations()
	out := make([]OperationInfo, 0, len(ops))
	for _, op := range ops {
		arity, _ := calculator.Arity(op) // ops come from the domain, so always valid
		out = append(out, OperationInfo{Name: string(op), Arity: arity})
	}
	writeJSON(w, http.StatusOK, out)
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// mapError translates domain errors into HTTP status + stable error code.
//
//   - 400 Bad Request: the request itself is malformed (client bug).
//   - 422 Unprocessable Entity: the request is well-formed but the math
//     is undefined (division by zero, sqrt of a negative, overflow).
func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, calculator.ErrUnknownOperation):
		return http.StatusBadRequest, "UNKNOWN_OPERATION"
	case errors.Is(err, calculator.ErrOperandCount):
		return http.StatusBadRequest, "INVALID_OPERAND_COUNT"
	case errors.Is(err, calculator.ErrInvalidOperand):
		return http.StatusBadRequest, "INVALID_OPERAND"
	case errors.Is(err, calculator.ErrDivisionByZero):
		return http.StatusUnprocessableEntity, "DIVISION_BY_ZERO"
	case errors.Is(err, calculator.ErrNegativeSqrt):
		return http.StatusUnprocessableEntity, "NEGATIVE_SQRT"
	case errors.Is(err, calculator.ErrNotFinite):
		return http.StatusUnprocessableEntity, "RESULT_NOT_FINITE"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR"
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, ErrorBody{Error: ErrorDetail{Code: code, Message: msg}})
}
