package collections

import "time"

// Comment corresponds to the structure of the Comment collection in mongodb
type Comment struct {
	ID        string    `json:"id,omitempty" bson:"_id"`
	IDimdb    string    `json:"idimdb,omitempty" bson:"idimdb"`
	UserID    string    `json:"userid,omitempty" bson:"userid,omitempty"`
	Content   string    `json:"content,omitempty" bson:"content,omitempty"`
	CreatedAt time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
}
