package dto

type ChangePasswordRequest struct {
	CurrentPassword         string `json:"current_password" validate:"required,min=6"`
	NewPassword             string `json:"new_password" validate:"required,min=6"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,min=6"`
}

type ValidateTokenUserResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ValidateTokenResponse struct {
	User ValidateTokenUserResponse `json:"user"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}
