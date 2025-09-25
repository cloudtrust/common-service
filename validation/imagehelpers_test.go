package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageMimeType(t *testing.T) {
	t.Run("JPEG image", func(t *testing.T) {
		var jpegImage = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
		mimeType, err := GetImageMimeType(jpegImage)
		assert.NoError(t, err)
		assert.Equal(t, "image/jpeg", mimeType)
	})
	t.Run("PNG image", func(t *testing.T) {
		var pngImage = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		mimeType, err := GetImageMimeType(pngImage)
		assert.NoError(t, err)
		assert.Equal(t, "image/png", mimeType)
	})
	t.Run("GIF image", func(t *testing.T) {
		var gifImage = []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}
		mimeType, err := GetImageMimeType(gifImage)
		assert.NoError(t, err)
		assert.Equal(t, "image/gif", mimeType)
	})
	t.Run("SVG image", func(t *testing.T) {
		var svgImage = []byte{0x3C, 0x73, 0x76, 0x67, 0x20, 0x78, 0x6D, 0x6C}
		mimeType, err := GetImageMimeType(svgImage)
		assert.NoError(t, err)
		assert.Equal(t, "image/svg+xml", mimeType)
	})
	t.Run("SVG image with <?", func(t *testing.T) {
		var svgImage = []byte{0x3C, 0x3F, 0x78, 0x6D, 0x6C, 0x20, 0x76, 0x65}
		mimeType, err := GetImageMimeType(svgImage)
		assert.NoError(t, err)
		assert.Equal(t, "image/svg+xml", mimeType)
	})
	t.Run("Unknown image type", func(t *testing.T) {
		var unknownImage = []byte{0x00, 0x01, 0x02, 0x03, 0x04}
		mimeType, err := GetImageMimeType(unknownImage)
		assert.NoError(t, err)
		assert.Equal(t, "application/octet-stream", mimeType)
	})
	t.Run("Too short image data", func(t *testing.T) {
		var shortImage = []byte{0xFF}
		mimeType, err := GetImageMimeType(shortImage)
		assert.Error(t, err)
		assert.Equal(t, "", mimeType)
	})
}
