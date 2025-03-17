package frequency

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestFrequency(t *testing.T) {
	testCase := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:  "Usual equation",
			input: "1,2,3,4,5,6,7,8,9",
			expected: `Количество чисел в списке: 9
Число 9 повторяется 1 раз. Частота повторения: 0.111
Число 5 повторяется 1 раз. Частота повторения: 0.111
Число 7 повторяется 1 раз. Частота повторения: 0.111
Число 8 повторяется 1 раз. Частота повторения: 0.111
Число 4 повторяется 1 раз. Частота повторения: 0.111
Число 6 повторяется 1 раз. Частота повторения: 0.111
Число 1 повторяется 1 раз. Частота повторения: 0.111
Число 2 повторяется 1 раз. Частота повторения: 0.111
Число 3 повторяется 1 раз. Частота повторения: 0.111
Мода: 1, 2, 3, 4, 5, 6, 7, 8, 9`,
			err: nil,
		},
		{
			name:  "Negative Str",
			input: "-1,-2,-3,-4,-5,-6,-7,-8,-9",
			expected: `Количество чисел в списке: 9
Число -6 повторяется 1 раз. Частота повторения: 0.111
Число -7 повторяется 1 раз. Частота повторения: 0.111
Число -9 повторяется 1 раз. Частота повторения: 0.111
Число -1 повторяется 1 раз. Частота повторения: 0.111
Число -2 повторяется 1 раз. Частота повторения: 0.111
Число -3 повторяется 1 раз. Частота повторения: 0.111
Число -4 повторяется 1 раз. Частота повторения: 0.111
Число -5 повторяется 1 раз. Частота повторения: 0.111
Число -8 повторяется 1 раз. Частота повторения: 0.111
Мода: -9, -8, -7, -6, -5, -4, -3, -2, -1`,
			err: nil,
		},
		{
			name:  "Fractional values Str",
			input: "0.1,1.23,2.5,0.1,3.14159,4.0,1.23,5.67,2.5,6.999,7.0001,0.1,8.8,9.12,1.23,3.14159,5.67,6.999,4.0,7.0001",
			expected: `Количество чисел в списке: 20
Число 0.100 повторяется 3 раз. Частота повторения: 0.150
Число 1.230 повторяется 3 раз. Частота повторения: 0.150
Число 5.670 повторяется 2 раз. Частота повторения: 0.100
Число 6.999 повторяется 2 раз. Частота повторения: 0.100
Число 9.120 повторяется 1 раз. Частота повторения: 0.050
Число 2.500 повторяется 2 раз. Частота повторения: 0.100
Число 3.142 повторяется 2 раз. Частота повторения: 0.100
Число 4 повторяется 2 раз. Частота повторения: 0.100
Число 7.000 повторяется 2 раз. Частота повторения: 0.100
Число 8.800 повторяется 1 раз. Частота повторения: 0.050
Мода: 0.100, 1.230`,
			err: nil,
		},
		{
			name:     "Error Parse Float",
			input:    "jnjnjnjlo,2,3,4,5,6,7,8,9",
			expected: "ошибка преобразования чисел: ошибка преобразования: strconv.ParseFloat: parsing \"jnjnjnjlo\": invalid syntax",
			err:      errors.New("ошибка преобразования чисел: ошибка преобразования: strconv.ParseFloat: parsing \"jnjnjnjlo\": invalid syntax"),
		},
		// ошибка сробатывает всегда на какой бы позиции не стояла чилсо
	}
	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			actual, err := CalculateFrequency(test.input)
			if test.err != nil {
				if err == nil || err.Error() != test.err.Error() {
					t.Errorf("expected error: '%v' got '%v'", test.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actual = strings.TrimSpace(strings.ReplaceAll(actual, "\r\n", "\n"))
			expected := strings.TrimSpace(strings.ReplaceAll(test.expected, "\r\n", "\n"))

			actualLines := strings.Split(actual, "\n")
			expectedLines := strings.Split(expected, "\n")

			sort.Strings(actualLines)
			sort.Strings(expectedLines)

			if !reflect.DeepEqual(actualLines, expectedLines) {
				t.Errorf("expected (unordered): %s\ngot (unordered): %s", strings.Join(expectedLines, "\n"), strings.Join(actualLines, "\n"))
			}

		})
	}
}
