package movie

import (
	"encoding/xml"
	"net/http"
	"os"
	"sort"
	"strconv"

	"../../../lib"
	coltypes "../../../mongodb/collections"
	"../../../mongodb/query"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const jackettAPIKey = "xxxx"
const jackettDomain = "xxxx"

type rss struct {
	Channel []struct {
		Item []struct {
			Title    string `xml:"title"`
			Size     int    `xml:"size"`
			MetaData []struct {
				Name  string `xml:"name,attr"`
				Value string `xml:"value,attr"`
			} `xml:"attr"`
		} `xml:"item"`
	} `xml:"channel"`
}

type magnetList struct {
	Title string `json:"title,omitempty"`
	Size  int    `json:"size,omitempty"`
	Seed  int    `json:"seed,omitempty"`
	Peer  int    `json:"peer,omitempty"`
	JWT   string `json:"jwt,omitempty"`
}

type torrentSorter []magnetList

func (a torrentSorter) Len() int      { return len(a) }
func (a torrentSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a torrentSorter) Less(i, j int) bool {
	return a[i].Seed > a[j].Seed
}

func getMagnet(indexer string, magnets *[]magnetList, filmDB coltypes.Movie) (int, string, error) {
	query := "https://" + jackettDomain + "/api/v2.0/indexers/" + indexer + "/results/torznab/api?apikey=" + jackettAPIKey + "&t=movie&q=" + filmDB.IDimdb
	resp, err := http.Get(query)
	if err != nil {
		return 500, "Call to external api failed", err
	}
	var itemsList rss

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&itemsList)
	if err != nil {
		return 500, "Extract data from xml failed", err
	}
	var tmp int
	var magnet string
	if len(itemsList.Channel) > 0 {
		for _, item := range itemsList.Channel[0].Item {
			var temp magnetList
			temp.Title = item.Title
			temp.Size = item.Size
			magnet = ""
			for _, meta := range item.MetaData {
				if meta.Name == "seeders" {
					tmp, err = strconv.Atoi(meta.Value)
					if err != nil {
						tmp = 0
					}
					temp.Seed = tmp
				}
				if meta.Name == "peers" {
					tmp, err = strconv.Atoi(meta.Value)
					if err != nil {
						tmp = 0
					}
					temp.Peer = tmp
				}
				if meta.Name == "magneturl" {
					magnet = meta.Value
				}
			}
			// Generate a jwt token to protect magnets
			tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
				jwt.MapClaims{
					"magnet":   magnet,
					"id":       filmDB.IDimdb,
					"duration": filmDB.Duration,
					"poster":   filmDB.Poster,
					"title":    filmDB.OriginalTitle,
				}).
				SignedString([]byte(os.Getenv("jwtSecret")))
			if err != nil {
				return 500, "Encode jwt failed", err
			}
			temp.JWT = tokenString
			if temp.Seed > 0 {
				*magnets = append(*magnets, temp)
			}
		}
	}
	return 0, "", nil
}

var (
	providers = []string{
		"yts",
		"kickasstorrent-kathow",
		"btdb",
	}
)

// GetTorrentMagnets is the route '/api/v1/torrents/{filmID}' with the method GET.
func GetTorrentMagnets(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(lib.MongoDB).(*mgo.Database)
	if !ok {
		lib.RespondWithErrorHTTP(w, 500, "Problem with database connection")
		return
	}
	vars := mux.Vars(r)
	filmID := vars["filmId"]
	// Get database data -> Rating/Comments
	dbMovie, err := query.FindMovie(bson.M{"_id": filmID}, db)
	if err != nil {
		if err.Error() == "Not Found" {
			lib.RespondWithErrorHTTP(w, 406, "Movie does not exists in the database")
			return
		}
		lib.RespondWithErrorHTTP(w, 500, "Failed to collect movie data in the database")
		return
	}
	var magnets []magnetList
	for _, provider := range providers {
		errCode, errStatus, err := getMagnet(provider, &magnets, dbMovie)
		if err != nil {
			lib.RespondWithErrorHTTP(w, errCode, errStatus)
			return
		}
	}
	sort.Sort(torrentSorter(magnets))
	if len(magnets) > 0 {
		lib.RespondWithJSON(w, http.StatusOK, magnets)
		return
	}
	lib.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "No torrents",
	})
}
