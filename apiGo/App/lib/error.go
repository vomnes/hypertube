package lib

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

// PrettyError create take a string as paramter to return a formated error
// that contains the error file and line
func PrettyError(err string) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.New("Unidentified: " + err)
	}
	lastIndex := strings.LastIndex(file, "server")
	if lastIndex == -1 {
		file = ""
	} else {
		fileBytes := []byte(file)
		fileBytes = fileBytes[lastIndex:]
		file = string(fileBytes)
	}
	return errors.New(file + ": l." + strconv.Itoa(line) + ": " + err)
}
