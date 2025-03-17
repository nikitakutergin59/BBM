package diskriminant

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func CalculateDiscriminant(a, b, c float64) float64 {
	return (b * b) - 4*(a*c)
}

func CalculateRoots(a, b, discriminant float64) ([]float64, error) {
	if discriminant < 0 {
		return nil, errors.New("дискриминант меньше нуля, нет действительных корней")
	} else if discriminant == 0 {
		root := (-b) / (2 * a)
		return []float64{root}, nil
	} else {
		sqrtDiscriminant := math.Sqrt(discriminant)
		root1 := ((-b) + sqrtDiscriminant) / (2 * a)
		root2 := ((-b) - sqrtDiscriminant) / (2 * a)
		return []float64{root1, root2}, nil
	}
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

func DiscriminantFromString(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) != 3 {
		return "", errors.New("неверное количество аргументов. Используйте: a b c")
	}

	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования a: %w", err)
	}
	b, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования b: %w", err)
	}
	c, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования c: %w", err)
	}

	var equation strings.Builder
	equation.WriteString("Ваше квадратное уравнение: ")
	if a-float64(int64(a)) == 0 {
		equation.WriteString(fmt.Sprintf("%dx²", int64(a)))
	} else {
		equation.WriteString(fmt.Sprintf("%sx²", formatFloat(a)))
	}
	if b-float64(int64(b)) == 0 {
		if b >= 0 {
			equation.WriteString(fmt.Sprintf("+%dx", int64(b)))
		} else {
			equation.WriteString(fmt.Sprintf("%dx", int64(b)))
		}
	} else {
		if b >= 0 {
			equation.WriteString(fmt.Sprintf("+%sx", formatFloat(b)))
		} else {
			equation.WriteString(fmt.Sprintf("%sx", formatFloat(b)))
		}
	}
	if c-float64(int64(c)) == 0 {
		if c >= 0 {
			equation.WriteString(fmt.Sprintf("+%d=0\n", int64(c)))
		} else {
			equation.WriteString(fmt.Sprintf("%d=0\n", int64(c)))
		}
	} else {
		if c >= 0 {
			equation.WriteString(fmt.Sprintf("+%s=0\n", formatFloat(c)))
		} else {
			equation.WriteString(fmt.Sprintf("%s=0\n", formatFloat(c)))
		}
	}

	discriminant := CalculateDiscriminant(a, b, c)
	roots, err := CalculateRoots(a, b, discriminant)
	if err != nil {
		return "", fmt.Errorf("ошибка вычисления корней: %w", err)
	}

	var result strings.Builder
	result.WriteString(equation.String()) // Добавляем отформатированное уравнение
	result.WriteString("Корни уравнения:\n")
	for i, root := range roots {
		if root-float64(int64(root)) != 0 {
			result.WriteString(fmt.Sprintf("Корень %d: %f\n", i+1, root))
		} else {
			result.WriteString(fmt.Sprintf("Корень %d: %d\n", i+1, int64(root)))
		}

	}
	return result.String(), nil
}
