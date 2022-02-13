package dto

import "mime/multipart"

type CreateCategoryDTO struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Icon        *multipart.FileHeader `form:"icon" binding:"required"`
}

type UpdateCategoryDTO struct {
	Name        string                `form:"name"`
	Description *string               `form:"description"`
	Icon        *multipart.FileHeader `form:"icon" binding:"required"`
}

type UpdateCategoryInput struct {
	Name        string
	Description *string
	Icon        *multipart.FileHeader
}
