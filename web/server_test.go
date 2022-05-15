package web

import (
	"context"
	"fmt"
	"github.com/maciej-kapusta/gomongo/config"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServerWithMongo(t *testing.T) {
	a := assert.New(t)

	timeout := 60 * time.Second
	startMongoCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	mongoC, err := tc.GenericContainer(startMongoCtx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(timeout),
		},
		Started: true,
	})

	ensure(t, err)

	defer func() { _ = mongoC.Terminate(context.TODO()) }()

	port, err := mongoC.MappedPort(context.TODO(), "27017")
	ensure(t, err)
	cfg := &config.Config{
		Port:     "8080",
		MongoUri: fmt.Sprintf("mongodb://localhost:%s", port),
		MongoDb:  "docs",
	}

	engine, err := SetupAll(cfg)
	ensure(t, err)
	ts := httptest.NewServer(engine)
	defer ts.Close()

	var addedId string
	t.Run("should post a document", func(t *testing.T) {
		doc1 := `{"name":"a","info":"b"}`

		resp, err := http.Post(
			fmt.Sprintf("%s/doc", ts.URL),
			"application/json",
			strings.NewReader(doc1),
		)
		a.Nil(err)
		a.Equal(200, resp.StatusCode)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		a.Nil(err)
		addedId = string(bodyBytes)
		a.NotEmpty(addedId)

	})

	t.Run("should get previously added document", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/doc/%s", ts.URL, addedId))
		ensure(t, err)
		a.Equal(200, resp.StatusCode)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		a.Nil(err)
		a.JSONEq(string(bodyBytes), `{"name":"a","info":"b"}`)
	})

}

func ensure(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
