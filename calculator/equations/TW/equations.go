package equations

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// создаём константы с типами токенов
const (
	TOKEN_NUMBER           = "NUMBER"
	TOKEN_VARIABLE         = "VARIABLE"
	TOKEN_FUNC             = "FUNC"
	TOKEN_OPERATOR         = "OPERATOR"
	TOKEN_PARENT_OPEN      = "OPEN"
	TOKEN_PARENT_CLOSE     = "CLOSE"
	TOKEN_CURLY_BRACE_OPEN = "CURLY_OPEN"
	TOKEN_POWER            = "POWER"
	TOKEN_MULT             = "MULT"
	TOKEN_UNARY_MINUS      = "UNARY_MINUS"
	TOKEN_UNKNOWN          = "UNKNOWN"
)

// структура к которой обращаеться функция токенизации
type Token struct {
	Type  string
	Value string
}

func Tokenize_BH(expression string) ([]Token, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	//регулярные выражения для разных типов токенов
	numberRegex := regexp.MustCompile(`^(\d+(\.\d*)?|\.\d+)`)
	variableRegex := regexp.MustCompile(`^(x|y)`)
	funcRegex := regexp.MustCompile(`^(sqrt|abs)`)
	argRegex := regexp.MustCompile(`^(\d+(\.\d*)?|\.\d+|\(([^()]+|\(([^()]*)\))*\))`) // регулярка для аргумента функции (число или выражение в скобках)
	operatorRegex := regexp.MustCompile(`^(\+|-|\*|/|=)`)
	powerRegex := regexp.MustCompile(`^\^(\d+(\.\d*)?|\.\d+|\(([^()]+|\(([^()]*)\))*\))`) // объединённые "^" с числами или варажениями в скобках
	openRegex := regexp.MustCompile(`^\(`)
	closeRegex := regexp.MustCompile(`^\)`)
	curlyOpenRegex := regexp.MustCompile(`^\{`)
	unknouRegex := regexp.MustCompile(`.`)

	token_BH := []Token{}
	unaryMinusNext := false

	for len(expression) > 0 {
		// токенизация типа числа
		if numberMatch := numberRegex.FindString(expression); numberMatch != "" {
			numberValue := numberMatch
			// проверяем, установлен ли флаг унарного минуса
			if unaryMinusNext {
				numberValue = "-" + numberValue // объединяем минус и число
				unaryMinusNext = false          // сбрасываем флаг
			}
			token_BH = append(token_BH, Token{TOKEN_NUMBER, numberValue})
			expression = expression[len(numberMatch):]

			// проверяем, следует ли за числом переменная (например, 2x)
			expression = strings.TrimSpace(expression) // удаляем пробелы после числа
			if len(expression) > 0 {
				// проверяем, что следующий символ - буква (переменная)
				if variableRegex.MatchString(string(expression[0])) {
					token_BH = append(token_BH, Token{TOKEN_MULT, "*"})
				}
				if openRegex.MatchString(string(expression[0])) {
					token_BH = append(token_BH, Token{TOKEN_MULT, "*"})
				}
			}
			continue
		}
		//токенизация типа перменные x,y
		if variableMatch := variableRegex.FindString(expression); variableMatch != "" {
			variableValue := variableMatch
			if unaryMinusNext {
				variableValue = "-" + variableValue
				unaryMinusNext = false
			}
			token_BH = append(token_BH, Token{TOKEN_VARIABLE, variableValue})
			expression = expression[len(variableMatch):]

			expression = strings.TrimSpace(expression)
			if len(expression) > 0 {
				if openRegex.MatchString(string(expression[0])) {
					token_BH = append(token_BH, Token{TOKEN_MULT, "*"})
				}
			}
			continue
		}
		//склеиваем функции с их аргументами
		if funcMatch := funcRegex.FindString(expression); funcMatch != "" {
			functionName := funcMatch
			expression = expression[len(funcMatch):]

			argMatch := argRegex.FindString(expression)
			if argMatch != "" {
				functionTokenValue := functionName + argMatch
				token_BH = append(token_BH, Token{TOKEN_FUNC, functionTokenValue})
				expression = expression[len(argMatch):]
				continue
			} else {
				return nil, fmt.Errorf("ожидался аргумент после функции '%s'", functionName)
			}
		}
		//токенизация типа операторы
		if operatorMatch := operatorRegex.FindString(expression); operatorMatch != "" {
			operator := operatorMatch
			expression = expression[len(operatorMatch):]

			if operator == "-" {
				// проверяем, является ли это унарным минусом или оператором вычитания
				if len(token_BH) == 0 || // Это начало выражения
					(len(token_BH) > 0 && token_BH[len(token_BH)-1].Type == TOKEN_PARENT_OPEN) || // после открывающей скобки
					(len(token_BH) > 0 && token_BH[len(token_BH)-1].Type == TOKEN_OPERATOR) || // после другого оператора
					(len(token_BH) > 0 && token_BH[len(token_BH)-1].Type == TOKEN_UNARY_MINUS) { // после унарного минуса
					// это унарный минус
					unaryMinusNext = true // устанавливаем флаг
					continue              // не добавляем токен UNARY_MINUS, а просто переходим к следующей итерации
				} else {
					// это оператор вычитания
					token_BH = append(token_BH, Token{TOKEN_OPERATOR, operator})
					continue
				}
			} else {
				// это другой оператор (+, *, /, =)
				token_BH = append(token_BH, Token{TOKEN_OPERATOR, operator})
				continue
			}
		}
		//склеиваем символ степени с последующим числом или выражением в скобках
		if strings.HasPrefix(expression, "^") {
			powerMatch := powerRegex.FindString(expression)
			if powerMatch != "" {
				// найдено число или выражение в скобках после "^"
				powerTokenValue := powerMatch
				token_BH = append(token_BH, Token{TOKEN_POWER, powerTokenValue})
				expression = expression[len(powerMatch):]
				continue
			} else {
				// после "^" нет ни числа, ни выражения в скобках, это ошибка
				return nil, fmt.Errorf("ожидалось число или выражение в скобках после оператора '^'")
			}
		}
		//токенизация открывающей скобки
		if openMatch := openRegex.FindString(expression); openMatch != "" {
			openValue := openMatch
			if unaryMinusNext {
				openValue = "-" + openValue
				unaryMinusNext = false
			}
			// проверяем, есть ли коэффициент перед скобкой
			if len(token_BH) > 0 {
				lastToken := token_BH[len(token_BH)-1]
				if lastToken.Type == TOKEN_NUMBER {
					// eсли перед скобкой число, значит коэффициент есть. Ничего не добавляем.
				} else {
					if lastToken.Type == TOKEN_OPERATOR {
						// eсли перед скобкой оператор добавляем "1*"
						token_BH = append(token_BH, Token{TOKEN_NUMBER, "1"})
						token_BH = append(token_BH, Token{TOKEN_MULT, "*"})
					} else if lastToken.Type == TOKEN_PARENT_CLOSE {
						//если перед скобкой закрывающая скобка то добавляем "*1*"
						token_BH = append(token_BH, Token{TOKEN_MULT, "*"})
						token_BH = append(token_BH, Token{TOKEN_NUMBER, "1"})
						token_BH = append(token_BH, Token{TOKEN_MULT, "*"})
					}
				}
			} else {
				// eсли это первая скобка в выражении, добавляем "1*"
				token_BH = append(token_BH, Token{TOKEN_NUMBER, "1"})
				token_BH = append(token_BH, Token{TOKEN_MULT, "*"})
			}
			token_BH = append(token_BH, Token{TOKEN_PARENT_OPEN, openValue})
			expression = expression[len(openMatch):]
			continue
		}
		//токенизация закрывающей скобки
		if closeMatch := closeRegex.FindString(expression); closeMatch != "" {
			token_BH = append(token_BH, Token{TOKEN_PARENT_CLOSE, closeMatch})
			expression = expression[len(closeMatch):]
			continue
		}
		//токенизация открывающей фигурной скобки
		if curlyOpenMatch := curlyOpenRegex.FindString(expression); curlyOpenMatch != "" {
			token_BH = append(token_BH, Token{TOKEN_CURLY_BRACE_OPEN, curlyOpenMatch})
			expression = expression[len(curlyOpenMatch):]
			continue
		}
		//если токен не подходит не под один из типов
		if unknouMatch := unknouRegex.FindString(expression); unknouMatch != "" {
			token_BH = append(token_BH, Token{TOKEN_UNKNOWN, unknouMatch})
			expression = expression[1:]
		}
	}
	//log.Println("распаршенное уравнение", token_BH)
	return token_BH, nil
}

