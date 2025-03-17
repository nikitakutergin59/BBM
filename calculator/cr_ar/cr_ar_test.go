package crar

import (
	"errors"
	"strings"
	"testing"
)

func TestCrAr(t *testing.T) {
	testCase := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:  "Usual equation",
			input: "1,2,3,4,5,6,7,8,9",
			expected: `Количество чисел: 9
Среднее арифметическое: 5
Дисперсия: 6.666666666666667
Стандартное отклонение: 2.581988897471611
Медиана: 5
Максимальное значение: 9
Минимальное значение: 1
Размах: 8 (значительный)`,
			err: nil,
		},
		{
			name:  "Negative Str",
			input: "-1,-2,-3,-4,-5,-6,-7,-8,-9",
			expected: `Количество чисел: 9
Среднее арифметическое: -5
Дисперсия: 6.666666666666667
Стандартное отклонение: 2.581988897471611
Медиана: -5
Максимальное значение: -1
Минимальное значение: -9
Размах: 8 (значительный)`,
			err: nil,
		},
		{
			name:  "Fractional values Str",
			input: "0.1,1.23,2.5,0.1,3.14159,4.0,1.23,5.67,2.5,6.999,7.0001,0.1,8.8,9.12,1.23,3.14159,5.67,6.999,4.0,7.0001",
			expected: `Количество чисел: 20
Среднее арифметическое: 4.02657
Дисперсия: 8.27148596204900
Стандартное отклонение: 2.87601911712162
Медиана: 3.57079
Максимальное значение: 9.12
Минимальное значение: 0.1
Размах: 9.0 (значительный)`,
			err: nil,
		},
		{
			name:     "Error Parse Float",
			input:    "jnjnjnjlo,2,3,4,5,6,7,8,9",
			expected: "ошибка вычисления статистики: ошибка преобразования числа: strconv.ParseFloat: parsing \"jnjnjnjlo\": invalid syntax",
			err:      errors.New("ошибка вычисления статистики: ошибка преобразования числа: strconv.ParseFloat: parsing \"jnjnjnjlo\": invalid syntax"),
		},
		// ошибка сробатывает всегда на какой бы позиции не стояла чилсо
	}
	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			actual, err := StatsTelegram(test.input)
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
				t.Errorf("expected %s\ngot\n%s", expected, actual)
			}
		})
	}
}
