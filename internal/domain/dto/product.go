package dto

type CreateProductDTO struct {
	Name        string   `form:"name" binding:"required"`
	Description string   `form:"description" binding:"required"`
	Price       float64  `form:"price" binding:"required"`
	CategoryID  string   `form:"category_id" binding:"required"`
	Images      []string `form:"images"`
	Weight      int64    `form:"weight" binding:"required"`
}

type UpdateProductDTO struct {
	Name        string    `form:"name" binding:"required"`
	Description *string   `form:"description" binding:"required"`
	Price       *float64  `form:"price" binding:"required"`
	CategoryID  string    `form:"category_id" binding:"required"`
	Images      *[]string `form:"images"`
	Weight      int64     `form:"weight" binding:"required"`
}

type UpdateProductInput struct {
	Name        string
	Description *string
	Price       *float64
	CategoryID  string
	Images      *[]string `form:"images"`
}
