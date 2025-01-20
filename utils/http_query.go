package utils

import (
	"fmt"
	"net/http"
	"rezvin-pro-bot/types/responses"
	"strconv"
)

func GetLimitQueryParam(r *http.Request, defaultLimit, maxLimit int) (int, *responses.Response) {
	limit := r.URL.Query().Get("limit")

	if limit == "" {
		return defaultLimit, nil
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return 0, responses.NewBadRequestResponse(fmt.Sprintf("Invalid `limit` format %s. Should be integer", limit))
	}

	if limitInt > maxLimit {
		return 0, responses.NewBadRequestResponse(fmt.Sprintf("Invalid `limit` %d. Should be from 0 to %d", limitInt, maxLimit))
	}

	return limitInt, nil
}

func GetOffsetQueryParam(r *http.Request, defaultOffset int) (int, *responses.Response) {
	offset := r.URL.Query().Get("offset")

	if offset == "" {
		return defaultOffset, nil
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return 0, responses.NewBadRequestResponse(fmt.Sprintf("Invalid `offset` format %s. Should be integer", offset))
	}

	if offsetInt < 0 {
		return 0, responses.NewBadRequestResponse(fmt.Sprintf("Invalid `offset` %d. Should be greater or equal 0", offsetInt))
	}

	return offsetInt, nil
}
