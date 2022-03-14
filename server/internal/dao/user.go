package dao

import (
	"ipashare/internal/model"
	"ipashare/pkg/tools"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newUser(db *gorm.DB) *user {
	u := &user{db}
	var users []model.User
	if u.db.Find(&users).Error == nil && users != nil && len(users) == 0 {
		salt := uuid.New().String()
		u.db.Create(&model.User{
			Username: "admin",
			Password: tools.MD5LowercaseEncode("e10adc3949ba59abbe56e057f20f883e" + salt),
			Salt:     salt,
		})
	}
	return u
}

type user struct {
	db *gorm.DB
}

var _ model.UserStore = (*user)(nil)

func (u *user) Query(username string) (*model.User, error) {
	var user model.User
	err := u.db.Where("username = ?", username).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *user) Update(oldUsername, newUsername, password string) error {
	salt := uuid.New().String()
	return u.db.Model(&model.User{}).
		Where("username = ?", oldUsername).
		Updates(map[string]interface{}{
			"username": newUsername,
			"password": tools.MD5LowercaseEncode(password + salt),
			"salt":     salt,
		}).Error
}
