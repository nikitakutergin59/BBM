package telegram

import (
	"fmt"
	//"log"
	"strings"

	"github.com/nikitakutergin59/calculator/equations/TW"
	"github.com/nikitakutergin59/calculator/equations/linear"
)

// функция для взаимодействия с телеграмм ботом
func EquationTelegram(expression string) (string, error) {
	tokens, err := equations.Tokenize_BH(expression)
	if err != nil {
		return "", fmt.Errorf("ошибка токенизации: %w", err)
	}
	// строим строку с исходным уравнением
	var eSD strings.Builder
	for _, token := range tokens {
		eSD.WriteString(token.Value)
	}
	equationString := eSD.String()

	// определяем тип уравнения (до раскрытия скобок!)
	equationType, err := equations.WhatTypeEquations(tokens)
	if err != nil {
		return "", fmt.Errorf("ошибка при определении типа уравнения: %w", err)
	}

	var simplifiedTokens []equations.Token // тут будет храниться результат упрощения
	var invertedEquationString string      // для конветации инвертированного в строку
	var simplifiedEquationString string    // для конвертации упрощённого в строку

	// инвертируем упрощаем уравнение (для линейных)
	if equationType == "Линейное" {
		// Инвертируем упрощенное уравнение
		invertedTokens, err := linear.InvertedEquations(tokens)
		if err != nil {
			return "", fmt.Errorf("ошибка инвертирования уравнения: %w", err)
		}
		var iESD strings.Builder
		for _, token := range invertedTokens {
			iESD.WriteString(token.Value)
		}
		invertedEquationString = iESD.String() + "=0"

		// раскрываем скобки в инвертированном уравнении
		simplifiedTokens, err = linear.OpenAllParent(invertedTokens)
		if err != nil {
			return "", fmt.Errorf("ошибка при раскрытии скобок: %w", err)
		}

		// cтроим строку с упрощенным уравнением
		var sESD strings.Builder
		for _, token := range simplifiedTokens {
			sESD.WriteString(token.Value)
		}
		simplifiedEquationString = sESD.String() + "=0"
	}

	// формируем строку ответа для бота
	response := fmt.Sprintf(
		"Уравнение: %s\n"+
			"Тип уравнения: %s\n"+
			"Инвертированное уравнение: %s\n"+
			"Упрощенное уравнение: %s\n",
		equationString,
		equationType,
		invertedEquationString,
		simplifiedEquationString,
	)

	return response, nil
}
