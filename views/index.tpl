<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Cat Voting App</title>
  <!-- Bootstrap CSS -->
  <link
    href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css"
    rel="stylesheet"
  >
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
  <div class="container mt-4">
    <!-- Top Navigation Buttons -->
    <div class="d-flex justify-content-center mb-4">
      <button class="btn btn-outline-danger mx-2" onclick="displaySection('voting')">Voting</button>
      <button class="btn btn-outline-secondary mx-2" onclick="displaySection('breeds')">Breeds</button>
      <button class="btn btn-outline-primary mx-2" onclick="displaySection('favorites')">Favs</button>
    </div>

    <!-- Voting Section -->
    <div id="image-container" class="section text-center">
      <img
        id="cat-image"
        src=""
        alt="Cat"
        class="img-fluid rounded border shadow"
        style="max-width: 500px;"
      />
    </div>

    <!-- Breed Information Section -->
    <div id="breed-info" class="mt-3" style="display: block;"></div>

    <!-- Breeds Dropdown -->
    <div class="text-center mt-3">
      <select
        id="breed-select"
        class="form-select w-auto mx-auto"
        style="display: none;"
        onchange="fetchBreedCats()"
      >
      </select>
    </div>

    <!-- Bottom Buttons -->
    <div class="text-center mt-3">
      <button class="btn btn-light mx-2" onclick="voteCat(1)">
        <i class="bi bi-hand-thumbs-up-fill text-success"></i>
      </button>
      <button class="btn btn-light mx-2" onclick="addToFavorites()">
        <i class="bi bi-heart-fill text-danger"></i>
      </button>
      <button class="btn btn-light mx-2" onclick="voteCat(0)">
        <i class="bi bi-hand-thumbs-down-fill text-danger"></i>
      </button>
    </div>

    <!-- Favorites Section -->
    <div
      id="favorites-section"
      class="section row mt-4 g-3 justify-content-center"
      style="display: none;"
    >
      <!-- Favorite images will be dynamically added here -->
    </div>
  </div>

  <!-- Bootstrap JS and Icons -->
  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js"></script>
  <script src="/static/js/main.js"></script>
</body>
</html>
