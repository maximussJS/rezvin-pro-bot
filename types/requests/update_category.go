package requests

import (
	"rezvin-pro-bot/models"
)

type UpdateCategoryRequest struct {
	Price float64 `json:"price" validate:"gt=0"`
	Name  string  `json:"name" validate:"min=3,max=255"`
}

func (r UpdateCategoryRequest) ToCategory() models.Category {
	return models.Category{
		Name:  r.Name,
		Price: r.Price,
	}
}
