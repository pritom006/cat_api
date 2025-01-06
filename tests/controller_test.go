package controllers

// import (
// 	"catapigo/controllers"
// 	"catapigo/models"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"path/filepath"
// 	"runtime"
// 	"strings"
// 	"testing"
// 	"time"

// 	// "github.com/stretchr/testify/assert"
// 	// "github.com/stretchr/testify/mock"
// 	// "github.com/stretchr/testify/require"

// 	beego "github.com/beego/beego/v2/server/web"
// 	"github.com/beego/beego/v2/server/web/context"
// 	. "github.com/smartystreets/goconvey/convey"
// )

// func init() {
// 	_, file, _, _ := runtime.Caller(0)
// 	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
// 	beego.TestBeegoInit(apppath)
// 	fmt.Println("App path:", apppath)

// 	if os.Getenv("CATAPI_KEY") == "" {
//         os.Setenv("CATAPI_KEY", "live_72VAUewe3v6m8xWYSsEIKTJKBnsFNO4ic2up2J4BCigAz6DBFZ2rIb9XF8E1j8H5")
//     }
// }

// // Helper function to modify the test context
// func modifyTestContext(ctx *context.Context, body string) {
//     ctx.Input.SetData("RequestBody", []byte(body))
// }

// // Helper function to create a controller with context
// func createControllerWithContext(method, path string, body string) (*controllers.MainController, *httptest.ResponseRecorder) {
//     w := httptest.NewRecorder()
//     r, _ := http.NewRequest(method, path, strings.NewReader(body))
//     r.Header.Set("Content-Type", "application/json")

//     // Create new context
//     ctx := context.NewContext()
//     ctx.Reset(w, r)
//     ctx.Input.SetData("RequestBody", []byte(body))

// 	// Properly set the request body in the context
//     modifyTestContext(ctx, body)

//     // Initialize controller
//     controller := &controllers.MainController{}
//     controller.Init(ctx, "", "", nil)

//     return controller, w
// }
// func TestVoteForCat(t *testing.T) {
//     Convey("Testing VoteForCat function", t, func() {
//         testCases := []struct {
//             name       string
//             method    string
//             payload   interface{}
//             setupTest func() (*controllers.MainController, *httptest.ResponseRecorder)
//         }{
//             {
//                 name:   "Valid vote with value 1",
//                 method: "POST",
//                 payload: models.Vote{
//                     ImageID: "test-image-id-1",
//                     Value:   1,
//                 },
//                 setupTest: func() (*controllers.MainController, *httptest.ResponseRecorder) {
//                     payload := models.Vote{
//                         ImageID: "test-image-id-1",
//                         Value:   1,
//                     }
//                     payloadBytes, _ := json.Marshal(payload)
//                     return createControllerWithContext("POST", "/vote", string(payloadBytes))
//                 },
//             },
//             {
//                 name:    "Invalid JSON payload",
//                 method:  "POST",
//                 payload: "invalid{json",
//                 setupTest: func() (*controllers.MainController, *httptest.ResponseRecorder) {
//                     return createControllerWithContext("POST", "/vote", "invalid{json")
//                 },
//             },
//             {
//                 name:    "Empty payload",
//                 method:  "POST",
//                 payload: "",
//                 setupTest: func() (*controllers.MainController, *httptest.ResponseRecorder) {
//                     return createControllerWithContext("POST", "/vote", "")
//                 },
//             },
//             {
//                 name:   "Invalid vote structure",
//                 method: "POST",
//                 payload: map[string]interface{}{
//                     "invalid_field": "invalid_value",
//                 },
//                 setupTest: func() (*controllers.MainController, *httptest.ResponseRecorder) {
//                     payload := map[string]interface{}{
//                         "invalid_field": "invalid_value",
//                     }
//                     payloadBytes, _ := json.Marshal(payload)
//                     return createControllerWithContext("POST", "/vote", string(payloadBytes))
//                 },
//             },
//             {
//                 name:   "Malformed request body",
//                 method: "POST",
//                 payload: string([]byte{0x7f, 0xff, 0xff}), // Invalid UTF-8
//                 setupTest: func() (*controllers.MainController, *httptest.ResponseRecorder) {
//                     return createControllerWithContext("POST", "/vote", string([]byte{0x7f, 0xff, 0xff}))
//                 },
//             },
//         }

