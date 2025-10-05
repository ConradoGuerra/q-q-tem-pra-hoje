# O que é que tem para hoje?

"q-q-tem-pra-hoje" (a rough translation of "what's for today?" in Portuguese) is a recipe management and recommendation system. It allows users to manage their ingredients, add recipes, and get recommendations for what to cook based on the ingredients they have.

## Features

*   **Ingredient Management:** Add, update, delete, and view ingredients.
*   **Recipe Management:** Create, delete, and view recipes.
*   **Recipe Recommendations:** Get recipe recommendations based on available ingredients.

## Technologies

*   **Backend:** Go
*   **Database:** PostgreSQL
*   **Containerization:** Docker

## Getting Started

### Prerequisites

*   [Go](https://golang.org/doc/install)
*   [Docker](https://docs.docker.com/get-docker/)

### Installation & Running

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/q-q-tem-pra-hoje.git
    cd q-q-tem-pra-hoje
    ```

2.  **Start the database:**
    ```bash
    docker-compose up -d
    ```

3.  **Run the application:**
    ```bash
    go run cmd/main.go
    ```

The application will be running at `http://localhost:8080`.

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management.

### Running Migrations

- **Apply migrations:** `migrate -path migrations -database "postgres://user:secret@localhost:6432/q-q-tem-pra-hj-db?sslmode=disable" up`
- **Rollback last migration:** `migrate -path migrations -database "postgres://user:secret@localhost:6432/q-q-tem-pra-hj-db?sslmode=disable" down 1`
- **Check migration status:** `migrate -path migrations -database "postgres://user:secret@localhost:6432/q-q-tem-pra-hj-db?sslmode=disable" version`

Replace the database URL with your actual connection string.

### Creating New Migrations

To create a new migration:

1. `migrate create -ext sql -dir migrations -seq add_new_table`
2. Edit the generated `up` and `down` SQL files.
3. Run `migrate up` to apply.

## API Endpoints

### Ingredients

*   `POST /ingredient`: Add a new ingredient.
    *   **Body:**
        ```json
        {
          "name": "Flour",
          "measureType": "g",
          "quantity": 500
        }
        ```
*   `GET /ingredient`: Get all ingredients.
*   `PATCH /ingredient/{id}`: Update an ingredient.
    *   **Body:**
        ```json
        {
          "name": "All-purpose Flour",
          "measureType": "g",
          "quantity": 1000
        }
        ```
*   `DELETE /ingredient?id={id}`: Delete an ingredient.

### Recipes

*   `POST /recipe`: Add a new recipe.
    *   **Body:**
        ```json
        {
          "name": "Pancakes",
          "ingredients": [
            {
              "name": "Flour",
              "measureType": "g",
              "quantity": 200
            },
            {
              "name": "Milk",
              "measureType": "ml",
              "quantity": 250
            },
            {
              "name": "Egg",
              "measureType": "unit",
              "quantity": 2
            }
          ]
        }
        ```
*   `GET /recipe`: Get all recipes.
*   `DELETE /recipe?id={id}`: Delete a recipe.

### Recommendations

*   `GET /recommendation`: Get recipe recommendations based on available ingredients.

## Testing

To run the tests, use the following command:

```bash
go test ./...
```
