package dto

type RegisterPlayerRequest struct {
	Name     string `json:"name" binding:"required" example:"Иван"`
	Email    string `json:"email" binding:"required,email" example:"test@mail.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
}