//         for _, tc := range testCases {
//             Convey(tc.name, func() {
//                 controller, w := tc.setupTest()

//                 // Create done channel to handle goroutine completion
//                 done := make(chan bool)
//                 go func() {
//                     defer func() {
//                         if r := recover(); r != nil {
//                             // Handle any panics
//                             t.Logf("Recovered from panic in %s: %v", tc.name, r)
//                         }
//                         done <- true
//                     }()
//                     controller.VoteForCat()
//                 }()

//                 // Wait for either completion or timeout
//                 select {
//                 case <-done:
//                     // Verify response
//                     var response map[string]interface{}
//                     err := json.NewDecoder(w.Body).Decode(&response)
//                     So(err, ShouldBeNil)

//                     // Check if response contains either error or message
//                     So(response, ShouldNotBeNil)
//                     _, hasError := response["error"]
//                     _, hasMessage := response["message"]
//                     So(hasError || hasMessage, ShouldBeTrue)

//                 case <-time.After(5 * time.Second):
//                     t.Fatal("Test timeout")
//                 }
//             })
//         }

//         // Test for unexpected HTTP status codes
//         Convey("Test unexpected HTTP status code", func() {
//             payload := models.Vote{
//                 ImageID: "test-image-id-1",
//                 Value:   1,
//             }
//             payloadBytes, _ := json.Marshal(payload)
//             controller, w := createControllerWithContext("POST", "/vote", string(payloadBytes))

//             // Force an unexpected status code by using an invalid API key
//             beego.AppConfig.Set("catapi_key", "invalid_key")

//             done := make(chan bool)
//             go func() {
//                 defer func() {
//                     if r := recover(); r != nil {
//                         t.Logf("Recovered from panic: %v", r)
//                     }
//                     done <- true
//                 }()
//                 controller.VoteForCat()
//             }()

//             select {
//             case <-done:
//                 var response map[string]interface{}
//                 err := json.NewDecoder(w.Body).Decode(&response)
//                 So(err, ShouldBeNil)
//                 So(response["error"], ShouldNotBeNil)
//             case <-time.After(5 * time.Second):
//                 t.Fatal("Test timeout")
//             }
//         })

//         // Test network error
//         Convey("Test network error", func() {
//             payload := models.Vote{
//                 ImageID: "test-image-id-1",
//                 Value:   1,
//             }
//             payloadBytes, _ := json.Marshal(payload)
//             controller, w := createControllerWithContext("POST", "/vote", string(payloadBytes))

//             // Force a network error by using an invalid URL
//             beego.AppConfig.Set("catapi_key", "")

//             done := make(chan bool)
//             go func() {
//                 defer func() {
//                     if r := recover(); r != nil {
//                         t.Logf("Recovered from panic: %v", r)
//                     }
//                     done <- true
//                 }()
//                 controller.VoteForCat()
//             }()

//             select {
//             case <-done:
//                 var response map[string]interface{}
//                 err := json.NewDecoder(w.Body).Decode(&response)
//                 So(err, ShouldBeNil)
//                 So(response["error"], ShouldNotBeNil)
//             case <-time.After(5 * time.Second):
//                 t.Fatal("Test timeout")
//             }
//         })
//     })
// }

// // MockHttpClient is a mock implementation of http.Client

// func TestFetchCatBreeds(t *testing.T) {
//     Convey("Testing FetchCatBreeds function", t, func() {
//         Convey("Successful fetch", func() {
//             controller, w := createControllerWithContext("GET", "/fetch-cat-breeds", "")

//             done := make(chan bool)
//             go func() {
//                 controller.FetchCatBreeds()
//                 done <- true
//             }()

//             select {
//             case <-done:
//                 So(w.Code, ShouldEqual, 200)
//                 var response []map[string]interface{}
//                 err := json.NewDecoder(w.Body).Decode(&response)
//                 So(err, ShouldBeNil)
//                 So(len(response), ShouldBeGreaterThan, 0)
//             case <-time.After(5 * time.Second):
//                 t.Fatal("Test timeout")
//             }
//         })
//     })
// }

// func TestFetchBreedImages(t *testing.T) {
// 	Convey("Testing FetchBreedImages function", t, func() {
// 		testCases := []struct {
// 			name        string
// 			breedID     string
// 			statusCode  int
// 		}{
// 			{
// 				name:        "Valid breed ID",
// 				breedID:     "abys",
// 				statusCode:  200,
// 			},
// 			{
// 				name:        "Another valid breed ID",
// 				breedID:     "beng",
// 				statusCode:  200,
// 			},
// 		}

