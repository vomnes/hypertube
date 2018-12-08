package collections

import "time"

type Rating struct {
	Average float64 `json:"average,omitempty" bson:"average,omitempty`
	Number  int     `json:"number,omitempty" bson:"number,omitempty`
}

type Subtitle struct {
	Language string `json:"lang,omitempty" bson:"lang,omitempty`
	Path     string `json:"path,omitempty" bson:"path,omitempty`
}

type Video struct {
	Status    int8       `json:"status" bson:"status"` // Not omitempty otherwise 0 value isn't inserted
	Path      string     `json:"path,omitempty" bson:"path,omitempty"`
	Stream    bool       `json:"stream,omitempty" bson:"stream,omitempty"`
	Subtitles []Subtitle `json:"subtitles,omitempty" bson:"subs,omitempty"`
}

type Watched struct {
	UserID    string    `json:"userid,omitempty" bson:"userid,omitempty"`
	WatchedAt time.Time `json:"watchedat,omitempty" bson:"watchedat,omitempty"`
}

type Movie struct {
	IDimdb        string            `json:"idimdb,omitempty" bson:"_id"`
	OriginalTitle string            `json:"title,omitempty" bson:"title,omitempty"`
	Titles        map[string]string `json:"titles,omitempty" bson:"language_title,omitempty"`
	Year          int               `json:"year,omitempty" bson:"year,omitempty"`
	Duration      int               `json:"duration,omitempty" bson:"duration,omitempty"`
	Genres        []string          `json:"genres,omitempty" bson:"genres,omitempty"`
	Rating        Rating            `json:"rating,omitempty" bson:"rating,omitempty"`
	Poster        string            `json:"poster,omitempty" bson:"poster,omitempty"`
	Video         Video             `json:"video" bson:"video"` // Not omitempty otherwise 0 status value isn't inserted
	WatchedBy     []Watched         `json:"watchedby,omitempty" bson:"watchedby,omitempty"`
}
