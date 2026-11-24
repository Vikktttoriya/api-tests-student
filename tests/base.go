package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// SetupSuite действия перед запуском сьюта с тестами
func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

// TearDownSuite действия после запуска сьюта с тестами
func TearDownSuite() {
	log.Println("Tear down suite")
}

// Precondition выводит текст в консоль с меткой [Precondition]
func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}

// Step выводит в консоль описание шага теста с меткой [Step]
func Step(format string, args ...any) {
	utils.LogWithLabelAndTimeStamp("Step", fmt.Sprintf(format, args...))
}

// LoginTestUser авторизует основного пользователя TestBoard с кредами из переменных окружения TEST_LOGIN и TEST_PASSWORD
func LoginTestUser(t *testing.T) models.LoginOkResponse {
	Precondition("Авторизация пользователя с кредами из переменных окружения")
	return auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}
