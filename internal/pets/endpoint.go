package pets

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type GetPetWithIDEndpoint struct {
	service *Service
	logger  *zap.Logger
}

type CreatePetEndpoint struct {
	service *Service
	logger  *zap.Logger
}

type UpdatePetEndpoint struct {
	service *Service
	logger  *zap.Logger
}

type DeletePetEndpoint struct {
	service *Service
	logger  *zap.Logger
}

type SearchPetsEndpoint struct {
	service *Service
	logger  *zap.Logger
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
func NewEndpoints(service *Service, logger *zap.Logger) Endpoints {
	return Endpoints{
		CreatePetEndpoint:    MakeCreatePetEndpoint(service, logger),
		UpdatePetEndpoint:    MakeUpdatePetEndpoint(service, logger),
		DeletePetEndpoint:    MakeDeletePetEndpoint(service, logger),
		GetPetWithIDEndpoint: MakeGetPetWithIDEndpoint(service, logger),
		SearchPetsEndpoint:   MakeSearchPetsEndpoint(service, logger),
	}
}

// MakeGetPetWithIDEndpoint create endpoint for get a pet with ID service.
func MakeGetPetWithIDEndpoint(srv *Service, logger *zap.Logger) *GetPetWithIDEndpoint {
	newNewEndpoint := GetPetWithIDEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeCreatePetEndpoint create endpoint for create pet service.
func MakeCreatePetEndpoint(srv *Service, logger *zap.Logger) *CreatePetEndpoint {
	newNewEndpoint := CreatePetEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeUpdatePetEndpoint create endpoint for update pet service.
func MakeUpdatePetEndpoint(srv *Service, logger *zap.Logger) *UpdatePetEndpoint {
	newNewEndpoint := UpdatePetEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeDeletePetEndpoint create endpoint for the delete pet service.
func MakeDeletePetEndpoint(srv *Service, logger *zap.Logger) *DeletePetEndpoint {
	newNewEndpoint := DeletePetEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeSearchPetsEndpoint pet endpoint to search pets with filters.
func MakeSearchPetsEndpoint(srv *Service, logger *zap.Logger) *SearchPetsEndpoint {
	newNewEndpoint := SearchPetsEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

func (g *GetPetWithIDEndpoint) Do(ctx context.Context, request any) (any, error) {
	petID, ok := request.(PetID)
	if !ok {
		g.logger.Error("invalid pet id", zap.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid pet id")
	}

	petFound, err := g.service.QueryByID(ctx, petID)
	if err != nil {
		g.logger.Error(
			"something went wrong trying to get a pet with the given id",
			zap.Error(err),
		)
	}

	g.logger.Debug("find pet by id endpoint", zap.String("result", fmt.Sprintf("%+v", petFound)))

	return newGetPetWithIDResult(petFound, err), nil
}

func (c *CreatePetEndpoint) Do(ctx context.Context, request any) (any, error) {
	newPet, ok := request.(*NewPet)
	if !ok {
		c.logger.Error("invalid new pet type", zap.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid new pet type")
	}

	newid, err := c.service.Create(ctx, *newPet)
	if err != nil {
		c.logger.Error(
			"something went wrong trying to create a pet with the given id",
			zap.Error(err),
		)
	}
	return newCreatePetResult(newid, err), nil
}

func (u *UpdatePetEndpoint) Do(ctx context.Context, request any) (any, error) {
	updatePet, ok := request.(*UpdatePet)
	if !ok {
		u.logger.Error("invalid update pet type", zap.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid update pet type")
	}

	err := u.service.Update(ctx, *updatePet)
	if err != nil {
		u.logger.Error(
			"something went wrong trying to update a pet with the given id",
			zap.Error(err),
		)
	}

	return newUpdatePetResult(err), nil
}

func (d *DeletePetEndpoint) Do(ctx context.Context, request any) (any, error) {
	petID, ok := request.(PetID)
	if !ok {
		d.logger.Error("invalid delete pet type", zap.String("received", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid pet id type")
	}

	err := d.service.Delete(ctx, petID)
	if err != nil {
		d.logger.Error(
			"something went wrong trying to delete a pet with the given id",
			zap.Error(err),
		)
	}
	return newDeletePetResult(err), nil

}

func (s *SearchPetsEndpoint) Do(ctx context.Context, request any) (any, error) {
	petFilters, ok := request.(QueryFilter)
	if !ok {
		s.logger.Error("invalid pet filters", zap.String("received", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid pet filters")
	}

	searchResult, err := s.service.Query(ctx, petFilters)
	if err != nil {
		s.logger.Error(
			"something went wrong trying to search pets with the given filter",
			zap.Error(err),
		)
	}

	s.logger.Debug("search pets endpoint", zap.String("result", fmt.Sprintf("%+v", searchResult)))

	return newSearchPetsDataResult(searchResult, err), nil
}
