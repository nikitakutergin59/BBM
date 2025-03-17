package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculator(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected float64
		err      error
	}{
		{
			name:     "Simple addition",
			input:    "2 + 3",
			expected: 5,
			err:      nil,
		},
		{
			name:     "Simple subtraction",
			input:    "5 - 2",
			expected: 3,
			err:      nil,
		},
		{
			name:     "Simple multiplication",
			input:    "4 * 6",
			expected: 24,
			err:      nil,
		},
		{
			name:     "Simple division",
			input:    "8 / 2",
			expected: 4,
			err:      nil,
		},
		{
			name:     "Expression with parentheses",
			input:    "(2 + 3) * 4",
			expected: 20,
			err:      nil,
		},
		{
			name:     "Expression with multiple operations",
			input:    "10 + 2 * 3 - 4 / 2",
			expected: 14,
			err:      nil,
		},
		{
			name:     "Complex expression with parentheses",
			input:    "((10 / 2) + 10) * 2",
			expected: 30,
			err:      nil,
		},
		{
			name:     "Division by zero",
			input:    "5 / 0",
			expected: 0,
			err:      errors.New("division by zero"),
		},
		{
			name:     "Mismatched parentheses - opening",
			input:    "(2 + 3",
			expected: 0,
			err:      errors.New("mismatched parentheses"),
		},
		{
			name:     "Mismatched parentheses - closing",
			input:    "2 + 3)",
			expected: 0,
			err:      errors.New("mismatched parentheses"),
		},
		{
			name:     "Invalid character",
			input:    "2 + a",
			expected: 0,
			err:      errors.New("invalid character: a"),
		},
		{
			name:     "Negative numbers",
			input:    "-2 + 5",
			expected: 3,
			err:      nil,
		},
		{
			name:     "Negative number behind +",
			input:    "2 + -3",
			expected: -1,
			err:      nil,
		},
		{
			name:     "Negative number behind -",
			input:    "2 --3",
			expected: 5,
			err:      nil,
		},
		{
			name:     "Negative number bihind *",
			input:    "2 *-3",
			expected: -6,
			err:      nil,
		},
		{
			name:     "Negative nuvber bihind /",
			input:    "2 /-3",
			expected: -0.6666666666666666,
			err:      nil,
		},
		{
			name:     "Float numbers",
			input:    "2.5 + 2.5",
			expected: 5,
			err:      nil,
		},
		{
			name:     "Float number as result",
			input:    "5 / 2",
			expected: 2.5,
			err:      nil,
		},
		{
			name:     "Negative result",
			input:    "2-5",
			expected: -3,
			err:      nil,
		},
		{
			name:     "Very big numbers",
			input:    "1000000000000 + 2000000000000",
			expected: 3000000000000,
			err:      nil,
		},
		{
			name:     "Very big float number",
			input:    "1000000000000.5 + 2000000000000.5",
			expected: 3000000000001,
			err:      nil,
		},
		{
			name:     "Floating point imprecision",
			input:    "0.1+0.2",
			expected: 0.3,
			err:      nil,
		},
		{
			name:     "Empty expression",
			input:    "",
			expected: 0,
			err:      errors.New("invalid expression"),
		},
		{
			name:     "Sqrt",
			input:    "sqrt(25)",
			expected: 5,
			err:      nil,
		},
		{
			name:     "Sqrt in expression",
			input:    "sqrt(4)+2",
			expected: 4,
			err:      nil,
		},
		{
			name:     "invalid sqrt",
			input:    "sqrt(9",
			expected: 0,
			err:      errors.New("missing closing parenthesis"),
		},
		{
			name:     "Zero Value in sqrt",
			input:    "sqrt()",
			expected: 0,
			err:      errors.New("invalid number in sqrt: strconv.ParseFloat: parsing \"\": invalid syntax"),
		},
		{
			name:     "Float %1f",
			input:    "0.1+0.1",
			expected: 0.2,
			err:      nil,
		},
		{
			name:     "Float %2f",
			input:    "0.01+0.01",
			expected: 0.02,
			err:      nil,
		},
		{
			name:     "Float %3f",
			input:    "0.001+0.001",
			expected: 0.002,
			err:      nil,
		},
		{
			name:     "Float %4f",
			input:    "0.0001+0.0001",
			expected: 0.0002,
			err:      nil,
		},
		{
			name:     "Float %5f",
			input:    "0.00001+0.00001",
			expected: 0.00002,
			err:      nil,
		},
		{
			name:     "Float %6f",
			input:    "0.000001+0.000001",
			expected: 0.000002,
			err:      nil,
		},
		{
			name:     "Float %7f",
			input:    "0.0000001+0.0000001",
			expected: 0.0000002,
			err:      nil,
		},
		{
			name:     "Float %8f",
			input:    "0.00000001+0.00000001",
			expected: 0.00000002,
			err:      nil,
		},
		{
			name:     "Float %9f",
			input:    "0.000000001+0.000000001",
			expected: 0.000000002,
			err:      nil,
		},
		{
			name:     "Float %10f",
			input:    "0.0000000001+0.0000000001",
			expected: 0.0000000002,
			err:      nil,
		},
		{
			name:     "Float division %1f",
			input:    "1/10",
			expected: 0.1,
			err:      nil,
		},
		{
			name:     "Float division %2f",
			input:    "1/100",
			expected: 0.01,
			err:      nil,
		},
		{
			name:     "Float division %3f",
			input:    "1/1000",
			expected: 0.001,
			err:      nil,
		},
		{
			name:     "Float division %4f",
			input:    "1/10000",
			expected: 0.0001,
			err:      nil,
		},
		{
			name:     "Float division %5f",
			input:    "1/100000",
			expected: 0.00001,
			err:      nil,
		},
		{
			name:     "Float division %6f",
			input:    "1/1000000",
			expected: 0.000001,
			err:      nil,
		},
		{
			name:     "Float division %7f",
			input:    "1/10000000",
			expected: 0.0000001,
			err:      nil,
		},
		{
			name:     "Float division %8f",
			input:    "1/100000000",
			expected: 0.00000001,
			err:      nil,
		},
		{
			name:     "Float division %9f",
			input:    "1/1000000000",
			expected: 0.000000001,
			err:      nil,
		},
		{
			name:     "Float division %10f",
			input:    "1/10000000000",
			expected: 0.0000000001,
			err:      nil,
		},
		{
			name:     "degree value",
			input:    "3**3",
			expected: 27,
			err:      nil,
		},
		{
			name:     "degree negative value",
			input:    "2**(-3)",
			expected: 0.125,
			err:      nil,
		},
		{
			name:     "degree onli negative value",
			input:    "-2**(-3)",
			expected: -0.125,
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Calc(tc.input)
			if tc.err != nil {
				if err == nil || err.Error() != tc.err.Error() {
					t.Errorf("expected error '%v', got '%v'", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(actual-tc.expected) > 1e-10 {
				t.Errorf("expected %f, got %f", tc.expected, actual)
			}
		})
	}
}
