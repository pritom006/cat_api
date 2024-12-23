// Variable to track the current cat image ID
let currentCatId = "";
let currentImageURL = ""; // To ensure the same image is retained in voting

// Fetch a random cat image and set it as the current image
function fetchCat(initial = false) {
  // Check if the image should be fetched or reused
  if (initial || !currentImageURL) {
    // Call the backend API to fetch a new cat image
    fetch("/fetch-new-cat", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((res) => {
        // Check if the response is OK
        if (!res.ok) {
          throw new Error("Failed to fetch a new cat image.");
        }
        return res.json();
      })
      .then((data) => {
        if (data && data.url && data.id) {
          // Update the image element and global variables
          const img = document.getElementById("cat-image");
          img.src = data.url;
          currentImageURL = data.url;
          currentCatId = data.id;
        } else {
          throw new Error("Invalid data format received from API.");
        }
      })
      .catch((error) => {
        console.error("Error fetching new cat image:", error);
        alert("An error occurred while fetching the cat image. Please try again.");
      });
  } else {
    // Reuse the existing image
    const img = document.getElementById("cat-image");
    img.src = currentImageURL;
  }
}

// Handle voting (like/dislike) for the current cat image
function voteCat(value) {
  fetch("/vote", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ image_id: currentCatId, value }),
  }).then(() => {
    fetchCat(true); // Load a new cat image after voting
  });
}

// Handle AddFavorite functionality with subID
function addToFavorites() {
  console.log("Current Cat ID:", currentCatId)
  fetch("/addToFavorites", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ image_id: currentCatId, sub_id: "test" }), // Ensure this matches backend expectations
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.message) {
        alert(data.message);
        setTimeout(fetchFavorites, 500); // Delay to allow backend processing
      } else {
        alert(data.error || "Failed to add to favorites. Try again.");
      }
    })
    .then(() => {
      fetchCat(true);
    })
    .catch((error) => {
      console.error("Error adding to favorites:", error);
    });
}

// Fetch favorites for the "favorites" section
// Utility function to toggle between flex and grid views
function toggleView(viewClass) {
  const imagesContainer = document.getElementById("images-container");

  if (viewClass === "grid-view") {
    imagesContainer.classList.remove("flex-view");
    imagesContainer.classList.add("grid-view");
  } else {
    imagesContainer.classList.remove("grid-view");
    imagesContainer.classList.add("flex-view");
  }
}

