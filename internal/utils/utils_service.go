package utils

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type IUtilsService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, hash string) bool
}

type utilsService struct {
	ctx   context.Context
	debug bool
}

func NewUtilsService(ctx context.Context, debug bool) IUtilsService {
	return &utilsService{
		ctx:   ctx,
		debug: debug,
	}
}

func (s *utilsService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *utilsService) CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
