package bezu

import (
	"errors"
	"strings"
	"testing"
)

func TestBezu(t *testing.T) {
	testCase := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:  "Usual equation",
			input: "1 7 -4 -28",
			expected: `Уравнение: 1x³+7x²-4x-28=0
Корни:
Корень 1: -7
Корень 2: -2
Корень 3: 2`,
			err: nil,
		},
		{
			name:  "Negative cf",
			input: "-1 -7 -4 -28",
			expected: `Уравнение: -1x³-7x²-4x-28=0
Корни:
Корень 1: -7`,
			err: nil,
		},
		{
			name:  "Fractional values cf",
			input: "0.052 1.488 -0.052 -1.488",
			expected: `Уравнение: 0.052x³+1.488x²-0.052x-1.488=0
Корни:
Корень 1: -28.615
Корень 2: -1
Корень 3: 1`,
			err: nil,
		},
		{
			name:     "Total Zero",
			input:    "0 0 0 0",
			expected: "ошибка вычисления: ошибка решения кубического уравнения: не удалось найти вещественный корень",
			err:      errors.New("ошибка вычисления: ошибка решения кубического уравнения: не удалось найти вещественный корень"),
		},
		{
			name:     "Lots of coefficients",
			input:    "1 7 -4 -28 11564515",
			expected: "неверное количество аргументов.  Используйте: a b c d",
			err:      errors.New("неверное количество аргументов.  Используйте: a b c d"),
		},
		{
			name:     "Error Parse Float",
			input:    "jqevjn ndwjkvnqoewp qnwvenqdowvq qkodeevnqokdvnqop",
			expected: "ошибка преобразования коэффициента 1: strconv.ParseFloat: parsing \"jqevjn\": invalid syntax",
			err:      errors.New("ошибка преобразования коэффициента 1: strconv.ParseFloat: parsing \"jqevjn\": invalid syntax"),
		},
		{
			name:     "Error Parse Float_1",
			input:    "jnvaonvadov 7 -4 -28",
			expected: "ошибка преобразования коэффициента 1: strconv.ParseFloat: parsing \"jnvaonvadov\": invalid syntax",
			err:      errors.New("ошибка преобразования коэффициента 1: strconv.ParseFloat: parsing \"jnvaonvadov\": invalid syntax"),
		},
		{
			name:     "Eror Parse Float_2",
			input:    "1 mwkmefok -4 -28",
			expected: "ошибка преобразования коэффициента 2: strconv.ParseFloat: parsing \"mwkmefok\": invalid syntax",
			err:      errors.New("ошибка преобразования коэффициента 2: strconv.ParseFloat: parsing \"mwkmefok\": invalid syntax"),
		},
		{
			name:     "Error Parse FLoat",
			input:    "1 7 mkfmkmk -28",
			expected: "ошибка преобразования коэффициента 3: strconv.ParseFloat: parsing \"mkfmkmk\": invalid syntax",
			err:      errors.New("ошибка преобразования коэффициента 3: strconv.ParseFloat: parsing \"mkfmkmk\": invalid syntax"),
		},
		{
			name:     "Error Parse Float",
			input:    "1 7 -4 mvpsdmvlpdm",
			expected: "ошибка преобразования коэффициента 4: strconv.ParseFloat: parsing \"mvpsdmvlpdm\": invalid syntax",
			err:      errors.New("ошибка преобразования коэффициента 4: strconv.ParseFloat: parsing \"mvpsdmvlpdm\": invalid syntax"),
		},
		{
			name:     "Error Cubic Equation",
			input:    "0 1 0 1",
			expected: "ошибка вычисления: ошибка решения кубического уравнения: не удалось найти вещественный корень",
			err:      errors.New("ошибка вычисления: ошибка решения кубического уравнения: не удалось найти вещественный корень"),
		},
		{
			name:     "Very Big Cf",
			input:    "1 -7 -4 1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expected: "ошибка преобразования коэффициента 4: strconv.ParseFloat: parsing \"1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\": value out of range",
			err:      errors.New("ошибка преобразования коэффициента 4: strconv.ParseFloat: parsing \"1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\": value out of range"),
		},
	}
	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			actual, err := BezuTelegram(test.input)
			if test.err != nil {
				if err == nil || err.Error() != test.err.Error() {
					t.Errorf("expected error '%v' got '%v'", test.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actual = strings.TrimSpace(strings.ReplaceAll(actual, "\r\n", "\n"))
			expected := strings.TrimSpace(strings.ReplaceAll(test.expected, "\r\n", "\n"))

			if actual != expected {
				t.Errorf("expected:\n%s\ngot:\n%s", expected, actual)
			}
		})
	}
}
