package e2e_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"q-q-tem-pra-hoje/internal/repository/postgres"
	controller "q-q-tem-pra-hoje/internal/server/controller/recipe"
	recipeService "q-q-tem-pra-hoje/internal/service/recipe"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupDatabase(t *testing.T) *sql.DB {
	err := godotenv.Load("../../../../.env")

	if err != nil {
		t.Fatalf("error loading .env files: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS recipes (id SERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE);")
	if err != nil {
		t.Fatalf("failed to create table recipes: %v", err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS recipes_ingredients (
        recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        measure_type TEXT NOT NULL,
        quantity INT NOT NULL,
        PRIMARY KEY (recipe_id,name));
    `)

	if err != nil {
		t.Fatalf("failed to create table recipes_ingredients: %v", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ingredients_storage (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            measure_type TEXT NOT NULL,
            quantity INT NOT NULL
        );
    `)

	if err != nil {
		t.Fatalf("failed to create the ingredients_storage table: %v", err)
	}
	return db
}

func TestRecipeController_Add(t *testing.T) {
	db := setupDatabase(t)
	t.Cleanup(func() {

		_, err := db.Exec(`DROP TABLE recipes_ingredients`)

		if err != nil {
			t.Fatalf("failed to drop table recipes_ingredients: %v", err)
		}

		_, err = db.Exec(`DROP TABLE recipes`)

		if err != nil {
			t.Fatalf("failed to drop table recipes: %v", err)
		}

		db.Close()

	})

	repository := postgres.NewRecipeManager(db)
	service := recipeService.NewRecipeService(repository)
	controller := controller.RecipeController{RecipeProvider: service}
	ts := httptest.NewServer(controller)

	defer ts.Close()

	t.Run("should create a recipe", func(t *testing.T) {

		body := `{"name":"Rice", "ingredients": [
        {"name": "Onion", "measureType":"unit","quantity":1},
        {"name": "Rice", "measureType":"mg","quantity":500},
        {"name": "Garlic", "measureType":"unit","quantity":2}
      ]}`

		resp, err := http.Post(ts.URL+"/recipe", "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to get ingredients: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
	})
}

// func TestRecipeController_GetRecommendation(t *testing.T) {
// 	db := setupDatabase(t)
// 	t.Cleanup(func() {
//
// 		_, err := db.Exec(`DROP TABLE ingredients_storage`)
//
// 		if err != nil {
// 			t.Fatalf("failed to drop table ingredients_storage: %v", err)
// 		}
//
// 		_, err = db.Exec(`DROP TABLE recipes_ingredients`)
//
// 		if err != nil {
// 			t.Fatalf("failed to drop table recipes_ingredients: %v", err)
// 		}
//
// 		_, err = db.Exec(`DROP TABLE recipes`)
//
// 		if err != nil {
// 			t.Fatalf("failed to drop table recipes: %v", err)
// 		}
//
// 		db.Close()
//
// 	})
//
// 	query := `INSERT INTO ingredients_storage(name, measure_type, quantity)
//             VALUES ($1, $2, $3), ($4, $5, $6);`
// 	_, err := db.Exec(query, "Onion", "unit", 1, "Rice", "mg", 500)
//
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	recipes := []recipe.Recipe{
// 		{Name: "Rice with Onion and Garlic", Ingredients: []ingredient.Ingredient{
// 			{Name: "Onion", MeasureType: "unit", Quantity: 1},
// 			{Name: "Rice", MeasureType: "mg", Quantity: 500},
// 			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
// 		}},
// 		{Name: "Rice with Garlic", Ingredients: []ingredient.Ingredient{
// 			{Name: "Rice", MeasureType: "mg", Quantity: 500},
// 			{Name: "Garlic", MeasureType: "unit", Quantity: 2},
// 		}},
// 		{Name: "Rice with Onion", Ingredients: []ingredient.Ingredient{
// 			{Name: "Onion", MeasureType: "unit", Quantity: 1},
// 			{Name: "Rice", MeasureType: "mg", Quantity: 500},
// 		}},
// 		{Name: "Fries", Ingredients: []ingredient.Ingredient{
// 			{Name: "Potato", MeasureType: "unit", Quantity: 2},
// 		}},
// 	}
//
// 	for _, recipe := range recipes {
// 		var recipeId int
// 		err := db.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id;", recipe.Name).Scan(&recipeId)
// 		if err != nil {
// 			t.Fatalf("failed to insert recipe %q: %v", recipe.Name, err)
// 		}
//
// 		for _, ing := range recipe.Ingredients {
// 			_, err := db.Exec(`
//                 INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
//                 VALUES ($1, $2, $3, $4)
//                 ON CONFLICT (recipe_id, name) DO NOTHING;
//             `, recipeId, ing.Name, ing.MeasureType, ing.Quantity)
// 			if err != nil {
// 				t.Fatalf("failed to insert ingredient %q for recipe %q: %v", ing.Name, recipe.Name, err)
// 			}
// 		}
// 	}
//
// 	recipeRepository := postgres.NewRecipeManager(db)
// 	ingredientRepository := postgres.NewIngredientStorageManager(db)
// 	recipeService := recipeService.NewRecipeService(recipeRepository)
// 	ingredientService := ingredientService.NewService(&ingredientRepository)
// 	controller := controller.RecipeController{RecipeProvider: recipeService, IngredientProvider: ingredientService}
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/recommendation", controller.GetRecommendation)
// 	ts := httptest.NewServer(mux)
//
// 	defer ts.Close()
//
// 	t.Run("should provide the recommendations", func(t *testing.T) {
//
// 		recipes := []recipe.Recipe{
// 			{Name: "Rice with Onion and Garlic", Ingredients: []ingredient.Ingredient{
// 				{Name: "Onion", MeasureType: "unit", Quantity: 1},
// 				{Name: "Rice", MeasureType: "mg", Quantity: 500},
// 				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
// 			}},
// 			{Name: "Rice with Garlic", Ingredients: []ingredient.Ingredient{
// 				{Name: "Rice", MeasureType: "mg", Quantity: 500},
// 				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
// 			}},
// 			{Name: "Rice with Onion", Ingredients: []ingredient.Ingredient{
// 				{Name: "Onion", MeasureType: "unit", Quantity: 1},
// 				{Name: "Rice", MeasureType: "mg", Quantity: 500},
// 			}},
// 			{Name: "Fries", Ingredients: []ingredient.Ingredient{
// 				{Name: "Potato", MeasureType: "unit", Quantity: 2},
// 			}},
// 		}
//
// 		resp, err := http.Get(ts.URL + "/recommendation")
// 		if err != nil {
// 			t.Fatalf("Failed to get recommendations: %v", err)
// 		}
//
// 		if resp.StatusCode != http.StatusOK {
// 			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
// 		}
//
// 		expectedRecommendations := []recipe.Recommendation{
// 			{Recommendation: 1, Recipe: recipes[2]},
// 			{Recommendation: 2, Recipe: recipes[0]},
// 			{Recommendation: 3, Recipe: recipes[1]},
// 			{Recommendation: 4, Recipe: recipes[3]},
// 		}
//
// 		if err != nil {
// 			t.Fatalf("error while encoding expectedRecommendationsJSON: %v", err)
// 		}
// 		var recommendations []recipe.Recommendation
// 		err = json.NewDecoder(resp.Body).Decode(&recommendations)
// 		if err != nil {
// 			t.Fatalf("error while decoding recommendations: %v", err)
// 		}
// 		assert.Equal(t, expectedRecommendations, recommendations)
//
// 	})
// }
