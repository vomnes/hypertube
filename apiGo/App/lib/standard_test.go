package lib

import "testing"

var stringInArrayTests = []struct {
	item_s      interface{} // input
	stringArray []string
	expected    bool   // expected result
	testContent string // test details
}{
	{"a", []string{"a", "b", "c"}, true, "Basic - string"},
	{"z", []string{"a", "b", "c"}, false, "Not in array"},
	{42, []string{"a", "b", "c"}, false, "Not a string or an array"},
	{[]string{"a"}, []string{"a", "b", "c"}, true, "Basic - []string"},
	{[]string{"z", "a"}, []string{"a", "b", "c"}, true, "Basic - []string"},
	{[]string{"a", "z"}, []string{"a", "b", "c"}, true, "Basic - []string"},
	{[]string{"y", "z"}, []string{"a", "b", "c"}, false, "Not in array"},
	{[]string{"y"}, []string{"a", "b", "c"}, false, "Not in array"},
}

func TestStringInArray(t *testing.T) {
	for _, tt := range stringInArrayTests {
		actual := StringInArray(tt.item_s, tt.stringArray)
		if actual != tt.expected {
			t.Errorf("StringInArray(%s, %s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.item_s, tt.stringArray, tt.expected, actual, tt.testContent)
		}
	}
}