// функция для определения типа уравнений
func WhatTypeEquations(token_BH []Token) (string, error) {
	hasPower := false // указатель на наличие степени
	for _, token := range token_BH {
		if token.Type == TOKEN_UNKNOWN {
			return "", fmt.Errorf("токен %s не распознан", token.Value)
		}
		if token.Type == TOKEN_POWER {
			hasPower = true
		}
	}
	if hasPower {
		return "Нелинейное", nil
	} else {
		return "Линейное", nil
	}
}

// вспомогательная функция меняет знаки на противоположные используеться во всех пакетах
func InvertedOperator(token_BH []Token) []Token {
	invertedToken := make([]Token, len(token_BH), 4096)
	for i, token := range token_BH {
		invertedToken[i] = token
		if token.Type == TOKEN_OPERATOR {
			if token.Value == "+" {
				invertedToken[i].Value = "-"
			} else if token.Value == "-" {
				invertedToken[i].Value = "+"
			}
		}
	}
	return invertedToken
}

// MultiplyInnerSlice вспомогательная функция умножает число на всё что было в скобках
func MultiplyInnerSlice(innerSlice []Token, multiplier float64) ([]Token, error) {
	log.Println("НАЧАЛЬНОЕ ЗНАЧЕНИЕ", innerSlice)
	multipliedSlice := make([]Token, 0)

	for i, token := range innerSlice {
		switch token.Type {
		case TOKEN_NUMBER:
			numValue, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("ошибка преобразования строки в число: %w", err)
			}
			multipliedValue := numValue * multiplier
			// Форматируем с точностью до 6 знаков после запятой (можно изменить)
			multipliedToken := Token{Type: TOKEN_NUMBER, Value: fmt.Sprintf("%.6f", multipliedValue)}
			multipliedSlice = append(multipliedSlice, multipliedToken)
		case TOKEN_VARIABLE:
			if i > 0 && innerSlice[i-1].Type != TOKEN_MULT && innerSlice[i-1].Type != TOKEN_OPERATOR {
				multipliedSlice = append(multipliedSlice, Token{Type: TOKEN_MULT, Value: "*"})
			}
			multipliedSlice = append(multipliedSlice, token)
		default:
			multipliedSlice = append(multipliedSlice, token)
		}
	}

	log.Println("токены", multipliedSlice)
	formatMultSlice, err := formatFloat(multipliedSlice)
	if err != nil {
		return nil, fmt.Errorf("ошибка formatFloat: %w", err)
	}
	log.Println("после конвертации", formatMultSlice)
	return formatMultSlice, nil
}

