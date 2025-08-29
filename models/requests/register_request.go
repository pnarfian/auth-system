package request

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName string `json:"lastName" validate:"required"`
	TelephoneNo string `json:"telephoneNo" validate:"required"`
	Email string `json:"email" validate:"required,email"` 
	Password  string `json:"password" validate:"required"`
}