package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		operands []float64
		want     float64
		wantErr  error
	}{
		// Basic operations
		{"add integers", Add, []float64{2, 3}, 5, nil},
		{"add negatives", Add, []float64{-2, -3}, -5, nil},
		{"add decimals", Add, []float64{0.5, 0.25}, 0.75, nil},
		{"subtract", Subtract, []float64{10, 4}, 6, nil},
		{"subtract to negative", Subtract, []float64{4, 10}, -6, nil},
		{"multiply", Multiply, []float64{6, 7}, 42, nil},
		{"multiply by zero", Multiply, []float64{123, 0}, 0, nil},
		{"divide", Divide, []float64{10, 4}, 2.5, nil},
		{"divide negative", Divide, []float64{-9, 3}, -3, nil},

		// Advanced operations
		{"power", Power, []float64{2, 10}, 1024, nil},
		{"power fractional exponent", Power, []float64{9, 0.5}, 3, nil},
		{"power negative exponent", Power, []float64{2, -2}, 0.25, nil},
		{"sqrt", Sqrt, []float64{16}, 4, nil},
		{"sqrt zero", Sqrt, []float64{0}, 0, nil},
		{"percentage", Percentage, []float64{15, 200}, 30, nil},
		{"percentage over 100", Percentage, []float64{150, 20}, 30, nil},

		// Edge cases
		{"divide by zero", Divide, []float64{1, 0}, 0, ErrDivisionByZero},
		{"zero divided by zero", Divide, []float64{0, 0}, 0, ErrDivisionByZero},
		{"sqrt negative", Sqrt, []float64{-4}, 0, ErrNegativeSqrt},
		{"overflow", Multiply, []float64{math.MaxFloat64, 10}, 0, ErrNotFinite},
		{"zero to negative power", Power, []float64{0, -1}, 0, ErrNotFinite},
		{"negative base fractional exponent", Power, []float64{-8, 0.5}, 0, ErrNotFinite},
		{"NaN operand", Add, []float64{math.NaN(), 1}, 0, ErrInvalidOperand},
		{"Inf operand", Add, []float64{math.Inf(1), 1}, 0, ErrInvalidOperand},

		// Operand count / unknown operation
		{"binary with one operand", Add, []float64{1}, 0, ErrOperandCount},
		{"binary with no operands", Divide, nil, 0, ErrOperandCount},
		{"unary with two operands", Sqrt, []float64{4, 2}, 0, ErrOperandCount},
		{"unknown operation", Operation("modulo"), []float64{1, 2}, 0, ErrUnknownOperation},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Calculate(tc.op, tc.operands...)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Calculate(%s, %v) = %v, want %v", tc.op, tc.operands, got, tc.want)
			}
		})
	}
}

func TestArity(t *testing.T) {
	for _, op := range Operations() {
		n, err := Arity(op)
		if err != nil {
			t.Fatalf("Arity(%s) returned error: %v", op, err)
		}
		want := 2
		if op == Sqrt {
			want = 1
		}
		if n != want {
			t.Errorf("Arity(%s) = %d, want %d", op, n, want)
		}
	}

	if _, err := Arity("nope"); !errors.Is(err, ErrUnknownOperation) {
		t.Errorf("expected ErrUnknownOperation, got %v", err)
	}
}
