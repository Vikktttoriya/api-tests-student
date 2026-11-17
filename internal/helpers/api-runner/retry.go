package api_runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"api-tests-template/internal/utils"
)

const (
	maxRequestAttempts = 5
	retryBaseDelay     = 500 * time.Millisecond
)

// httpClient HTTP-клиент для всех запросов тестов: повторяет запросы, отклонённые nginx с кодом 503
var httpClient = &http.Client{Transport: retryOnUnavailableTransport{base: http.DefaultTransport}}

// retryOnUnavailableTransport повторяет запрос, если nginx ответил 503 HTML-страницей ошибки
// (превышен лимит запросов или сервис временно недоступен)
type retryOnUnavailableTransport struct {
	base http.RoundTripper
}

func (rt retryOnUnavailableTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	body, err := readAndCloseBody(request)
	if err != nil {
		return nil, err
	}

	for attempt := 1; ; attempt++ {
		attemptRequest := request.Clone(request.Context())
		if body != nil {
			attemptRequest.Body = io.NopCloser(bytes.NewReader(body))
		}

		response, err := rt.base.RoundTrip(attemptRequest)
		if err != nil || !isRejectedByProxy(response) || attempt == maxRequestAttempts {
			return response, err
		}
		_ = response.Body.Close()

		delay := time.Duration(attempt) * retryBaseDelay
		utils.LogWithLabelAndTimeStamp("Retry", fmt.Sprintf("%s %s: 503 от nginx, повтор %d из %d через %s",
			request.Method, request.URL.Path, attempt, maxRequestAttempts-1, delay))
		time.Sleep(delay)
	}
}

// isRejectedByProxy проверяет, что ответ — 503 от nginx, а не ответ самого сервиса (сервис отвечает JSON)
func isRejectedByProxy(response *http.Response) bool {
	return response.StatusCode == http.StatusServiceUnavailable &&
		!strings.Contains(response.Header.Get("Content-Type"), "application/json")
}

// readAndCloseBody читает тело запроса целиком, чтобы его можно было отправить повторно
func readAndCloseBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}
	defer request.Body.Close()

	return io.ReadAll(request.Body)
}
