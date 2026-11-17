package fixtures

import (
	"testing"

	"api-tests-template/internal/constants/apiErrors"
	"api-tests-template/internal/constants/tokens"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	userFirstName = "Автотест"
	userLastName  = "Покупатель"
)

type UnauthorizedCase struct {
	Name            string
	Token           string
	ExpectedMessage string
}

// UnauthorizedCases возвращает варианты запроса без авторизации: без токена и с невалидным токеном
func UnauthorizedCases() []UnauthorizedCase {
	return []UnauthorizedCase{
		{Name: "без токена", Token: tokens.Missing, ExpectedMessage: apiErrors.MessageAuthorizationRequired},
		{Name: "с невалидным токеном", Token: tokens.Invalid, ExpectedMessage: apiErrors.MessageAuthorizationInvalid},
	}
}

// NewUser регистрирует нового пользователя с уникальным email и возвращает его токен и данные
func NewUser(t *testing.T) authModels.LoginOkResponse {
	return auth.Register(t, authModels.RegisterRequest{
		Email:     utils.RandomEmail(),
		Password:  utils.RandomPassword(),
		FirstName: userFirstName,
		LastName:  userLastName,
	})
}
