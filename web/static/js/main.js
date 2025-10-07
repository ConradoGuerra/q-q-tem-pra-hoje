function showLoading() {
  document.getElementById("loadingOverlay").style.visibility = "visible";
}

function hideLoading() {
  document.getElementById("loadingOverlay").style.visibility = "hidden";
}

function openModal() {
  document.getElementById("editModal").style.display = "flex";
}

function closeModal() {
  document.getElementById("editModal").style.display = "none";
}

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
let allRecipes = [];
let currentRecipePage = 1;
const recipesPerPage = 10;
let allRecommendations = [];
let currentRecPage = 1;
const recsPerPage = 10;

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

  currentRecipePage = 1; // Reset to first page
  await getRecipes();
  hideLoading();
}

async function getRecipes() {
  const response = await fetch("http://localhost:8080/recipe");
  allRecipes = await response.json();
  displayRecipes(currentRecipePage);
}

function displayRecipes(page) {
  const container = document.getElementById("recipes");
  const start = (page - 1) * recipesPerPage;
  const end = start + recipesPerPage;
  const pageRecipes = allRecipes.slice(start, end);

  container.innerHTML = pageRecipes
    .map(
      (recipe) => `<div class="recipe-card">
                    <h3>${recipe.Name}</h3>
                    ${recipe.Ingredients.length ? "<div>Requires:</div><div class=\"badges-container\">" : ""}
                    ${recipe.Ingredients.map(
                      (ing) =>
                        `<div class="ingredient-badge">${ing.Name}: ${ing.Quantity} ${ing.MeasureType}</div>`,
                    ).join("")}
                    ${recipe.Ingredients.length ? "</div>" : ""}
                    <button class="btn-danger" onclick="deleteRecipe('${recipe.Id}')">Delete</button>
                </div>
            `,
    )
    .join("");

  // Pagination
  const totalPages = Math.ceil(allRecipes.length / recipesPerPage);
  let pagination = '<div class="pagination">';
  if (page > 1) pagination += `<button onclick="changeRecipePage(${page - 1})">Prev</button>`;
  if (page < totalPages) pagination += `<button onclick="changeRecipePage(${page + 1})">Next</button>`;
  pagination += `</div>`;
  container.innerHTML += pagination;
}

function changeRecipePage(page) {
  currentRecipePage = page;
  displayRecipes(page);
}

async function deleteRecipe(id) {
  if (!confirm("Are you sure you want to delete this recipe?")) {
    return;
  }

  showLoading();
  await fetch(`http://localhost:8080/recipe?id=${id}`, {
    method: "DELETE",
  });

  currentRecipePage = 1; // Reset to first page
  await getRecipes();
  hideLoading();
}

async function getRecommendations() {
  showLoading();
  const response = await fetch("http://localhost:8080/recommendation");
  allRecommendations = await response.json();
  currentRecPage = 1; // Reset to first page
  displayRecommendations(currentRecPage);
  hideLoading();
}

function displayRecommendations(page) {
  const container = document.getElementById("recommendations");
  const start = (page - 1) * recsPerPage;
  const end = start + recsPerPage;
  const pageRecs = allRecommendations.slice(start, end);

  container.innerHTML = pageRecs
    .map(
      (recipe) => `<div class="recipe-card">
                    <h3>${recipe.Recipe.Name}</h3>
                    <div>Requires:</div>
                    <div class="badges-container">
                    ${recipe.Recipe.Ingredients.map(
                      (ing) =>
                        `<div class="ingredient-badge">${ing.Name}: ${ing.Quantity} ${ing.MeasureType}</div>`,
                    ).join("")}
                    </div>
                </div>
            `,
    )
    .join("");

  // Pagination
  const totalPages = Math.ceil(allRecommendations.length / recsPerPage);
  let pagination = '<div class="pagination">';
  if (page > 1) pagination += `<button onclick="changeRecPage(${page - 1})">Prev</button>`;
  if (page < totalPages) pagination += `<button onclick="changeRecPage(${page + 1})">Next</button>`;
  pagination += `</div>`;
  container.innerHTML += pagination;
}

function changeRecPage(page) {
  currentRecPage = page;
  displayRecommendations(page);
}

// Tab switching
document.querySelectorAll('.tab').forEach(tab => {
  tab.addEventListener('click', () => {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
    tab.classList.add('active');
    document.getElementById(tab.dataset.tab + '-tab').classList.add('active');
  });
});

// Event listeners
document.getElementById('addIngredientBtn').addEventListener('click', addIngredient);
document.getElementById('createRecipeBtn').addEventListener('click', createRecipe);
document.getElementById('getRecommendationsBtn').addEventListener('click', getRecommendations);
document.getElementById('addRecipeIngredientBtn').addEventListener('click', addRecipeIngredient);
document.getElementById('updateIngredientBtn').addEventListener('click', updateIngredient);
document.getElementById('closeModalBtn').addEventListener('click', closeModal);

getIngredients();
getRecipes();
