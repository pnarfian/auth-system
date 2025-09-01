package interfaces

import (
	"auth-system/models"
	request "auth-system/models/requests"
)

type UseCase interface {
	Register(data *request.RegisterRequest) (error)
	Login(data *request.LoginRequest) (string, error)
	Logout(userID int) (error)
	Forgot(data *request.ForgotRequest) (error)
	Reset(data *request.ResetRequest, token string) (error)
	Delete(userID int) (error)
	ValidatePassword(password string) (bool)
}

type Repository interface {
	GetUser(userID int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	InsertUser(data *models.User) (error)
	UpdateUser(data *models.User) (error)
	DeleteUser(data *models.User) (error)
	GetToken(tokenID int) (*models.Access_Token, error)
	InsertToken(data *models.Access_Token) (int, error)
	RevokeToken(userID int) (error)
	InsertResetToken(data *models.ResetToken) (error)
	GetResetToken(token string) (*models.ResetToken, error)
	UpdateResetToken(data *models.ResetToken) (error)
}