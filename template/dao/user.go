package dao

import (
	"golang.org/x/crypto/bcrypt"
)

func (d *Dao) GetUserIsIncognito(emailAes string) (isIncognito bool) {
	var user model.AuthUser
	u, err := user.GetUser(d.db, emailAes)
	if err != nil {
		isIncognito = false
		return
	} else {
		return u.IsIncognito
	}
}
func (d *Dao) GetRoleWhileLogin(email string) (err error, role string) {
	var user model.AuthUser
	u, err := user.GetUser(d.db, email)
	if err != nil {
		return err, ""
	}
	return nil, u.Role
}
func (d *Dao) IsEmailRegistered(email string) (err error, isRegistered bool) {
	var user model.AuthUser
	u, err := user.GetUser(d.db, email)
	if u.Email == "" {
		return err, false
	}
	return nil, true
}
func (d *Dao) IsEmailAndPasswordMatch(email string, pwd string) (isEmailAndPasswordMatch bool) {
	var user model.AuthUser
	u, _ := user.GetUser(d.db, email)
	result := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
	if result == nil {
		isEmailAndPasswordMatch = true
	} else {
		isEmailAndPasswordMatch = false
	}
	return
}
func (d *Dao) GetUser(emailAes string) (u model.AuthUser, err error) {
	var user model.AuthUser
	u, err = user.GetUser(d.db, emailAes)
	return
}
func (d *Dao) Register(user *model.AuthUser) (err error) {
	var u model.AuthUser
	err = u.CreateUser(d.db, user)
	return
}

func (d *Dao) ActiveEmail(emailAes string) (err error) {
	var user model.AuthUser
	err = user.ActivateEmail(d.db, emailAes)
	return
}

func (d *Dao) UpdateUser(u *model.AuthUser) (err error) {
	var user model.AuthUser
	err = user.UpdateUserInfo(d.db, u)
	return
}

func (d *Dao) ResetPassword(emailAes, passwordHash string) (err error) {
	var user model.AuthUser
	err = user.ResetUserPassword(d.db, emailAes, passwordHash)
	return
}
