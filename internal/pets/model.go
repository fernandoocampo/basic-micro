package pets

import (
	"fmt"

	"github.com/google/uuid"
)

// PetID defines pet id.
type PetID string

// OrderByField defines fields you can use to order queries.
type OrderByField string

// NewPet contains data to request the creation of a new pet.
type NewPet struct {
	Name string `json:"name"`
}

// UpdatePet contains data to request the update of a new pet.
type UpdatePet struct {
	ID   PetID  `json:"id"`
	Name string `json:"name"`
}

// Pet contains pet data.
type Pet struct {
	ID   PetID  `json:"id"`
	Name string `json:"name"`
}

// ValidationError define pet validation logic.
type ValidationError struct {
	Messages []string
}

// QueryFilter contains data for query filters.
type QueryFilter struct {
	PetName     string
	OrderBy     OrderByField
	PageNumber  uint8
	RowsPerPage uint8
}

// GetPetWithIDResult standard roesponse for get a Pet with an ID.
type GetPetWithIDResult struct {
	Pet *Pet
	Err string
}

// CreatePetResult standard response for create Pet.
type CreatePetResult struct {
	ID  PetID
	Err string
}

// UpdatePetResult standard response for updating a pet.
type UpdatePetResult struct {
	Err string
}

// DeletePetResult standard response for deleting a pet.
type DeletePetResult struct {
	Err string
}

// SearchPetsResult contains search pets result data.
type SearchPetsResult struct {
	Pets        []Pet
	Total       int
	Page        uint8
	RowsPerPage uint8
}

// SearchPetsDataResult standard roespnse for get a Pet with an ID.
type SearchPetsDataResult struct {
	SearchResult SearchPetsResult
	Err          string
}

const (
	// EmptyPetID is the pet id that empty or nil.
	EmptyPetID        = PetID("")
	EmptyOrderByField = OrderByField("")

	PageNumberDefault  = uint8(1)
	RowsPerPageDefault = uint8(10)
)

// order by field possible values
const (
	Name OrderByField = "Name"
)

func (e *ValidationError) addErrorMessage(message string) {
	e.Messages = append(e.Messages, message)
}

func (e *ValidationError) Error() string {
	// TODO improve this part
	return fmt.Sprintf("invalid pet data: %+v", e.Messages)
}

func newPetID() PetID {
	return PetID(uuid.New().String())
}

func buildNewPet(newPet NewPet) Pet {
	return Pet{
		ID:   newPetID(),
		Name: newPet.Name,
	}
}

func validPetToUpdate(pet UpdatePet) error {
	err := new(ValidationError)

	if pet.ID == "" {
		err.addErrorMessage("pet id cannot be empty")
	}

	if pet.Name == "" {
		err.addErrorMessage("pet name cannot be empty")
	}

	if len(err.Messages) > 0 {
		return err
	}

	return nil
}

func (q QueryFilter) isInvalid() bool {
	return q.PetName == ""
}

func (q *QueryFilter) fillDefaultValues() {
	if q.OrderBy == EmptyOrderByField {
		q.OrderBy = Name
	}

	if q.PageNumber == 0 {
		q.PageNumber = PageNumberDefault
	}

	if q.RowsPerPage == 0 {
		q.RowsPerPage = RowsPerPageDefault
	}
}

// newGetPetWithIDResult create a new GetPetWithIDResult
func newGetPetWithIDResult(pet *Pet, err error) GetPetWithIDResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return GetPetWithIDResult{
		Pet: pet,
		Err: errmessage,
	}
}

// newCreatePetResult create a new CreatePetResponse
func newCreatePetResult(id PetID, err error) CreatePetResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return CreatePetResult{
		ID:  id,
		Err: errmessage,
	}
}

// newUpdatePetResult udpate a new UpdatePetResponse
func newUpdatePetResult(err error) UpdatePetResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return UpdatePetResult{
		Err: errmessage,
	}
}

// newDeletePetResult udpate a new DeletePetResponse
func newDeletePetResult(err error) DeletePetResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return DeletePetResult{
		Err: errmessage,
	}
}

// newSearchPetsResult create a new SearchPetsResult
func newSearchPetsDataResult(result SearchPetsResult, err error) SearchPetsDataResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return SearchPetsDataResult{
		SearchResult: result,
		Err:          errmessage,
	}
}

func (p PetID) String() string {
	return string(p)
}
