package stl_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/philips-software/go-dip-api/stl"
	"github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	muxSTL      *http.ServeMux
	serverSTL   *httptest.Server
	token       string
	tokenSource oauth2.TokenSource
	client      *stl.Client
)

func setup(t *testing.T) (func(), error) {
	muxSTL = http.NewServeMux()
	serverSTL = httptest.NewServer(muxSTL)

	token = "44d20214-7879-4e35-923d-f9d4e01c9746"
	tokenSource = oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})

	var err error
	client, err = stl.NewClient(tokenSource, &stl.Config{
		STLAPIURL: serverSTL.URL,
	})
	if !assert.Nil(t, err) {
		t.Fatalf("invalid STL client")
		return func() {}, err
	}

	return func() {
		serverSTL.Close()
	}, nil
}

func TestDebug(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	client, err = stl.NewClient(tokenSource, &stl.Config{
		STLAPIURL: serverSTL.URL,
		DebugLog:  tmpfile,
	})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer client.Close()
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}() // clean up

	var query struct {
		App stl.AppResource `graphql:"applicationResource(id: $id, name: $name)"`
	}
	err = client.Query(context.Background(), &query, map[string]interface{}{
		"id":   graphql.Int(1),
		"name": graphql.String("name"),
	})
	assert.NotNil(t, err)

	fi, err := tmpfile.Stat()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if fi.Size() == 0 {
		t.Errorf("Expected something to be written to DebugLog")
	}
}
