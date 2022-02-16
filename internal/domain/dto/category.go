package dto

import "mime/multipart"

type ValidationCategoryDTO struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Icon        *multipart.FileHeader `form:"icon" binding:"required"`
}

type ValidationUpdateCategoryDTO struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
}

type CreateCategoryDTO struct {
	Name        string
	Description string
	Icon        string
}

type UpdateCategoryDTO struct {
	Name        string
	Description string
	Icon        string
}

type UpdateCategoryInput struct {
	Name        string
	Description string
	Icon        string
}
