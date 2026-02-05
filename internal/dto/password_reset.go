package dto

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetVerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

type PasswordResetEmailRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Token            string `json:"token" binding:"required"`
	ResetURL         string `json:"reset_url" binding:"omitempty,url"`
	ExpiresInMinutes int    `json:"expires_in_minutes" binding:"omitempty,gte=1,lte=1440"`
}
