package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testImageSize = 64

var testImageColors = []color.RGBA{
	{R: 220, G: 60, B: 60, A: 255},
	{R: 60, G: 160, B: 80, A: 255},
	{R: 60, G: 100, B: 220, A: 255},
	{R: 230, G: 180, B: 40, A: 255},
}

// GeneratePNGFiles создаёт count однотонных PNG-картинок разного цвета во временной папке теста
// и возвращает пути к ним. Папка удаляется автоматически после завершения теста
func GeneratePNGFiles(t *testing.T, count int) []string {
	t.Helper()

	dir := t.TempDir()
	paths := make([]string, 0, count)
	for i := range count {
		img := image.NewRGBA(image.Rect(0, 0, testImageSize, testImageSize))
		fill := image.NewUniform(testImageColors[i%len(testImageColors)])
		draw.Draw(img, img.Bounds(), fill, image.Point{}, draw.Src)

		filePath := filepath.Join(dir, fmt.Sprintf("photo_%d.png", i+1))
		file, err := os.Create(filePath)
		require.NoError(t, err, "Не удалось создать файл картинки")
		require.NoError(t, png.Encode(file, img), "Не удалось записать картинку")
		require.NoError(t, file.Close(), "Не удалось закрыть файл картинки")

		paths = append(paths, filePath)
	}

	return paths
}
