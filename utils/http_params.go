package utils

import (
	"github.com/julienschmidt/httprouter"
	"rezvin-pro-bot/types/responses"
	"strconv"
)

func GetIdParam(p httprouter.Params) (uint, *responses.Response) {
	idString := p.ByName("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		return 0, responses.NewBadRequestResponse("Invalid `id` format. Should be integer greater than 0")
	}

	return uint(id), nil
}

func GetCategoryIdParam(p httprouter.Params) (uint, *responses.Response) {
	idString := p.ByName("categoryId")

	id, err := strconv.Atoi(idString)

	if err != nil {
		return 0, responses.NewBadRequestResponse("Invalid `categoryId` format. Should be integer greater than 0")
	}

	return uint(id), nil
}
