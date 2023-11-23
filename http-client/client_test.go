package http_client_test

import (
	"github.com/kpearce2430/keputils/http-client"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Do stuff BEFORE the tests!")
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}

func TestGetDefaultClient(t *testing.T) {

	const defaultURL200 = "https://www.yahoo.com"                // "https://httpstat.us/200"
	const defaultURL404 = "https://www.yahoo.com/blah/blah/blah" //"https://httpstat.us/404"

	client := http_client.GetDefaultClient(10, true)
	resp, err := client.R().Get(defaultURL200)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != "200 OK" {
		t.Logf("%+v\n", resp.Status)
		t.Fatal("Response Not 200")
	}

	resp, err = client.R().Get(defaultURL404)
	if err != nil {
		t.Fatal(err)
	}

	if resp.IsError() {
		t.Logf("Error: %s", resp.Status)
	} else {
		t.Error("Expecting an 404 Not Found Error")
	}

	if resp.Status != "404 Not Found" {
		t.Logf("%+v\n", resp.Status)
		t.Fatal("Response Not 404 Not Found")
	}

}
