package lib

import (
	"errors"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Oauth struct {
	Token    string
	Provider string
}

type DataJWT struct {
	Duration       time.Duration
	ISS            string
	Sub            string
	UserID         string
	Username       string
	FirstName      string
	LastName       string
	ProfilePicture string
	Oauth          Oauth
}

func GenerateJWT(settings DataJWT) (string, error) {
	now := time.Now().Local()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":             now.Unix(),
		"exp":             now.Add(settings.Duration).Unix(),
		"iss":             settings.ISS,
		"sub":             settings.Sub,
		"userId":          settings.UserID,
		"username":        settings.Username,
		"firstname":       settings.FirstName,
		"lastname":        settings.LastName,
		"profile_picture": settings.ProfilePicture,
		"oauth": map[string]string{
			"token":    settings.Oauth.Token,
			"provider": settings.Oauth.Provider,
		},
	})
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseJWT check the JWT validity and get the data from the JWT
func AnalyseJWT(strToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return JWTSecret, nil
	})
	if err != nil {
		if ve, yes := err.(*jwt.ValidationError); yes {
			if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
				return nil, errors.New("Token expired")
			}
		}
		log.Println(PrettyError("[JWT] Not a valid JSON Web Token - " + err.Error()))
		return nil, errors.New("Not a valid JSON Web Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Not a valid token")
	}
	return claims, nil
}
