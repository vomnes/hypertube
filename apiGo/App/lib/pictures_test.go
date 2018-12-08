package lib

import (
	"os"
	"testing"
)

var pictureTests = []struct {
	typeImage     string
	imageBase64   string
	errorExpected bool
	errorContent  string
	testContent   string
}{
	{"image/jpeg", "/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAALCAAQABABAREA/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMECP/EAB4QAAEFAQEBAQEAAAAAAAAAAAQDBQYHCAIJARQS/9oACAEBAAA/AN17njdi6FtO8qT3puJvJpfL26WBopqu+cL1XZkPn5r55wzbXDKRaMeOloAcgb4lX7jeMdZhHxBxAPmtfQ6Qp/mfyRfgIwnHbEz/AGtR9J4L3E3i0np7ccnbLorzrC9V1nEIAfHfOeCaweFqtjoMuPDYG2YQQSjY69isaLaA3zKwJlIlPhMgFJ4Pv9PeacC3lqas7e1HnPKhku5xtqCCO2lJ8vXjFYQRuR/Q7Fs9HhhiLO8qvLxGjZDXar6lwMkgEzFfelSv19AAmt8ukqbL3fmOtKh1HnLVZkRG2Xp+du2a58vYbFXwY+VPOfGMDFmZizOzKszzJTI3YqzEj2Mqgayif2iV+vhwBC//2Q==", false, "", "Generate jpeg file succeed"},
	{"image/png", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk6Pn/HwAEqwKMiLo2HgAAAABJRU5ErkJggg==", false, "", "Generate png file succeed"},
	{"image/png", "/9j/4AAQSkZJRgABAQIAIwAjAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAALCAAQABABAREA/8QAFwAAAwEAAAAAAAAAAAAAAAAAAQMECP/EAB4QAAEFAQEBAQEAAAAAAAAAAAQDBQYHCAIJARQS/9oACAEBAAA/AN17njdi6FtO8qT3puJvJpfL26WBopqu+cL1XZkPn5r55wzbXDKRaMeOloAcgb4lX7jeMdZhHxBxAPmtfQ6Qp/mfyRfgIwnHbEz/AGtR9J4L3E3i0np7ccnbLorzrC9V1nEIAfHfOeCaweFqtjoMuPDYG2YQQSjY69isaLaA3zKwJlIlPhMgFJ4Pv9PeacC3lqas7e1HnPKhku5xtqCCO2lJ8vXjFYQRuR/Q7Fs9HhhiLO8qvLxGjZDXar6lwMkgEzFfelSv19AAmt8ukqbL3fmOtKh1HnLVZkRG2Xp+du2a58vYbFXwY+VPOfGMDFmZizOzKszzJTI3YqzEj2Mqgayif2iV+vhwBC//2Q==", true, "Corrupted file [image/png] | Error: png: invalid format: not a PNG file", "Generate png file failed - Corrupted file"},
	{"image/jpg", "/9j", true, "Corrupted file [image/jpg] | Error: unexpected EOF", "Generate png file failed - Corrupted file"},
	{"image/jpeg", "==", true, "Corrupted file [image/jpeg] | Error: unexpected EOF", "Generate jpeg file failed - Corrupted file"},
	{"image/random", "==", true, "Image type [image/random] not accepted, support only png, jpg and jpeg images", "File type doesn't exists"},
}

func TestGeneratePicture(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("Failed to get the root path name - EncodeBase64")
	}
	for _, tt := range pictureTests {
		pathString, err := Base64ToImageFile(tt.imageBase64, tt.typeImage, path, "/storage/tests/", GetRandomString(42))
		if err != nil {
			if tt.errorExpected {
				if err.Error() != tt.errorContent {
					t.Errorf("GeneratePicture - Expected the error '%s' has '%s' - Test type: \033[31m%s\033[0m", tt.errorContent, err.Error(), tt.testContent)
				}
			} else {
				t.Error(err)
			}
		} else {
			if tt.errorExpected {
				t.Errorf("GeneratePicture - Expected the error '%s' has 'nil' - Test type: \033[31m%s\033[0m", tt.errorContent, tt.testContent)
			}
		}
		// Check : File created
		empty, err := FileExists(path + pathString)
		if err != nil {
			t.Error(err)
		}
		if !empty {
			t.Errorf("GeneratePicture - The file doesn't exists, it hasn't been created - Test type: \033[31m%s\033[0m", tt.testContent)
		}
	}
	err = os.RemoveAll(path + "/storage/tests/")
	if err != nil {
		t.Error(err)
	}
}

var base64Tests = []struct {
	base64        string
	errorContent  string
	typeImage     string
	imageBase64   string
	testContent   string
	errorExpected bool
}{
	{"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk6Pn/HwAEqwKMiLo2HgAAAABJRU5ErkJggg==", "", "image/png", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk6Pn/HwAEqwKMiLo2HgAAAABJRU5ErkJggg==", "Extract image type and base64 succeed", false},
	{"data:image/random;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk6Pn/HwAEqwKMiLo2HgAAAABJRU5ErkJggg==", "", "image/random", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk6Pn/HwAEqwKMiLo2HgAAAABJRU5ErkJggg==", "Extract image/random and base64", false},
	{"data:image/failed;iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk6Pn/HwAEqwKMiLo2HgAAAABJRU5ErkJggg==", "Base64 doesn't match with the pattern 'data:image/[...];base64,[...]'", "", "", "Extract invalid pattern", true},
	{"", "Base64 can't be empty", "", "", "Empty base64 field", true},
}

func TestBase64ToImageFile(t *testing.T) {
	for _, tt := range base64Tests {
		typeImage, imageBase64, err := ExtractDataPictureBase64(tt.base64)
		if typeImage != tt.typeImage {
			t.Errorf("Base64ToImageFile - Test type: \033[31m%s\033[0m - File type extracted doesn't match, expect %s has %s", tt.testContent, tt.typeImage, typeImage)
		}
		if imageBase64 != tt.imageBase64 {
			t.Errorf("Base64ToImageFile - Test type: \033[31m%s\033[0m - Image base64 extracted doesn't match, expect %s has %s", tt.testContent, tt.imageBase64, imageBase64)
		}
		if tt.errorExpected && err != nil && err.Error() != tt.errorContent {
			t.Errorf("Base64ToImageFile - Test type: \033[31m%s\033[0m - Expected the error '%s' has '%s'", tt.testContent, tt.errorContent, err.Error())
		}
		if tt.errorExpected && err == nil {
			t.Errorf("Base64ToImageFile - Test type: \033[31m%s\033[0m - Expected the error '%s' has 'nil'", tt.testContent, tt.errorContent)
		}
	}
}
