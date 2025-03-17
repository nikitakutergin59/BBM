package bezu

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// теорема безу
// делители свобоного члена

// Polynomial - структура для представления многочлена
type Polynomial struct {
	Coefficients []float64
}

// Value - вычисляет значение многочлена в точке x
func (p Polynomial) Value(x float64) float64 {
	value := 0.0
	for i, coeff := range p.Coefficients {
		value += coeff * math.Pow(x, float64(len(p.Coefficients)-1-i))
	}
	return value
}

// Derivative - вычисляет производную многочлена в точке x
func (p Polynomial) Derivative(x float64) float64 {
	derivative := 0.0
	for i, coeff := range p.Coefficients {
		if i < len(p.Coefficients)-1 {
			derivative += coeff * float64(len(p.Coefficients)-1-i) * math.Pow(x, float64(len(p.Coefficients)-2-i))
		}
	}
	return derivative
}

func NewtonMethod(polynomial Polynomial, initialGuess float64, tolerance float64, maxIterations int) (float64, error) {
	x := initialGuess
	for i := 0; i < maxIterations; i++ {
		value := polynomial.Value(x)
		derivative := polynomial.Derivative(x)

		if derivative == 0 {
			return 0, fmt.Errorf("производная равна нулю в точке x = %f", x)
		}

		nextX := x - value/derivative
		if math.Abs(nextX-x) < tolerance {
			return nextX, nil
		}
		x = nextX
	}
	return 0, fmt.Errorf("метод не сошелся после %d итераций", maxIterations)
}

// hornerScheme - схема Горнера для многочлена
func hornerScheme(coefficients []float64, x float64) []float64 {
	result := make([]float64, len(coefficients)-1)
	result[len(coefficients)-2] = coefficients[len(coefficients)-1]
	for i := len(coefficients) - 2; i > 0; i-- {
		result[i-1] = coefficients[i] + result[i]*x
	}
	return result
}

// Вычисление корней кубического уравнения
func solveCubicReal(polynomial Polynomial) ([]float64, error) {
	// Поиск первого корня методом Ньютона (изменены начальные приближения).
	roots := make([]float64, 0, 3)
	initialGuesses := []float64{-3, -2, -1, 0, 1, 2, 3}

	for _, guess := range initialGuesses {
		root, err := NewtonMethod(polynomial, guess, 0.0001, 1000)
		if err != nil {
			// Если метод Ньютона не сошёлся, пропускаем это начальное приближение.
			continue
		}
		// Проверяем, что полученный корень действительный (не NaN и не бесконечность)
		if !math.IsNaN(root) && !math.IsInf(root, 0) {
			roots = append(roots, root)
		}
	}

	if len(roots) == 0 {
		return nil, fmt.Errorf("не удалось найти вещественный корень")
	}

	// Используем первый найденный корень для дальнейших вычислений.
	root := roots[0]
	reducedCoefficients := hornerScheme(polynomial.Coefficients, root)

	if len(reducedCoefficients) < 2 {
		return roots, nil //Если не квадратное уравнение, то уже нашли один корень
	}

	a, b, c := reducedCoefficients[0], reducedCoefficients[1], reducedCoefficients[2]
	discriminant := b*b - 4*a*c

	if discriminant >= 0 {
		sqrtDiscriminant := math.Sqrt(discriminant)
		root2 := (-b + sqrtDiscriminant) / (2 * a)
		root3 := (-b - sqrtDiscriminant) / (2 * a)
		roots = append(roots, root2, root3)
	}

	// Удаляем дубликаты и округляем до 3 знаков после запятой
	uniqueRoots := make([]float64, 0)
	seen := make(map[float64]bool)
	for _, v := range roots {
		roundedV := math.Round(v*1000) / 1000 // Округление до 3 знаков после запятой
		if !seen[roundedV] {
			seen[roundedV] = true
			uniqueRoots = append(uniqueRoots, roundedV)
		}
	}
	sort.Float64s(uniqueRoots)

	return uniqueRoots, nil
}

func BezuCalculate(polynomial Polynomial) ([]float64, float64, error) {
	//Здесь происходит само вычисление
	roots, err := solveCubicReal(polynomial)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка решения кубического уравнения: %w", err)
	}

	//Находим приблизительный корень методом Ньютона
	root, err := NewtonMethod(polynomial, 1, 0.0001, 10000)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка метода Ньютона: %w", err)
	}
	return roots, root, nil

}

func formatFloat(num float64) string {
	s := strconv.FormatFloat(num, 'f', -1, 64)
	parts := strings.Split(s, ".")
	if len(parts) == 2 {
		decimalPart := parts[1]
		precision := 0
		for i, r := range decimalPart {
			if r != '0' {
				precision = i + 1
			}
		}
		if precision > 0 {
			return fmt.Sprintf("%."+fmt.Sprint(precision)+"f", num)
		}
	}
	return s
}

func BezuTelegram(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) != 4 {
		return "", errors.New("неверное количество аргументов.  Используйте: a b c d")
	}

	coefficients := make([]float64, 4)
	for i := 0; i < 4; i++ {
		coeff, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			return "", fmt.Errorf("ошибка преобразования коэффициента %d: %w", i+1, err)
		}
		coefficients[i] = coeff
	}

	polynomial := Polynomial{Coefficients: coefficients}
	roots, _, err := BezuCalculate(polynomial)
	if err != nil {
		return "", fmt.Errorf("ошибка вычисления: %w", err)
	}

	var result strings.Builder
	result.WriteString("Уравнение: ")
	degree := len(coefficients) - 1 // Степень многочлена
	degreeMap := map[int]string{
		1: "\u00B9",
		2: "\u00B2",
		3: "\u00B3",
	}

	for i, coeff := range coefficients {
		if coeff == 0 {
			continue
		}
		if i > 0 {
			if coeff >= 0 {
				result.WriteString("+")
			}
		}
		result.WriteString(formatFloat(coeff))

		if degree-i > 0 { // Если степень x больше 0 (не константа)
			result.WriteString("x")
			if degree-i > 1 {
				if val, ok := degreeMap[degree-i]; ok {
					result.WriteString(val)
				}
			}
		}
	}
	result.WriteString("=0\n")
	result.WriteString("Корни:\n")
	for i, root := range roots {
		result.WriteString(fmt.Sprintf("Корень %d: %s\n", i+1, formatFloat(root)))
	}
	return result.String(), nil
}
