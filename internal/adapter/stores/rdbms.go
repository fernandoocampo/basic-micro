package stores

import (
	"context"
	"log/slog"

	"github.com/fernandoocampo/basic-micro/internal/pets"
)

type Setup struct {
	Logger *slog.Logger
}

// Store handles logic to persist data from this microservice.
type Store struct {
	logger *slog.Logger
}

func NewStore(setup Setup) *Store {
	newStore := Store{
		logger: setup.Logger,
	}

	return &newStore
}

func (s *Store) Save(ctx context.Context, newPet pets.Pet) error {
	s.logger.Info("Saving new pet in database")
	return nil
}

func (s *Store) Update(ctx context.Context, pet pets.UpdatePet) error {
	s.logger.Info("Updating new pet in database")
	return nil
}

func (s *Store) Delete(ctx context.Context, pet pets.Pet) error {
	s.logger.Info("Deleting new pet in database")
	return nil
}

func (s *Store) Query(ctx context.Context, filter pets.QueryFilter) (pets.SearchPetsResult, error) {
	s.logger.Info("Querying pets in database")
	result := pets.SearchPetsResult{
		Pets: []pets.Pet{
			{
				ID:   pets.PetID("56016eaf-5e15-44db-839c-ef4f7f9df437"),
				Name: "Drila",
			},
			{
				ID:   pets.PetID("ec665f5e-da4e-4f51-bc4c-310dd7cc9590"),
				Name: "Michael",
			},
		},
		Total:       2,
		Page:        filter.PageNumber,
		RowsPerPage: filter.RowsPerPage,
	}
	return result, nil
}

func (s *Store) QueryByID(ctx context.Context, id pets.PetID) (*pets.Pet, error) {
	s.logger.Info("Querying pet by id in database")
	pet := pets.Pet{
		ID:   pets.PetID("56016eaf-5e15-44db-839c-ef4f7f9df437"),
		Name: "Drila",
	}

	s.logger.Info("a pet was found", slog.Any("pet", pet))

	return &pet, nil
}
