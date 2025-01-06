let currentCatId = "";
let currentImageURL = ""; 


function toggleSpinner(show) {
  const spinner = document.getElementById("spinner");
  spinner.style.display = show ? "block" : "none";
}

// Fetch a random cat image and set it as the current image
function fetchCat(initial = false) {
  
  // Check if the image should be fetched or reused
  if (initial || !currentImageURL) {
    toggleSpinner(true);
    fetch("/fetch-new-cat", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((res) => {
        if (!res.ok) {
          throw new Error("Failed to fetch a new cat image.");
        }
        return res.json();
      })
      .then((data) => {
        if (data && data.url && data.id) {
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
      })
      .finally(() => {
        toggleSpinner(false); 
      });
  } else {
    const img = document.getElementById("cat-image");
    img.src = currentImageURL;
  }
}

// Handle voting (like/dislike) for the current cat image
function voteCat(value) {
  toggleSpinner(true)
  fetch("/vote", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ image_id: currentCatId, value }),
  }).then(() => {
    fetchCat(true); 
  })
  .catch((error) => {
    console.error("Error voting on cat:", error);
  })
  .finally(() => {
    toggleSpinner(false); 
  });
}



function addToFavorites() {
  toggleSpinner(true);
  console.log("Adding to favorites, Cat ID:", currentCatId);
  
  fetch("/addToFavorites", {
      method: "POST",
      headers: {
          "Content-Type": "application/json",
      },
      body: JSON.stringify({ 
          image_id: currentCatId,
          sub_id: "test"
      }),
  })
  .then(response => response.json())
  .then(data => {
      console.log("Add to favorites response:", data);
      if (data.message) {
          alert(data.message);
          fetchFavorites(); // Refresh favorites list
      } else if (data.error) {
          throw new Error(data.error);
      }
  })
  .then(() => {
      fetchCat(true);
  })
  .catch(error => {
      console.error("Error adding to favorites:", error);
      alert("Failed to add to favorites: " + error.message);
  })
  .finally(() => {
      toggleSpinner(false);
  });
}




function fetchFavorites() {
  toggleSpinner(true);
  console.log("Fetching favorites...");

  fetch("/favorites", {
      method: "GET",
      headers: {
          "Content-Type": "application/json"
      },
  })
  .then(response => response.json())
  .then(data => {
      console.log("Favorites data:", data);
      const container = document.getElementById("favorites-section");
      const imagesContainer = document.getElementById("images-container");
      
      // Clear existing images
      imagesContainer.innerHTML = "";

      // Handle error response
      if (data.error) {
          throw new Error(data.error);
      }

      // Handle empty response
      if (!Array.isArray(data) || data.length === 0) {
          imagesContainer.innerHTML = "<p>No favorites found</p>";
          container.style.display = "block";
          return;
      }

      // Process each favorite
      data.forEach(item => {
          if (!item.image || !item.image.url) {
              console.warn("Invalid favorite item:", item);
              return;
          }

          const imgWrapper = document.createElement("div");
          imgWrapper.className = "image-wrapper";

          const img = document.createElement("img");
          img.src = item.image.url;
          img.alt = "Favorite Cat";
          img.className = "favorite-image";
          
          // Error handling for images
          img.onerror = () => {
              img.src = "placeholder-image-url.jpg"; // Add a placeholder image URL
              img.alt = "Failed to load image";
          };

          imgWrapper.appendChild(img);
          imagesContainer.appendChild(imgWrapper);
      });

      container.style.display = "block";
  })
  .catch(error => {
      console.error("Error fetching favorites:", error);
      const imagesContainer = document.getElementById("images-container");
      imagesContainer.innerHTML = `<p>Error loading favorites: ${error.message}</p>`;
  })
  .finally(() => {
      toggleSpinner(false);
  });
}



function toggleView(view) {
  const imagesContainer = document.getElementById("images-container");
  const flexButton = document.getElementById("flex-view-button");
  const gridButton = document.getElementById("grid-view-button");

  console.log("Toggle View Called:", view); 
  if (view === "flex-view") {
    imagesContainer.className = "flex-view"; 
    flexButton.classList.add("active");
    gridButton.classList.remove("active");
  } else if (view === "grid-view") {
    imagesContainer.className = "grid-view"; 
    gridButton.classList.add("active");
    flexButton.classList.remove("active");
  }
}

// Attach event listeners to toggle buttons
document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("flex-view-button").addEventListener("click", () => toggleView("flex-view"));
  document.getElementById("grid-view-button").addEventListener("click", () => toggleView("grid-view"));
  toggleView("flex-view");
});



// Track the current section being displayed

let currentSection = "";

function displaySection(section) {
  document.querySelectorAll(".section").forEach((sec) => {
    sec.style.display = "none"; 
  });

  const breedSelect = document.getElementById("breed-select");
  const breedInfoDiv = document.getElementById("breed-info");
  const sliderContainer = document.getElementById("breed-image-slider");
  const favoritesSection = document.getElementById("favorites-section");
  const votingIcons = document.getElementById("voting-icons");

  // Clear breed data when leaving the "breeds" section
  if (currentSection === "breeds" && section !== "breeds") {
    breedSelect.style.display = "none";
    breedSelect.innerHTML = ""; 
    breedInfoDiv.innerHTML = ""; 
    sliderContainer.style.display = "none"; 
    sliderContainer.innerHTML = ""; 
  }

  if (section === "breeds") {
    votingIcons.style.display = "none"; 
    fetchBreeds(); 
    breedSelect.style.display = "block"; 
    document.getElementById("image-container").style.display = "none"; 
  } else if (section === "voting") {
    votingIcons.style.display = "block"; 
    document.getElementById("image-container").style.display = "block";
    favoritesSection.style.display = "none"; 

    if (currentSection === "breeds" || currentSection === "favorites") {
      fetchCat(true); 
    } else {
      fetchCat(false); 
    }
  } else if (section === "favorites") {
    votingIcons.style.display = "none"; 
    fetchFavorites();
    favoritesSection.style.display = "flex"; 
  }

  // Update the current section
  currentSection = section;
}



// Fetch all available breeds and populate the dropdown
function fetchBreeds() {
  toggleSpinner(true);
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
    })
    .catch((error) => {
      console.error("Error fetching breeds:", error);
      alert("An error occurred while fetching breeds. Please try again.");
    })
    .finally(() => {
      toggleSpinner(false); 
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

    toggleSpinner(true);
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
      })
      .catch((error) => {
        console.error("Error fetching breed images:", error);
        alert("An error occurred while fetching breed images. Please try again.");
      })
      .finally(() => {
        toggleSpinner(false); 
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
