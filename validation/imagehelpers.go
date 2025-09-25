package validation

import cerrors "github.com/cloudtrust/common-service/v2/errors"

func GetImageMimeType(imageData []byte) (string, error) {
	var mimeType string
	// Check first bytes of the file to determine mime type (https://en.wikipedia.org/wiki/List_of_file_signatures)
	if len(imageData) > 2 {
		switch {
		case imageData[0] == 0xFF && imageData[1] == 0xD8:
			mimeType = "image/jpeg"
		case imageData[0] == 0x89 && imageData[1] == 0x50:
			mimeType = "image/png"
		case imageData[0] == 0x47 && imageData[1] == 0x49:
			mimeType = "image/gif"
		case (imageData[0] == 0x3C && imageData[1] == 0x73) || (imageData[0] == 0x3C && imageData[1] == 0x3F):
			mimeType = "image/svg+xml"
		default:
			mimeType = "application/octet-stream"
		}
	} else {
		return "", cerrors.Error{Message: "image data too short to determine mime type"}
	}
	return mimeType, nil
}