// Fetch favorites for the "favorites" section
function fetchFavorites() {
  fetch("/favorites", {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((res) => res.json())
    .then((data) => {
      const container = document.getElementById("favorites-section");
      container.innerHTML = ""; // Clear previous favorites

      // Add toggle buttons for views
      const buttonContainer = document.createElement("div");
      buttonContainer.style.marginBottom = "10px";

      const flexButton = document.createElement("button");
      flexButton.textContent = "Flex View";
      flexButton.style.marginRight = "10px";
      flexButton.onclick = () => toggleView("flex-view");

      const gridButton = document.createElement("button");
      gridButton.textContent = "Grid View";
      gridButton.onclick = () => toggleView("grid-view");

      buttonContainer.appendChild(flexButton);
      buttonContainer.appendChild(gridButton);
      container.appendChild(buttonContainer);

      // Add favorite images
      if (data.length === 0) {
        container.innerHTML += "<p>No favorites yet!</p>";
        return;
      }

      const imagesContainer = document.createElement("div");
      imagesContainer.id = "images-container"; // Wrapper for images
      imagesContainer.className = "flex-view"; // Default to flex view

      data.forEach((item) => {
        const img = document.createElement("img");
        img.src = item.image.url;
        img.alt = "Favorite Cat";
        imagesContainer.appendChild(img);
      });

      container.appendChild(imagesContainer);
    })
    .catch((error) => {
      console.error("Error fetching favorites:", error);
    });
}

// Track the current section being displayed
let currentSection = "";

function displaySection(section) {
  document.querySelectorAll(".section").forEach((sec) => {
    sec.style.display = "none"; // Hide all sections
  });

  const breedSelect = document.getElementById("breed-select");
  const breedInfoDiv = document.getElementById("breed-info");
  const sliderContainer = document.getElementById("breed-image-slider");

  // Clear breed data when leaving the "breeds" section
  if (currentSection === "breeds" && section !== "breeds") {
    breedSelect.style.display = "none";
    breedSelect.innerHTML = ""; // Clear dropdown
    breedInfoDiv.innerHTML = ""; // Clear breed info
    sliderContainer.style.display = "none"; // Hide slider
    sliderContainer.innerHTML = ""; // Clear slider content
  }

  const votingIcons = document.getElementById("voting-icons");

  if (section === "breeds") {
    votingIcons.style.display = "none"; // Hide voting icons in breed section
    fetchBreeds(); // Fetch breeds when entering the "breeds" section
    breedSelect.style.display = "block"; // Show breed dropdown
    document.getElementById("image-container").style.display = "none"; // Hide voting image container
  } else if (section === "voting") {
    votingIcons.style.display = "block"; // Show voting icons
    document.getElementById("image-container").style.display = "block";

    // Only fetch a new cat if coming from breeds or favorites
    if (currentSection === "breeds" || currentSection === "favorites") {
      fetchCat(true); // Force a new cat image
    } else {
      fetchCat(false); // Retain the same image
    }
  } else if (section === "favorites") {
    //votingIcons.style.display = none; // Show voting icons
    fetchFavorites();
    document.getElementById("favorites-section").style.display = "flex";
  }

  // Update the current section
  currentSection = section;
}




// Fetch all available breeds and populate the dropdown
function fetchBreeds() {
  fetch("/fetch-breeds", {
    method: "GET",
    headers: {
      "Content-Type": "application/json"
    },
  })
    .then((res) => res.json())
    .then((data) => {
      const select = document.getElementById("breed-select");
      select.style.display = "block";
      select.innerHTML = "<option>Select Breed</option>";

      window.breedData = data

      data.forEach((breed) => {
        const option = document.createElement("option");
        option.value = breed.id;
        option.textContent = breed.name;
        select.appendChild(option);
      });

      if (data.length > 0) {
        select.value = data[0].id;
        fetchBreedCats();
      }
    });
}

// Fetch cat images for the selected breed
function fetchBreedCats() {
  const breed = document.getElementById("breed-select").value;
  const breedInfoDiv = document.getElementById("breed-info");
  const sliderContainer = document.getElementById("breed-image-slider");

  if (breed !== "Select Breed") {
    const breedData = window.breedData.find((b) => b.id === breed);
    breedInfoDiv.innerHTML = `
      <p><strong>${breedData.name}</strong> (${breedData.origin})</p>
      <p>${breedData.description}</p>
      <p><a href="${breedData.wikipedia_url}" target="_blank">Wikipedia</a></p>
    `;

    fetch(`/fetch-breed-images?breed_id=${breed}`)
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          console.error(data.error);
          return;
        }
        sliderContainer.innerHTML = `
          <div class="carousel-container">
            <div id="carousel-${breed}" class="carousel-wrapper">
              ${data
                .map(
                  (img) =>
                    `<div class="carousel-img"><img src="${img.url}" alt="Breed Image"></div>`
                )
                .join("")}
            </div>
            <div id="dots-${breed}" class="dots-container">
              ${data
                .map((_, i) => `<div class="dot ${i === 0 ? "active" : ""}" data-index="${i}"></div>`)
                .join("")}
            </div>
          </div>
        `;
        initializeCarouselForBreed(breed, data.length);
        sliderContainer.style.display = "block";
      });
  } else {
    sliderContainer.style.display = "none";
    breedInfoDiv.innerHTML = "";
  }
}

function initializeCarouselForBreed(breedId, imageCount) {
  const carousel = document.getElementById(`carousel-${breedId}`);
  const dots = document.querySelectorAll(`#dots-${breedId} .dot`);
  let currentIndex = 0;


  function updateCarousel(index) {
    const carousel = document.getElementById(`carousel-${breedId}`);
    carousel.style.transform = `translateX(-${index * 100}%)`; // Horizontal slide
    dots.forEach((dot, i) => dot.classList.toggle("active", i === index));
  }

  let interval = setInterval(() => {
    currentIndex = (currentIndex + 1) % imageCount;
    updateCarousel(currentIndex);
  }, 3000);

  dots.forEach((dot, index) => {
    dot.addEventListener("click", () => {
      clearInterval(interval);
      currentIndex = index;
      updateCarousel(currentIndex);
      interval = setInterval(() => {
        currentIndex = (currentIndex + 1) % imageCount;
        updateCarousel(currentIndex);
      }, 3000);
    });
  });

  updateCarousel(0);
}


// Initialize the application with the first cat image
window.onload = () => fetchCat(true);
