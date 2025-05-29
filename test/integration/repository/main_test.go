package respository_integration_test

import (
	"database/sql"
	"os"
	"q-q-tem-pra-hoje/internal/domain/ingredient"
	"q-q-tem-pra-hoje/internal/domain/recipe"
	"q-q-tem-pra-hoje/internal/testutil"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	dsn, teardown := testutil.SetupTestDB()

	db := testutil.Connect(dsn)
	testutil.SetDB(db)
	testutil.RunMigrations(db)
	code := m.Run()

	defer teardown()

	os.Exit(code)
}


func cleanUpTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE ingredients_storage RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE recipes RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatal(err)
	}
}

func createDataset(t *testing.T, db *sql.DB) {
	query := `INSERT INTO ingredients_storage(id, name, measure_type, quantity) 
            VALUES (1, $1, $2, $3), (2, $4, $5, $6);`

	_, err := db.Exec(query, "onion", "unit", 10, "garlic", "unit", 10)

	if err != nil {
		t.Fatal(err)
	}

	testRecipes := []recipe.Recipe{
		{Name: "Rice with Onion and Garlic",
			Ingredients: []ingredient.Ingredient{
				{Name: "Onion", MeasureType: "unit", Quantity: 1},
				{Name: "Rice", MeasureType: "mg", Quantity: 500},
				{Name: "Garlic", MeasureType: "unit", Quantity: 2},
			},
		},
		{
			Name: "Tomato Soup",
			Ingredients: []ingredient.Ingredient{
				{Name: "Tomato", MeasureType: "unit", Quantity: 4},
				{Name: "Water", MeasureType: "ml", Quantity: 500},
				{Name: "Salt", MeasureType: "mg", Quantity: 10},
			},
		},
	}

	for _, tr := range testRecipes {
		var recipeId int
		err := db.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id", tr.Name).Scan(&recipeId)
		if err != nil {
			t.Fatal(err)
		}

		for _, ing := range tr.Ingredients {
			_, err := db.Exec(`
                INSERT INTO recipes_ingredients (recipe_id, name, measure_type, quantity)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (recipe_id, name) DO NOTHING;
            `, recipeId, ing.Name, ing.MeasureType, ing.Quantity)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
