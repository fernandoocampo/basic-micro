package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/fernandoocampo/basic-micro/internal/pets"
)

type GetPetWithIDEncoder struct {
	logger *slog.Logger
}

type SearchPetsEncoder struct {
	logger *slog.Logger
}

type CreatePetEncoder struct {
	logger *slog.Logger
}

type UpdatePetEncoder struct {
	logger *slog.Logger
}

type DeletePetEncoder struct {
	logger *slog.Logger
}

type PetEncoders struct {
	GetByIDEncoder *GetPetWithIDEncoder
	SearchEncoder  *SearchPetsEncoder
	CreateEncoder  *CreatePetEncoder
	UpdateEncoder  *UpdatePetEncoder
	DeleteEncoder  *DeletePetEncoder
}

var (
	errUnableToEncodeResult = errors.New("unable to encode the result")
)

func NewPetEncoders(logger *slog.Logger) PetEncoders {
	newEncoders := PetEncoders{
		GetByIDEncoder: NewGetPetWithIDEncoder(logger),
		SearchEncoder:  NewSearchPetsEncoder(logger),
		CreateEncoder:  NewCreatePetEncoder(logger),
		UpdateEncoder:  NewUpdatePetEncoder(logger),
		DeleteEncoder:  NewDeletePetEncoder(logger),
	}

	return newEncoders
}

func NewGetPetWithIDEncoder(logger *slog.Logger) *GetPetWithIDEncoder {
	newEncoder := GetPetWithIDEncoder{
		logger: logger,
	}

	return &newEncoder
}

func NewSearchPetsEncoder(logger *slog.Logger) *SearchPetsEncoder {
	newEncoder := SearchPetsEncoder{
		logger: logger,
	}

	return &newEncoder
}

func NewCreatePetEncoder(logger *slog.Logger) *CreatePetEncoder {
	newEncoder := CreatePetEncoder{
		logger: logger,
	}

	return &newEncoder
}

func NewUpdatePetEncoder(logger *slog.Logger) *UpdatePetEncoder {
	newEncoder := UpdatePetEncoder{
		logger: logger,
	}

	return &newEncoder
}

func NewDeletePetEncoder(logger *slog.Logger) *DeletePetEncoder {
	newEncoder := DeletePetEncoder{
		logger: logger,
	}

	return &newEncoder
}

func (c *CreatePetEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(pets.CreatePetResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to pets.CreatePetResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build create pet response")
	}

	err := encodeResultWithJSON(w, toCreatePetResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode create pet result: %w", err)
	}

	return nil
}

func (u *UpdatePetEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(pets.UpdatePetResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to pets.UpdatePetResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build update pet response")
	}

	err := encodeResultWithJSON(w, toUpdatePetResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode update pet result: %w", err)
	}

	return nil
}

func (u *DeletePetEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(pets.DeletePetResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to pets.DeletePetResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build delete pet response")
	}

	err := encodeResultWithJSON(w, toDeletePetResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode delete pet result: %w", err)
	}

	return nil
}

func (g *GetPetWithIDEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(pets.GetPetWithIDResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to pets.GetPetWithIDResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build get pet response")
	}

	err := encodeResultWithJSON(w, toGetPetWithIDResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode get pet by id result: %w", err)
	}

	return nil
}

func (s *SearchPetsEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(pets.SearchPetsDataResult)
	if !ok {
		log.Println("level", "ERROR", "msg", "cannot transform to pets.SearchPetsDataResult", "received", fmt.Sprintf("%T", response))
		return errors.New("cannot build search pets response")
	}

	err := encodeResultWithJSON(w, toSearchPetsResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode search pets result: %w", err)
	}

	return nil
}

func encodeResultWithJSON(w http.ResponseWriter, message Result) error {
	w.Header().Set("Content-Type", "application/json")

	if message.Failed() {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		return errUnableToEncodeResult
	}

	return nil
}
