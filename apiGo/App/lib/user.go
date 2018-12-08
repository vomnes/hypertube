package lib

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/rangetable"
)

const (
	// UsernameMinLength corresponds to the minimum character length of the username
	UsernameMinLength = 1
	// UsernameMaxLength corresponds to the maximum character length of the username
	UsernameMaxLength = 254
	// EmailAddressMaxLength corresponds to the maximum character length of the email address
	EmailAddressMaxLength = 254
	// PasswordMinLength corresponds to the minimum character length of a password
	PasswordMinLength = 4
)

// IsValidUsername check if the string parameter is a valid username
// Check length maximum and minimum and authorized characters (a-zA-Z0-9.-_)
// Return a boolean
func IsValidUsername(s string) bool {
	if len(s) < UsernameMinLength || len(s) > UsernameMaxLength {
		return false
	}
	if !regexp.MustCompile(`^[\p{L}0-9\.\-_]+$`).MatchString(s) {
		return false
	}
	return true
}

// IsValidFirstLastName check if the string parameter is a valid lastname or firstname
// Check length maximum and minimum and authorized characters (a-zA-Z -)
// Return a boolean
func IsValidFirstLastName(s string) bool {
	if len(s) < 1 || len(s) > UsernameMaxLength {
		return false
	}
	if !regexp.MustCompile(`^[\p{L}\-]+$`).MatchString(s) {
		return false
	}
	return true
}

// IsValidEmailAddress check if the string parameter is a valid email address
// Return a boolean
func IsValidEmailAddress(s string) bool {
	reEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !reEmail.MatchString(s) {
		return false
	}
	return true
}

// IsValidPassword check if the string parameter is a valid password
// A valid password must have a minimum length of PasswordMinLength,
// no white space at all, at least 1 digit OR 1 special char (@#$%^&+=)
// Return a boolean
func IsValidPassword(s string) bool {
	var hasDigit, hasSpecial bool
	var specialChars = []*unicode.RangeTable{rangetable.New('@', '#', '$', '%', '^', '&', '+', '=')}
	for _, s := range s {
		if unicode.IsSpace(s) {
			return false
		}
		if unicode.IsDigit(s) {
			hasDigit = true
		}
		if unicode.IsOneOf(specialChars, s) {
			hasSpecial = true
		}
	}
	if len(s) < PasswordMinLength || (!hasDigit && !hasSpecial) {
		return false
	}
	return true
}

// IsValidText check if the string parameter is a text
// Check length maximum and authorized characters (a-zA-Z0-9 .,?!&-_*-+@#$%;)
// Return a boolean
func IsValidText(s string, lengthMax int) bool {
	if len(s) > lengthMax {
		return false
	}
	if !regexp.MustCompile(`^[\p{L}0-9\ \.\,\?\!\&\-\_\*\-\+\@\#\$\%\;]+$`).MatchString(s) {
		return false
	}
	return true
}

// IsOnlyLowercaseLetters check if the string parameter has only lowercase letters
// Check length maximum and authorized characters (a-z)
// Return a boolean
func IsOnlyLowercaseLetters(s string, lengthMax int) bool {
	if len(s) > lengthMax {
		return false
	}
	if !regexp.MustCompile(`^[a-z]+$`).MatchString(s) {
		return false
	}
	return true
}

// IsValidCommonName check if the string parameter is a common name
// Check length maximum and authorized characters (a-zA-Z0-9.-_ )
// Return a boolean
func IsValidCommonName(s string) bool {
	if len(s) < 1 || len(s) > UsernameMaxLength {
		return false
	}
	if !regexp.MustCompile(`^[\p{L}0-9\.\-_\ ]+$`).MatchString(s) {
		return false
	}
	return true
}

// IsValidDate check if the string parameter is a valid date dd/mm/yyyy
// Return a boolean and status error
func IsValidDate(s string) (bool, error) {
	if len(s) > len("mm/dd/yyyy") {
		return false, nil
	}
	if !regexp.MustCompile(`^[0-9\/]+$`).MatchString(s) {
		return false, nil
	}
	if strings.Count(s, "/") != 2 {
		return false, nil
	}
	var day, month, year string
	idx := strings.Index(s, "/")
	if idx != -1 {
		day = s[:idx]
	}
	lastIdx := strings.LastIndex(s, "/")
	if lastIdx != -1 {
		year = s[lastIdx+1:]
	}
	month = Strsub(s, idx+1, 2)
	var monthLimit = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return false, err
	}
	if monthInt < 1 || monthInt > 12 {
		return false, nil
	}
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return false, err
	}
	if dayInt < 1 || dayInt > monthLimit[monthInt-1] {
		return false, nil
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return false, err
	}
	if yearInt <= 999 || yearInt > 9999 {
		return false, nil
	}
	return true, nil
}

// IsValidTag check if the string parameter is a valid tag
// Check length maximum and minimum and authorized characters (a-z0-9-_)
// Return a boolean
func IsValidTag(s string) bool {
	if len(s) < 1 || len(s) > UsernameMaxLength {
		return false
	}
	if !regexp.MustCompile(`^[a-z0-9\-_]+$`).MatchString(s) {
		return false
	}
	return true
}

// IsValidTag check if the string parameter is a valid IP Addres
// Return a boolean
func IsValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	if !regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`).
		MatchString(ipAddress) {
		return false
	}
	return true
}

// GetAge take a *time.Time as parameter and return the date age in a int
func GetAge(date *time.Time) int {
	if date == nil {
		return 0
	}
	return int(time.Since(*date).Hours() / 8760)
}

func IsValidBase64Picture(base64 string) bool {
	if regexp.MustCompile(`data:image\/([a-zA-Z]*);base64,([^\"]*)`).MatchString(base64) {
		return true
	}
	return false
}
