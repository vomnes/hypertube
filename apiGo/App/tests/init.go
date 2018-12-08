package tests

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/kylelemons/godebug/pretty"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	mgo "gopkg.in/mgo.v2"
)

var (
	// MongoDB corresponds to the test database
	MongoDB *mgo.Database
	// MailjetClient corresponds to the test mailjet
	MailjetClient *mailjet.Client
	// TimeTest allows to round about time for tests
	TimeTest = time.Now()
)

// InitTimeTest allows to round about time for tests
func InitTimeTest() {
	cfg := pretty.CompareConfig
	cfg.Formatter[reflect.TypeOf(time.Time{})] = func(t time.Time) string {
		if t.Nanosecond() == 0 {
			return fmt.Sprint(t)
		}
		diff := t.Sub(TimeTest)
		if diff.Nanoseconds() < 0 {
			diff = -diff
		}
		if diff.Nanoseconds() < 50000 {
			return "Now rounded to 0.5 secondes"
		}
		return fmt.Sprintf("%d-%d-%d %d:%d:%d.%s\n", TimeTest.Year(), TimeTest.Month(), TimeTest.Day(),
			TimeTest.Hour(), TimeTest.Minute(), TimeTest.Second(), string(strconv.Itoa(TimeTest.Nanosecond())[0]))
	}
}

// DbClean delete all of rows of the tables in the test database and from redis
func DbClean() {
	if MongoDB == nil {
		log.Panic("Connection to MongoDB database failed")
	}
	collections := []string{
		"users",
		"movies",
		"comments",
	}
	for _, coll := range collections {
		if _, err := MongoDB.C(coll).RemoveAll(nil); err != nil {
			log.Fatal(err)
		}
	}
}
