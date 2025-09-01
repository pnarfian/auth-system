package repositories

import (
	"auth-system/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (Repository) {
	return Repository{db: db}
}

func (r Repository) GetUser(userID int) (*models.User, error) {
	var user *models.User

	err := r.db.Where("id = ?", userID).Where("is_deleted = false").Find(&user).Error

	return user, err
}

func (r Repository) GetUserByUsername(username string) (*models.User, error) {
	var user *models.User

	err := r.db.Where("username = ?", username).Where("is_deleted = false").Find(&user).Error

	return user, err
}

func (r Repository) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User

	err := r.db.Where("email = ?", email).Where("is_deleted = false").Find(&user).Error

	return user, err
}
 
func (r Repository) InsertUser(data *models.User) (error) {
	err := r.db.Create(data).Error

	return err
}

func (r Repository) UpdateUser(data *models.User) (error) {
	err := r.db.Save(data).Error

	return err
}

func (r Repository) DeleteUser(data *models.User) (error) {
	data.IsDeleted = true
	err := r.db.Save(data).Error

	return err
}

func (r Repository) GetToken(tokenID int) (*models.Access_Token, error) {
	var token *models.Access_Token
	err := r.db.Where("id = ?", tokenID).Find(&token).Error

	return token, err
}

func (r Repository) InsertToken(data *models.Access_Token) (int, error) {
	err := r.db.Create(data).Error

	if err != nil {
		return 0, err
	}

	return int(data.ID), nil
}

func (r Repository) RevokeToken(userID int) (error) {
	var tokens []*models.Access_Token

	err := r.db.Where("user_id = ?", userID).Find(&tokens).Error

	if err != nil {
		return err
	}

	for _, currToken := range tokens {
		currToken.Revoked = true
		err = r.db.Save(currToken).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func (r Repository) InsertResetToken(data *models.ResetToken) (error) {
	err := r.db.Create(data).Error

	return err
}

func (r Repository) GetResetToken(token string) (*models.ResetToken, error) {
	var resetToken *models.ResetToken

	err := r.db.Where("token = ?", token).Where("is_used = false").Find(&resetToken).Error

	return resetToken, err
}

func (r Repository) UpdateResetToken(data *models.ResetToken) (error) {
	data.IsUsed = true

	err := r.db.Save(data).Error

	return err
}