package lib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)

type key int

const (
	// MongoDB key is used as value in order to store database in the context
	MongoDB key = iota
	// MongoDB key is used as value in order to store database session in the context
	MongoDBSession
	// MailJet key is used as value in order to store mailjet connection
	MailJet
	// UserID key is used as value in order to store userId from JSON Web Token in the context
	UserID
	// Username key is used as value in order to store username from JSON Web Token in the context
	Username
	// FirstName key is used as value in order to store firstname from JSON Web Token in the context
	FirstName
	// LastName key is used as value in order to store lastname from JSON Web Token in the context
	LastName
	// ProfilePicture key is used as value in order to store profile picture from JSON Web Token in the context
	ProfilePicture
)

var (
	LocaleAvailable = []string{
		"en", // English - Default
		"fr", // French
		"it", // Italian
	}
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_-"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GetRandomString create a random string with a length of n characters
// with the characters include in letterBytes
func GetRandomString(n int) string {
	b := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// StringInArray take a string or an array of strings and a array of string as parameter
// If the first argument is an array of string, the function return true if
// at list one of the elements array in the array of the second argument
// Return true if the string is in the array of string else false
func StringInArray(a interface{}, list []string) bool {
	var elements []string
	typeA := reflect.TypeOf(a)
	if typeA.String() == "string" {
		elements = []string{a.(string)}
	}
	if typeA.String() == "[]string" {
		elements = a.([]string)
	}
	if len(elements) == 0 {
		return false
	}
	for _, b := range list {
		for _, elem := range elements {
			if b == elem {
				return true
			}
		}
	}
	return false
}

func Strsub(input string, start int, end int) string {
	var output string
	if start < 0 || end < 0 {
		return ""
	}
	for i := start; i < start+end; i++ {
		output += string(input[i])
	}
	return output
}

func TrimStringFromString(s, sub string) string {
	if idx := strings.Index(s, sub); idx != -1 {
		return s[:idx]
	}
	return s
}

func TrimStringFromLastString(s, sub string) string {
	if idx := strings.LastIndexAny(s, sub); idx != -1 {
		return s[:idx]
	}
	return s
}

func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// UniqueTimeToken generates a unique base64 token using a key and time.Now() as string
func UniqueTimeToken(key string) string {
	now := time.Now()
	data := []byte(key + "&" + now.String())
	return base64.StdEncoding.EncodeToString(data)
}

func InterfaceToByte(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func SWAPStrings(str1, str2 *string) {
	*str1, *str2 = *str2, *str1
}
