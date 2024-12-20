
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
    .then(()=> {
      fetchCat(true);
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


// Show the relevant section 
// function displaySection(section) {
//   document.querySelectorAll(".section").forEach((sec) => {
//     sec.style.display = "none"; // Hide all sections
//   });

//   const breedSelect = document.getElementById("breed-select");
//   breedSelect.style.display = "none"; // Hide breed dropdown by default

//   if (section === "voting") {
//     document.getElementById("image-container").style.display = "block";
//     fetchCat(false); // Ensure the same image is shown unless voting
//   } else if (section === "breeds") {
//     fetchBreeds();
//     breedSelect.style.display = "block"; // Show breed dropdown
//     document.getElementById("image-container").style.display = "block";
//   } else if (section === "favorites") {
//     fetchFavorites();
//     document.getElementById("favorites-section").style.display = "flex";
//   }
// }


let currentSection = ""; // Track the currently active section
function displaySection(section) {
  document.querySelectorAll(".section").forEach((sec) => {
    sec.style.display = "none"; // Hide all sections
  });

  const breedSelect = document.getElementById("breed-select");
  breedSelect.style.display = "none"; // Hide breed dropdown by default

  const breedInfoDiv = document.getElementById("breed-info");
  if (section !== "breeds") {
    breedInfoDiv.innerHTML = ""; // Clear the description and Wikipedia URL
  }

  if (section === "voting") {
    document.getElementById("image-container").style.display = "block";

    // Only fetch a new cat if coming from breeds or favorites
    if (currentSection === "breeds" || currentSection === "favorites") {
      fetchCat(true); // Force a new cat image
    } else {
      fetchCat(false); // Retain the same image
    }
  } else if (section === "breeds") {
    fetchBreeds();
    breedSelect.style.display = "block"; // Show breed dropdown
    document.getElementById("image-container").style.display = "block";
  } else if (section === "favorites") {
    fetchFavorites();
    document.getElementById("favorites-section").style.display = "flex";
  }

  // Update the current section
  currentSection = section;
}
// Show the relevant section
// function displaySection(section) {
//   document.querySelectorAll(".section").forEach((sec) => {
//     sec.style.display = "none"; // Hide all sections
//   });

//   const breedSelect = document.getElementById("breed-select");
//   breedSelect.style.display = "none"; // Hide breed dropdown by default

//   if (section === "voting") {
//     document.getElementById("image-container").style.display = "block";

//     // Only fetch a new cat if coming from breeds or favorites
//     if (currentSection === "breeds" || currentSection === "favorites") {
//       fetchCat(true); // Force a new cat image
//     } else {
//       fetchCat(false); // Retain the same image
//     }
//   } else if (section === "breeds") {
//     fetchBreeds();
//     breedSelect.style.display = "block"; // Show breed dropdown
//     document.getElementById("image-container").style.display = "block";
//   } else if (section === "favorites") {
//     fetchFavorites();
//     document.getElementById("favorites-section").style.display = "flex";
//   }

//   // Update the current section
//   currentSection = section;
// }



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

      window.breedData = data

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
  const breedInfoDiv = document.getElementById("breed-info");

  if (breed !== "Select Breed") {
    fetch(`https://api.thecatapi.com/v1/images/search?breed_ids=${breed}`)
      .then((res) => res.json())
      .then((data) => {
        const img = document.getElementById("cat-image");
        img.src = data[0].url;
        currentCatId = data[0].id;
      });

      // Display breed Description and Wikipedia url
      const breeds = window.breedData.find((b) => b.id === breed);
      if (breeds) {
        breedInfoDiv.innerHTML = `
        <p><strong>Description:</strong> ${breeds.description}</p>
        <p><strong>Learn more:</strong> <a href="${breeds.wikipedia_url}" target="_blank">Wikipedia</a></p>
      `;
      } else {
        breedInfoDiv.innerHTML = "<p>No additional information available for this breed.</p>";
      } 
  } else {
    breedInfoDiv.innerHTML = "";
  }
}

// Initialize the application with the first cat image
window.onload = () => fetchCat(true);
