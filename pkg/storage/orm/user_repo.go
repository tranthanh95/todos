package orm

import (
	"errors"
	"fmt"

	"todo-lists/pkg/logger"
	"todo-lists/pkg/user"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type userRepo struct {
	db  *gorm.DB
	log logger.LogInfoFormat
}

func NewUserRepo(db *gorm.DB, log logger.LogInfoFormat) user.Repository {
	return &userRepo{db, log}
}

func (u *userRepo) Delete(id string) error {
	u.log.Debugf("deleting the user with id : %s", id)

	if u.db.Delete(&user.User{}, "user_id = ?", id).Error != nil {
		errMsg := fmt.Sprintf("error while deleting the user with id : %s", id)
		u.log.Errorf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (u *userRepo) GetAll() ([]*user.User, error) {
	u.log.Debug("get all the users")

	users := make([]*user.User, 0)
	err := u.db.Find(&users).Error
	if err != nil {
		u.log.Debug("no single user found")
		return nil, err
	}
	return users, nil
}

func (u *userRepo) GetByID(id string) (*user.User, error) {
	u.log.Debugf("get user details by id : %s", id)

	user := &user.User{}
	err := u.db.Where("user_id = ?", id).First(&user).Error
	if err != nil {
		u.log.Errorf("user not found with id : %s, reason : %v", id, err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) Store(usr *user.User) error {
	fmt.Println("tao user")
	fmt.Println("asdasd", usr.Email)
	user := &user.User{}
	err := u.db.Where("email = ?", usr.Email).First(&user).Error

	fmt.Println("asdasd ", user.ID)
	fmt.Println("asdasd ", user.Email)

	fmt.Println("err ", err.Error())

	if err != nil && err.Error() != "record not found" {
		fmt.Println("dkmmm")
		return errors.New("please signin")
	}

	if user.Email != "" {
		return errors.New("user exist")
	}

	fmt.Println("tao user voi password ", user.Password)

	u.log.Debugf("creating the user with email : %v", usr.Email)
	// Generates a hashed version of our password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(fmt.Sprintf("error hashing password: %v", err))
	}
	usr.Password = string(hashedPass)
	
	fmt.Println("asdasdas", string(hashedPass))
	err = u.db.Create(&usr).Error
	if err != nil {
		u.log.Errorf("error while creating the user, reason : %v", err)
		return err
	}
	return nil
}

func (u *userRepo) Update(usr *user.User) error {
	u.log.Debugf("updating the user, user_id : %v", usr.ID)

	err := u.db.Model(&usr).Updates(user.User{FirstName: usr.FirstName, LastName: usr.LastName, Password: usr.Password}).Error
	if err != nil {
		u.log.Errorf("error while updating the user, reason : %v", err)
		return err
	}
	return nil
}