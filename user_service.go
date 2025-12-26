package main

import (
	"context"
	"errors"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (uservice *UserService) GetAllUsers(ctx context.Context) ([]User, error) {
	return uservice.repo.GetAll(ctx)
}

func (uservice *UserService) GetUserById(ctx context.Context, id int) (*User, error) {
	if id < 1 {
		return nil, ErrIdMustBeGtZero
	}
	return uservice.repo.GetById(ctx, id)
}

func (uservice *UserService) CreateNewUser(ctx context.Context, name string, password string) (int, error) {
	if len(name) < 1 {
		return 0, ErrIdMustBeGtZero
	}
	if len(password) < 6 {
		return 0, ErrPasswordMustBeGt6
	}

	encryptPassword, err := Encrypter(password)
	if err != nil {
		return 0, err
	}

	newUser := User{
		Name:     name,
		Password: encryptPassword,
	}
	return uservice.repo.Create(ctx, newUser)
}

func (uservice *UserService) RenameUser(ctx context.Context, id int, newName string) error {
	if len(newName) < 1 {
		return ErrLenNameIsZero
	}
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	err := uservice.repo.UpdateName(ctx, id, newName)
	if err != nil {
		return err
	}
	return nil
}

func (uservice *UserService) ChangeUsersPass(ctx context.Context, id int, oldPass string, newPass string) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	if len(oldPass) < 6 {
		return errors.New("old password must be greater than 6")
	}
	if len(newPass) < 6 {
		return errors.New("new password must be greater than 6")
	}

	if oldPass == newPass {
		return ErrNewPasswordIsSame
	}

	user, err := uservice.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	if !CompareHashAndPassword(user.Password, oldPass) {
		return ErrOldPasswordIsWrong
	}

	newHashPass, err := Encrypter(newPass)
	if err != nil {
		return err
	}
	user.Password = newHashPass
	return uservice.repo.UpdatePassword(ctx, id, newHashPass)

}

func (uservice *UserService) DeleteUser(ctx context.Context, id int) error {
	if id < 1 {
		return ErrIdMustBeGtZero
	}
	return uservice.repo.Delete(ctx, id)
}
