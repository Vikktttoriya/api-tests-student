package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonModels "api-tests-template/internal/models"
)

// ErrorResponseEquals проверяет код и текст ошибки в ответе API
func ErrorResponseEquals(t *testing.T, actual commonModels.ErrorResponse, expectedError, expectedMessage string) {
	t.Helper()

	require.Equalf(t, expectedError, actual.Error, "Некорректное значение error в ответе (message: %q)", actual.Message)
	require.Equal(t, expectedMessage, actual.Message, "Некорректное значение message в ответе")
	logCheck("получена ожидаемая ошибка %s: %s", actual.Error, actual.Message)
}