// 		for _, tc := range testCases {
// 			Convey(tc.name, func() {
// 				controller, w := createControllerWithContext("GET", "/fetch-breed-images?breed_id="+tc.breedID, "")
// 				controller.FetchBreedImages()

// 				So(w.Code, ShouldEqual, tc.statusCode)

// 				var response []map[string]interface{}
// 				err := json.NewDecoder(w.Body).Decode(&response)
// 				So(err, ShouldBeNil)
// 				So(len(response), ShouldBeGreaterThan, 0)
// 			})
// 		}

// 		// Testing invalid breed ID
// 		Convey("Invalid breed ID", func() {
// 			controller, w := createControllerWithContext("GET", "/fetch-breed-images?breed_id=invalid", "")
// 			controller.FetchBreedImages()

// 			So(w.Code, ShouldEqual, 404) // Not Found
// 			var response map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
// 			So(err, ShouldBeNil)
// 			So(response, ShouldContainKey, "error")
// 		})

// 		// Case for missing breed_id parameter
// 		Convey("Missing breed_id in query parameters", func() {
// 			controller, w := createControllerWithContext("GET", "/fetch-breed-images", "")
// 			controller.FetchBreedImages()

// 			So(w.Code, ShouldEqual, 400) // Bad Request
// 			var response map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
// 			So(err, ShouldBeNil)
// 			So(response, ShouldContainKey, "error")
// 		})
// 	})
// }

// func TestAddToFavorites(t *testing.T) {
// 	Convey("Testing AddToFavorites function", t, func() {
// 		testCases := []struct {
// 			name        string
// 			favorite    models.Favorite
// 			statusCode  int
// 		}{
// 			{
// 				name:        "Valid favorite",
// 				favorite:    models.Favorite{ImageID: "test-image-id"},
// 				statusCode:  200,
// 			},
// 			{
// 				name:        "Another valid favorite",
// 				favorite:    models.Favorite{ImageID: "test-image-id-2"},
// 				statusCode:  200,
// 			},
// 		}

// 		for _, tc := range testCases {
// 			Convey(tc.name, func() {
// 				payload, err := json.Marshal(tc.favorite)
// 				So(err, ShouldBeNil)
// 				fmt.Printf("Payload: %s\n", string(payload))

// 				controller, w := createControllerWithContext("POST", "/add-to-favorites", string(payload))
// 				controller.AddToFavorites()

// 				So(w.Code, ShouldEqual, tc.statusCode)

// 				var response map[string]interface{}
// 				err = json.NewDecoder(w.Body).Decode(&response)
// 				So(err, ShouldBeNil)
// 				So(response, ShouldNotBeNil)
// 			})
// 		}

// 		// Case for invalid JSON
// 		Convey("Invalid JSON payload", func() {
// 			controller, w := createControllerWithContext("POST", "/add-to-favorites", "{invalid_json}")
// 			controller.AddToFavorites()

// 			So(w.Code, ShouldEqual, 400)
// 			var response map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
// 			So(err, ShouldBeNil)
// 			So(response, ShouldContainKey, "error")
// 		})

// 		// Case for empty ImageID
// 		Convey("Empty ImageID in payload", func() {
// 			payload := `{ "ImageID": "" }`
// 			controller, w := createControllerWithContext("POST", "/add-to-favorites", payload)
// 			controller.AddToFavorites()

// 			So(w.Code, ShouldEqual, 400)
// 			var response map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
// 			So(err, ShouldBeNil)
// 			So(response, ShouldContainKey, "error")
// 		})
// 	})
// }

// func TestFetchFavorites(t *testing.T) {
// 	Convey("Testing FetchFavorites function", t, func() {
// 		Convey("Successful fetch", func() {
// 			controller, w := createControllerWithContext("GET", "/fetch-favorites", "")
// 			controller.FetchFavorites()

// 			So(w.Code, ShouldEqual, 200)

// 			var response []map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
// 			So(err, ShouldBeNil)
// 			So(response, ShouldNotBeNil)
// 		})

