package request

type ResetRequest struct {
	NewPassword  string `json:"newPassword" validate:"required"`
}
