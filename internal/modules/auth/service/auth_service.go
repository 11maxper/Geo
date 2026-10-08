package service

import (
	"project/internal/infrastructure/cript"
	"project/internal/infrastructure/jwt"
	"project/internal/models"
	"project/internal/modules/auth/storage"

	"go.uber.org/zap"
)

type Auther interface {
	Register(creds models.Credentials) (string, error)
	Login(creds models.Credentials) map[string]string
}

type AuthService struct {
	tokenAuth *jwt.Jwt
	hasher    cript.Bcripter
	storage   storage.Auther
	logger    *zap.Logger
}

func NewAuthService(tokenAuth *jwt.Jwt, hasher cript.Bcripter, storage storage.Auther, logger *zap.Logger) Auther {
	return &AuthService{
		tokenAuth: tokenAuth,
		hasher:    hasher,
		storage:   storage,
		logger:    logger,
	}
}

// Register регистрирует пользователя, хешируют пароль
func (a *AuthService) Register(creds models.Credentials) (string, error) {
	// хешируем пароль
	hash, err := a.hasher.HashPassword(creds.Password)
	if err != nil {
		a.logger.Error("HashPassword()", zap.Error(err))
		return "", err
	}

	// сохраняем пользователя и хеш пароля в БД
	err = a.storage.Create(creds.Username, hash)
	if err != nil {
		a.logger.Error("storage.Create()", zap.Error(err))
		return "error register", err
	}

	return "user register", nil
}

// Login проверяет логин и пароль пользователя, выдает токен
func (a *AuthService) Login(creds models.Credentials) map[string]string {
	// проверка на наличие пользователя
	user, err := a.storage.FindByName(creds.Username)
	if err != nil {
		a.logger.Error("storage.FindByName()", zap.Error(err))
		return map[string]string{"error": "user not found"}
	}

	// сравнение введенного пароля с паролем из БД
	err = a.hasher.CompareHashAndPassword(user.PasswordHash, creds.Password)
	if err != nil {
		a.logger.Error("CompareHashAndPassword()", zap.Error(err))
		return map[string]string{"error": "invalid password"}
	}

	//получение токена
	tokenString := a.tokenAuth.GetTokenString(creds.Username)

	// выдача токена пользователю
	return map[string]string{"token": tokenString}
}
