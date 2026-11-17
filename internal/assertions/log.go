package assertions

import (
	"fmt"

	"api-tests-template/internal/utils"
)

// logCheck выводит в консоль результат успешной проверки с меткой [Check]
func logCheck(format string, args ...any) {
	utils.LogWithLabelAndTimeStamp("Check", fmt.Sprintf(format, args...))
}