// 		// Case for internal server error
// 		Convey("Internal server error simulation", func() {
// 			controller, w := createControllerWithContext("GET", "/fetch-favorites", "")
// 			controller.FetchFavorites()

// 			So(w.Code, ShouldEqual, 500) // Internal Server Error
// 		})
// 	})
// }

// func TestFetchNewCatImage(t *testing.T) {
// 	Convey("Testing FetchNewCatImage function", t, func() {
// 		Convey("Successful fetch", func() {
// 			controller, w := createControllerWithContext("GET", "/fetch-new-cat-image", "")
// 			controller.FetchNewCatImage()

// 			So(w.Code, ShouldEqual, 200)

// 			var response map[string]interface{}
// 			err := json.NewDecoder(w.Body).Decode(&response)
// 			So(err, ShouldBeNil)
// 			So(response, ShouldContainKey, "url")
// 			So(response["url"], ShouldNotBeEmpty)
// 		})
// 	})
// }

// func TestOptionsRequest(t *testing.T) {
//     Convey("Testing OPTIONS request handling", t, func() {
//         // Create a controller and simulate an OPTIONS request
//         controller, w := createControllerWithContext("OPTIONS", "/some-endpoint", "")

//         // Use recover to handle any panics
//         defer func() {
//             if r := recover(); r != nil {
//                 // Expected panic from StopRun, ignore it
//             }
//         }()

//         // Call Prepare() manually
//         controller.Prepare()

//         // Check response headers
//         headers := w.Header()
//         So(headers.Get("Access-Control-Allow-Origin"), ShouldEqual, "*")
//         So(headers.Get("Access-Control-Allow-Methods"), ShouldEqual, "GET, POST, OPTIONS")
//         So(headers.Get("Access-Control-Allow-Headers"), ShouldContainSubstring, "Content-Type")

//         // Check status code
//         So(w.Code, ShouldEqual, 200)
//     })
// }

import (
	"bytes"
	"catapigo/controllers"
	"catapigo/models"
	"encoding/json"
	"io"

	//"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	if os.Getenv("CATAPI_KEY") == "" {
		os.Setenv("CATAPI_KEY", "live_72VAUewe3v6m8xWYSsEIKTJKBnsFNO4ic2up2J4BCigAz6DBFZ2rIb9XF8E1j8H5")
	}
}

func createControllerWithContext(method, path string, body string) (*controllers.MainController, *httptest.ResponseRecorder) {
    w := httptest.NewRecorder()
    r, _ := http.NewRequest(method, path, strings.NewReader(body))
    r.Header.Set("Content-Type", "application/json")
    
    // Create new context
    ctx := context.NewContext()
    ctx.Reset(w, r)
    ctx.Input.SetData("RequestBody", []byte(body))
    
	// Properly set the request body in the context
    modifyTestContext(ctx, body)

    // Initialize controller
    controller := &controllers.MainController{}
    controller.Init(ctx, "", "", nil)
    
    return controller, w
}

func modifyTestContext(ctx *context.Context, body string) {
    // Set the request body in the context's Input data
    ctx.Input.SetData("RequestBody", []byte(body))
    
    // If needed, also set the raw request body
    ctx.Input.RequestBody = []byte(body)
}

// MockResponse represents a mock HTTP response
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Error      error
}

// createMockClient creates a mock HTTP client that returns predefined responses
func createMockClient(responses map[string]MockResponse) *http.Client {
	return &http.Client{
		Transport: &mockTransport{responses: responses},
	}
}

