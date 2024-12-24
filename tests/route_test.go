package tests

import (
	"catapigo/controllers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	// Initialize Beego Router
	beego.Router("/", &controllers.MainController{}, "get:ServeFrontend")
	beego.Router("/fetch-breeds", &controllers.MainController{}, "get:FetchCatBreeds")
	beego.Router("/vote", &controllers.MainController{}, "post:VoteForCat")
	beego.Router("/favorites", &controllers.MainController{}, "get:FetchFavorites")
	beego.Router("/addToFavorites", &controllers.MainController{}, "post:AddToFavorites")
	beego.Router("/fetch-new-cat", &controllers.MainController{}, "get:FetchNewCatImage")
	beego.Router("/fetch-breed-images", &controllers.MainController{}, "get:FetchBreedImages")

	// Test ServeFrontend
	t.Run("Test ServeFrontend", func(t *testing.T) {
		// Step 1: Create a new HTTP request for the `/` route (the endpoint for ServeFrontend)
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err) // Ensure no error while creating the request

		// Step 2: Create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// Step 3: Serve the HTTP request using Beego's handler
		web.BeeApp.Handlers.ServeHTTP(rr, req)

		// Step 4: Assert the status code is 200 OK
		assert.Equal(t, http.StatusOK, rr.Code)

		// Step 5: Assert the response body contains content from the `index.tpl`
		// For example, check if a generic HTML tag or text from the template exists
		// Assuming index.tpl contains some basic HTML, such as <html> or <body>
		assert.Contains(t, rr.Body.String(), "<html>") // Adjust if necessary

		// Alternatively, you can check for other expected content in the response.
		// For example, if you know your template has a specific string like "Welcome" or "Cat API", you can verify that:
		// assert.Contains(t, rr.Body.String(), "Welcome to the Cat API")
	})

	// Test FetchCatBreeds
	t.Run("Test FetchCatBreeds", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/fetch-breeds", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// You can further check the response body if needed (it should contain breed data)
	})

	// Test VoteForCat
	t.Run("Test VoteForCat", func(t *testing.T) {
		// Assuming Vote has the following structure: {"image_id": "imageId", "value": 1}
		vote := `{"image_id": "1234", "value": 1}`
		req, err := http.NewRequest("POST", "/vote", strings.NewReader(vote))
		req.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Check if the response is successful, you can also assert response body for success message
	})

	// Test FetchFavorites
	t.Run("Test FetchFavorites", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/favorites", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Check response body contains the favorites
	})

	// Test AddToFavorites
	t.Run("Test AddToFavorites", func(t *testing.T) {
		favorite := `{"image_id": "1234"}`
		req, err := http.NewRequest("POST", "/addToFavorites", strings.NewReader(favorite))
		req.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Check if the response is successful, you can also assert response body for success message
	})

	// Test FetchNewCatImage
	t.Run("Test FetchNewCatImage", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/fetch-new-cat", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Check if the response body contains a new cat image
	})

	// Test FetchBreedImages
	t.Run("Test FetchBreedImages", func(t *testing.T) {
		// Example breed_id: "abc123"
		req, err := http.NewRequest("GET", "/fetch-breed-images?breed_id=abc123", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Check if the response contains images related to the breed
	})
}
