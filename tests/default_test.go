

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"

    "github.com/beego/beego/v2/core/logs"

	_ "catapigo/routers"

	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}


// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}




// package test

// import (
// 	"bytes"
// 	"catapigo/controllers"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"path/filepath"
// 	"runtime"
// 	"strings"
// 	"testing"

// 	"github.com/beego/beego/v2/core/logs"
// 	beego "github.com/beego/beego/v2/server/web"
// 	"github.com/smartystreets/goconvey/convey" // Import GoConvey for testing
// 	_ "catapigo/routers"
// )

// // Initialize test setup
// func init() {
// 	_, file, _, _ := runtime.Caller(0)
// 	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
// 	beego.TestBeegoInit(apppath)
// }

// // Test for basic routes and methods
// func TestBeegoBasicRoute(t *testing.T) {
// 	r, _ := http.NewRequest("GET", "/", nil)
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	logs.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

// 	convey.Convey("Subject: Test Station Endpoint\n", t, func() {
// 		convey.Convey("Status Code Should Be 200", func() {
// 			convey.So(w.Code, convey.ShouldEqual, 200)
// 		})
// 		convey.Convey("The Result Should Not Be Empty", func() {
// 			convey.So(w.Body.Len(), convey.ShouldBeGreaterThan, 0)
// 		})
// 	})
// }

