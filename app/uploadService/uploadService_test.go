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

func TestOnlyLetters(t *testing.T) {
	s := "123abc$%D^ E "
	expect := "abcDE"
	result := OnlyLetters(s)
	if result != expect {
		t.Errorf("OnlyLetters test failed! expected %s, received %s ", expect, result)
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

func TestIsValidMimeType2(t *testing.T) {
	if IsValidMimeType("") {
		t.Errorf("IsValidMimeType returns valid for an empty string")
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

func TestParseTags(t *testing.T) {
	s := "appLE, banana ,bird, ap2ple, cherry"
	expect := []string{"apple", "banana", "bird", "cherry"}

	result := ParseTags(s)

	for i, r := range result {
		if r != expect[i] {
			t.Errorf("ParseTags test failed expected: %s, received: %s", expect, result)
			break
		}
	}
}

func TestParseTags2(t *testing.T) {
	s := ""
	result := ParseTags(s)

	if len(result) > 0 {
		t.Errorf("ParseTages test failed, expected empty slice, received: %s of length: %v", result, len(result))
	}

}
