package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"./lib"
	"./mongodb"
	"github.com/gocolly/colly"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Video struct {
	Status int8 `json:"status" bson:"status"` // Not omitempty otherwise 0 value isn't inserted
}

type Rating struct {
	Average float64 `json:"average,omitempty" bson:"average,omitempty"`
	Number  int     `json:"number,omitempty" bson:"number,omitempty"`
}

type Movie struct {
	IDimdb        string            `json:"idimdb" bson:"_id"`                      // title.basics.tsv.gz - 0
	OriginalTitle string            `json:"title,omitempty" bson:"title,omitempty"` // title.basics.tsv.gz - 3
	Titles        map[string]string `json:"titles,omitempty" bson:"language_title,omitempty"`
	Year          int               `json:"year,omitempty" bson:"year,omitempty"`         // title.basics.tsv.gz - 5
	Duration      int               `json:"duration,omitempty" bson:"duration,omitempty"` // title.basics.tsv.gz - 7
	Genres        []string          `json:"genres,omitempty" bson:"genres,omitempty"`     // title.basics.tsv.gz - 8
	Rating        Rating            `json:"rating,omitempty" bson:"rating,omitempty"`
	Poster        string            `json:"poster,omitempty" bson:"poster,omitempty"`
	Video         Video             `json:"video" bson:"video"` // Not omitempty otherwise 0 status value isn't inserted
}

const (
	YearMin     = 1950
	DurationMin = 60
	VotesMin    = 10000
)

func createTSVReader(fileName string) (*bufio.Reader, io.ReadCloser, *gzip.Reader, error) {
	fmt.Print("Get file with HTTP call from 'https://datasets.imdbws.com/'")
	// Call the file location
	resp, err := http.Get("https://datasets.imdbws.com/" + fileName)
	if err != nil {
		return nil, nil, nil, err
	}
	if resp.StatusCode == 404 {
		return nil, nil, nil, err
	}
	// Make the gzip data readable
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	// Use bufio to be able to read line by line
	reader := bufio.NewReader(gzipReader)
	fmt.Println(" - File [OK]")
	return reader, resp.Body, gzipReader, nil
}

func readLine(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return []string{}, err
	}
	line = line[:len(line)-1]
	return strings.Split(line, "\t"), nil
}

func basicFile(reader *bufio.Reader, movies *[]Movie) error {
	var err error
	var items, categories []string
	var year, duration int64
	for {
		items, err = readLine(reader)
		if err != nil {
			return err
		}
		if items[0] != "tconst" &&
			items[1] == "movie" &&
			len(items) == 9 {
			// Clean data fields with \N
			for i, item := range items {
				if strings.Contains(item, "\\N") {
					items[i] = ""
				}
				if i == 5 || i == 7 {
					if item == "" {
						items[i] = "0"
					}
				}
			}
			// Handle year field
			year, err = strconv.ParseInt(items[5], 10, 32)
			if err != nil {
				year = 0
			}
			// Don't store movie with no year
			if year < YearMin {
				continue
			}
			// Handle duration field
			duration, err = strconv.ParseInt(items[7], 10, 32)
			if err != nil {
				duration = 0
			}
			if duration < DurationMin {
				continue
			}
			err = nil
			categories = strings.Split(items[8], ",")
			if len(categories) == 0 || (len(categories) == 1 && categories[0] == "") {
				continue
			}
			*movies = append(*movies, Movie{
				IDimdb:        items[0],
				OriginalTitle: items[3],
				Year:          int(year),
				Duration:      int(duration),
				Genres:        categories,
				Video: Video{
					Status: -1,
				},
			})
		}
	}
}

func ratingFile(reader *bufio.Reader, movies *[]Movie) error {
	var err error
	var items []string
	data := make(map[string]Rating)
	var average float64
	var number int64
	for {
		items, err = readLine(reader)
		if err != nil {
			break
		}
		if len(items) != 3 {
			continue
		}
		average, err = strconv.ParseFloat(items[1], 32)
		if err != nil {
			average = 0
		}
		number, err = strconv.ParseInt(items[2], 10, 32)
		if err != nil {
			number = 0
		}
		data[items[0]] = Rating{
			Average: math.Floor(average*10) / 10,
			Number:  int(number),
		}
	}
	// Match ratings with movies
	for i, movie := range *movies {
		(*movies)[i].Rating = data[movie.IDimdb]
	}
	return err
}

func titlesFile(reader *bufio.Reader, movies *[]Movie) error {
	var err error
	var items []string
	var language string
	data := make(map[string]map[string]string)
	for {
		items, err = readLine(reader)
		if err != nil {
			break
		}
		language = strings.ToLower(items[3])
		if len(items) != 8 ||
			!lib.StringInArray(language, []string{"gb", "us", "fr", "it", "es"}) {
			continue
		}
		if data[items[0]] == nil {
			data[items[0]] = make(map[string]string)
		}
		data[items[0]][language] = items[2]
	}
	// Match titles with movies
	for i, movie := range *movies {
		(*movies)[i].Titles = data[movie.IDimdb]
	}
	return err
}

type fileToStruct func(reader *bufio.Reader, movies *[]Movie) error