// // Test for FetchFavorites endpoint
// func TestFetchFavoritesEndpoint(t *testing.T) {
// 	r, _ := http.NewRequest("GET", "/favorites", nil)
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	convey.Convey("Subject: Test FetchFavorites Endpoint\n", t, func() {
// 		convey.Convey("Status Code Should Be 200", func() {
// 			convey.So(w.Code, convey.ShouldEqual, 200)
// 		})
// 		convey.Convey("Response Should Contain JSON", func() {
// 			var response []map[string]interface{}
// 			err := json.Unmarshal(w.Body.Bytes(), &response)
// 			convey.So(err, convey.ShouldBeNil)
// 			convey.So(len(response), convey.ShouldBeGreaterThan, 0)
// 		})
// 	})
// }

// // Test for FetchNewCatImage endpoint
// func TestFetchNewCatImageEndpoint(t *testing.T) {
// 	r, _ := http.NewRequest("GET", "/fetch-new-cat", nil)
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	convey.Convey("Subject: Test FetchNewCatImage Endpoint\n", t, func() {
// 		convey.Convey("Status Code Should Be 200", func() {
// 			convey.So(w.Code, convey.ShouldEqual, 200)
// 		})
// 		convey.Convey("Response Should Contain JSON", func() {
// 			var response map[string]interface{}
// 			err := json.Unmarshal(w.Body.Bytes(), &response)
// 			convey.So(err, convey.ShouldBeNil)
// 			convey.So(response, convey.ShouldContainKey, "url")
// 		})
// 	})
// }

// // Test for AddToFavorites endpoint (POST request)
// func TestAddToFavoritesEndpoint(t *testing.T) {
// 	favorite := map[string]string{"image_id": "abc123"}

// 	body, _ := json.Marshal(favorite)
// 	r, _ := http.NewRequest("POST", "/addToFavorites", bytes.NewBuffer(body))
// 	r.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	convey.Convey("Subject: Test AddToFavorites Endpoint\n", t, func() {
// 		convey.Convey("Status Code Should Be 200", func() {
// 			convey.So(w.Code, convey.ShouldEqual, 200)
// 		})
// 		convey.Convey("Response Should Contain Success Message", func() {
// 			var response map[string]string
// 			err := json.Unmarshal(w.Body.Bytes(), &response)
// 			convey.So(err, convey.ShouldBeNil)
// 			convey.So(response["message"], convey.ShouldEqual, "Added to favorites successfully")
// 		})
// 	})
// }

// // Test for VoteForCat endpoint (POST request)
// func TestVoteForCatEndpoint(t *testing.T) {
// 	vote := map[string]interface{}{
// 		"image_id": "abc123",
// 		"value":    1,
// 	}

// 	body, _ := json.Marshal(vote)
// 	r, _ := http.NewRequest("POST", "/vote", bytes.NewBuffer(body))
// 	r.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)

// 	convey.Convey("Subject: Test VoteForCat Endpoint\n", t, func() {
// 		convey.Convey("Status Code Should Be 200", func() {
// 			convey.So(w.Code, convey.ShouldEqual, 200)
// 		})
// 		convey.Convey("Response Should Contain Success Message", func() {
// 			var response map[string]string
// 			err := json.Unmarshal(w.Body.Bytes(), &response)
// 			convey.So(err, convey.ShouldBeNil)
// 			convey.So(response["message"], convey.ShouldEqual, "Vote submitted successfully")
// 		})
// 	})
// }

// // Test Prepare and ServeFrontend for controllers
// func TestControllerPrepareAndServeFrontend(t *testing.T) {
// 	// Test Prepare
// 	t.Run("TestPrepare", func(t *testing.T) {
// 		ctrl := &controllers.MainController{}
// 		w := httptest.NewRecorder()
// 		r := httptest.NewRequest("OPTIONS", "/", nil)
// 		ctrl.Ctx.Reset(w, r)
// 		ctrl.Prepare()

// 		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
// 			t.Error("CORS header not set correctly")
// 		}
// 	})

// 	// Test ServeFrontend
// 	t.Run("TestServeFrontend", func(t *testing.T) {
// 		ctrl := &controllers.MainController{}
// 		w := httptest.NewRecorder()
// 		r := httptest.NewRequest("GET", "/", nil)
// 		ctrl.Ctx.Reset(w, r)
// 		ctrl.ServeFrontend()

// 		if ctrl.TplName != "index.tpl" {
// 			t.Errorf("Expected template name 'index.tpl', got %s", ctrl.TplName)
// 		}
// 	})
// }

// // Test FetchCatBreeds functionality
// func TestFetchCatBreedsFunctionality(t *testing.T) {
// 	t.Run("SuccessfulFetch", func(t *testing.T) {
// 		ctrl := &controllers.MainController{}
// 		w := httptest.NewRecorder()
// 		r := httptest.NewRequest("GET", "/fetch-breeds", nil)
// 		ctrl.Ctx.Reset(w, r)
// 		ctrl.FetchCatBreeds()

// 		var response []map[string]interface{}
// 		err := json.NewDecoder(w.Body).Decode(&response)
// 		if err != nil && !strings.Contains(w.Body.String(), "error") {
// 			t.Errorf("Failed to decode response: %v", err)
// 		}
// 	})
// }

// // Test FetchBreedImages functionality
// func TestFetchBreedImagesFunctionality(t *testing.T) {
// 	t.Run("WithValidBreedID", func(t *testing.T) {
// 		ctrl := &controllers.MainController{}
// 		w := httptest.NewRecorder()
// 		r := httptest.NewRequest("GET", "/fetch-breed-images?breed_id=test", nil)
// 		ctrl.Ctx.Reset(w, r)
// 		ctrl.FetchBreedImages()

// 		var response []map[string]interface{}
// 		err := json.NewDecoder(w.Body).Decode(&response)
// 		if err != nil && !strings.Contains(w.Body.String(), "error") {
// 			t.Errorf("Failed to decode response: %v", err)
// 		}
// 	})
// }

// // Test VoteForCat functionality
// func TestVoteForCatFunctionality(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		payload  []byte
// 		wantErr  bool
// 		errMsg   string
// 	}{
// 		{
// 			name:     "ValidVote",
// 			payload:  []byte(`{"image_id":"test123","value":1}`),
// 			wantErr:  false,
// 		},
// 		{
// 			name:     "InvalidJSON",
// 			payload:  []byte(`{invalid json}`),
// 			wantErr:  true,
// 			errMsg:   "Invalid input",
// 		},
// 		{
// 			name:     "EmptyVote",
// 			payload:  []byte(`{}`),
// 			wantErr:  true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := &controllers.MainController{}
// 			w := httptest.NewRecorder()
// 			r := httptest.NewRequest("POST", "/vote", bytes.NewBuffer(tt.payload))
// 			ctrl.Ctx.Reset(w, r)
// 			ctrl.Ctx.Input.RequestBody = tt.payload
// 			ctrl.VoteForCat()

// 			var response map[string]string
// 			err := json.NewDecoder(w.Body).Decode(&response)
// 			if err != nil {
// 				t.Errorf("Failed to decode response: %v", err)
// 				return
// 			}

// 			if tt.wantErr && response["error"] == "" {
// 				t.Error("Expected error but got success")
// 			}
// 			if tt.errMsg != "" && response["error"] != tt.errMsg {
// 				t.Errorf("Expected error message %s, got %s", tt.errMsg, response["error"])
// 			}
// 		})
// 	}
// }
