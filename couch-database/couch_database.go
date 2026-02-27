package couch_database

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kelseyhightower/envconfig"
	couchdbclient "github.com/kpearce2430/keputils/couchdb-client"
	"github.com/kpearce2430/keputils/http-client"
	"github.com/segmentio/encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	errInvalidStatusResponse = errors.New("invalid status response")
)

type DatabaseStore[T interface{}] struct {
	databaseConfig *DatabaseConfig
	httpClient     *http.Client
}

func (ds DatabaseStore[T]) callCouchDB(method string, u *url.URL, body []byte, qArgs ...string) (int, []byte, error) {
	request := http.Request{
		Method: method,
		URL:    u,
		Header: make(http.Header),
	}

	if len(body) > 0 {
		request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	if len(qArgs) > 0 {
		if len(qArgs)%2 == 1 {
			logrus.Error("qargs must be even")
			return -1, []byte{}, errors.New("qargs must be even")
		}
		q := request.URL.Query()
		for i := 0; i < len(qArgs); i += 2 {
			q.Add(qArgs[i], qArgs[i+1])
		}
		request.URL.RawQuery = q.Encode()
	}

	request.SetBasicAuth(ds.databaseConfig.Username, ds.databaseConfig.Password)

	response, err := ds.httpClient.Do(&request)
	if err != nil {
		logrus.Error(err.Error())
		return -1, []byte{}, err
	}

	if response == nil {
		logrus.Error("response is nil")
		return -1, []byte{}, errors.New("response is nil")
	}

	if response.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logrus.Error(err.Error())
			}
		}(response.Body)
	}

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(err.Error())
		return response.StatusCode, []byte{}, err
	}
	return response.StatusCode, respBody, nil
}

func (ds DatabaseStore[T]) DocumentURL(key string) (*url.URL, error) {
	return url.Parse(ds.databaseConfig.DocumentURL(key))
}

func GetDataStoreByDatabaseName[T interface{}](databaseName string) (*DatabaseStore[T], error) {
	var dbConfig DatabaseConfig
	err := envconfig.Process("", &dbConfig)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	dbConfig.DatabaseName = databaseName
	datastore := NewDataStore[T](&dbConfig)
	return &datastore, nil
}

func NewDataStore[T interface{}](config *DatabaseConfig) DatabaseStore[T] {
	return DatabaseStore[T]{config, http_client.GetDefaultClient(10)}
}

func DataStore[T interface{}](prefix string) DatabaseStore[T] {
	dbConfig, _ := NewDatabaseConfig(prefix)
	dbStore := NewDataStore[T](dbConfig)
	return dbStore
}

func New[T interface{}](name string, url string, user string, pswd string) DatabaseStore[T] {
	dbConfig := DatabaseConfig{name, url, user, pswd}
	dbStore := DatabaseStore[T]{&dbConfig, http_client.GetDefaultClient(10)}
	return dbStore
}

// CreateCouchDBServer I put this here so that other test packages can use it.
func CreateCouchDBServer(ctx context.Context) (testcontainers.Container, error) {
	env := make(map[string]string)

	env["COUCHDB_USER"] = "admin"
	env["COUCHDB_PASSWORD"] = "password"

	req := testcontainers.ContainerRequest{
		Image:        "couchdb:3.3.2",
		ExposedPorts: []string{"5984/tcp"},
		WaitingFor:   wait.ForListeningPort("5984/tcp"),
		Env:          env,
	}

	couchDBServer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return couchDBServer, nil
}

func (ds DatabaseStore[T]) CouchDBUp() bool {
	u, err := url.Parse(fmt.Sprintf("%s/_up", ds.databaseConfig.CouchDBUrl))
	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	statusCode, body, err := ds.callCouchDB(http.MethodGet, u, []byte{})
	if statusCode != http.StatusOK {
		return false
	}

	var result couchdbclient.CouchDBStatus
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	if result.Status == "ok" {
		return true
	}
	return false
}

func (ds DatabaseStore[T]) GetConfig() string {
	return fmt.Sprintf("%s : %s : %s :%s",
		ds.databaseConfig.DatabaseName, ds.databaseConfig.CouchDBUrl,
		ds.databaseConfig.Username, ds.databaseConfig.Password)
}