func extractData(fileName string, movies *[]Movie, extractor fileToStruct) error {
	fmt.Print("Start to extract data from " + fileName + "\n")
	reader, file, gzipReader, err := createTSVReader(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	defer gzipReader.Close()
	err = extractor(reader, movies)
	if err != io.EOF {
		return err
	}
	fmt.Print("Finish to extract data from " + fileName + "\n")
	return nil
}

func filterData(movies []Movie, db *mgo.Database) ([]Movie, error) {
	// Get movies already inserted in the database
	pipe := db.C("movies").Pipe([]bson.M{{"$group": bson.M{"_id": nil, "idsimdb": bson.M{"$push": "$_id"}}}})
	var currentMovies struct {
		IDsIMDB []string `bson:"idsimdb"`
	}
	err := pipe.One(&currentMovies)
	var currentMoviesIDs []string
	if err != nil {
		if err != mgo.ErrNotFound {
			return []Movie{}, err
		}
	} else {
		currentMoviesIDs = currentMovies.IDsIMDB
	}
	var data []Movie
	// Get ids []string of the current movies in the database
	// in order to filter them
	for _, movie := range movies {
		if !lib.StringInArray(movie.IDimdb, currentMoviesIDs) &&
			len(movie.Titles) > 0 &&
			movie.Rating != (Rating{}) &&
			movie.Rating.Number > VotesMin {
			data = append(data, movie)
		}
	}
	return data, nil
}

type SafePosters struct {
	URL map[string]string
	mux sync.Mutex
	wg  sync.WaitGroup
}

func trimLastString(s, sub string) string {
	if idx := strings.Index(s, sub); idx != -1 {
		return s[:idx]
	}
	return s
}

func (posters *SafePosters) scrapePoster(IDimdb string, total int, count *int) {
	c := colly.NewCollector()

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		if e.Attr("class") != "poster" {
			return
		}
		link := e.Attr("src")
		posters.mux.Lock()
		posters.URL[IDimdb] = trimLastString(link, "_V1_") + "_V1_.jpg"
		posters.mux.Unlock()
		c.Visit(e.Request.AbsoluteURL(link))
		defer posters.wg.Done()
		(*count)++
		fmt.Printf("Scraped: % 3.2f%%\r", float32(*count)*100.0/float32(total))
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://www.imdb.com/title/" + IDimdb + "/reviews")
}

func structToJSONFile(data interface{}, filename string) error {
	posterJSONBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Remove(path + filename)
	if err != nil && err.Error() != "remove "+path+filename+": no such file or directory" {
		return err
	}
	fmt.Println("Update/Create " + path + filename)
	f, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(posterJSONBody)
	if err != nil {
		return err
	}
	return nil
}

func JSONToInterface(data interface{}, filename string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := os.Open(path + filename)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(data)
	if err != nil {
		if err.Error() != "EOF" {
			return err
		}
	}
	return nil
}

func getPosters(ids []string, movies *[]Movie) map[string]string {
	fmt.Print("Start to get posters from https://www.imdb.com/\n")
	posters := SafePosters{
		URL: make(map[string]string),
	}
	total := len(ids)
	count := 0
	posters.wg.Add(total)
	for index, id := range ids {
		go posters.scrapePoster(id, total, &count)
		if index%5 == 0 {
			time.Sleep(1 * time.Second)
		}
	}
	posters.wg.Wait()
	fmt.Println("Finish")
	// Update movies with the new posters
	for i, movie := range *movies {
		if lib.StringInArray(movie.IDimdb, ids) {
			(*movies)[i].Poster = posters.URL[movie.IDimdb]
		}
	}
	return posters.URL
}

func mergeMapsString(ms ...map[string]string) map[string]string {
	res := make(map[string]string)
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func main() {
	// Etablish a connection with mongodb
	var dbName = os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "hypertube"
	}
	dbSession, errStr := mongodb.MongoDBConn(dbName)
	if errStr != "" {
		fmt.Println(errStr)
		return
	}
	db := dbSession.DB(dbName)
	var movies []Movie
	// -> Basic
	err := extractData("title.basics.tsv.gz", &movies, basicFile)
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}
	// -> Rating/Votes
	err = extractData("title.ratings.tsv.gz", &movies, ratingFile)
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}
	// -> Title Languages
	err = extractData("title.akas.tsv.gz", &movies, titlesFile)
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}
	movies, err = filterData(movies, db)
	if err != nil {
		fmt.Printf("MongoDB Find > Failed!: %v\n", err)
		return
	}
	nbMovies := len(movies)
	fmt.Printf("Need to insert %d movies\n", nbMovies)
	var oldPosters map[string]string
	// Read the json poster file in order to know if we need to scrap again the website
	err = JSONToInterface(&oldPosters, "/posters.json")
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}
	// Match json posters with new movies
	var needToScrapPosterIds []string
	for i, movie := range movies {
		if oldPosters[movie.IDimdb] != "" {
			movies[i].Poster = oldPosters[movie.IDimdb]
		} else {
			needToScrapPosterIds = append(needToScrapPosterIds, movie.IDimdb)
		}
	}
	if len(needToScrapPosterIds) < 25 {
		fmt.Println("New posters:", needToScrapPosterIds)
	} else {
		fmt.Println("New posters: Too much to display")
	}
	var newPosters map[string]string
	if len(needToScrapPosterIds) != 0 {
		newPosters = getPosters(needToScrapPosterIds, &movies)
	}
	// If new posters then update the json poster file
	posters := mergeMapsString(oldPosters, newPosters)
	// Generate a file with only the posters to avoid the scraping
	err = structToJSONFile(posters, "/posters.json")
	if err != nil {
		fmt.Printf(" > Failed!: %v\n", err)
		return
	}
	// Insert movies in the database
	fmt.Printf("Insert %d movies in the collection movies of the database\n", nbMovies)
	for index, movie := range movies {
		err := db.C("movies").Insert(movie)
		if err != nil {
			fmt.Printf("Insert data failed: IDimdb[%s] -> %v\n", movie.IDimdb, err)
		}
		fmt.Printf("Inserted: % 3.2f%%\r", float32(index)*100.0/float32(nbMovies))
	}
	fmt.Println("\rInserted: 100.00%")
}
