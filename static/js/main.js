// let currentCatId = "";

// // Fetch a random cat image and set it as the current image
// function fetchCat(initial = false) {
//   if (initial) {
//     fetch("https://api.thecatapi.com/v1/images/search")
//       .then((res) => res.json())
//       .then((data) => {
//         const img = document.getElementById("cat-image");
//         img.src = data[0].url;
//         currentCatId = data[0].id;
//       });
//   }
// }

// // Handle voting (like/dislike) for the current cat image
// function voteCat(value) {
//   fetch("/vote", {
//     method: "POST",
//     headers: { "Content-Type": "application/json" },
//     body: JSON.stringify({ image_id: currentCatId, value }),
//   }).then(() => {
//     fetchCat(true); // Load a new cat image after voting
//   });
// }

// // Add the current cat image to favorites
// function addToFavorites() {
//   fetch("https://api.thecatapi.com/v1/favourites", {
//     method: "POST",
//     headers: {
//       "Content-Type": "application/json",
//       "x-api-key": "YOUR_CAT_API_KEY",
//     },
//     body: JSON.stringify({ image_id: currentCatId }),
//   })
//     .then((response) => response.json())
//     .then(() => {
//       alert("Added to favorites!");
//     })
//     .catch((error) => {
//       console.error("Error adding to favorites:", error);
//     });
// }

// // Fetch all favorite cat images
// function fetchFavorites() {
//   fetch("/favorites")
//     .then((res) => res.json())
//     .then((data) => {
//       const container = document.getElementById("favorites-section");
//       container.innerHTML = "";
//       container.style.display = "flex";
//       container.style.flexWrap = "wrap";
//       container.style.overflowY = "scroll";
//       data.forEach((item) => {
//         const img = document.createElement("img");
//         img.src = item.image.url;
//         img.alt = "Favorite Cat";
//         img.style.width = "300px";
//         img.style.height = "auto";
//         img.style.margin = "10px";
//         container.appendChild(img);
//       });
//     })
//     .catch((error) => {
//       console.error("Error fetching favorites:", error);
//     });
// }

// // Show the relevant section
// // function showSection(section) {
// //   document.querySelectorAll(".section").forEach((sec) => {
// //     sec.style.display = "none";
// //   });
// //   document.getElementById(`${section}-section`).style.display = "block";
// // }

// let lastSection = "home"; // Track the last visited section

// // Show a specific section of the application (voting, breeds, favorites)
// function showSection(section) {
//   const breedSelect = document.getElementById("breed-select");
//   breedSelect.style.display = "none"; // Hide the breed select dropdown by default

//   if (section === "voting") {
//     if (lastSection === "breeds" || lastSection === "favorites") {
//       fetchCat(true); // Refresh the cat image only if coming from Breeds or Favorites
//     }
//   } else if (section === "breeds") {
//     fetchBreeds();
//   } else if (section === "favorites") {
//     fetchFavorites();
//   }

//   lastSection = section; // Update the last visited section
// }

// // Fetch all cat breeds and populate the dropdown
// function fetchBreeds() {
//   fetch("/fetch-breeds")
//     .then((res) => res.json())
//     .then((data) => {
//       const select = document.getElementById("breed-select");
//       select.style.display = "block";
//       select.innerHTML = "<option>Select Breed</option>";
//       data.forEach((breed) => {
//         const option = document.createElement("option");
//         option.value = breed.id;
//         option.textContent = breed.name;
//         select.appendChild(option);
//       });
//     });
// }

// // Fetch a cat image based on the selected breed
// function fetchBreedCats() {
//   const breed = document.getElementById("breed-select").value;
//   if (breed !== "Select Breed") {
//     fetch(`https://api.thecatapi.com/v1/images/search?breed_ids=${breed}`)
//       .then((res) => res.json())
//       .then((data) => {
//         const img = document.getElementById("cat-image");
//         img.src = data[0].url;
//         currentCatId = data[0].id;
//       });
//   }
// }

// // Initialize the application with the first cat image
// window.onload = () => fetchCat(true);


// Variable to track the current cat image ID
let currentCatId = "";
let currentImageURL = ""; // To ensure the same image is retained in voting

