package pets_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fernandoocampo/basic-micro/internal/pets"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	// Given
	newPet := pets.NewPet{
		Name: "drila",
	}

	storerMock := newStorerMock()

	settings := pets.ServiceSetup{
		Storer: storerMock,
		Logger: newLogger(),
	}

	service := pets.NewService(settings)

	ctx := context.TODO()

	// When
	petID, err := service.Create(ctx, newPet)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, storerMock.ids[0], petID)
}

func TestCreateButError(t *testing.T) {
	t.Parallel()

	// Given
	newPet := pets.NewPet{
		Name: "drila",
	}

	expectedPetID := pets.PetID("")
	expectedError := errors.New("unable to save pet in the repository")

	storerMock := newStorerMock(
		withError(errors.New("error saving")),
	)

	settings := pets.ServiceSetup{
		Storer: storerMock,
		Logger: newLogger(),
	}

	service := pets.NewService(settings)

	ctx := context.TODO()

	// When
	petID, err := service.Create(ctx, newPet)

	// Then
	assert.Error(t, err)
	assert.Empty(t, storerMock.ids)
	assert.Equal(t, expectedPetID, petID)
	assert.Equal(t, expectedError, err)
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	// Given
	updatePet := pets.UpdatePet{
		ID:   pets.PetID("858455b7-e182-4122-a1b6-132c64d2f77b"),
		Name: "drila",
	}

	storerMock := newStorerMock()

	settings := pets.ServiceSetup{
		Storer: storerMock,
		Logger: newLogger(),
	}

	service := pets.NewService(settings)

	ctx := context.TODO()

	// When
	err := service.Update(ctx, updatePet)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, updatePet, storerMock.updatedPet)
}

func TestUpdateButError(t *testing.T) {
	t.Parallel()

	// Given
	updatePet := pets.UpdatePet{
		ID:   pets.PetID("858455b7-e182-4122-a1b6-132c64d2f77b"),
		Name: "drila",
	}

	expectedError := errors.New("unable to update pet in the repository")

	storerMock := newStorerMock(
		withError(errors.New("error updating")),
	)

	settings := pets.ServiceSetup{
		Storer: storerMock,
		Logger: newLogger(),
	}

	service := pets.NewService(settings)

	ctx := context.TODO()

	// When
	err := service.Update(ctx, updatePet)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	// Given
	petID := pets.PetID("858455b7-e182-4122-a1b6-132c64d2f77b")
	foundPet := pets.Pet{
		ID:   pets.PetID("858455b7-e182-4122-a1b6-132c64d2f77b"),
		Name: "drila",
	}

	storerMock := newStorerMock(
		withFoundPet(foundPet),
	)

	settings := pets.ServiceSetup{
		Storer: storerMock,
		Logger: newLogger(),
	}

	service := pets.NewService(settings)

	ctx := context.TODO()

	// When
	err := service.Delete(ctx, petID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, foundPet, storerMock.foundPet)
}

func TestQuery(t *testing.T) {
}

func TestQueryByID(t *testing.T) {
}

type storerMock struct {
	pets.Storer
	err        error
	ids        []pets.PetID
	deletedIDs []pets.PetID

	updatedPet pets.UpdatePet
	deletedPet pets.Pet
	foundPet   pets.Pet
}

func newStorerMock(options ...func(*storerMock)) *storerMock {
	newStorerMock := storerMock{
		ids: make([]pets.PetID, 0),
	}

	for _, opt := range options {
		opt(&newStorerMock)
	}

	return &newStorerMock
}

func withError(err error) func(*storerMock) {
	return func(s *storerMock) {
		s.err = err
	}
}

func withFoundPet(pet pets.Pet) func(*storerMock) {
	return func(s *storerMock) {
		s.foundPet = pet
	}
}

func (s *storerMock) Save(ctx context.Context, newPet pets.Pet) error {
	if s.err != nil {
		return s.err
	}

	s.ids = append(s.ids, newPet.ID)

	return nil
}

func (s *storerMock) Update(ctx context.Context, pet pets.UpdatePet) error {
	if s.err != nil {
		return s.err
	}

	s.updatedPet = pet

	return nil
}

func (s *storerMock) Delete(ctx context.Context, pet pets.Pet) error {
	if s.err != nil {
		return s.err
	}

	s.deletedPet = pet

	return nil
}

func (s *storerMock) QueryByID(ctx context.Context, id pets.PetID) (*pets.Pet, error) {
	if s.err != nil {
		return nil, s.err
	}

	return &s.foundPet, nil
}

func newLogger() *zap.Logger {
	return zap.NewExample()
}
