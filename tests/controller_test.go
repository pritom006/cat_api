package controllers

import (
	"catapigo/controllers"
	"catapigo/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/require"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	fmt.Println("App path:", apppath)
}

// Helper function to create a controller with context
func createControllerWithContext(method, path string, body string) (*controllers.MainController, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.NewContext()
	ctx.Reset(w, r)
	
	controller := &controllers.MainController{}
	controller.Init(ctx, "", "", nil)
	return controller, w
}

func TestVoteForCat(t *testing.T) {
	Convey("Testing VoteForCat function", t, func() {
		testCases := []struct {
			name        string
			payload     models.Vote
			statusCode  int
		}{
			{
				name: "Valid vote with value 1",
				payload: models.Vote{
					ImageID: "test-image-id-1",
					Value:   1,
				},
				statusCode:  200,
			},
			{
				name: "Valid vote with value 0",
				payload: models.Vote{
					ImageID: "test-image-id-2",
					Value:   0,
				},
				statusCode:  200,
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				payload, err := json.Marshal(tc.payload)
				So(err, ShouldBeNil)

				controller, w := createControllerWithContext("POST", "/vote", string(payload))
				controller.VoteForCat()

				So(w.Code, ShouldEqual, tc.statusCode)

				var response map[string]interface{}
				err = json.NewDecoder(w.Body).Decode(&response)
				So(err, ShouldBeNil)
				So(response, ShouldNotBeNil)
			})
		}

		// Additional case for invalid method
		Convey("Unsupported HTTP method", func() {
			controller, w := createControllerWithContext("GET", "/vote", "")
			controller.VoteForCat()
			So(w.Code, ShouldEqual, 405) // Method Not Allowed
		})
	})
}

// MockHttpClient is a mock implementation of http.Client

func TestFetchCatBreeds(t *testing.T) {
	Convey("Testing FetchCatBreeds function", t, func() {
		Convey("Successful fetch", func() {
			controller, w := createControllerWithContext("GET", "/fetch-cat-breeds", "")
			controller.FetchCatBreeds()

			So(w.Code, ShouldEqual, 200)

			var response []map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(len(response), ShouldBeGreaterThan, 0)

			// Verify response structure
			if len(response) > 0 {
				So(response[0], ShouldContainKey, "id")
				So(response[0], ShouldContainKey, "name")
			}
		})

		// Case for internal server error simulation (e.g., DB failure)
		Convey("Internal server error simulation", func() {
			// Simulate a failure in your data fetching logic
			// For example, mock your service to return an error
			controller, w := createControllerWithContext("GET", "/fetch-cat-breeds", "")
			controller.FetchCatBreeds()

			So(w.Code, ShouldEqual, 500) // Internal Server Error
		})
	})
}

func TestFetchBreedImages(t *testing.T) {
	Convey("Testing FetchBreedImages function", t, func() {
		testCases := []struct {
			name        string
			breedID     string
			statusCode  int
		}{
			{
				name:        "Valid breed ID",
				breedID:     "abys",
				statusCode:  200,
			},
			{
				name:        "Another valid breed ID",
				breedID:     "beng",
				statusCode:  200,
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				controller, w := createControllerWithContext("GET", "/fetch-breed-images?breed_id="+tc.breedID, "")
				controller.FetchBreedImages()

				So(w.Code, ShouldEqual, tc.statusCode)

				var response []map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				So(err, ShouldBeNil)
				So(len(response), ShouldBeGreaterThan, 0)
			})
		}

		// Testing invalid breed ID
		Convey("Invalid breed ID", func() {
			controller, w := createControllerWithContext("GET", "/fetch-breed-images?breed_id=invalid", "")
			controller.FetchBreedImages()

			So(w.Code, ShouldEqual, 404) // Not Found
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})

		// Case for missing breed_id parameter
		Convey("Missing breed_id in query parameters", func() {
			controller, w := createControllerWithContext("GET", "/fetch-breed-images", "")
			controller.FetchBreedImages()

			So(w.Code, ShouldEqual, 400) // Bad Request
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})
	})
}

func TestAddToFavorites(t *testing.T) {
	Convey("Testing AddToFavorites function", t, func() {
		testCases := []struct {
			name        string
			favorite    models.Favorite
			statusCode  int
		}{
			{
				name:        "Valid favorite",
				favorite:    models.Favorite{ImageID: "test-image-id"},
				statusCode:  200,
			},
			{
				name:        "Another valid favorite",
				favorite:    models.Favorite{ImageID: "test-image-id-2"},
				statusCode:  200,
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				payload, err := json.Marshal(tc.favorite)
				So(err, ShouldBeNil)
				fmt.Printf("Payload: %s\n", string(payload))

				controller, w := createControllerWithContext("POST", "/add-to-favorites", string(payload))
				controller.AddToFavorites()

				So(w.Code, ShouldEqual, tc.statusCode)

				var response map[string]interface{}
				err = json.NewDecoder(w.Body).Decode(&response)
				So(err, ShouldBeNil)
				So(response, ShouldNotBeNil)
			})
		}

		// Case for invalid JSON
		Convey("Invalid JSON payload", func() {
			controller, w := createControllerWithContext("POST", "/add-to-favorites", "{invalid_json}")
			controller.AddToFavorites()

			So(w.Code, ShouldEqual, 400)
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})

		// Case for empty ImageID
		Convey("Empty ImageID in payload", func() {
			payload := `{ "ImageID": "" }`
			controller, w := createControllerWithContext("POST", "/add-to-favorites", payload)
			controller.AddToFavorites()

			So(w.Code, ShouldEqual, 400)
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})
	})
}

func TestFetchFavorites(t *testing.T) {
	Convey("Testing FetchFavorites function", t, func() {
		Convey("Successful fetch", func() {
			controller, w := createControllerWithContext("GET", "/fetch-favorites", "")
			controller.FetchFavorites()

			So(w.Code, ShouldEqual, 200)

			var response []map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldNotBeNil)
		})

		// Case for internal server error
		Convey("Internal server error simulation", func() {
			controller, w := createControllerWithContext("GET", "/fetch-favorites", "")
			controller.FetchFavorites()

			So(w.Code, ShouldEqual, 500) // Internal Server Error
		})
	})
}

func TestFetchNewCatImage(t *testing.T) {
	Convey("Testing FetchNewCatImage function", t, func() {
		Convey("Successful fetch", func() {
			controller, w := createControllerWithContext("GET", "/fetch-new-cat-image", "")
			controller.FetchNewCatImage()

			So(w.Code, ShouldEqual, 200)

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "url")
			So(response["url"], ShouldNotBeEmpty)
		})
	})
}

func TestOptionsRequest(t *testing.T) {
    Convey("Testing OPTIONS request handling", t, func() {
        // Create a controller and simulate an OPTIONS request to the endpoint
        controller, w := createControllerWithContext("OPTIONS", "/some-endpoint", "")

        // Call Prepare() manually to test the OPTIONS request handling code
        controller.Prepare()

        // Check if the response status is 200 (OK) as expected by the OPTIONS handling logic
        So(w.Code, ShouldEqual, 200)

        // Check if the response body is empty or has expected CORS headers
        So(w.Body.String(), ShouldBeEmpty) // Adjust based on the actual response
    })
}