package services

import (
	"auth-system/interfaces"
	"auth-system/models"
	"crypto/rand"
	"encoding/hex"
	"net/smtp"
	"time"
)

type EmailService struct {
	server         string
	port           string
	senderEmail    string
	senderPassword string
	repo interfaces.Repository
}

func NewEmailService(server string, port string, email string, pass string, r interfaces.Repository) EmailService {
	return EmailService{
		server:         server,
		port:           port,
		senderEmail:    email,
		senderPassword: pass,
		repo: r,
	}
}

func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 15)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Encode the random bytes into a hexadecimal string
	return hex.EncodeToString(bytes), nil
}

func (e EmailService) SendEmail(user *models.User) (error) {
	recipientEmail := user.Email

	token, err := GenerateSecureToken()
	if err != nil {
		return err
	}

	resetToken := &models.ResetToken {
		Token: token,
		UserID: user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}
	
	err = e.repo.InsertResetToken(resetToken)
	if err != nil {
		return err
	}

	url := "http://localhost:8080/user/reset?token=" + token

	subject := "Subject: Password Recovery\r\n"
	body := "To reset your password please click the following link:\n" + url
	message := []byte(subject + "\r\n" + body)

	auth := smtp.PlainAuth("", e.senderEmail, e.senderPassword, e.server)

	err = smtp.SendMail(e.server+":"+e.port, auth, e.senderEmail, []string{recipientEmail}, message)
	return err
}