func formatFloat(multipliedSlise []Token) ([]Token, error) {
	formatMultipliedSlise := make([]Token, 0, len(multipliedSlise))
	for _, formatMS := range multipliedSlise {
		log.Println("рабочая переменная", formatMS)
		switch formatMS.Type {
		case TOKEN_NUMBER:
			parts := strings.Split(formatMS.Value, ".")
			log.Printf("parts=%s", parts)
			if len(parts) == 2 {
				decimalPart := parts[1]
				precision := 0
				for i, r := range decimalPart {
					log.Printf("r=%d", r)
					log.Printf("decinalParts=%s", decimalPart)
					if r != '0' {
						precision = i + 1
						log.Printf("precision=%d", precision)
					}
				}
				floatFormatMSV, err := strconv.ParseFloat(formatMS.Value, 64)
				if err != nil {
					return nil, fmt.Errorf("ошибка конвертации %v", err)
				}
				strFormatMS := fmt.Sprintf("%."+fmt.Sprint(precision)+"f", floatFormatMSV)
				formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_NUMBER, Value: strFormatMS})
			}
		case TOKEN_VARIABLE:
			formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_VARIABLE, Value: formatMS.Value})
		case TOKEN_FUNC:
			formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_FUNC, Value: formatMS.Value})
		case TOKEN_OPERATOR:
			formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_OPERATOR, Value: formatMS.Value})
		case TOKEN_PARENT_OPEN:
			formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_PARENT_OPEN, Value: formatMS.Value})
		case TOKEN_PARENT_CLOSE:
			formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_PARENT_OPEN, Value: formatMS.Value})
		case TOKEN_CURLY_BRACE_OPEN:
			formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_CURLY_BRACE_OPEN, Value: formatMS.Value})
		case TOKEN_POWER:
			//пока просто добавляем в слайс без изменнений
			formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_POWER, Value: formatMS.Value})
		case TOKEN_MULT:
			formatMultipliedSlise = append(formatMultipliedSlise, Token{Type: TOKEN_MULT, Value: formatMS.Value})
		case TOKEN_UNKNOWN:
			return nil, fmt.Errorf("неизвестный токен: %+v", formatMS.Value)
		}
	}
	return formatMultipliedSlise, nil
}
