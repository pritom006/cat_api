package tests

import (
	"catapigo/controllers"
	"catapigo/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// Helper function to create a controller with context
func createControllerWithContext(method, path string, body string) (*controllers.MainController, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, path, strings.NewReader(body))
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
			expectError bool
		}{
			{
				name: "Valid vote",
				payload: models.Vote{
					ImageID: "test-image-id",
					Value:   1,
				},
				expectError: false,
			},
			{
				name: "Invalid vote value",
				payload: models.Vote{
					ImageID: "test-image-id",
					Value:   3,
				},
				expectError: true,
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				// Marshall the payload
				payload, _ := json.Marshal(tc.payload)

				// Create the request and controller
				controller, w := createControllerWithContext("POST", "/vote", string(payload))
				// Run the controller function
				controller.VoteForCat()

				// Decode the response
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				So(err, ShouldBeNil)

				if tc.expectError {
					So(response, ShouldContainKey, "error")
				} else {
					So(response, ShouldContainKey, "message")
					So(response["message"], ShouldEqual, "Vote submitted successfully")
				}
			})
		}
	})
}

func TestFetchCatBreeds(t *testing.T) {
	Convey("Testing FetchCatBreeds function", t, func() {
		controller, w := createControllerWithContext("GET", "/fetch-cat-breeds", "")

		// Run the controller function
		controller.FetchCatBreeds()

		// Decode the response
		var response []map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		So(err, ShouldBeNil)

		// Assert that the response is an array
		So(len(response), ShouldBeGreaterThan, 0)
	})
}

func TestFetchBreedImages(t *testing.T) {
	Convey("Testing FetchBreedImages function", t, func() {
		testCases := []struct {
			breedID     string
			expectError bool
		}{
			{
				breedID:     "abys",
				expectError: false,
			},
			{
				breedID:     "unknown-breed",
				expectError: true,
			},
		}

		for _, tc := range testCases {
			Convey("Fetching images for breed ID: "+tc.breedID, func() {
				// Create the request and controller
				controller, w := createControllerWithContext("GET", "/fetch-breed-images?breed_id="+tc.breedID, "")

				// Run the controller function
				controller.FetchBreedImages()

				// Decode the response
				var response []map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)

				if tc.expectError {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(len(response), ShouldBeGreaterThan, 0)
				}
			})
		}
	})
}

func TestAddToFavorites(t *testing.T) {
	Convey("Testing AddToFavorites function", t, func() {
		testCases := []struct {
			favorite    models.Favorite
			expectError bool
		}{
			{
				favorite:    models.Favorite{ImageID: "test-image-id"},
				expectError: false,
			},
			{
				favorite:    models.Favorite{},
				expectError: true,
			},
		}

		for _, tc := range testCases {
			Convey("Adding to favorites: "+tc.favorite.ImageID, func() {
				// Marshall the payload
				payload, _ := json.Marshal(tc.favorite)

				// Create the request and controller
				controller, w := createControllerWithContext("POST", "/add-to-favorites", string(payload))
				// Run the controller function
				controller.AddToFavorites()

				// Decode the response
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				So(err, ShouldBeNil)

				if tc.expectError {
					So(response, ShouldContainKey, "error")
				} else {
					So(response, ShouldContainKey, "message")
					So(response["message"], ShouldEqual, "Added to favorites successfully")
				}
			})
		}
	})
}

func TestFetchFavorites(t *testing.T) {
	Convey("Testing FetchFavorites function", t, func() {
		controller, w := createControllerWithContext("GET", "/fetch-favorites", "")

		// Run the controller function
		controller.FetchFavorites()

		// Decode the response
		var response []map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		So(err, ShouldBeNil)

		// Assert that the response is an array
		So(len(response), ShouldBeGreaterThan, 0)
	})
}


func TestFetchNewCatImage(t *testing.T) {
	Convey("Testing FetchNewCatImage function", t, func() {
		controller, w := createControllerWithContext("GET", "/fetch-new-cat-image", "")

		// Run the controller function
		controller.FetchNewCatImage()

		// Decode the response
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		So(err, ShouldBeNil)

		// Assert that the response contains a valid image URL
		So(response, ShouldContainKey, "url")
	})
}

// func TestControllerEndpoints(t *testing.T) {
// 	Convey("Testing Controller Endpoints", t, func() {
// 		Convey("Testing Prepare Method", func() {
// 			_, w := createControllerWithContext("OPTIONS", "/", "")
// 			// Manually set the headers that are normally set by Prepare
// 			w.Header().Set("Access-Control-Allow-Origin", "*")
// 			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
// 			So(w.Header().Get("Access-Control-Allow-Origin"), ShouldEqual, "*")
// 			So(w.Header().Get("Access-Control-Allow-Methods"), ShouldEqual, "GET, POST, OPTIONS")
// 			So(w.Header().Get("Access-Control-Allow-Headers"), ShouldContainSubstring, "Content-Type")
// 		})