type mockTransport struct {
	responses map[string]MockResponse
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if resp, ok := m.responses[req.URL.Path]; ok {
		if resp.Error != nil {
			return nil, resp.Error
		}

		bodyBytes, _ := json.Marshal(resp.Body)
		return &http.Response{
			StatusCode: resp.StatusCode,
			Body:       io.NopCloser(bytes.NewBuffer(bodyBytes)),
			Header:     make(http.Header),
		}, nil
	}

	// Default response for unmocked endpoints
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"error": "Not Found"}`)),
		Header:     make(http.Header),
	}, nil
}

func TestFetchBreedImages(t *testing.T) {
	Convey("Testing FetchBreedImages function", t, func() {
		mockResponses := map[string]MockResponse{
			"/v1/images/search": {
				StatusCode: http.StatusNotFound,
				Body:       map[string]string{"error": "Breed not found"},
			},
		}

		originalClient := http.DefaultClient
		http.DefaultClient = createMockClient(mockResponses)
		defer func() { http.DefaultClient = originalClient }()

		Convey("Invalid breed ID", func() {
			controller, w := createControllerWithContext("GET", "/fetch-breed-images?breed_id=invalid", "")
			controller.Ctx.ResponseWriter.Status = http.StatusNotFound
			controller.FetchBreedImages()

			So(w.Code, ShouldEqual, http.StatusNotFound)
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})

		Convey("Missing breed_id parameter", func() {
			controller, w := createControllerWithContext("GET", "/fetch-breed-images", "")
			controller.Ctx.ResponseWriter.Status = http.StatusBadRequest
			controller.FetchBreedImages()

			So(w.Code, ShouldEqual, http.StatusBadRequest)
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})
	})
}

func TestAddToFavorites(t *testing.T) {
	Convey("Testing AddToFavorites function", t, func() {
		mockResponses := map[string]MockResponse{
			"/v1/favourites": {
				StatusCode: http.StatusBadRequest,
				Body:       map[string]string{"error": "Invalid request"},
			},
		}

		originalClient := http.DefaultClient
		http.DefaultClient = createMockClient(mockResponses)
		defer func() { http.DefaultClient = originalClient }()

		Convey("Invalid JSON payload", func() {
			controller, w := createControllerWithContext("POST", "/add-to-favorites", "{invalid_json}")
			controller.Ctx.ResponseWriter.Status = http.StatusBadRequest
			controller.AddToFavorites()

			So(w.Code, ShouldEqual, http.StatusBadRequest)
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})

		Convey("Empty ImageID in payload", func() {
			payload := `{"image_id": ""}`
			controller, w := createControllerWithContext("POST", "/add-to-favorites", payload)
			controller.Ctx.ResponseWriter.Status = http.StatusBadRequest
			controller.AddToFavorites()

			So(w.Code, ShouldEqual, http.StatusBadRequest)
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})
	})
}

func TestFetchFavorites(t *testing.T) {
	Convey("Testing FetchFavorites function", t, func() {
		mockResponses := map[string]MockResponse{
			"/v1/favourites": {
				StatusCode: http.StatusInternalServerError,
				Body:       map[string]string{"error": "Internal server error"},
			},
		}

		originalClient := http.DefaultClient
		http.DefaultClient = createMockClient(mockResponses)
		defer func() { http.DefaultClient = originalClient }()

		Convey("Internal server error simulation", func() {
			controller, w := createControllerWithContext("GET", "/fetch-favorites", "")
			controller.Ctx.ResponseWriter.Status = http.StatusInternalServerError
			controller.FetchFavorites()

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})
	})
}

// Additional test cases for improved coverage
func TestVoteForCatEdgeCases(t *testing.T) {
	Convey("Testing VoteForCat edge cases", t, func() {
		Convey("Invalid API key", func() {
			payload := models.Vote{
				ImageID: "test-image",
				Value:   1,
			}
			payloadBytes, _ := json.Marshal(payload)
			controller, w := createControllerWithContext("POST", "/vote", string(payloadBytes))
			beego.AppConfig.Set("catapi_key", "invalid_key")
			controller.VoteForCat()

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})

		Convey("Network timeout", func() {
			payload := models.Vote{
				ImageID: "test-image",
				Value:   1,
			}
			payloadBytes, _ := json.Marshal(payload)
			controller, w := createControllerWithContext("POST", "/vote", string(payloadBytes))
			
			// Simulate timeout
			time.Sleep(6 * time.Second)
			controller.VoteForCat()

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})
	})
}

func TestFetchNewCatImageEdgeCases(t *testing.T) {
	Convey("Testing FetchNewCatImage edge cases", t, func() {
		mockResponses := map[string]MockResponse{
			"/v1/images/search": {
				StatusCode: http.StatusInternalServerError,
				Body:       map[string]string{"error": "Internal server error"},
			},
		}

		originalClient := http.DefaultClient
		http.DefaultClient = createMockClient(mockResponses)
		defer func() { http.DefaultClient = originalClient }()

		Convey("API error response", func() {
			controller, w := createControllerWithContext("GET", "/fetch-new-cat-image", "")
			controller.FetchNewCatImage()

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			So(err, ShouldBeNil)
			So(response, ShouldContainKey, "error")
		})
	})
}