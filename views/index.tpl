<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Cat Voting App</title>
  <link rel="stylesheet" href="/static/css/style.css">
  <script src="/static/js/reload.min.js"></script>
</head>
<body>
  <div class="container">
    <!-- Top Navigation Buttons -->
    <div class="top-buttons">
      <button onclick="displaySection('voting')">Voting</button>
      <button onclick="displaySection('breeds')">Breeds</button>
      <button onclick="displaySection('favorites')">Favorites</button>
    </div>

    <!-- Voting Section -->
    <div id="image-container" class="section">
      <img id="cat-image" src="" alt="Cat" />
    </div>

    <!-- Bottom Buttons -->
    <div class="bottom-buttons">
      <button onclick="voteCat(1)">👍</button>
      <button onclick="addToFavorites()">❤️</button>
      <button onclick="voteCat(0)">👎</button>
    </div>

    <!-- Breeds Dropdown -->
    <select id="breed-select" class="section" style="display: none;" onchange="fetchBreedCats()"></select>

    <!-- Favorites Section -->
    <div id="favorites-section" class="section" style="display: none; flex-wrap: wrap;">
      <!-- Favorite images will be dynamically added here -->
    </div>
  </div>

  <script src="/static/js/main.js"></script>
</body>
</html>
