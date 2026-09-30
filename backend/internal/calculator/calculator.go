// Package calculator contains the pure arithmetic domain logic.
//
// It has no knowledge of HTTP or JSON, which keeps it trivial to test and
// reuse from any transport (REST, CLI, gRPC...).
package calculator

import (
	"errors"
	"fmt"
	"math"
)

// Operation identifies a supported arithmetic operation.
type Operation string

const (
	Add        Operation = "add"
	Subtract   Operation = "subtract"
	Multiply   Operation = "multiply"
	Divide     Operation = "divide"
	Power      Operation = "power"
	Sqrt       Operation = "sqrt"
	Percentage Operation = "percentage"
)

// Domain errors. Callers use errors.Is to map them to transport-level codes.
var (
	ErrUnknownOperation = errors.New("unknown operation")
	ErrOperandCount     = errors.New("wrong number of operands")
	ErrDivisionByZero   = errors.New("division by zero")
	ErrNegativeSqrt     = errors.New("square root of a negative number")
	ErrInvalidOperand   = errors.New("operand must be a finite number")
	ErrNotFinite        = errors.New("result is not a finite number")
)

// spec describes how many operands an operation needs and how to compute it.
type spec struct {
	arity int
	fn    func(a, b float64) (float64, error)
}

var operations = map[Operation]spec{
	Add:      {2, func(a, b float64) (float64, error) { return a + b, nil }},
	Subtract: {2, func(a, b float64) (float64, error) { return a - b, nil }},
	Multiply: {2, func(a, b float64) (float64, error) { return a * b, nil }},
	Divide: {2, func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	}},
	Power: {2, func(a, b float64) (float64, error) { return math.Pow(a, b), nil }},
	Sqrt: {1, func(a, _ float64) (float64, error) {
		if a < 0 {
			return 0, ErrNegativeSqrt
		}
		return math.Sqrt(a), nil
	}},
	// Percentage returns a% of b (e.g. 15% of 200 = 30).
	Percentage: {2, func(a, b float64) (float64, error) { return a / 100 * b, nil }},
}

// Arity returns the number of operands an operation requires.
func Arity(op Operation) (int, error) {
	s, ok := operations[op]
	if !ok {
		return 0, fmt.Errorf("%w: %q", ErrUnknownOperation, op)
	}
	return s.arity, nil
}

// Operations lists every supported operation in a stable order.
func Operations() []Operation {
	return []Operation{Add, Subtract, Multiply, Divide, Power, Sqrt, Percentage}
}

// Calculate applies op to the given operands.
//
// Binary operations need exactly two operands; unary ones (sqrt) need one.
// Extra operands are rejected rather than silently ignored so client bugs
// surface early.
func Calculate(op Operation, operands ...float64) (float64, error) {
	s, ok := operations[op]
	if !ok {
		return 0, fmt.Errorf("%w: %q", ErrUnknownOperation, op)
	}
	if len(operands) != s.arity {
		return 0, fmt.Errorf("%w: %s expects %d operand(s), got %d",
			ErrOperandCount, op, s.arity, len(operands))
	}
	for _, v := range operands {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, ErrInvalidOperand
		}
	}

	var a, b float64
	a = operands[0]
	if s.arity == 2 {
		b = operands[1]
	}

	result, err := s.fn(a, b)
	if err != nil {
		return 0, err
	}
	// Guard against overflow (1e308 * 10) or undefined results (0^-1).
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrNotFinite
	}
	return result, nil
}
