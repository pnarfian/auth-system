package usecases

import (
	"auth-system/interfaces"
	"auth-system/models"
	request "auth-system/models/requests"
	"auth-system/services"
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo interfaces.Repository
	secretKey []byte
	email services.EmailService
	client *redis.Client
	ctx context.Context
}

func NewUseCase(r interfaces.Repository, s string, e services.EmailService, c *redis.Client, ctx context.Context) (UseCase) {
	return UseCase{repo: r, secretKey: []byte(s), email: e, client: c, ctx: ctx}
}

func (u UseCase) Register(data *request.RegisterRequest) (error) {
	user, _ := u.repo.GetUserByEmail(data.Email)
	if user.Username != "" {
		return errors.New("user already exists")
	}

	user, _ = u.repo.GetUserByUsername(data.Username)
	if user.Username != "" {
		return errors.New("user already exists")
	}

	if !u.ValidatePassword(data.Password) {
		return errors.New("invalid password")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)

	if err != nil {
		return err
	}

	user = &models.User{
		Username: data.Username,
		FirstName: data.FirstName,
		LastName: data.LastName,
		TelephoneNo: data.TelephoneNo,
		Email: data.Email,
		Password: string(password),
		CreatedAt: time.Now(),
		UpdateAt: time.Now(), 
	}

	err = u.repo.InsertUser(user)

	return err
}

func (u UseCase) Login(data *request.LoginRequest) (string, error) {
	user, err := u.repo.GetUserByUsername(data.Username)
	if err != nil {
		return "", err
	} else if user.Username == ""{
		return "", errors.New("incorrect username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return "", errors.New("incorrect password")
	}

	accessToken := &models.Access_Token{
		UserID: user.ID,
		Expires_at: time.Now().Add(time.Hour * 6),
	}

	tokenID, err := u.repo.InsertToken(accessToken)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": tokenID,
			"exp": time.Now().Add(time.Hour * 6).Unix(),
	})

	tokenString, err := token.SignedString(u.secretKey)
	if err != nil {
		return "", err
	}

	err = u.client.Set(u.ctx, tokenString, user.ID, time.Minute * 30).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u UseCase) Logout(userID int) (error) {
	err := u.repo.RevokeToken(userID)

	return err
}

func (u UseCase) Forgot(data *request.ForgotRequest) (error) {
	user, err := u.repo.GetUserByUsername(data.Username)
	if err != nil {
		return err
	}

	if user.Username == "" {
		return nil
	}

	err = u.email.SendEmail(user)

	return err
}

func (u UseCase) Reset(data *request.ResetRequest, token string) (error) {
	resetToken, err := u.repo.GetResetToken(token)
	if err != nil {
		return errors.New("invalid token")
	}
	
	if resetToken.ExpiresAt.Before(time.Now()) || resetToken.IsUsed {
		return errors.New("invalid token")
	}

	user, err := u.repo.GetUser(int(resetToken.UserID))
	if err != nil{
		return err
	}
	if user.Username == "" {
		return errors.New("user not found")
	}

	if !u.ValidatePassword(data.NewPassword) {
		return errors.New("invalid password")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 10)
	if err != nil {
		return err
	}

	user.Password = string(password)
	err = u.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	err = u.repo.UpdateResetToken(resetToken)
	return err
}

func (u UseCase) Delete(userID int) (error) {
	user, err := u.repo.GetUser(userID)
	if err != nil {
		return err
	}

	err = u.repo.RevokeToken(userID)
	if err != nil {
		return err
	}

	err = u.repo.DeleteUser(user)
	return err
}

func (u UseCase) ValidatePassword(password string) (bool) {
	if len(password) < 8 || len(password) > 64 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?~]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}