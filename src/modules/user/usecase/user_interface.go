package usecase

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/dto"
)

type UserUsecase interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req dto.RegisterRequest) error
}
