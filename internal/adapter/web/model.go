package web

import "github.com/fernandoocampo/basic-micro/internal/pets"

// Result standard result for the service
type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

// Pet contains pet data.
type Pet struct {
	ID string `json:"id"`
	// Name pet's name.
	Name string `json:"name"`
}

// NewPet contains the expected data for a new pet.
type NewPet struct {
	Name string `json:"name"`
}

// UpdatePet contains the expected data to update an pet.
type UpdatePet struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreatePetResponse standard response for create Pet
type CreatePetResponse struct {
	ID  string `json:"id"`
	Err string `json:"err,omitempty"`
}

// GetPetWithIDResponse standard response for get a Pet with an ID.
type GetPetWithIDResponse struct {
	Pet *Pet   `json:"pet"`
	Err string `json:"err,omitempty"`
}

// SearchPetsResponse standard response for searching pets with filters.
type SearchPetsResponse struct {
	Pets *SearchPetsResult `json:"result"`
	Err  string            `json:"err,omitempty"`
}

// SearchPetFilter contains filters to search pets
type SearchPetFilter struct {
	// Name pet's name.
	Name string
	// Order by field
	OrderBy string
	// Page page to query
	Page uint8
	// rows per page
	PageSize uint8
}

// SearchPetsResult contains search pets result data.
type SearchPetsResult struct {
	Pets     []Pet `json:"pets"`
	Total    int   `json:"total"`
	Page     uint8 `json:"page"`
	PageSize uint8 `json:"page_size"`
}

// toPet transforms new pet to a pet object.
func toPet(pet *pets.Pet) *Pet {
	if pet == nil {
		return nil
	}
	webPet := Pet{
		ID:   pet.ID.String(),
		Name: pet.Name,
	}
	return &webPet
}

// toSearchPetResult transforms new pet to a pet object.
func toSearchPetResult(result *pets.SearchPetsResult) *SearchPetsResult {
	if result == nil {
		return nil
	}
	petsFound := make([]Pet, 0)
	for _, v := range result.Pets {
		petFound := toPet(&v)
		petsFound = append(petsFound, *petFound)
	}
	webPet := SearchPetsResult{
		Pets:     petsFound,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.RowsPerPage,
	}
	return &webPet
}

func (r Result) NotSuccess() bool {
	return !r.Success
}

func (r Result) ThereAreErrors() bool {
	return len(r.Errors) > 0
}

func (r Result) Failed() bool {
	return (r.NotSuccess() && r.ThereAreErrors())
}

// toPet transforms new pet to a pet object.
func (n *NewPet) toPet() *pets.NewPet {
	if n == nil {
		return nil
	}
	petDomain := pets.NewPet{
		Name: n.Name,
	}
	return &petDomain
}

// toPet transforms udpate pet to a pet object.
func (u *UpdatePet) toPet() *pets.UpdatePet {
	if u == nil {
		return nil
	}
	petDomain := pets.UpdatePet{
		ID:   pets.PetID(u.ID),
		Name: u.Name,
	}
	return &petDomain
}

func toCreatePetResponse(petResult pets.CreatePetResult) Result {
	var message Result
	if petResult.Err == "" {
		message.Success = true
		message.Data = petResult.ID
	}
	if petResult.Err != "" {
		message.Errors = []string{petResult.Err}
	}
	return message
}

func toUpdatePetResponse(petResult pets.UpdatePetResult) Result {
	var message Result
	if petResult.Err == "" {
		message.Success = true
	}
	if petResult.Err != "" {
		message.Errors = []string{petResult.Err}
	}
	return message
}

func toDeletePetResponse(petResult pets.DeletePetResult) Result {
	var message Result
	if petResult.Err == "" {
		message.Success = true
	}
	if petResult.Err != "" {
		message.Errors = []string{petResult.Err}
	}
	return message
}

func toGetPetWithIDResponse(petResult pets.GetPetWithIDResult) Result {
	var message Result
	newPet := toPet(petResult.Pet)
	if petResult.Err == "" {
		message.Success = true
		message.Data = newPet
	}
	if petResult.Err != "" {
		message.Errors = []string{petResult.Err}
	}
	return message
}

func toSearchPetsResponse(petResult pets.SearchPetsDataResult) Result {
	var message Result

	if petResult.Err == "" {
		message.Success = true
		message.Data = toSearchPetResult(&petResult.SearchResult)
	}
	if petResult.Err != "" {
		message.Errors = []string{petResult.Err}
	}
	return message
}

func (s SearchPetFilter) toSearchPetFilter() pets.QueryFilter {
	return pets.QueryFilter{
		PetName:     s.Name,
		PageNumber:  s.Page,
		RowsPerPage: s.PageSize,
		OrderBy:     pets.OrderByField(s.OrderBy),
	}
}
