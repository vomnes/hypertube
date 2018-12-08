package mail

import (
	"log"
	"os"
	"testing"

	"../../../lib"
	mdb "../../../mongodb"
	"../../../tests"
)

func TestMain(m *testing.M) {
	var err string
	dbsession, err := mdb.MongoDBConn("")
	if err != "" {
		log.Fatal(err)
	}
	dbsession.Copy()
	defer dbsession.Close() // cleaning up
	tests.MongoDB = dbsession.DB("db_hypertube_tests")
	tests.MailjetClient = lib.MailJetConn()
	tests.DbClean()
	ret := m.Run()
	tests.DbClean()
	os.Exit(ret)
}
