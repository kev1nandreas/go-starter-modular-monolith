package repository

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/entity"
)

type UserRepository interface {
	FindByUsername(username string) (*entity.User, error)
	Create(user *entity.User) error
}
