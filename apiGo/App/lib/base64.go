package lib

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

func Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func Base64Encode(d []byte) string {
	return base64.StdEncoding.EncodeToString(d)
}

func ExtractBase64Struct(base64 string, data interface{}) error {
	byteData, err := Base64Decode(base64)
	if err != nil {
		return errors.New("[Base64] Failed to decode search parameters in header " + err.Error())
	}
	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return errors.New("[Unmarshal] Failed to unmarshal search parameters in header " + err.Error())
	}
	return nil
}
