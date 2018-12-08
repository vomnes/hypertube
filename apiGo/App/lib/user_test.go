package lib

import "testing"

var usernameTests = []struct {
	username    string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"abcdef", true, "Only lowercase letters"},
	{"123456789", true, "Only digit character"},
	{"ABCDEF", true, "Only uppercase letters"},
	{"abcdef-123456789_ABC.DEF", true, "Lowercase, uppercase, digit, -, _, ."},
	{"", false, "Too short"},
	{"abcdef-123456789_ABCDEF.vomnes.vomnes.vomnes.vomnes.vomnes.vomnesabcdef-123456789_ABCDEF.vomnes.vomnes.vomnes.vomnes.vomnes.vomnesabcdef-123456789_ABCDEF.vomnes.vomnes.vomnes.vomnes.vomnes.vomnesabcdef-123456789_ABCDEF.vomnes.vomnes.vomnes.vomnes.vomnes.vomnesabcdef-123456789_ABCDEF.vomnes.vomnes.vomnes.vomnes.vomnes.vomnes", false, "Too long"},
	{"vomnes->#$%^&*()_)(*&^%)", false, "Forbidden characters"},
	{" abcABC123 ", false, "Space before and after"},
	{"abcABéùàC123", true, "Valid with diacritic"},
}

func TestIsValidUsername(t *testing.T) {
	for _, tt := range usernameTests {
		actual := IsValidUsername(tt.username)
		if actual != tt.expected {
			t.Errorf("IsValidUsername(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.username, tt.expected, actual, tt.testContent)
		}
	}
}

var nameTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"abcdef", true, "Only lowercase letters"},
	{"ABCDEF", true, "Only uppercase letters"},
	{"abcdefABCDEF", true, "Lowercase, uppercase characters"},
	{"abcdef-ABCDEF", true, "Lowercase, uppercase and separator (HyphenMinus) characters"},
	{"abcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcd", true, "Limit max length"},
	{"", false, "Too short"},
	{"123456789", false, "Only digit character"},
	{"abcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdeabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdeabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdeabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcde", false, "Too long - Over max length"},
	{"vomnes.", false, "Forbidden characters '.'"},
	{"vomnes_", false, "Forbidden characters '_'"},
	{"vomnes%", false, "Forbidden characters '%'"},
	{"vomnes<", false, "Forbidden characters '<'"},
	{"vomnes>", false, "Forbidden characters '>'"},
	{"vomnes=/*-+/=%^&*()", false, "Forbidden characters"},
	{" abcABC123 ", false, "Space before and after"},
	{"abcABéùà", true, "Valid with diacritic"},
}

