// package controllers


// import (
// 	"catapigo/controllers"
// 	//"catapigo/routers"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	beego "github.com/beego/beego/v2/server/web"
// 	"github.com/stretchr/testify/assert"
// )

// func TestRouter(t *testing.T) {
// 	// Initialize Beego Router
// 	beego.Router("/fetch-breeds", &controllers.MainController{}, "get:FetchCatBreeds")
// 	beego.Router("/vote", &controllers.MainController{}, "post:VoteForCat")
// 	beego.Router("/favorites", &controllers.MainController{}, "get:FetchFavorites")
// 	beego.Router("/addToFavorites", &controllers.MainController{}, "post:AddToFavorites")
// 	beego.Router("/fetch-new-cat", &controllers.MainController{}, "get:FetchNewCatImage")
// 	beego.Router("/fetch-breed-images", &controllers.MainController{}, "get:FetchBreedImages")

// 	// Initialize Beego application (if not already initialized)
// 	if beego.BeeApp == nil {
// 		beego.Run()
// 	}

	
// 	// Test FetchCatBreeds
// 	t.Run("Test FetchCatBreeds", func(t *testing.T) {
// 		req, err := http.NewRequest("GET", "/fetch-breeds", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		beego.BeeApp.Handlers.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		// Further checks can be added to verify the actual content
// 	})

// 	// Test VoteForCat
// 	t.Run("Test VoteForCat", func(t *testing.T) {
// 		vote := `{"image_id": "1234", "value": 1}`
// 		req, err := http.NewRequest("POST", "/vote", strings.NewReader(vote))
// 		req.Header.Set("Content-Type", "application/json")
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		beego.BeeApp.Handlers.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		// Check if the response is successful, you can also assert response body for success message
// 	})

// 	// Test FetchFavorites
// 	t.Run("Test FetchFavorites", func(t *testing.T) {
// 		req, err := http.NewRequest("GET", "/favorites", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		beego.BeeApp.Handlers.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		// Check response body contains the favorites
// 	})

// 	// Test AddToFavorites
// 	t.Run("Test AddToFavorites", func(t *testing.T) {
// 		favorite := `{"image_id": "1234"}`
// 		req, err := http.NewRequest("POST", "/addToFavorites", strings.NewReader(favorite))
// 		req.Header.Set("Content-Type", "application/json")
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		beego.BeeApp.Handlers.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		// Check if the response is successful, you can also assert response body for success message
// 	})

// 	// Test FetchNewCatImage
// 	t.Run("Test FetchNewCatImage", func(t *testing.T) {
// 		req, err := http.NewRequest("GET", "/fetch-new-cat", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		beego.BeeApp.Handlers.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		// Check if the response body contains a new cat image
// 	})

// 	// Test FetchBreedImages
// 	t.Run("Test FetchBreedImages", func(t *testing.T) {
// 		req, err := http.NewRequest("GET", "/fetch-breed-images?breed_id=abc123", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		beego.BeeApp.Handlers.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		// Check if the response contains images related to the breed
// 	})
// }


package routers

import (
	//"catapigo/controllers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name           string
	method         string
	path           string
	payload        string
	expectedStatus int
	headers        map[string]string
}

func executeRequest(t *testing.T, tc testCase) *httptest.ResponseRecorder {
	req, err := http.NewRequest(tc.method, tc.path, strings.NewReader(tc.payload))
	assert.NoError(t, err)

	// Set default content type for POST requests
	if tc.method == "POST" && tc.headers == nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set custom headers if provided
	if tc.headers != nil {
		for key, value := range tc.headers {
			req.Header.Set(key, value)
		}
	}

	rr := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(rr, req)
	return rr
}

func TestRouterInitialization(t *testing.T) {
	// Test if all routes are properly registered
	expectedRoutes := []struct {
		path   string
		method string
		name   string
	}{
		{"/", "GET", "Frontend"},
		{"/fetch-breeds", "GET", "Cat Breeds"},
		{"/vote", "POST", "Vote"},
		{"/favorites", "GET", "Favorites"},
		{"/addToFavorites", "POST", "Add to Favorites"},
		{"/fetch-new-cat", "GET", "New Cat Image"},
		{"/fetch-breed-images", "GET", "Breed Images"},
	}

	for _, route := range expectedRoutes {
		t.Run(route.name+" route", func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, req)
			
			// For routes that exist, the status should not be 404
			assert.NotEqual(t, 404, w.Code, "Route should be registered: %s %s", route.method, route.path)
		})
	}
}

func TestRouteHandlers(t *testing.T) {
	testCases := []testCase{
		{
			name:           "Serve Frontend",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Fetch Cat Breeds",
			method:         "GET",
			path:           "/fetch-breeds",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Vote For Cat",
			method:         "POST",
			path:           "/vote",
			payload:        `{"image_id":"test-id","value":1}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Fetch Favorites",
			method:         "GET",
			path:           "/favorites",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Add To Favorites",
			method:         "POST",
			path:           "/addToFavorites",
			payload:        `{"image_id":"test-id"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Fetch New Cat Image",
			method:         "GET",
			path:           "/fetch-new-cat",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Fetch Breed Images",
			method:         "GET",
			path:           "/fetch-breed-images?breed_id=test",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := executeRequest(t, tc)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			
			if tc.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil && tc.path != "/" { // Skip JSON check for frontend route
					assert.NoError(t, err, "Response should be valid JSON")
				}
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	testCases := []testCase{
		{
			name:           "Invalid Route",
			method:         "GET",
			path:           "/invalid-route",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid Method",
			method:         "DELETE",
			path:           "/fetch-breeds",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Missing Breed ID",
			method:         "GET",
			path:           "/fetch-breed-images",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Content Type",
			method:         "POST",
			path:           "/vote",
			payload:        `{"image_id":"test-id","value":1}`,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Invalid JSON",
			method:         "POST",
			path:           "/vote",
			payload:        `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := executeRequest(t, tc)
			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}