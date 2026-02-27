package couch_database_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	couchdatabase "github.com/kpearce2430/keputils/couch-database"
	"github.com/stretchr/testify/assert"
)

type TestDocument struct {
	Id    string `json:"_id,omitempty"`
	Rev   string `json:"_rev,omitempty"`
	Name  string
	Value int64
}

var url string

func TestMain(m *testing.M) {
	ctx := context.Background()

	couchDBServer, _ := couchdatabase.CreateCouchDBServer(ctx)
	defer func() {
		if err := couchDBServer.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	ip, err := couchDBServer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err := couchDBServer.MappedPort(ctx, "5984")
	if err != nil {
		log.Fatal(err)
	}

	url = fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	log.Println(url)

	_ = os.Setenv("COUCHDB_DATABASE", "tester")
	_ = os.Setenv("COUCHDB_URL", url)
	_ = os.Setenv("COUCHDB_USER", "admin")
	_ = os.Setenv("COUCHDB_PASSWORD", "password")

	m.Run()

}

func TestDatabaseConfig(t *testing.T) {
	dbConfig, err := couchdatabase.NewDatabaseConfig("")

	assert.Nil(t, err, fmt.Sprintf("%+v", err))

	assert.Equal(t, "tester", dbConfig.DatabaseName, "database name mismatch")
	assert.Equal(t, url, dbConfig.CouchDBUrl, "database url mismatch")
	assert.Equal(t, "admin", dbConfig.Username, "couchdb username mismatch")
	assert.Equal(t, "password", dbConfig.Password, "couchdb password mismatch")

	databaseStore := couchdatabase.NewDataStore[TestDocument](dbConfig)

	assert.NotNil(t, databaseStore, "database store is nil")

	if databaseStore.DatabaseCreate() != true {
		t.Fatal("Error creating a database")
	}

	log.Printf("Database created")

	testDocument := TestDocument{Name: "name", Value: 1}

	_, err = databaseStore.DocumentCreate("key", &testDocument)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDataStore(t *testing.T) {
	err := os.Setenv("MY_COUCHDB_DATABASE", "junk")
	if err != nil {
		t.Error(err)
		return
	}
	t.Setenv("MY_COUCHDB_URL", url)
	t.Setenv("MY_COUCHDB_USER", "admin")
	t.Setenv("MY_COUCHDB_PASSWORD", "password")

	databaseStore := couchdatabase.DataStore[TestDocument]("MY")

	assert.NotNil(t, databaseStore, "database store is nil")

	if databaseStore.DatabaseCreate() != true {
		t.Fatal("Error creating a database")
	}

	log.Printf("Database created")

	testDocument := TestDocument{Name: "name", Value: 1}

	_, err = databaseStore.DocumentCreate("key", &testDocument)

	if err != nil {
		t.Fatal(err)
	}

}

func TestDatabaseStore_CouchDBUp(t *testing.T) {
	databaseStore := couchdatabase.New[TestDocument]("name", url, "admin", "password")
	if databaseStore.CouchDBUp() == true {
		log.Println("Datastore Couch DB is Up")
	} else {
		t.Fatal("Couchdb not up")
	}
}

func TestCouchDBUp(t *testing.T) {
	databaseStore := couchdatabase.New[TestDocument]("name", url, "admin", "password")

	if databaseStore.DatabaseCreate() != true {
		t.Fatal("Error creating a database")
	}

	t.Log("Database created")

	testDocument := TestDocument{Name: "name", Value: 1}

	revision, err := databaseStore.DocumentCreate("key", &testDocument)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Document created revision: %s", revision)

	getDocument, err := databaseStore.DocumentGet("key")

	if err != nil {
		t.Fatal(err)
	}

	if getDocument == nil {
		t.Fatal("Document not found")
	}

	t.Logf("%s, %s, %s, %d", getDocument.Id, getDocument.Rev, getDocument.Name, getDocument.Value)

	getDocument.Name = "New Name"

	revision, err = databaseStore.DocumentUpdate(getDocument.Id, getDocument.Rev, getDocument)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Document updated new revision: %s", revision)

	getDocument, err = databaseStore.DocumentGet("key")

	if err != nil {
		t.Fatal(err)
	}

	if getDocument == nil {
		t.Fatal("Document not found")
	}

	couchDatabaseInfo, err := databaseStore.DatabaseExists()
	assert.Nil(t, err, "Database Exists returns error")
	t.Log("db info>", couchDatabaseInfo)

	log.Printf("{ %s, %s, %s, %d }", getDocument.Id, getDocument.Rev, getDocument.Name, getDocument.Value)

	revision, err = databaseStore.DocumentDelete("key", getDocument.Rev)

	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Document deleted new revision: %s", revision)
	t.Log("all done")
}

func TestNotFound(t *testing.T) {

	databaseStore := couchdatabase.New[TestDocument]("name", url, "admin", "password")
	dbInfo, err := databaseStore.DatabaseExists()
	if err != nil {
		t.Log(err.Error())
	}

	if dbInfo == nil {
		if !databaseStore.DatabaseCreate() {
			t.Log("Unable to create database")
			t.FailNow()
		}
	}

	doc, err := databaseStore.DocumentGet("any-key")
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if doc != nil {
		t.Log(doc)
		t.FailNow()
	}

}

func TestGetDataStoreByDatabaseName(t *testing.T) {
	dbName, ok := os.LookupEnv("COUCHDB_DATABASE")
	if !ok {
		t.Error("COUCHDB_DATABASE not set")
		return
	}
	db, err := couchdatabase.GetDataStoreByDatabaseName[TestDocument](dbName)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if db == nil {
		t.Error("db is nil")
	}
}
