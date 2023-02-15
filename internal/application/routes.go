package application

import (
	"net/http"

	"github.com/fernandoocampo/basic-micro/internal/adapter/web"
	"github.com/fernandoocampo/basic-micro/internal/pets"
	"github.com/gorilla/mux"
)

type petsRouter struct {
	router    *mux.Router
	endpoints pets.Endpoints
	decoders  web.PetDecoders
	encoders  web.PetEncoders
}

func newPetsRouter(petsRouter petsRouter) http.Handler {
	petsRouter.router.Methods(http.MethodPost).Path("/pets").Handler(
		web.NewHandler().
			WithEndpoint(petsRouter.endpoints.CreatePetEndpoint).
			WithDecoder(petsRouter.decoders.CreateDecoder).
			WithEncoder(petsRouter.encoders.CreateEncoder),
	)

	petsRouter.router.Methods(http.MethodPut).Path("/pets").Handler(
		web.NewHandler().
			WithEndpoint(petsRouter.endpoints.UpdatePetEndpoint).
			WithDecoder(petsRouter.decoders.UpdateDecoder).
			WithEncoder(petsRouter.encoders.UpdateEncoder),
	)

	petsRouter.router.Methods(http.MethodDelete).Path("/pets/{id}").Handler(
		web.NewHandler().
			WithEndpoint(petsRouter.endpoints.DeletePetEndpoint).
			WithDecoder(petsRouter.decoders.DeleteDecoder).
			WithEncoder(petsRouter.encoders.DeleteEncoder),
	)

	petsRouter.router.Methods(http.MethodGet).Path("/pets/{id}").Handler(
		web.NewHandler().
			WithEndpoint(petsRouter.endpoints.GetPetWithIDEndpoint).
			WithDecoder(petsRouter.decoders.GetByIDDecoder).
			WithEncoder(petsRouter.encoders.GetByIDEncoder),
	)

	petsRouter.router.Methods(http.MethodGet).Path("/pets").Handler(
		web.NewHandler().
			WithEndpoint(petsRouter.endpoints.SearchPetsEndpoint).
			WithDecoder(petsRouter.decoders.SearchDecoder).
			WithEncoder(petsRouter.encoders.SearchEncoder),
	)

	return petsRouter.router
}
