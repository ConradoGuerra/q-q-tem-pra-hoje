-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS recipes_ingredients;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS ingredients_storage;