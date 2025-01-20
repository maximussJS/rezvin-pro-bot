package requests

import (
	"rezvin-pro-bot/models"
)

type CreateCategoryRequest struct {
	Price float64 `json:"price" validate:"required,gt=0"`
	Name  string  `json:"name" validate:"required,min=3,max=255"`
}

func (r CreateCategoryRequest) ToCategory() models.Category {
	return models.Category{
		Name:  r.Name,
		Price: r.Price,
	}
}
