package linear

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nikitakutergin59/calculator/equations/TW" // директория TW реализует пакет equations
)

// инвертируем уравнение
func InvertedEquations(tokens []equations.Token) ([]equations.Token, error) {
	equationsIndex := -1
	for i, token := range tokens {
		if token.Type == equations.TOKEN_OPERATOR && token.Value == "=" {
			equationsIndex = i
			break
		}
	}
	if equationsIndex == -1 {
		return nil, fmt.Errorf("оператор '=' не найден")
	}

	leftSide := tokens[:equationsIndex]
	rightSide := tokens[equationsIndex+1:]
	//log.Println("левая часть", leftSide)
	//log.Println("правая часть", rightSide)
	if len(rightSide) > 0 {
		if rightSide[0].Type == equations.TOKEN_NUMBER && !strings.HasPrefix(rightSide[0].Value, "-") {
			rightSide = append([]equations.Token{{Type: equations.TOKEN_OPERATOR, Value: "+"}}, rightSide...)
		}
		if rightSide[0].Type == equations.TOKEN_PARENT_OPEN && !strings.HasPrefix(rightSide[0].Value, "-") {
			rightSide = append([]equations.Token{{Type: equations.TOKEN_OPERATOR, Value: "+"}}, rightSide...)
		}
		if rightSide[0].Type == equations.TOKEN_VARIABLE && !strings.HasPrefix(rightSide[0].Value, "-") {
			rightSide = append([]equations.Token{{Type: equations.TOKEN_OPERATOR, Value: "+"}}, rightSide...)
		}
	}

	//интвертируем знаки в правой части
	invertedRightSide := make([]equations.Token, len(rightSide))
	for j, token := range rightSide {
		invertedToken := token
		switch token.Type {
		case equations.TOKEN_OPERATOR: // иневертация операторов
			switch token.Value {
			case "+":
				invertedToken.Value = "-"
			case "-":
				invertedToken.Value = "+"
			default:
				invertedToken = token //оставляем как есть
			}
		case equations.TOKEN_MULT:
			invertedToken = token // оставляем как есть
		case equations.TOKEN_NUMBER: // инвертация чисел
			if strings.HasPrefix(invertedToken.Value, "-") {
				invertedToken.Value = strings.Replace(invertedToken.Value, "-", "+", 1)
			} else {
				invertedToken.Value = strings.Replace(invertedToken.Value, "+", "-", 1)
			}
		case equations.TOKEN_PARENT_OPEN: // инвертация скобок
			if strings.HasPrefix(invertedToken.Value, "-") {
				invertedToken.Value = strings.Replace(invertedToken.Value, "-", "+", 1)
			} else {
				invertedToken.Value = strings.Replace(invertedToken.Value, "+", "-", 1)
			}
		case equations.TOKEN_VARIABLE: // инвертация переменных
			if strings.HasPrefix(invertedToken.Value, "-") {
				invertedToken.Value = strings.Replace(invertedToken.Value, "-", "+", 1)
			} else {
				invertedToken.Value = strings.Replace(invertedToken.Value, "+", "-", 1)
			}
		default:
			invertedToken = token // просто копируем токен, если это переменная или скобка
		}
		invertedRightSide[j] = invertedToken
	}
	// форматирование результатов
	result_inverted := make([]equations.Token, 0, len(leftSide)+len(invertedRightSide)+1)
	result_inverted = append(result_inverted, leftSide...)
	//log.Println("инвертация", result_inverted)
	result_inverted = append(result_inverted, invertedRightSide...)
	//log.Println("инвертация_1", result_inverted)
	return result_inverted, nil
}

// функция для раскрытия скобко
func OpenParent(result_inverted []equations.Token) ([]equations.Token, error) {
	for i, openTheDoor := range result_inverted {
		log.Printf("i=%d", i)
		if openTheDoor.Type == equations.TOKEN_PARENT_OPEN && openTheDoor.Value == "(" {
			closeIndex := -1
			openParent := 1
			for j := i + 1; j < len(result_inverted); j++ {
				if result_inverted[j].Type == equations.TOKEN_PARENT_OPEN && result_inverted[j].Value == "(" {
					openParent++ // встретили ещё одну открывающую скобку
				} else if result_inverted[j].Type == equations.TOKEN_PARENT_CLOSE && result_inverted[j].Value == ")" {
					openParent--
					if openParent == 0 {
						closeIndex = j
						log.Printf("closeIndex=%d", closeIndex)
						break
					}
				}
			}
			if closeIndex == -1 { //проверка на наличие закрывающей скобки
				return nil, fmt.Errorf("не найдена закрывающая скобка")
			}
			// извлекаем срез в скобках
			parentSlise := result_inverted[i+1 : closeIndex]
			//раскрываем скобки рекурсией
			processedParentSlise, err := OpenParent(parentSlise)
			if err != nil {
				return nil, err
			}
			// проверяем есть ли число перед скобками
			multiplier := 0.0
			multiplierIndex := -1
			if i > 0 {
				log.Println("i > 0 result_inverted[i-2]", result_inverted[i-2])
				if i >= 2 && result_inverted[i-2].Type == equations.TOKEN_NUMBER {
					multiplierValue, err := strconv.ParseFloat(result_inverted[i-2].Value, 64)
					if err != nil {
						return nil, fmt.Errorf("ошибка при преобразовании числа: %w", err)
					}
					multiplier = multiplierValue
					log.Println("поступает в фуункции из equations", multiplier)
					multiplierIndex = i - 0
					log.Println("инвекс", multiplierIndex)
				}
			}

			// умножение содержимого скобок на число перед скобкой
			multipliedInnerSlice, err := equations.MultiplyInnerSlice(processedParentSlise, multiplier)
			if err != nil {
				return nil, fmt.Errorf("ошибка MultiplyInnerSlice: %w", err)
			}

			// инвертация знаков
			if i > 0 && result_inverted[i-1].Type == equations.TOKEN_OPERATOR && result_inverted[i-1].Value == "-" {
				multipliedInnerSlice = equations.InvertedOperator(multipliedInnerSlice)
			}

			if multiplierIndex != -1 {
				// удаление множителя
				result_inverted = append(result_inverted[:multiplierIndex], result_inverted[i:]...)
				// Важно: скорректировать i после удаления!
				i = multiplierIndex // Теперь i указывает на место, куда будет вставлен срез
			}

			var tail []equations.Token
			if closeIndex+1 < len(result_inverted) {
				tail = result_inverted[closeIndex+1:]
			} else {
				tail = []equations.Token{}
			}
			open_result_inverted := append(result_inverted[:i], append(multipliedInnerSlice, tail...)...)
			return open_result_inverted, nil
		}
	}
	return result_inverted, nil
}

// функция проверят есть ли ещё скобки в уравнение если есть то перезапускает OpenParent
func OpenAllParent(all_open_result_inverted []equations.Token) ([]equations.Token, error) {
	result_inverted := all_open_result_inverted
	for {
		allLen := len(result_inverted)

		var err error
		result_inverted, err = OpenParent(result_inverted)
		if err != nil {
			return nil, err
		}

		if len(result_inverted) == allLen {
			break
		}
	}
	//log.Println("OpenAllParent", result_inverted)
	return result_inverted, nil
}
