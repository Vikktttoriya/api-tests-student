package tests

import (
	"fmt"
	"log"

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
