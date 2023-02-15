package web_test

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/fernandoocampo/basic-micro/internal/adapter/web"
	"github.com/fernandoocampo/basic-micro/internal/pets"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetPetWithIDDecoder(t *testing.T) {
	// Given
	var emptyBody []byte
	ctx := context.TODO()
	logger := newDummyLogger()
	decoder := web.NewGetPetWithIDDecoder(logger)
	givenPetID := "e65d36b3-ca19-4c33-8f59-917ab7399b44"

	getPetWithIDRequest := createHTTPRequest(t, emptyBody, http.MethodGet, "http://anyhost/pets/"+givenPetID)
	getPetWithIDRequest = mux.SetURLVars(getPetWithIDRequest, map[string]string{
		"id": givenPetID,
	})

	expectedRequest := pets.PetID("e65d36b3-ca19-4c33-8f59-917ab7399b44")

	// When
	got, err := decoder.Decode(ctx, getPetWithIDRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, got)
}

func TestSearchPetsDecoder(t *testing.T) {
	// Given
	var emptyBody []byte
	ctx := context.TODO()
	logger := newDummyLogger()
	pageSize := "15"
	pageNumber := "1"
	orderBy := "name"
	givenSearchName := "drila"
	decoder := web.NewSearchPetsDecoder(logger)

	searchPetsRequest := createHTTPRequest(t, emptyBody, http.MethodGet, "http://anyhost/pets")
	requestQuery := url.Values{}
	requestQuery.Add("page", pageNumber)
	requestQuery.Add("pagesize", pageSize)
	requestQuery.Add("name", givenSearchName)
	requestQuery.Add("orderby", orderBy)
	searchPetsRequest.URL.RawQuery = requestQuery.Encode()

	expectedFilter := pets.QueryFilter{
		PetName:     "drila",
		PageNumber:  1,
		RowsPerPage: 15,
		OrderBy:     "name",
	}

	// When
	got, err := decoder.Decode(ctx, searchPetsRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedFilter, got)
}

func TestCreatePetDecoder(t *testing.T) {
	// Given
	givenCreateBody := []byte(`{"name":"drila"}`)
	ctx := context.TODO()
	logger := newDummyLogger()
	decoder := web.NewCreatePetDecoder(logger)
	createPetRequest := createHTTPRequest(t, givenCreateBody, http.MethodPost, "http://anyhost/pets")
	expectedCreateRequest := &pets.NewPet{
		Name: "drila",
	}
	// When
	got, err := decoder.Decode(ctx, createPetRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedCreateRequest, got)
}

func TestUpdatePetDecoder(t *testing.T) {
	// Given
	givenUpdateBody := []byte(`{"id":"388df4d7-75a4-4690-af0d-32a73899fdc3","name":"drila"}`)
	ctx := context.TODO()
	logger := newDummyLogger()
	decoder := web.NewUpdatePetDecoder(logger)
	updatePetRequest := createHTTPRequest(t, givenUpdateBody, http.MethodPut, "http://anyhost/pets")
	expectedUpdateRequest := &pets.UpdatePet{
		ID:   "388df4d7-75a4-4690-af0d-32a73899fdc3",
		Name: "drila",
	}

	// When
	got, err := decoder.Decode(ctx, updatePetRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedUpdateRequest, got)
}

func TestDeletePetDecoder(t *testing.T) {
	// Given
	var emptyBody []byte
	ctx := context.TODO()
	logger := newDummyLogger()
	decoder := web.NewDeletePetDecoder(logger)
	givenPetID := "e65d36b3-ca19-4c33-8f59-917ab7399b44"

	deletePetRequest := createHTTPRequest(t, emptyBody, http.MethodDelete, "http://anyhost/pets/"+givenPetID)
	deletePetRequest = mux.SetURLVars(deletePetRequest, map[string]string{
		"id": givenPetID,
	})

	expectedRequest := pets.PetID("e65d36b3-ca19-4c33-8f59-917ab7399b44")

	// When
	got, err := decoder.Decode(ctx, deletePetRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, got)
}

func createHTTPRequest(t *testing.T, body []byte, httpMethod, url string) *http.Request {
	t.Helper()

	newHTTPRequest, err := http.NewRequest(
		http.MethodGet,
		url,
		bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("unexpected error creating request: %s", err)
	}

	return newHTTPRequest
}

func newDummyLogger() *zap.Logger {
	return zap.NewExample()
}
