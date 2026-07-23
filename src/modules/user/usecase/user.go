package usecase

import (
	"errors"

	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/dto"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/entity"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/repository"
	"golang.org/x/crypto/bcrypt"

	"github.com/kev1nandreas/go-rest-api-template/src/shared/auth"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return &dto.LoginResponse{Token: token}, nil
}

func (u *userUsecase) Register(req dto.RegisterRequest) error {
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return errors.New("could not hash password")
	}

	newUser := &entity.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := u.repo.Create(newUser); err != nil {
		return errors.New("could not save user")
	}

	return nil
}
