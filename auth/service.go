package auth

import (
	"log/slog"
	"time"

	"hrms/models"
	"hrms/users"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  users.UserRepository
	jwtSecret string
	log       *slog.Logger
}

func NewAuthService(userRepo users.UserRepository, log *slog.Logger, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		log:       log,
	}
}

func (as *AuthService) Register(user *UserAuth) error {
	userFound, err := as.userRepo.FindByEmail(user.Email)
	if userFound != nil {
		return ErrUserExists
	}
	if err != nil {
		return err
	}
	hashedPassword, err := as.hashPassword(user.Password)
	if err != nil {
		return err
	}
	userDB := &models.User{
		Email:    user.Email,
		Password: hashedPassword,
	}
	if err := as.userRepo.SaveUser(userDB); err != nil {
		return ErrSaveUser
	}

	return nil
}

func (as *AuthService) Login(user *UserAuth) (string, error) {
	userDB, err := as.userRepo.FindByEmail(user.Email)
	passIsValid := as.passwordValid(userDB.Password, user.Password)

	if userDB != nil && passIsValid {
		return as.generateToken(userDB.ID)
	} else if userDB == nil && err == nil || userDB != nil && !passIsValid {
		return "", ErrInvalidCredentials
	} else {
		return "", err
	}
}

func (as *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(as.jwtSecret))
	if err != nil {
		as.log.Error("Failed to generate token", "error", err.Error())
		return "", ErrGenerateToken
	}

	return token, nil
}

func (as *AuthService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		as.log.Error("Failed to hash password", "error", err.Error())
		return "", ErrHashPassword
	}
	return string(hashedBytes), nil
}

func (as *AuthService) passwordValid(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		as.log.Error("Failed to compare password", "error", err.Error())
	}
	return err == nil
}
