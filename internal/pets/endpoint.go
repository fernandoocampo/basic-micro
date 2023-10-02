package pets

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type GetPetWithIDEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type CreatePetEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type UpdatePetEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type DeletePetEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type SearchPetsEndpoint struct {
	service *Service
	logger  *slog.Logger
}

// Endpoints is a wrapper for endpoints
type Endpoints struct {
	GetPetWithIDEndpoint *GetPetWithIDEndpoint
	CreatePetEndpoint    *CreatePetEndpoint
	UpdatePetEndpoint    *UpdatePetEndpoint
	DeletePetEndpoint    *DeletePetEndpoint
	SearchPetsEndpoint   *SearchPetsEndpoint
}

// NewEndpoints Create the endpoints for pets application.
func NewEndpoints(service *Service, logger *slog.Logger) Endpoints {
	return Endpoints{
		CreatePetEndpoint:    MakeCreatePetEndpoint(service, logger),
		UpdatePetEndpoint:    MakeUpdatePetEndpoint(service, logger),
		DeletePetEndpoint:    MakeDeletePetEndpoint(service, logger),
		GetPetWithIDEndpoint: MakeGetPetWithIDEndpoint(service, logger),
		SearchPetsEndpoint:   MakeSearchPetsEndpoint(service, logger),
	}
}

// MakeGetPetWithIDEndpoint create endpoint for get a pet with ID service.
func MakeGetPetWithIDEndpoint(srv *Service, logger *slog.Logger) *GetPetWithIDEndpoint {
	newNewEndpoint := GetPetWithIDEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeCreatePetEndpoint create endpoint for create pet service.
func MakeCreatePetEndpoint(srv *Service, logger *slog.Logger) *CreatePetEndpoint {
	newNewEndpoint := CreatePetEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeUpdatePetEndpoint create endpoint for update pet service.
func MakeUpdatePetEndpoint(srv *Service, logger *slog.Logger) *UpdatePetEndpoint {
	newNewEndpoint := UpdatePetEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeDeletePetEndpoint create endpoint for the delete pet service.
func MakeDeletePetEndpoint(srv *Service, logger *slog.Logger) *DeletePetEndpoint {
	newNewEndpoint := DeletePetEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeSearchPetsEndpoint pet endpoint to search pets with filters.
func MakeSearchPetsEndpoint(srv *Service, logger *slog.Logger) *SearchPetsEndpoint {
	newNewEndpoint := SearchPetsEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

func (g *GetPetWithIDEndpoint) Do(ctx context.Context, request any) (any, error) {
	petID, ok := request.(PetID)
	if !ok {
		g.logger.Error("invalid pet id", slog.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid pet id")
	}

	petFound, err := g.service.QueryByID(ctx, petID)
	if err != nil {
		g.logger.Error(
			"querying pet with the given id",
			slog.String("id", petID.String()),
			slog.String("error", err.Error()),
		)
	}

	g.logger.Debug("find pet by id endpoint", slog.String("result", fmt.Sprintf("%+v", petFound)))

	return newGetPetWithIDResult(petFound, err), nil
}

func (c *CreatePetEndpoint) Do(ctx context.Context, request any) (any, error) {
	newPet, ok := request.(*NewPet)
	if !ok {
		c.logger.Error("invalid new pet type", slog.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid new pet type")
	}

	newid, err := c.service.Create(ctx, *newPet)
	if err != nil {
		c.logger.Error(
			"creating pet",
			slog.String("new_pet", fmt.Sprintf("%+v", newPet)),
			slog.String("error", err.Error()),
		)
	}
	return newCreatePetResult(newid, err), nil
}

func (u *UpdatePetEndpoint) Do(ctx context.Context, request any) (any, error) {
	updatePet, ok := request.(*UpdatePet)
	if !ok {
		u.logger.Error("invalid update pet type", slog.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid update pet type")
	}

	err := u.service.Update(ctx, *updatePet)
	if err != nil {
		u.logger.Error(
			"updating a pet with the given id",
			slog.String("error", err.Error()),
		)
	}

	return newUpdatePetResult(err), nil
}

func (d *DeletePetEndpoint) Do(ctx context.Context, request any) (any, error) {
	petID, ok := request.(PetID)
	if !ok {
		d.logger.Error("invalid delete pet type", slog.String("received", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid pet id type")
	}

	err := d.service.Delete(ctx, petID)
	if err != nil {
		d.logger.Error(
			"deleting pet with the given id",
			slog.String("id", petID.String()),
			slog.String("error", err.Error()),
		)
	}
	return newDeletePetResult(err), nil

}

func (s *SearchPetsEndpoint) Do(ctx context.Context, request any) (any, error) {
	petFilters, ok := request.(QueryFilter)
	if !ok {
		s.logger.Error("invalid pet filters", slog.String("received", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid pet filters")
	}

	searchResult, err := s.service.Query(ctx, petFilters)
	if err != nil {
		s.logger.Error(
			"querying pets with the given filter",
			slog.String("filters", fmt.Sprintf("%+v", petFilters)),
			slog.String("error", err.Error()),
		)
	}

	s.logger.Debug("search pets endpoint", slog.String("result", fmt.Sprintf("%+v", searchResult)))

	return newSearchPetsDataResult(searchResult, err), nil
}
