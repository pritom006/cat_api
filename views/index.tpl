<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Cat Voting App</title>
  <!-- Bootstrap CSS -->
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
  <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons/font/bootstrap-icons.css" rel="stylesheet">
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
  <div class="container mt-4 mx-auto" style="max-width: 800px;">
    <div class="card shadow rounded-3 p-4">
      <!-- Top Navigation Buttons -->
      <div class="d-flex justify-content-center mb-4">
        <button class="btn btn-outline-danger mx-2 d-flex align-items-center" onclick="displaySection('voting')">
          <i class="bi bi-arrow-repeat me-2"></i> Voting
        </button>
        <button class="btn btn-outline-secondary mx-2 d-flex align-items-center" onclick="displaySection('breeds')">
          <i class="bi bi-search me-2"></i> Breeds
        </button>
        <button class="btn btn-outline-primary mx-2 d-flex align-items-center" onclick="displaySection('favorites')">
          <i class="bi bi-heart me-2"></i> Favs
        </button>
      </div>

      <!-- Voting Section -->
      <div id="image-container" class="section text-center">
        <img id="cat-image" src="" alt="Cat" class="img-fluid rounded border shadow" style="max-width: 500px;" />
      </div>

      <!-- Breed Information Section -->
      <div id="breed-info" class="mt-3"></div>

      <!-- Breeds Dropdown -->
      <div class="text-center mt-3">
        <select id="breed-select" class="form-select w-auto mx-auto" style="display: none;" onchange="fetchBreedCats()"></select>
      </div>

      <!-- Image Slider for Breeds Section -->
      <div id="breed-image-slider" class="mt-4" style="display: none;"></div>

      <!-- Bottom Buttons -->
      <div id="voting-icons" class="text-center mt-3">
        <button class="btn btn-light-custom mx-2" onclick="voteCat(1)">
          <i class="bi bi-hand-thumbs-up"></i>
        </button>
        <button class="btn btn-light-custom mx-2" onclick="addToFavorites()">
          <i class="bi bi-heart"></i>
        </button>
        <button class="btn btn-light-custom mx-2" onclick="voteCat(0)">
          <i class="bi bi-hand-thumbs-down"></i>
        </button>
      </div>

      <!-- Favorites Section -->
      <div id="favorites-section" class="section row mt-4 g-3 justify-content-center" style="display: none;">
        <!-- Favorite images will be dynamically added here -->
      </div>
    </div>
  </div>

  <!-- Bootstrap JS and Icons -->
  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js"></script>
  <script src="/static/js/main.js"></script>
</body>
</html>
