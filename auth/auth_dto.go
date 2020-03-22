package auth

type LoginDto struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginOauthDto struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}