package pets

import (
	"context"
	"errors"
	"fmt"

	"log/slog"
)

// Storer defines persistence behavior
type Storer interface {
	Save(ctx context.Context, newPet Pet) error
	Update(ctx context.Context, pet UpdatePet) error
	Delete(ctx context.Context, pet Pet) error
	Query(ctx context.Context, filter QueryFilter) (SearchPetsResult, error)
	// QueryByID find and return a pet with the given id.
	// If pet does not exist it returns a nil pet and nil error.
	QueryByID(ctx context.Context, id PetID) (*Pet, error)
}

// ServiceSetup contains service metadata.
type ServiceSetup struct {
	Storer Storer
	Logger *slog.Logger
}

// Service implements pets business logic.
type Service struct {
	storer Storer
	logger *slog.Logger
}

var (
	errSavePet    = errors.New("unable to save pet in the repository")
	errQueryPet   = errors.New("unable to query pet")
	errQueryPets  = errors.New("unable to query pets")
	errDeletePet  = errors.New("unable to delete pet")
	errUpdatePet  = errors.New("unable to update pet in the repository")
	errEmptyPetID = errors.New("pet id cannot be empty")
)

// NewService create a new pets service.
func NewService(settings ServiceSetup) *Service {
	newService := Service{
		logger: settings.Logger,
		storer: settings.Storer,
	}

	return &newService
}

// Create create a pet and store it in a database.
func (s *Service) Create(ctx context.Context, newPet NewPet) (PetID, error) {
	s.logger.Info("starting to create a new pet")
	pet := buildNewPet(newPet)

	err := s.storer.Save(ctx, pet)
	if err != nil {
		s.logger.Error("creating pet", "error", err)

		return EmptyPetID, errSavePet
	}

	s.logger.Debug(
		"pet was created",
		slog.String("id", pet.ID.String()),
	)

	return pet.ID, nil
}

// Update update a pet in a database.
func (s *Service) Update(ctx context.Context, pet UpdatePet) error {
	s.logger.Debug("starting to update pet")
	err := validPetToUpdate(pet)
	if err != nil {
		return fmt.Errorf("unable to update pet: %w", err)
	}

	err = s.storer.Update(ctx, pet)
	if err != nil {
		s.logger.Error("updating pet", "error", err)

		return errUpdatePet
	}

	return nil
}

func (s *Service) QueryByID(ctx context.Context, id PetID) (*Pet, error) {
	s.logger.Debug("starting query pet by id")
	if id == EmptyPetID {
		return nil, errEmptyPetID
	}

	pet, err := s.storer.QueryByID(ctx, id)
	if err != nil {
		s.logger.Error(
			"querying pet with id",
			"error", err,
			slog.String("id", id.String()))

		return nil, errQueryPet
	}

	return pet, nil
}

// Delete detele a pet from database.
func (s *Service) Delete(ctx context.Context, id PetID) error {
	pet, err := s.QueryByID(ctx, id)
	if err != nil {
		return errDeletePet
	}

	if pet == nil {
		s.logger.Info(
			"unable to delete pet cause it does not exist",
			slog.String("id", fmt.Sprintf("%+v", id)),
		)

		return nil
	}

	err = s.storer.Delete(ctx, *pet)
	if err != nil {
		s.logger.Error("deleting pet",
			"error", err,
			slog.String("id", id.String()))

		return errDeletePet
	}

	return nil
}

func (s *Service) Query(ctx context.Context, filter QueryFilter) (SearchPetsResult, error) {
	if filter.isInvalid() {
		return SearchPetsResult{}, nil
	}

	filter.fillDefaultValues()

	result, err := s.storer.Query(ctx, filter)
	if err != nil {
		s.logger.Error(
			"querying pets by filter",
			"error", err,
			slog.String("filter", fmt.Sprintf("%+v", filter)))

		return SearchPetsResult{}, errQueryPets
	}

	return result, nil
}
