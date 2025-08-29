package request

type ForgotRequest struct {
	Username string `json:"username" validate:"required"`
}