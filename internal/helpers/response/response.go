package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ReadBodyWithStatus читает и закрывает тело ответа, проверяет статус-код.
// При несовпадении статуса в сообщение попадает тело ответа, чтобы была видна причина падения
func ReadBodyWithStatus(t *testing.T, raw *http.Response, expectedStatusCode int) []byte {
	t.Helper()
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")
	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, body)

	return body
}

// ParseWithStatus проверяет статус-код ответа и разбирает JSON-тело в модель T
func ParseWithStatus[T any](t *testing.T, raw *http.Response, expectedStatusCode int) T {
	t.Helper()

	body := ReadBodyWithStatus(t, raw, expectedStatusCode)

	var result T
	require.NoErrorf(t, json.Unmarshal(body, &result), "Не удалось разобрать тело ответа: %s", body)

	return result
}
