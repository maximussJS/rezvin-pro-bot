package bot

import (
	"fmt"
	"rezvin-pro-bot/src/constants"
	"rezvin-pro-bot/src/types"
	"strconv"
	"strings"
)

func AddParamsToQueryString(prefix string, params *types.Params) string {
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
		paramPairs = append(paramPairs, fmt.Sprintf("pid=%d", params.ProgramId))
	}
	if params.UserId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("uid=%d", params.UserId))
	}
	if params.ExerciseId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("eid=%d", params.ExerciseId))
	}
	if params.UserProgramId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("upid=%d", params.UserProgramId))
	}
	if params.UserResultId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("urid=%d", params.UserResultId))
	}
	if params.MeasureId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("mid=%d", params.MeasureId))
	}
	if params.UserMeasureId != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("umid=%d", params.UserMeasureId))
	}
	if params.Limit != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("l=%d", params.Limit))
	}
	if params.Offset != 0 {
		paramPairs = append(paramPairs, fmt.Sprintf("o=%d", params.Offset))
	}
	if params.Reps != constants.Zero {
		paramPairs = append(paramPairs, fmt.Sprintf("r=%d", params.Reps))
	}

	if len(paramPairs) == 0 {
		return prefix
	}

	queryBuilder.WriteString(strings.Join(paramPairs, "&"))

	return queryBuilder.String()
}

func ParseParamsFromQueryString(queryStr string) (*types.Params, error) {
	params := types.NewEmptyParams()

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
		case "pid":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid programId: %v", err)
			}
			params.ProgramId = uint(parsedValue)
		case "uid":
			parsedValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid userId: %v", err)
			}
			params.UserId = parsedValue
		case "eid":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid exerciseId: %v", err)
			}
			params.ExerciseId = uint(parsedValue)
		case "upid":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid userProgramId: %v", err)
			}
			params.UserProgramId = uint(parsedValue)
		case "urid":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid userResultId: %v", err)
			}
			params.UserResultId = uint(parsedValue)
		case "mid":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid measureId: %v", err)
			}
			params.MeasureId = uint(parsedValue)
		case "umid":
			parsedValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid userMeasureId: %v", err)
			}
			params.UserMeasureId = uint(parsedValue)
		case "l":
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid limit: %v", err)
			}
			params.Limit = parsedValue
		case "o":
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid offset: %v", err)
			}
			params.Offset = parsedValue
		case "r":
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid reps: %v", err)
			}
			params.Reps = constants.Reps(parsedValue)
		default:
			return nil, fmt.Errorf("unknown key: %s", key)
		}
	}

	return params, nil
}
