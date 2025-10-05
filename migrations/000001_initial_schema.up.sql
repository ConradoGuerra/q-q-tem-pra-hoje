-- Initial schema for ingredients, recipes, and recipes_ingredients tables
CREATE TABLE IF NOT EXISTS ingredients_storage (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    measure_type TEXT NOT NULL,
    quantity INT NOT NULL
);

CREATE TABLE IF NOT EXISTS recipes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS recipes_ingredients (
    recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    measure_type TEXT NOT NULL,
    quantity INT NOT NULL,
    PRIMARY KEY (recipe_id, name)
);
