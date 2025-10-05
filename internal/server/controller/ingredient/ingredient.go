package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"strconv"
)

var (
	ErrInvalidRequestBody    = errors.New("invalid or missing fields in request body")
	ErrMethodNotAllowed      = errors.New("method not allowed")
	ErrInvalidId             = errors.New("invalid id parameter")
	ErrMissingId             = errors.New("id parameter is required")
	ErrInternalServerError   = errors.New("internal server error")
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type IngredientController struct {
	service ingredient.IngredientStorageProvider
}

func NewIngredientController(service ingredient.IngredientStorageProvider) *IngredientController {
	if service == nil {
		panic("ingredient service cannot be nil")
	}
	return &IngredientController{service: service}
}

func (ic *IngredientController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ic.Add(w, r)
	case http.MethodGet:
		ic.GetAll(w, r)
	case http.MethodPatch:
		ic.Update(w, r)
	case http.MethodDelete:
		ic.Delete(w, r)
	default:
		ic.respondWithError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	}
}

type IngredientInput struct {
	Name        string `json:"name"`
	MeasureType string `json:"measureType"`
	Quantity    int    `json:"quantity"`
}

func (ic *IngredientController) Add(w http.ResponseWriter, r *http.Request) {
	var input IngredientInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ic.respondWithError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	if input.Name == "" || input.MeasureType == "" {
		ic.respondWithError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	ing := ingredient.NewIngredient(nil, input.Name, input.MeasureType, input.Quantity)

	if err := ic.service.Add(ing); err != nil {
		ic.respondWithError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	ic.respondWithJSON(w, http.StatusCreated, nil)
}

func (ic *IngredientController) GetAll(w http.ResponseWriter, r *http.Request) {
	ingredients, err := ic.service.FindIngredients()
	if err != nil {
		ic.respondWithError(w, http.StatusInternalServerError, errors.New("failed to retrieve ingredients"))
		return
	}

	ic.respondWithJSON(w, http.StatusOK, ingredients)
}

func (ic *IngredientController) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if id == "" || err != nil {
		ic.respondWithError(w, http.StatusBadRequest, ErrMissingId)
		return
	}
	var input IngredientInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ic.respondWithError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	if input.Name == "" || input.MeasureType == "" {
		ic.respondWithError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	updatedIngredient := ingredient.NewIngredient(&intId, input.Name, input.MeasureType, input.Quantity)

	if err := ic.service.Update(updatedIngredient); err != nil {
		ic.respondWithError(w, http.StatusInternalServerError, errors.New("failed to update ingredient"))
		return
	}

	ic.respondWithJSON(w, http.StatusOK, nil)
}

func (ic *IngredientController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := ic.extractIdFromPath(r.URL.Query().Get("id"))

	if err != nil {
		ic.respondWithError(w, http.StatusBadRequest, ErrInvalidId)
		return
	}

	if err := ic.service.Delete(uint(id)); err != nil {
		ic.respondWithError(w, http.StatusInternalServerError, errors.New("failed to delete ingredient"))
		return
	}

	ic.respondWithJSON(w, http.StatusNoContent, nil)
}

func (ic *IngredientController) respondWithError(w http.ResponseWriter, code int, err error) {
	ic.respondWithJSON(w, code, Response{Message: err.Error()})
}

func (ic *IngredientController) respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
func (ic *IngredientController) extractIdFromPath(idString string) (uint, error) {
	if idString == "" {
		return 0, ErrMissingId
	}

	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		return 0, ErrInvalidId
	}
	return uint(id), nil

}