func (ds DatabaseStore[T]) DatabaseExists() (*CouchDatabaseInfo, error) {
	createDatabaseURL, err := url.Parse(fmt.Sprintf("%s/%s", ds.databaseConfig.CouchDBUrl, ds.databaseConfig.DatabaseName))
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	statusCode, body, err := ds.callCouchDB(http.MethodGet, createDatabaseURL, []byte{})
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	switch statusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		logrus.Info("Database does not exist")
		return nil, nil
	default:
		logrus.Error("Invalid status response:", statusCode)
		return nil, errInvalidStatusResponse
	}

	var couchDatabaseInfo CouchDatabaseInfo
	err = json.Unmarshal(body, &couchDatabaseInfo)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return &couchDatabaseInfo, nil
}

func (ds DatabaseStore[T]) DatabaseCreate() bool {
	createDatabaseURL, err := url.Parse(fmt.Sprintf("%s/%s", ds.databaseConfig.CouchDBUrl, ds.databaseConfig.DatabaseName))
	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	statusCode, body, err := ds.callCouchDB(http.MethodPut, createDatabaseURL, []byte{})
	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	switch statusCode {
	case http.StatusOK, http.StatusCreated:
	default:
		logrus.Error("Invalid response:", statusCode)
		return false

	}

	var couchDBResponse couchdbclient.CouchDBResponse
	err = json.Unmarshal(body, &couchDBResponse)
	if err != nil {
		logrus.Error(err.Error())
		return false
	}
	logrus.Info("Created:", couchDBResponse)
	return couchDBResponse.Ok
}

func (ds DatabaseStore[T]) DocumentCreate(key string, document *T) (string, error) {
	documentUrl, err := ds.DocumentURL(key)
	if err != nil {
		return "", errors.New("invalid document url")
	}

	data, err := json.Marshal(document)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}

	statusCode, body, err := ds.callCouchDB(http.MethodPut, documentUrl, data)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}

	switch statusCode {
	case http.StatusOK, http.StatusCreated:
	default:
		logrus.Error("Invalid response:", statusCode)
		return "", errInvalidStatusResponse
	}

	var couchDBResponse couchdbclient.CouchDBResponse
	err = json.Unmarshal(body, &couchDBResponse)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}
	return couchDBResponse.Rev, nil
}

func (ds DatabaseStore[T]) DocumentGet(key string) (*T, error) {
	documentUrl, err := ds.DocumentURL(key)
	if err != nil {
		logrus.Error("could not create document url for key:", key)
		return nil, errors.New("could not create document url")
	}

	statusCode, body, err := ds.callCouchDB(http.MethodGet, documentUrl, []byte{})
	if err != nil {
		return nil, err
	}

	var responseDocument T
	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &responseDocument)
		if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}
	case http.StatusNotFound:
		return nil, nil
	default:
		logrus.Error("Invalid status response:", statusCode)
		return nil, errInvalidStatusResponse
	}

	return &responseDocument, nil
}

func (ds DatabaseStore[T]) DocumentUpdate(key string, revision string, document *T) (string, error) {
	documentUrl, err := ds.DocumentURL(key)
	if err != nil {
		logrus.Error("could not create document url for key:", err.Error())
		return "", err
	}

	data, err := json.Marshal(document)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}

	statusCode, body, err := ds.callCouchDB(http.MethodPut, documentUrl, data, "_rev", revision)
	if err != nil {
		return "", err
	}

	switch statusCode {
	case http.StatusOK, http.StatusCreated:
		var couchDBResponse couchdbclient.CouchDBResponse
		err = json.Unmarshal(body, &couchDBResponse)
		if err != nil {
			logrus.Error(err.Error())
			return "", err
		}
		return couchDBResponse.Rev, nil
	}
	logrus.Error("Invalid status response:", statusCode)
	return "", errInvalidStatusResponse

}

func (ds DatabaseStore[T]) DocumentDelete(key string, revision string) (string, error) {
	documentUrl, err := ds.DocumentURL(key)
	if err != nil {
		logrus.Error("could not create document url for key:", err.Error())
		return "", err
	}

	statusCode, body, err := ds.callCouchDB(http.MethodDelete, documentUrl, []byte{}, "rev", revision)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}

	switch statusCode {
	case http.StatusOK, http.StatusCreated:
		var couchDBResponse couchdbclient.CouchDBResponse
		err = json.Unmarshal(body, &couchDBResponse)
		if err != nil {
			logrus.Error(err.Error())
			return "", err
		}
		return couchDBResponse.Rev, nil
	}

	logrus.Error("Invalid status response:", statusCode)
	return "", errInvalidStatusResponse
}
