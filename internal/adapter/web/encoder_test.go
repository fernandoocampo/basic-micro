package web_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fernandoocampo/basic-micro/internal/adapter/web"
	"github.com/fernandoocampo/basic-micro/internal/pets"
	"github.com/stretchr/testify/assert"
)

func TestEncodeCreatePet(t *testing.T) {
	// Given
	givenEndpointResult := pets.CreatePetResult{
		ID:  pets.PetID("82853922-4481-4a95-8691-30f36c61e45a"),
		Err: "",
	}

	expectedEncodedResult := web.Result{
		Success: true,
		Errors:  nil,
		Data:    "82853922-4481-4a95-8691-30f36c61e45a",
	}

	encoder := web.NewCreatePetEncoder(newDummyLogger())

	ctx := context.TODO()
	recorder := httptest.NewRecorder()

	// When
	err := encoder.Encode(ctx, recorder, givenEndpointResult)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, expectedEncodedResult, createWebResult(t, recorder.Body, pets.EmptyPetID))
}

func TestEncodeGetPetWithID(t *testing.T) {
	// Given
	givenEndpointResult := pets.GetPetWithIDResult{
		Pet: &pets.Pet{
			ID:   pets.PetID("82853922-4481-4a95-8691-30f36c61e45a"),
			Name: "drila",
		},
		Err: "",
	}

	expectedEncodedResult := web.Result{
		Success: true,
		Errors:  nil,
		Data: &web.Pet{
			ID:   "82853922-4481-4a95-8691-30f36c61e45a",
			Name: "drila",
		},
	}

	encoder := web.NewGetPetWithIDEncoder(newDummyLogger())

	ctx := context.TODO()
	recorder := httptest.NewRecorder()

	// When
	err := encoder.Encode(ctx, recorder, givenEndpointResult)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, expectedEncodedResult, createWebResult(t, recorder.Body, &web.Pet{}))
}

func createWebResult(t *testing.T, body io.Reader, data any) web.Result {
	t.Helper()

	var result web.Result
	result.Data = data
	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)

		return result
	}

	return result
}
