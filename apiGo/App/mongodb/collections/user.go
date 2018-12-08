package collections

type AccountType struct {
	Level int8   `json:"level,omitempty" bson:"level,omitempty"`
	Type  string `json:"type,omitempty" bson:"type,omitempty"`
}

type User struct {
	ID             string      `json:"_id,omitempty" bson:"_id,omitempty"`
	Username       string      `json:"username,omitempty" bson:"username,omitempty"`
	FirstName      string      `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName       string      `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Password       string      `json:"password,omitempty" bson:"password,omitempty"`
	ProfilePicture string      `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"`
	Email          string      `json:"email,omitempty" bson:"email,omitempty"`
	AccountType    AccountType `json:"account_type,omitempty" bson:"account_type,omitempty"`
	RandomToken    string      `json:"random_token,omitempty" bson:"random_token,omitempty"`
	Locale         string      `json:"locale,omitempty" bson:"locale,omitempty"`
	MoviesWatched  *[]int      `json:"movies_watched,omitempty" bson:"movies_watched,omitempty"`
}
