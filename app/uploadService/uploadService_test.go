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
	goodName := "GoodName"
	result := CleanFilename(goodName)
	if result != goodName {
		t.Errorf("CleanFilename test failed exptectd: %s recieved: %s", goodName, result)
	}
}