// 		Convey("Testing ServeFrontend", func() {
// 			controller, _ := createControllerWithContext("GET", "/", "")
// 			controller.ServeFrontend()
// 			So(controller.TplName, ShouldEqual, "index.tpl")
// 		})

// 		Convey("Testing FetchCatBreeds", func() {
// 			controller, w := createControllerWithContext("GET", "/breeds", "")
// 			controller.FetchCatBreeds()

// 			var response []map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
			
// 			if err != nil {
// 				var errorResp map[string]string
// 				err = json.NewDecoder(w.Body).Decode(&errorResp)
// 				So(err, ShouldBeNil)
// 			}
// 		})

// 		Convey("Testing FetchBreedImages", func() {
// 			controller, w := createControllerWithContext("GET", "/breed-images?breed_id=beng", "")
// 			controller.FetchBreedImages()

// 			var response []map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
			
// 			if err != nil {
// 				var errorResp map[string]string
// 				err = json.NewDecoder(w.Body).Decode(&errorResp)
// 				So(err, ShouldBeNil)
// 			}
// 		})

// 		Convey("Testing VoteForCat", func() {
// 			testCases := []struct {
// 				name        string
// 				payload     models.Vote
// 				expectError bool
// 			}{
// 				{
// 					name: "Valid vote",
// 					payload: models.Vote{
// 						ImageID: "test-image-id",
// 						Value:   1,
// 					},
// 					expectError: false,
// 				},
// 				{
// 					name: "Invalid vote value",
// 					payload: models.Vote{
// 						ImageID: "test-image-id",
// 						Value:   3,
// 					},
// 					expectError: true,
// 				},
// 			}

// 			for _, tc := range testCases {
// 				Convey(tc.name, func() {
// 					// Marshall the payload
// 					payload, _ := json.Marshal(tc.payload)
					
// 					// Create the request and controller
// 					controller, w := createControllerWithContext("POST", "/vote", string(payload))
// 					// Run the controller function
// 					controller.VoteForCat()

// 					// Wait for the response (use time.Sleep to wait for the goroutine)
// 					// This ensures the asynchronous process completes before the test asserts
// 					// (This is generally not recommended in production code, but can help in tests).
// 					// Ideally, the test should be more sophisticated to handle the async logic better.
// 					// In production, better ways to sync go routines like channels can be used.
					
// 					// Decode the response
// 					var response map[string]string
// 					err := json.NewDecoder(w.Body).Decode(&response)
// 					So(err, ShouldBeNil)

// 					// Log the response for debugging purposes
// 					t.Log("Response body:", w.Body.String())

// 					// Assert that the response contains the correct key
// 					if tc.expectError {
// 						So(response, ShouldContainKey, "error")
// 					} else {
// 						So(response, ShouldContainKey, "message")
// 						So(response["message"], ShouldEqual, "Vote submitted successfully")
// 					}
// 				})
// 			}
// 		})

// 		Convey("Testing AddToFavorites", func() {
// 			testCases := []struct {
// 				name        string
// 				payload     string
// 				expectError bool
// 			}{
// 				{
// 					name:        "Valid favorite",
// 					payload:     `{"image_id": "test-image-id"}`,
// 					expectError: false,
// 				},
// 				{
// 					name:        "Missing image_id",
// 					payload:     `{}`,
// 					expectError: true,
// 				},
// 				{
// 					name:        "Invalid JSON",
// 					payload:     `{invalid-json}`,
// 					expectError: true,
// 				},
// 			}

// 			for _, tc := range testCases {
// 				Convey(tc.name, func() {
// 					controller, w := createControllerWithContext("POST", "/favorites", tc.payload)
// 					controller.AddToFavorites()

// 					var response map[string]string
// 					err := json.NewDecoder(w.Body).Decode(&response)
// 					So(err, ShouldBeNil)

// 					if tc.expectError {
// 						So(response, ShouldContainKey, "error")
// 					} else {
// 						So(response, ShouldContainKey, "message")
// 					}
// 				})
// 			}
// 		})

// 		Convey("Testing FetchFavorites", func() {
// 			controller, w := createControllerWithContext("GET", "/favorites", "")
// 			controller.FetchFavorites()

// 			var response []map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
			
// 			if err != nil {
// 				var errorResp map[string]string
// 				err = json.NewDecoder(w.Body).Decode(&errorResp)
// 				So(err, ShouldBeNil)
// 			}
// 		})

// 		Convey("Testing FetchNewCatImage", func() {
// 			controller, w := createControllerWithContext("GET", "/new-cat", "")
// 			controller.FetchNewCatImage()

// 			var response map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
			
// 			if err != nil {
// 				var errorResp map[string]string
// 				err = json.NewDecoder(w.Body).Decode(&errorResp)
// 				So(err, ShouldBeNil)
// 			}
// 		})
// 	})
// }
