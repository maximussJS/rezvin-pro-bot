package validate_data

import (
	"fmt"
	"strconv"
)

func ValidateWeightAnswer(text string) (int, error) {
	weight, err := strconv.Atoi(text)

	if err != nil {
		return 0, fmt.Errorf("введіть число від 0 до 1000")
	}

	if weight < 0 || weight > 1000 {
		return 0, fmt.Errorf("введіть число від 0 до 1000")
	}

	return weight, nil
}
