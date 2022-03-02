package dao

import (
	"supersign/internal/model"
	"supersign/pkg/tools"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newUser(db *gorm.DB) *user {
	u := &user{db}
	_, err := u.Query("admin")
	if err == gorm.ErrRecordNotFound {
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

func (u *user) Update(username, password string) error {
	salt := uuid.New().String()
	return u.db.Model(&model.User{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"password": tools.MD5LowercaseEncode(password + salt),
			"salt":     salt,
		}).Error
}
