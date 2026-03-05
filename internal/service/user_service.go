package service

import (
	"amass/internal/models"
	"amass/internal/repository"
	"amass/internal/utils"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IUserService interface {
	Create(u *models.UserLogin) error
	Login(u *models.UserLogin) (*models.LoginResponse, error)
	GetUser(username string) (*models.User, error)
}

type userService struct {
	ctx      context.Context
	debug    bool
	userRepo repository.IUserRepository
	utils    utils.IUtilsService
}

func NewUserService(
	ctx context.Context, debug bool, repo AllRepository, utilsService utils.IUtilsService,
) IUserService {
	return &userService{
		ctx:      ctx,
		debug:    debug,
		userRepo: repo.IUserRepository,
		utils:    utilsService,
	}
}

func (s *userService) Create(u *models.UserLogin) error {
	if u == nil {
		return fmt.Errorf("user cannot be nil")
	}

	// ปกติใช้ secret key ในการ config KEY สำหรับการเข้ารหัสและถอดรหัส แต่ในที่นี้จะใช้ os.Getenv("KEY") แทนเพื่อความง่ายในการทดสอบ
	encryptedPassword, err := s.utils.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %v", err)
	}

	user := models.User{
		Username:     u.Username,
		PasswordHash: encryptedPassword,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.userRepo.Create(&user)
}

func (s *userService) Login(u *models.UserLogin) (*models.LoginResponse, error) {

	userLogin := &models.UserLogin{
		Username: u.Username,
	}

	user, err := s.userRepo.GetUser(userLogin)
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}

	encryptedPassword := s.utils.CheckPassword(u.Password, user.PasswordHash)

	if !encryptedPassword {
		return nil, fmt.Errorf("invalid password")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("KEY")
	}
	if jwtSecret == "" {
		return nil, fmt.Errorf("jwt secret is not configured")
	}

	jwtExpireHours := 24
	if rawExpire := os.Getenv("JWT_EXPIRES_HOURS"); rawExpire != "" {
		value, parseErr := strconv.Atoi(rawExpire)
		if parseErr != nil || value <= 0 {
			return nil, fmt.Errorf("invalid JWT_EXPIRES_HOURS")
		}
		jwtExpireHours = value
	}

	now := time.Now()
	expireAt := now.Add(time.Duration(jwtExpireHours) * time.Hour)

	claims := jwt.RegisteredClaims{
		Subject:   user.Username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expireAt),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate jwt: %v", err)
	}

	return &models.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expireAt.Unix(),
	}, nil
}

func (s *userService) GetUser(username string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	userReq := &models.UserLogin{Username: username}
	user, err := s.userRepo.GetUser(userReq)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	return user, nil
}
