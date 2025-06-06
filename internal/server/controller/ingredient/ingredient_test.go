package controller_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	controller "q-q-tem-pra-hoje/internal/server/controller/ingredient"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockIngredientService struct {
	addFunc             func(ingredient.Ingredient) error
	findIngredientsFunc func() ([]ingredient.Ingredient, error)
	updateFunc          func(ingredient.Ingredient) error
	lastIngredient      ingredient.Ingredient
	deleteFunc          func(uint) error
	lastDeletedId       uint
}

func (m *MockIngredientService) Add(ing ingredient.Ingredient) error {
	m.lastIngredient = ing
	if m.addFunc != nil {
		return m.addFunc(ing)
	}
	return nil
}

func (m *MockIngredientService) FindIngredients() ([]ingredient.Ingredient, error) {
	if m.findIngredientsFunc != nil {
		return m.findIngredientsFunc()
	}
	return []ingredient.Ingredient{}, nil
}

func (m *MockIngredientService) Update(ing ingredient.Ingredient) error {
	m.lastIngredient = ing
	if m.updateFunc != nil {
		return m.updateFunc(ing)
	}
	return nil
}

func (m *MockIngredientService) Delete(id uint) error {
	m.lastDeletedId = id
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func TestIngredientController_ServeHTTP(t *testing.T) {
	t.Run("it should return method not allowed", func(t *testing.T) {
		mockService := &MockIngredientService{}

		ctrl := controller.NewIngredientController(mockService)

		req := httptest.NewRequest(http.MethodPut, "/ingredient", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ctrl.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

		assert.JSONEq(t, `{"message":"method not allowed"}`, w.Body.String())
	})
}

func TestIngredientController_Add(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		mockAddFunc    func(ingredient.Ingredient) error
		expectedStatus int
		expectedBody   string
		validateMock   func(*testing.T, *MockIngredientService)
	}{
		{
			name:        "Valid ingredient",
			requestBody: `{"name":"Salt","measureType":"unit","quantity":1}`,
			mockAddFunc: func(ing ingredient.Ingredient) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			validateMock: func(t *testing.T, m *MockIngredientService) {
				assert.Equal(t, "Salt", m.lastIngredient.Name)
				assert.Equal(t, "unit", m.lastIngredient.MeasureType)
				assert.Equal(t, 1, m.lastIngredient.Quantity)
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"name":`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid or missing fields in request body"}`,
		},
		{
			name:           "Empty required fields",
			requestBody:    `{"name":"","measureType":"","quantity":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid or missing fields in request body"}`,
		},
		{
			name:        "Service error",
			requestBody: `{"name":"Salt","measureType":"unit","quantity":1}`,
			mockAddFunc: func(ing ingredient.Ingredient) error {
				return errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockIngredientService{
				addFunc: tc.mockAddFunc,
			}
			ctrl := controller.NewIngredientController(mockService)

			req := httptest.NewRequest(http.MethodPost, "/ingredient", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ctrl.Add(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}

			if tc.validateMock != nil {
				tc.validateMock(t, mockService)
			}
		})
	}
}

func TestIngredientController_GetAll(t *testing.T) {
	testCases := []struct {
		name           string
		mockFindFunc   func() ([]ingredient.Ingredient, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful retrieval",
			mockFindFunc: func() ([]ingredient.Ingredient, error) {
				id1 := int(1)
				id2 := int(2)
				return []ingredient.Ingredient{
					{Id: &id1, Name: "onion", Quantity: 20, MeasureType: "unit"},
					{Id: &id2, Name: "garlic", Quantity: 2, MeasureType: "unit"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"Id":1,"Name":"onion","MeasureType":"unit","Quantity":20},{"Id":2,"Name":"garlic","MeasureType":"unit","Quantity":2}]`,
		},
		{
			name: "Service error",
			mockFindFunc: func() ([]ingredient.Ingredient, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"failed to retrieve ingredients"}`,
		},
		{
			name: "Empty ingredients list",
			mockFindFunc: func() ([]ingredient.Ingredient, error) {
				return []ingredient.Ingredient{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockIngredientService{
				findIngredientsFunc: tc.mockFindFunc,
			}
			ctrl := controller.NewIngredientController(mockService)

			req := httptest.NewRequest(http.MethodGet, "/ingredient", nil)
			w := httptest.NewRecorder()

			ctrl.GetAll(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestIngredientController_Update(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		mockUpdateFunc func(ingredient.Ingredient) error
		expectedStatus int
		expectedBody   string
		validateMock   func(*testing.T, *MockIngredientService)
	}{
		{
			name:        "Valid update",
			requestBody: `{"name":"Salt","measureType":"unit","quantity":3}`,
			mockUpdateFunc: func(ing ingredient.Ingredient) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			validateMock: func(t *testing.T, m *MockIngredientService) {
				assert.Equal(t, "Salt", m.lastIngredient.Name)
				assert.Equal(t, "unit", m.lastIngredient.MeasureType)
				assert.Equal(t, 3, m.lastIngredient.Quantity)
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"name":`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid or missing fields in request body"}`,
		},
		{
			name:           "Empty required fields",
			requestBody:    `{"name":"","measureType":"","quantity":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid or missing fields in request body"}`,
		},
		{
			name:        "Service error",
			requestBody: `{"name":"Salt","measureType":"unit","quantity":1}`,
			mockUpdateFunc: func(ing ingredient.Ingredient) error {
				return errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"failed to update ingredient"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockIngredientService{
				updateFunc: tc.mockUpdateFunc,
			}
			ctrl := controller.NewIngredientController(mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("PATCH /ingredient/{id}", func(w http.ResponseWriter, r *http.Request) {
				ctrl.Update(w, r)
			})

			reqURL := "/ingredient/1"
			req := httptest.NewRequest(http.MethodPatch, reqURL, bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}

			if tc.validateMock != nil {
				tc.validateMock(t, mockService)
			}
		})
	}
}

func TestIngredientController_Delete(t *testing.T) {
	testCases := []struct {
		name           string
		urlPath        string
		mockDeleteFunc func(uint) error
		expectedStatus int
		expectedBody   string
		validateMock   func(*testing.T, *MockIngredientService)
	}{
		{
			name:    "Successful deletion",
			urlPath: "/ingredient?id=42",
			mockDeleteFunc: func(id uint) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
			validateMock: func(t *testing.T, m *MockIngredientService) {
				assert.Equal(t, uint(42), m.lastDeletedId)
			},
		},
		{
			name:           "Missing Id parameter",
			urlPath:        "/ingredient",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid id parameter"}`,
		},
		{
			name:           "Invalid Id format",
			urlPath:        "/ingredient?id=abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid id parameter"}`,
		},
		{
			name:    "Service error",
			urlPath: "/ingredient?id=42",
			mockDeleteFunc: func(id uint) error {
				return errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"failed to delete ingredient"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockIngredientService{
				deleteFunc: tc.mockDeleteFunc,
			}
			ctrl := controller.NewIngredientController(mockService)

			req := httptest.NewRequest(http.MethodDelete, tc.urlPath, nil)
			w := httptest.NewRecorder()

			ctrl.Delete(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}

			if tc.validateMock != nil {
				tc.validateMock(t, mockService)
			}
		})
	}
}
