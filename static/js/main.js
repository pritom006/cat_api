let currentCatId = "";

function fetchCat() {
  fetch("https://api.thecatapi.com/v1/images/search")
    .then(res => res.json())
    .then(data => {
      const img = document.getElementById("cat-image");
      img.src = data[0].url;
      currentCatId = data[0].id;
    });
}

function voteCat(value) {
  fetch("/vote", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ image_id: currentCatId, value })
  }).then(() => alert("Vote submitted!"));
}

function addToFavorites() {
  fetch("https://api.thecatapi.com/v1/favourites", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "x-api-key": "YOUR_CAT_API_KEY"
    },
    body: JSON.stringify({ image_id: currentCatId })
  }).then(() => alert("Added to favorites!"));
}

function fetchFavorites() {
  fetch("/favorites")
    .then(res => res.json())
    .then(data => {
      const container = document.getElementById("favorites-section");
      container.innerHTML = "";
      data.forEach(item => {
        const img = document.createElement("img");
        img.src = item.image.url;
        img.width = 100;
        container.appendChild(img);
      });
    });
}

function showSection(section) {
  if (section === "voting") fetchCat();
  if (section === "breeds") fetchBreeds();
}

function fetchBreeds() {
  fetch("/fetch-breeds")
    .then(res => res.json())
    .then(data => {
      const select = document.getElementById("breed-select");
      select.style.display = "block";
      select.innerHTML = "<option>Select Breed</option>";
      data.forEach(breed => {
        const option = document.createElement("option");
        option.value = breed.id;
        option.textContent = breed.name;
        select.appendChild(option);
      });
    });
}

function fetchBreedCats() {
  const breed = document.getElementById("breed-select").value;
  fetch(`https://api.thecatapi.com/v1/images/search?breed_ids=${breed}`)
    .then(res => res.json())
    .then(data => {
      document.getElementById("cat-image").src = data[0].url;
      currentCatId = data[0].id;
    });
}

window.onload = fetchCat;
