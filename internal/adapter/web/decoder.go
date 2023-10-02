package web

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fernandoocampo/basic-micro/internal/pets"
	"github.com/gorilla/mux"
)

type GetPetWithIDDecoder struct {
	logger *slog.Logger
}

type SearchPetsDecoder struct {
	logger *slog.Logger
}

type CreatePetDecoder struct {
	logger *slog.Logger
}

type UpdatePetDecoder struct {
	logger *slog.Logger
}

type DeletePetDecoder struct {
	logger *slog.Logger
}

type PetDecoders struct {
	GetByIDDecoder *GetPetWithIDDecoder
	SearchDecoder  *SearchPetsDecoder
	CreateDecoder  *CreatePetDecoder
	UpdateDecoder  *UpdatePetDecoder
	DeleteDecoder  *DeletePetDecoder
}

func NewPetDecoders(logger *slog.Logger) PetDecoders {
	newDecoders := PetDecoders{
		GetByIDDecoder: NewGetPetWithIDDecoder(logger),
		SearchDecoder:  NewSearchPetsDecoder(logger),
		CreateDecoder:  NewCreatePetDecoder(logger),
		UpdateDecoder:  NewUpdatePetDecoder(logger),
		DeleteDecoder:  NewDeletePetDecoder(logger),
	}

	return newDecoders
}

func NewGetPetWithIDDecoder(logger *slog.Logger) *GetPetWithIDDecoder {
	newDecoder := GetPetWithIDDecoder{
		logger: logger,
	}

	return &newDecoder
}

func NewSearchPetsDecoder(logger *slog.Logger) *SearchPetsDecoder {
	newDecoder := SearchPetsDecoder{
		logger: logger,
	}

	return &newDecoder
}

func NewCreatePetDecoder(logger *slog.Logger) *CreatePetDecoder {
	newDecoder := CreatePetDecoder{
		logger: logger,
	}

	return &newDecoder
}

func NewUpdatePetDecoder(logger *slog.Logger) *UpdatePetDecoder {
	newDecoder := UpdatePetDecoder{
		logger: logger,
	}

	return &newDecoder
}

func NewDeletePetDecoder(logger *slog.Logger) *DeletePetDecoder {
	newDecoder := DeletePetDecoder{
		logger: logger,
	}

	return &newDecoder
}

func (g *GetPetWithIDDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	petIDParam, ok := v["id"]
	if !ok {
		return nil, errors.New("pet ID was not provided")
	}
	return pets.PetID(petIDParam), nil
}

func (g *DeletePetDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	petIDParam, ok := v["id"]
	if !ok {
		return nil, errors.New("pet ID was not provided")
	}
	return pets.PetID(petIDParam), nil
}

func (s *SearchPetsDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	filterRequest := SearchPetFilter{
		Page:     1,
		PageSize: 10,
	}

	filters := r.URL.Query()

	if v, ok := filters["name"]; ok {
		filterRequest.Name = v[0]
	}

	if v, ok := filters["page"]; ok {
		page, err := strconv.Atoi(v[0])
		if err != nil {
			log.Println("level", "ERROR", "invalid page parameter, it must be an integer", "error", err)
			page = 1
		}
		filterRequest.Page = uint8(page)
	}
	if v, ok := filters["pagesize"]; ok {
		pageSize, err := strconv.Atoi(v[0])
		if err != nil {
			log.Println("level", "ERROR", "invalid page size parameter, it must be an integer", "error", err)
			pageSize = 10
		}
		filterRequest.PageSize = uint8(pageSize)
	}

	if v, ok := filters["orderby"]; ok {
		filterRequest.OrderBy = v[0]
	}

	filter := filterRequest.toSearchPetFilter()

	return filter, nil
}

func (c *CreatePetDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	log.Println("level", "DEBUG", "msg", "decoding new pet request")
	var req NewPet
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("level", "ERROR", "new pet request could not be decoded. Request: %q because of: %s", string(body), err.Error())
		return nil, err
	}

	log.Println("level", "DEBUG", "msg", "pet request was decoded", "request", req)

	domainPet := req.toPet()

	return domainPet, nil
}

func (u *UpdatePetDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	log.Println("level", "DEBUG", "msg", "decoding update pet request")
	var req UpdatePet
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("level", "ERROR", "update pet request could not be decoded. Request: %q because of: %s", string(body), err.Error())
		return nil, err
	}

	log.Println("level", "DEBUG", "msg", "pet request was decoded", "request", req)

	domainPet := req.toPet()

	return domainPet, nil
}
