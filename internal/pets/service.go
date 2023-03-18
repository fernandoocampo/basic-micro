package pets

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
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
	Logger *zap.Logger
}

// Service implements pets business logic.
type Service struct {
	storer Storer
	logger *zap.Logger
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
	pet := buildNewPet(newPet)

	err := s.storer.Save(ctx, pet)
	if err != nil {
		s.logger.Error("unable to create pet", zap.Error(err))

		return EmptyPetID, errSavePet
	}

	s.logger.Debug(
		"pet was created",
		zap.String("id", pet.ID.String()),
	)

	return pet.ID, nil
}

// Update update a pet in a database.
func (s *Service) Update(ctx context.Context, pet UpdatePet) error {
	err := validPetToUpdate(pet)
	if err != nil {
		return fmt.Errorf("unable to update pet: %w", err)
	}

	err = s.storer.Update(ctx, pet)
	if err != nil {
		s.logger.Error("unable to update pet", zap.Error(err))

		return errUpdatePet
	}

	return nil
}

func (s *Service) QueryByID(ctx context.Context, id PetID) (*Pet, error) {
	if id == EmptyPetID {
		return nil, errEmptyPetID
	}

	pet, err := s.storer.QueryByID(ctx, id)
	if err != nil {
		s.logger.Error(
			"unable to query pet by id",
			zap.String("id", fmt.Sprintf("%+v", id)),
			zap.Error(err))

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
			zap.String("id", fmt.Sprintf("%+v", id)),
		)

		return nil
	}

	err = s.storer.Delete(ctx, *pet)
	if err != nil {
		s.logger.Error("unable to delete pet",
			zap.String("id", fmt.Sprintf("%+v", id)),
			zap.Error(err))

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
			"unable to query pets by filter",
			zap.String("filter", fmt.Sprintf("%+v", filter)),
			zap.Error(err))

		return SearchPetsResult{}, errQueryPet
	}

	return result, nil
}