func TestIsValidFirstLastName(t *testing.T) {
	for _, tt := range nameTests {
		actual := IsValidFirstLastName(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidFirstLastName(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var emailAddressTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"valentin.omnes@gmail.com", true, "Valid email address"},
	{"vomnes@student.42.fr", true, "Valid email address"},
	{"true@student-42-fr", true, "Valid email address"},
	{"true_true@student.fr", true, "'_' before '@'"},
	{"f@s.fr", true, "Short email address"},
	{"ç$€§/az@gmail.com", false, "Illegal characters before '@'"},
	{"false@student_42_fr", false, "Illegal characters '_'"},
	{"false@student<42.fr", false, "Illegal characters '<'"},
	{"false@student>42.fr", false, "Illegal characters '>'"},
	{"false@student@42.fr", false, "Illegal characters '@'"},
	{"false@student*42.fr", false, "Illegal characters '*'"},
	{"studentstudentstudentstudentstudentstudentstudentstudentstudentstudent" +
		"studentstudent@studentstudentstudentstudentstudentstudentstudentstudent" +
		"studentstudentstudentstudent.fr", false, "Too long email address"},
}

func TestIsValidEmailAddress(t *testing.T) {
	for _, tt := range emailAddressTests {
		actual := IsValidEmailAddress(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidEmailAddress(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var passwordTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"abcABC123", true, "Valid password"},
	{"abcABC@", true, "Valid password - @"},
	{"abcABC#", true, "Valid password - #"},
	{"abcABC$", true, "Valid password - $"},
	{"abcABC%", true, "Valid password - %"},
	{"abcABC^", true, "Valid password - ^"},
	{"abcABC&", true, "Valid password - &"},
	{"abcABC+", true, "Valid password - +"},
	{"abcABC=", true, "Valid password - ="},
	{"abcABC@42", true, "Valid password - @ + Digit"},
	{"abcABC§", false, "Invalid password - §"},
	{"abcABCDEF", false, "Only letters"},
	{"abc", false, "Too short"},
	{"     \t      ", false, "Only space"},
	{"", false, "Empty"},
}

func TestIsValidPassword(t *testing.T) {
	for _, tt := range passwordTests {
		actual := IsValidPassword(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidPassword(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var textTests = []struct {
	str         string // input
	length      int    // length max
	expected    bool   // expected result
	testContent string // test details
}{
	{"Hello world ? 42 .,?!&-_*-+@#$%", 255, true, "Valid"},
	{"Hello world ? 42", 255, true, "Valid"},
	{"aaaaa bbbbb ccccc", 255, true, "Only lowercases"},
	{"AAAAA BBBBB CCCCC", 255, true, "Only uppercases"},
	{"123456789", 255, true, "Only digits"},
	{"Hello world ? 42 .,?!&-_*-+@#$%", 5, false, "Too long"},
	{"Hello world ? ", 255, true, "Valid"},
	{"Hello world ? ≤", 255, false, "Invalid char ≤"},
	{"Hello world ? <", 255, false, "Invalid char <"},
	{"Hello world ? >", 255, false, "Invalid char >"},
	{"Hello world ? &", 255, true, "Valid char &"},
	{"I&#39;m Valentin Omnes", 255, true, "Valid with escaped char"},
	{"abcABéùàC123", 255, true, "Valid with diacritic"},
}

func TestIsValidText(t *testing.T) {
	for _, tt := range textTests {
		actual := IsValidText(tt.str, tt.length)
		if actual != tt.expected {
			t.Errorf("IsValidText(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var letterTests = []struct {
	str         string // input
	length      int    // length max
	expected    bool   // expected result
	testContent string // test details
}{
	{"hello", 32, true, "Valid"},
	{"hello", 2, false, "Too long"},
	{"hello123456789", 32, false, "Contains digit"},
	{"Hello", 32, false, "Contains uppercase"},
	{"Hello?@", 32, false, "Contains uppercase"},
}

func TestIsOnlyLowercaseLetters(t *testing.T) {
	for _, tt := range letterTests {
		actual := IsOnlyLowercaseLetters(tt.str, tt.length)
		if actual != tt.expected {
			t.Errorf("IsOnlyLowercaseLetters(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var dateTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"06/03/1995", true, "Valid"},
	{"6/03/1995", true, "Valid - Day"},
	{"12/12/1995", true, "Invalid day"},
	{"06/13/1995333", false, "Invalid length"},
	{"06.12.1995", false, "Invalid characters '.'"},
	{"06-12-1995", false, "Invalid characters '-'"},
	{"06/03/199a", false, "Invalid characters 'a'"},
	{"06/1995", false, "Invalid number of '/'"},
	{"06/13/1995", false, "Invalid month"},
	{"29/02/2017", false, "Limited to 28 days in February"},
	{"31/04/2017", false, "Limited to 30 days in April"},
	{"30/04/2017", true, "Limited to 30 days in April"},
}

func TestIsValidDate(t *testing.T) {
	for _, tt := range dateTests {
		actual, _ := IsValidDate(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidDate(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var tagTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"test", true, "Valid"},
	{"test09", true, "Valid with digit"},
	{"test-a", true, "Valid with '-'"},
	{"test_b", true, "Valid with '_'"},
	{"", false, "Empty"},
	{"testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest", false, "Too long"},
	{"This", false, "Invalid character uppercase"},
	{"this is false", false, "Invalid character space"},
}

func TestIsValidTag(t *testing.T) {
	for _, tt := range tagTests {
		actual := IsValidTag(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidTag(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var CommonNameTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"aB12-_. ", true, "Valid"},
	{"37.169.43. 146", true, "Valid"},
	{"validContent", true, "Valid"},
	{"<h1>Title</h1>", false, "HTML"},
	{"Title§è!çà)", false, "Invalid characters"},
	{"abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabc", false, "Too long"},
	{"abcABéùàC123", true, "Valid with diacritic"},
}

func TestIsValidCommonName(t *testing.T) {
	for _, tt := range CommonNameTests {
		actual := IsValidCommonName(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidCommonName(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var IP4Tests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"37.169.43.146", true, "Valid"},
	{"192.168.a.0", false, "Invalid type"},
	{"FE80:0000:0000:0000:0202:B3FF:FE1E:8329", false, "Invalid type"},
	{"FE80::0202:B3FF:FE1E:8329", false, "Invalid type"},
	{"http://[2001:db8:0:1]:80", false, "Invalid type"},
}

func TestIsValidIP4(t *testing.T) {
	for _, tt := range IP4Tests {
		actual := IsValidIP4(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidIP4(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}
