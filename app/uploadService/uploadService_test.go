package uploadService

import "testing"

func TestCleanFilename1(t *testing.T) {
	badString := "bad%file^&name0047"
	expect1 := "badfilename0047"
	result := CleanFilename(badString)
	if result != expect1 {
		t.Errorf("CleanFilename test failed exptectd: %s recieved: %s", expect1, result)
	}
}

func TestCleanFilename2(t *testing.T) {
	goodName := "Good_Name"
	result := CleanFilename(goodName)
	if result != goodName {
		t.Errorf("CleanFilename test failed exptectd: %s recieved: %s", goodName, result)
	}
}

func TestIsValidMimeType(t *testing.T) {
	valid := []string{"image/png", "image/jpeg", "video/mp4"}
	notValid := []string{"image/foo", "video/mp3", "something"}
	for _, v := range valid {
		if !IsValidMimeType(v) {
			t.Errorf("%s should be a valid mime type", v)
			break
		}
	}

	for _, nv := range notValid {
		if IsValidMimeType(nv) {
			t.Errorf("%s is NOT a valid mime type", nv)
		}
	}
}

func TestGetRootDir(t *testing.T) {
	image := "image/png"
	video := "video/mp4"
	inValid := "something/something"

	if GetRootDir(image) != "images" {
		t.Errorf("%s should return images dir", image)
	}
	if GetRootDir(video) != "videos" {
		t.Errorf("%s should return video dir", video)
	}

	if GetRootDir(inValid) != "unknown" {
		t.Errorf("%s should return unknown dir", inValid)
	}
}
