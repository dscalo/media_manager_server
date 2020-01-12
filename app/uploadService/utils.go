package uploadService

import (
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func CleanFilename(name string) string {
	i := 0
	bytes := []byte(name)
	for _, b := range bytes {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') || (b == '_') || (b == '.') {
			bytes[i] = b
			i++
		}
	}
	return string(bytes[:i])
}

func OnlyLetters(s string) string {
	i := 0
	bytes := []byte(s)
	for _, b := range bytes {
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') {
			bytes[i] = b
			i++
		}
	}
	return string(bytes[:i])
}

func IsValidMimeType(mime string) bool {
	valid := true
	switch strings.TrimSpace(mime) {
	case "image/jpeg":
	case "image/gif":
	case "image/png":
	case "video/mpeg":
	case "video/ogg":
	case "video/mp4":
	default:
		valid = false
	}
	return valid
}

func GetRootDir(mimeType string) string {
	dir := "unknown"

	switch {
	case strings.Contains(mimeType, "video"):
		dir = "videos"
	case strings.Contains(mimeType, "image"):
		dir = "images"
	default:
		{
		}
	}
	return dir
}

func GetMimeType(f multipart.File) string {
	buffer := make([]byte, 512)
	_, err := f.Read(buffer)
	if err != nil {
		return ""
	}
	return http.DetectContentType(buffer)
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ParseTags(tagString string) []string {
	// used to eliminate dupes
	set := make(map[string]struct{})

	var tags []string
	if tagString == "" {
		return tags
	}

	for _, s := range strings.Split(tagString, ",") {
		cleaned := OnlyLetters(s)
		lower := strings.ToLower(cleaned)
		if _, ok := set[lower]; ok {
			continue
		} else {
			set[lower] = struct{}{}
			tags = append(tags, lower)
		}
	}
	return tags
}
