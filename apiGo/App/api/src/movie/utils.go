package movie

import (
	coltypes "../../../mongodb/collections"
)

func hasBeenWatched(userID string, watchedBy []coltypes.Watched) bool {
	for _, watcher := range watchedBy {
		if watcher.UserID == userID {
			return true
		}
	}
	return false
}

func decodeVideoStatus(status int8) string {
	switch status {
	case -1:
		return "not downloaded"
	case 0:
		return "waiting"
	case 1:
		return "downloading"
	case 2:
		return "ready"
	}
	return "not downloaded"
}

func encodeVideoStatus(status string) []int8 {
	switch status {
	case "not downloaded":
		return []int8{-1}
	case "waiting":
		return []int8{0, 1}
	case "downloading":
		return []int8{0, 1}
	case "ready":
		return []int8{2}
	}
	return []int8{2}
}

type categoryMovie struct {
	ID        string   `json:"id"`
	Poster    string   `json:"poster"`
	Title     string   `json:"title"`
	Rating    float64  `json:"rating"`
	IsWatched bool     `json:"is_watched"`
	Status    string   `json:"status"`
	Genres    []string `json:"genres,omitempty"`
	Year      int      `json:"year,omitempty"`
	Duration  int      `json:"duration,omitempty"`
}

type show struct {
	genres, year, duration bool
}

type formatMovie struct {
	movies           []coltypes.Movie
	userID, language string
	show             show
}

func (f *formatMovie) formatListMovies() []categoryMovie {
	var response []categoryMovie
	var thisTitle string
	for _, movie := range f.movies {
		// Handle differents language for title
		thisTitle = movie.Titles[f.language]
		if thisTitle == "" {
			thisTitle = movie.OriginalTitle
		}
		item := categoryMovie{
			ID:        movie.IDimdb,
			Poster:    movie.Poster,
			Title:     thisTitle,
			Rating:    movie.Rating.Average,
			IsWatched: hasBeenWatched(f.userID, movie.WatchedBy),
			Status:    decodeVideoStatus(movie.Video.Status),
		}
		if f.show.genres {
			item.Genres = movie.Genres
		}
		if f.show.year {
			item.Year = movie.Year
		}
		if f.show.duration {
			item.Duration = movie.Duration
		}
		response = append(response, item)
	}
	return response
}
