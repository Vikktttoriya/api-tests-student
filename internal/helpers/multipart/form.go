package multipart

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// BuildForm собирает тело multipart/form-data из текстовых полей и файлов и возвращает его вместе с Content-Type запроса.
// Content-Type каждого файла определяется по его содержимому (image/png, image/bmp, text/plain...)
func BuildForm(t *testing.T, fields map[string]string, fileFieldName string, filePaths []string) (body, contentType string) {
	t.Helper()

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	for name, value := range fields {
		require.NoErrorf(t, writer.WriteField(name, value), "Не удалось записать поле %s", name)
	}

	for _, filePath := range filePaths {
		content, err := os.ReadFile(filepath.Clean(filePath))
		require.NoErrorf(t, err, "Не удалось прочитать файл %s", filePath)

		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name=%q; filename=%q`, fileFieldName, filepath.Base(filePath)))
		header.Set("Content-Type", http.DetectContentType(content))

		part, err := writer.CreatePart(header)
		require.NoErrorf(t, err, "Не удалось добавить файл %s в форму", filePath)
		_, err = part.Write(content)
		require.NoErrorf(t, err, "Не удалось записать файл %s в форму", filePath)
	}

	require.NoError(t, writer.Close(), "Не удалось завершить multipart-форму")

	return buffer.String(), writer.FormDataContentType()
}
