package handlers

import (
	"encoding/json"
	"net/http"
	"project/internal/models"
	"project/internal/modules/auth/service"
)

type Auther interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type AuthController struct {
	service service.Auther
}

func NewAuthController(service service.Auther) Auther {
	return &AuthController{
		service: service,
	}
}

// RegisterHandler — регистрация нового пользователя.
//
// @Summary Регистрация нового пользователя
// @Description Создаёт нового пользователя с указанными логином и паролем.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "Данные пользователя"
// @Success 201 "Пользователь успешно зарегистрирован"
// @Failure 400 "Некорректный запрос или пользователь уже существует"
// @Failure 500 "Ошибка сервера"
// @Router /register [post]
func (a *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := a.service.Register(creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// LoginHandler — аутентификация пользователя и выдача JWT.
//
// @Summary Авторизация пользователя
// @Description Проверяет логин и пароль, возвращает JWT-токен для доступа к защищённым эндпоинтам.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "Имя пользователя и пароль"
// @Failure 200 "Пользователь не существует или пароль не совпадает"
// @Failure 400 "Некорректный запрос"
// @Failure 500 "Ошибка сервера"
// @Router /login [post]
func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := a.service.Login(creds)
	if _, ok := message["error"]; ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
		return
	}

	json.NewEncoder(w).Encode(message)
}
