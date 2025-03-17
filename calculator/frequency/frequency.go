package frequency

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// частота
type FrequencyResult struct {
	Value         float64
	ValueQuantity int
}

// Подсчет частоты
func countFrequency(array []float64) []FrequencyResult {
	frequencyMap := make(map[float64]int)
	for _, value := range array {
		frequencyMap[value]++
	}
	var results []FrequencyResult
	for value, quantity := range frequencyMap {
		results = append(results, FrequencyResult{value, quantity})
	}
	return results
}

// Преобразование строки в массив чисел
func ParseNumbers(s string) ([]float64, error) {
	var numbers []float64
	valuesStr := strings.Split(s, ",")
	for _, valueStr := range valuesStr {
		value, err := strconv.ParseFloat(strings.TrimSpace(valueStr), 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования: %w", err)
		}
		numbers = append(numbers, value)
	}
	return numbers, nil
}

// Нахождение моды (с учетом нескольких мод с максимальной частотой)
func FindModa(array []float64) []float64 {
	results := countFrequency(array)
	if len(results) == 0 {
		return []float64{} // Пустой массив, если входной массив пуст
	}

	maxQuantity := 0
	for _, result := range results {
		if result.ValueQuantity > maxQuantity {
			maxQuantity = result.ValueQuantity
		}
	}

	var moda []float64
	for _, result := range results {
		if result.ValueQuantity == maxQuantity {
			moda = append(moda, result.Value)
		}
	}

	// Только если несколько чисел имеют максимальную частоту, тогда возвращается весь список мод
	sort.Float64s(moda)

	return moda

}

// Форматирование вывода
func CalculateFrequency(input string) (string, error) {
	numbers, err := ParseNumbers(input)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования чисел: %w", err)
	}

	results := countFrequency(numbers)
	numCount := len(numbers)
	moda := FindModa(numbers)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Количество чисел в списке: %d\n", numCount))

	for _, result := range results {
		ratio := float64(result.ValueQuantity) / float64(numCount)
		if result.Value-float64(int64(result.Value)) != 0 {
			sb.WriteString(fmt.Sprintf("Число %.3f повторяется %d раз. Частота повторения: %.3f\n", result.Value, result.ValueQuantity, ratio))
		} else {
			sb.WriteString(fmt.Sprintf("Число %d повторяется %d раз. Частота повторения: %.3f\n", int64(result.Value), result.ValueQuantity, ratio))
		}

	}

	sb.WriteString("Мода: ")
	for i, val := range moda {
		if val-float64(int64(val)) != 0 {
			sb.WriteString(fmt.Sprintf("%.3f", val))
		} else {
			sb.WriteString(fmt.Sprintf("%d", int64(val)))
		}

		if i < len(moda)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("\n")
	return sb.String(), nil
}
