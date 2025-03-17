package diskriminant

import (
	"errors"
	"strings"
	"testing"
)

func TestDiskriminante(t *testing.T) {
	testCase := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:  "Usual equation",
			input: "1 2 -3",
			expected: `Ваше квадратное уравнение: 1x²+2x-3=0
Корни уравнения:
Корень 1: 1
Корень 2: -3`,
			err: nil,
		},
		{
			name:  "Negative cf",
			input: "-1 -5 -6",
			expected: `Ваше квадратное уравнение: -1x²-5x-6=0
Корни уравнения:
Корень 1: -3
Корень 2: -2`,
			err: nil,
		},
		{
			name:  "Fractional values cf",
			input: "0.052 0.052 -1.488",
			expected: `Ваше квадратное уравнение: 0.052x²+0.052x-1.488=0
Корни уравнения:
Корень 1: 4.872652
Корень 2: -5.872652`,
			err: nil,
		},
		{
			name:     "There are no valid roots",
			input:    "-5 -6 -7",
			expected: "ошибка вычисления корней: дискриминант меньше нуля, нет действительных корней",
			err:      errors.New("ошибка вычисления корней: дискриминант меньше нуля, нет действительных корней"),
		},
		{
			name:     "Incorrect number of arguments",
			input:    "5 -6",
			expected: "неверное количество аргументов. Используйте: a b c",
			err:      errors.New("неверное количество аргументов. Используйте: a b c"),
		},
		{
			name:     "Empty line",
			input:    "",
			expected: "неверное количество аргументов. Используйте: a b c",
			err:      errors.New("неверное количество аргументов. Используйте: a b c"),
		},
		{
			name:     "Error Parse Float",
			input:    "rver ewqg qewfefqw",
			expected: "ошибка преобразования a: strconv.ParseFloat: parsing \"rver\": invalid syntax",
			err:      errors.New("ошибка преобразования a: strconv.ParseFloat: parsing \"rver\": invalid syntax"),
		},
		{
			name:     "Error Parse Float_1",
			input:    "1, 52, 1488",
			expected: "ошибка преобразования a: strconv.ParseFloat: parsing \"1,\": invalid syntax",
			err:      errors.New("ошибка преобразования a: strconv.ParseFloat: parsing \"1,\": invalid syntax"),
		},
		{
			name:     "Error Parse Float_2",
			input:    "gnhfgdhjf 5 6",
			expected: "ошибка преобразования a: strconv.ParseFloat: parsing \"gnhfgdhjf\": invalid syntax",
			err:      errors.New("ошибка преобразования a: strconv.ParseFloat: parsing \"gnhfgdhjf\": invalid syntax"),
		},
		{
			name:     "Error Parse FLoat_3",
			input:    "5 dfwhd 6",
			expected: "ошибка преобразования b: strconv.ParseFloat: parsing \"dfwhd\": invalid syntax",
			err:      errors.New("ошибка преобразования b: strconv.ParseFloat: parsing \"dfwhd\": invalid syntax"),
		},
		{
			name:     "Error Parse Float_4",
			input:    "5 6 fbhhwsfd",
			expected: "ошибка преобразования c: strconv.ParseFloat: parsing \"fbhhwsfd\": invalid syntax",
			err:      errors.New("ошибка преобразования c: strconv.ParseFloat: parsing \"fbhhwsfd\": invalid syntax"),
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			actual, err := DiscriminantFromString(test.input)
			if test.err != nil {
				if err == nil || err.Error() != test.err.Error() {
					t.Errorf("expected error '%v', got '%v'", test.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			//нормализация строк
			actual = strings.TrimSpace(strings.ReplaceAll(actual, "\r\n", "\n"))
			expected := strings.TrimSpace(strings.ReplaceAll(test.expected, "\r\n", "\n"))

			if actual != expected {
				t.Errorf("expected:\n%s\ngot:\n%s", expected, actual)
			}
		})
	}
}