// Fetch a random cat image and set it as the current image
function fetchCat(initial = false) {
  if (initial || !currentImageURL) {
    fetch("https://api.thecatapi.com/v1/images/search")
      .then((res) => res.json())
      .then((data) => {
        const img = document.getElementById("cat-image");
        img.src = data[0].url;
        currentImageURL = data[0].url;
        currentCatId = data[0].id;
      });
  } else {
    // If currentImageURL exists, just set it without fetching again
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

// Add the current cat image to favorites
// function addToFavorites() {
//   fetch("/addToFavorites", { // Change this to call your backend endpoint that will handle adding to favorites
//     method: "POST",
//     headers: {
//       "Content-Type": "application/json",
//     },
//     body: JSON.stringify({ image_id: currentCatId }),
//   })
//     .then((response) => {
//       if (response.ok) {
//         alert("Added to favorites!");
//       } else {
//         alert("Failed to add to favorites. Try again.");
//       }
//     })
//     .catch((error) => {
//       console.error("Error adding to favorites:", error);
//     });
// }



// Fetch all favorite cat images and display them
// function fetchFavorites() {
//   fetch("/favorites", {
//     method: "GET",
//     headers: {
//       "Content-Type": "application/json" 
//       //"x-api-key": "live_72VAUewe3v6m8xWYSsEIKTJKBnsFNO4ic2up2J4BCigAz6DBFZ2rIb9XF8E1j8H5",
//     },
//   })
//     .then((res) => res.json())
//     .then((data) => {
//       const container = document.getElementById("favorites-section");
//       container.innerHTML = ""; // Clear previous favorites

//       if (data.length === 0) {
//         container.innerHTML = "<p>No favorites yet!</p>";
//         return;
//       }

//       data.forEach((item) => {
//         const img = document.createElement("img");
//         img.src = item.image.url;
//         img.alt = "Favorite Cat";
//         img.style.width = "300px";
//         img.style.height = "auto";
//         img.style.margin = "10px";
//         container.appendChild(img);
//       });
//     })
//     .catch((error) => {
//       console.error("Error fetching favorites:", error);
//     });
// }


function addToFavorites() {
  console.log(currentCatId)
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
        alert(data.message); // Success message
        setTimeout(fetchFavorites, 500); // Delay to allow backend processing
      } else {
        alert(data.error || "Failed to add to favorites. Try again.");
      }
    })
    .catch((error) => {
      console.error("Error adding to favorites:", error);
    });
}


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

      if (data.length === 0) {
        container.innerHTML = "<p>No favorites yet!</p>";
        return;
      }

      data.forEach((item) => {
        const img = document.createElement("img");
        img.src = item.image.url;
        img.alt = "Favorite Cat";
        img.style.width = "300px";
        img.style.height = "auto";
        img.style.margin = "10px";
        container.appendChild(img);
      });
    })
    .catch((error) => {
      console.error("Error fetching favorites:", error);
    });
}


// Show the relevant section (renamed to displaySection)
function displaySection(section) {
  document.querySelectorAll(".section").forEach((sec) => {
    sec.style.display = "none"; // Hide all sections
  });

  const breedSelect = document.getElementById("breed-select");
  breedSelect.style.display = "none"; // Hide breed dropdown by default

  if (section === "voting") {
    document.getElementById("image-container").style.display = "block";
    fetchCat(false); // Ensure the same image is shown unless voting
  } else if (section === "breeds") {
    fetchBreeds();
    breedSelect.style.display = "block"; // Show breed dropdown
    document.getElementById("image-container").style.display = "block";
  } else if (section === "favorites") {
    fetchFavorites();
    document.getElementById("favorites-section").style.display = "flex";
  }
}

// Fetch all available breeds and populate the dropdown
function fetchBreeds() {
  fetch("/fetch-breeds", {
    method: "GET",
    headers: {
      "Content-Type": "application/json" 
      //"x-api-key": "live_72VAUewe3v6m8xWYSsEIKTJKBnsFNO4ic2up2J4BCigAz6DBFZ2rIb9XF8E1j8H5",
    },
  })
    .then((res) => res.json())
    .then((data) => {
      const select = document.getElementById("breed-select");
      select.style.display = "block";
      select.innerHTML = "<option>Select Breed</option>";
      data.forEach((breed) => {
        const option = document.createElement("option");
        option.value = breed.id;
        option.textContent = breed.name;
        select.appendChild(option);
      });
    });
}

// Fetch cat images for the selected breed
function fetchBreedCats() {
  const breed = document.getElementById("breed-select").value;
  if (breed !== "Select Breed") {
    fetch(`https://api.thecatapi.com/v1/images/search?breed_ids=${breed}`)
      .then((res) => res.json())
      .then((data) => {
        const img = document.getElementById("cat-image");
        img.src = data[0].url;
        currentCatId = data[0].id;
      });
  }
}

// Initialize the application with the first cat image
window.onload = () => fetchCat(true);
