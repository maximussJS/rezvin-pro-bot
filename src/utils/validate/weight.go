package validate_data

import (
	"fmt"
	"strconv"
	"strings"
)

func ValidateWeightAnswer(text string) (int, error) {
	weight, err := strconv.Atoi(strings.TrimSpace(text))

	if err != nil {
		return 0, fmt.Errorf("введіть число від 0 до 1000")
	}

	if weight < 0 || weight > 1000 {
		return 0, fmt.Errorf("введіть число від 0 до 1000")
	}

	return weight, nil
}

func ValidateValueAnswer(text string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(text), 64)

	if err != nil {
		return 0, fmt.Errorf("введіть число від 0.0 до 200.0")
	}

	if value < 0 || value > 200 {
		return 0, fmt.Errorf("введіть число від 0.0 до 200.0")
	}

	return value, nil
}

func ValidateStringAnswer(text string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("введіть текст")
	}

	return text, nil
}
