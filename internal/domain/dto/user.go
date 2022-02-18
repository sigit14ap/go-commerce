package dto

type CreateUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

type UpdateUserDTO struct {
	Name string `json:"name"`
}

type UpdateUserInput struct {
	Name string `json:"name"`
}

type SignUpDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}
