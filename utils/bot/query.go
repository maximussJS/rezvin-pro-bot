package bot

import (
	"fmt"
	"rezvin-pro-bot/types/bot"
	"strconv"
	"strings"
)

func AddParamsToQueryString(prefix string, params *bot_types.Params) string {
	if params == nil {
		return prefix
	}

	separator := "?"
	if strings.Contains(prefix, "?") {
		separator = "&"
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(prefix)
	queryBuilder.WriteString(separator)

	var paramPairs []string

	if params.ProgramId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("programId=%d", params.ProgramId))
	}
	if params.UserId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("userId=%d", params.UserId))
	}
	if params.ExerciseId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("exerciseId=%d", params.ExerciseId))
	}
	if params.UserProgramId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("userProgramId=%d", params.UserProgramId))
	}
	if params.UserExerciseRecordId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("userExerciseRecordId=%d", params.UserExerciseRecordId))
	}
	if params.Limit != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("limit=%d", params.Limit))
	}
	if params.Offset != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("offset=%d", params.Offset))
	}

	if len(paramPairs) == 0 {
		return prefix
	}

	queryBuilder.WriteString(strings.Join(paramPairs, "&"))

	return queryBuilder.String()
}

func ParseParamsFromQueryString(queryStr string) (*bot_types.Params, error) {
	params := bot_types.NewEmptyParams()

	parts := strings.SplitN(queryStr, "?", 2)
	if len(parts) != 2 {
		return params, nil
	}

	queryPart := parts[1]

	pairs := strings.Split(queryPart, "&")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		if key == "" || value == "" {
			continue
		}

		switch key {
		case "programId":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid programId: %v", err)
			}
			params.ProgramId = uint(parsedValue)
		case "userId":
			parsedValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid userId: %v", err)
			}
			params.UserId = parsedValue
		case "exerciseId":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid exerciseId: %v", err)
			}
			params.ExerciseId = uint(parsedValue)
		case "userProgramId":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid userProgramId: %v", err)
			}
			params.UserProgramId = uint(parsedValue)
		case "userExerciseRecordId":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid userExerciseRecordId: %v", err)
			}
			params.UserExerciseRecordId = uint(parsedValue)
		case "limit":
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid limit: %v", err)
			}
			params.Limit = parsedValue
		case "offset":
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid offset: %v", err)
			}
			params.Offset = parsedValue
		default:
			return nil, fmt.Errorf("unknown key: %s", key)
		}
	}

	return params, nil
}
