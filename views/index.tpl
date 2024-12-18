<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Cat Voting App</title>
  <link rel="stylesheet" href="/static/css/style.css">
  <script src="/static/js/reload.min.js"></script>
</head>
<body>
  <div class="container">
    <div class="top-buttons">
      <button onclick="showSection('voting')">Voting</button>
      <button onclick="showSection('breeds')">Breeds</button>
      <button onclick="fetchFavorites()">Favs</button>
    </div>
    <div id="image-container">
      <img id="cat-image" src="" alt="Cat" />
    </div>
    <div class="bottom-buttons">
      <button onclick="voteCat(1)">👍</button>
      <button onclick="addToFavorites()">❤️</button>
      <button onclick="voteCat(0)">👎</button>
    </div>
    <select id="breed-select" style="display:none" onchange="fetchBreedCats()"></select>
    <div id="favorites-section"></div>
  </div>

  <script src="/static/js/main.js"></script>
</body>
</html>
