package comment

import (
	"log"
	"net/http"
	"os"
	"testing"

	mdb "../../../mongodb"
	"../../../tests"
	"github.com/gorilla/mux"
)

func testServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/comment/{filmId}", Comments)
	return r
}

func TestMain(m *testing.M) {
	var err string
	dbsession, err := mdb.MongoDBConn("")
	if err != "" {
		log.Fatal(err)
	}
	dbsession.Copy()
	defer dbsession.Close() // cleaning up
	tests.MongoDB = dbsession.DB("db_hypertube_tests")
	tests.InitTimeTest()
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}
