function showLoading() {
  document.getElementById("loadingOverlay").style.visibility = "visible";
}

function hideLoading() {
  document.getElementById("loadingOverlay").style.visibility = "hidden";
}

// Modal functions
function openModal() {
  document.getElementById("editModal").style.display = "flex";
}

function closeModal() {
  document.getElementById("editModal").style.display = "none";
}

// Ingredients
async function addIngredient() {
  showLoading();
  const ingredient = {
    name: document.getElementById("ingredientName").value,
    measureType: document.getElementById("measureType").value,
    quantity: parseInt(document.getElementById("quantity").value),
  };

  await fetch("http://localhost:8080/ingredient", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(ingredient),
  });

  // Clear form
  document.getElementById("ingredientName").value = "";
  document.getElementById("quantity").value = "";
  document.getElementById("measureType").value = "";

  await getIngredients();
  hideLoading();
}

async function getIngredients() {
  const response = await fetch("http://localhost:8080/ingredient");
  let ingredients = await response.json();
  if (!ingredients) {
    const list = document.getElementById("ingredientList");
    list.innerHTML = "";
    return;
  }
  const list = document.getElementById("ingredientList");
  list.innerHTML = ingredients
    .map(
      (i) =>
        `<div class="ingredient-item">
                    ${i.Name}
                    <span class="ingredient-badge">${i.Quantity} ${i.MeasureType}</span>
                    <div class="ingredient-actions">
                      <button class="btn-warning" onclick="openEditModal('${i.Id}', '${i.Name}', ${i.Quantity}, '${i.MeasureType}')">Edit</button>
                      <button class="btn-danger" onclick="deleteIngredient('${i.Id}')">Delete</button>
                    </div>
                </div>`,
    )
    .join("");
}

function openEditModal(id, name, quantity, measureType) {
  document.getElementById("editIngredientId").value = id;
  document.getElementById("editIngredientName").value = name;
  document.getElementById("editQuantity").value = quantity;
  document.getElementById("editMeasureType").value = measureType;
  openModal();
}

async function updateIngredient() {
  showLoading();
  const id = document.getElementById("editIngredientId").value;
  const ingredient = {
    name: document.getElementById("editIngredientName").value,
    quantity: parseInt(document.getElementById("editQuantity").value),
    measureType: document.getElementById("editMeasureType").value,
  };

  await fetch(`http://localhost:8080/ingredient/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(ingredient),
  });

  closeModal();
  await getIngredients();
  hideLoading();
}

async function deleteIngredient(id) {
  if (!confirm("Are you sure you want to delete this ingredient?")) {
    return;
  }

  showLoading();
  await fetch(`http://localhost:8080/ingredient?id=${id}`, {
    method: "DELETE",
  });

  await getIngredients();
  hideLoading();
}

let recipeIngredients = [];

function addRecipeIngredient() {
  const newInput = document.createElement("div");
  newInput.className = "input-group";
  newInput.innerHTML = `
                <div class="ingredient-input">
                    <input type="text" class="ingName" placeholder="Ingredient name">
                    <input type="number" class="ingQuantity" placeholder="Quantity">
                    <select id="ingMeasure">
                      <option value="unit">Unit</option>
                      <option value="g">Gram</option>
                      <option value="mg">Milligram</option>
                    </select>
                </div>
            `;
  document.getElementById("recipeIngredients").appendChild(newInput);
}

async function createRecipe() {
  showLoading();
  const ingredients = [];
  const recipeName = document.getElementById("recipeName").value;

  document.querySelectorAll(".ingredient-input").forEach((div) => {
    ingredients.push({
      Name: div.querySelector(".ingName").value,
      Quantity: +div.querySelector(".ingQuantity").value,
      MeasureType: div.querySelector("#ingMeasure").value,
    });
  });

  await fetch("http://localhost:8080/recipe", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name: recipeName, ingredients }),
  });

  document.getElementById("recipeName").value = "";
  document.getElementById("recipeIngredients").innerHTML = `
          <div class="input-group">
            <label>Ingredient</label>
            <div class="ingredient-input">
              <input type="text" class="ingName" placeholder="Ingredient name">
              <input type="number" class="ingQuantity" placeholder="Quantity">
              <select id="ingMeasure">
                <option value="unit">Unit</option>
                <option value="g">Gram</option>
                <option value="mg">Milligram</option>
              </select>
            </div>
          </div>
        `;

  await getRecipes();
  hideLoading();
}

async function getRecipes() {
  const response = await fetch("http://localhost:8080/recipe");
  const recipes = await response.json();

  const container = document.getElementById("recipes");

  container.innerHTML = recipes
    .map(
      (recipe) => `<div class="recipe-card">
                    <h3>${recipe.Name}</h3>
                    <div>Requires:</div>
                    ${recipe.Ingredients.map(
                      (ing) =>
                        `<div class="ingredient-badge">${ing.Name}: ${ing.Quantity} ${ing.MeasureType}</div>`,
                    ).join("")}
                    <button class="btn-danger" onclick="deleteRecipe('${recipe.Id}')">Delete</button>
                </div>
            `,
    )
    .join("");
}

async function deleteRecipe(id) {
  if (!confirm("Are you sure you want to delete this recipe?")) {
    return;
  }

  showLoading();
  await fetch(`http://localhost:8080/recipe?id=${id}`, {
    method: "DELETE",
  });

  await getRecipes();
  hideLoading();
}

async function getRecommendations() {
  showLoading();
  const response = await fetch("http://localhost:8080/recommendation");
  const recipes = await response.json();

  const container = document.getElementById("recommendations");

  container.innerHTML = recipes
    .map(
      (recipe) => `<div class="recipe-card">
                    <h3>${recipe.Recipe.Name}</h3>
                    <div>Requires:</div>
                    ${recipe.Recipe.Ingredients.map(
                      (ing) =>
                        `<div class="ingredient-badge">${ing.Name}: ${ing.Quantity} ${ing.MeasureType}</div>`,
                    ).join("")}
                </div>
            `,
    )
    .join("");

  hideLoading();
}

getIngredients();
getRecipes();
