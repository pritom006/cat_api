package controllers

import (
	"bytes"
	"catapigo/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}




func (c *MainController) Prepare() {
    c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
    c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    c.Ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
    
    // Handle OPTIONS requests
    if c.Ctx.Request.Method == "OPTIONS" {
        c.Ctx.Output.SetStatus(200)
        c.StopRun()
    }
}

func (c *MainController) ServeFrontend() {
    c.TplName = "index.tpl" // This will render the index.tpl file
}




// FetchCatBreeds fetches all available cat breeds from TheCatAPI
func (c *MainController) FetchCatBreeds() {
	apiKey, _ := beego.AppConfig.String("catapi_key")
	url := "https://api.thecatapi.com/v1/breeds"

	resultChan := make(chan []map[string]interface{})
	errChan := make(chan error)

	go func() {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		var result []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			errChan <- err
			return
		}

		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		c.Data["json"] = result
	case err := <-errChan:
		c.Data["json"] = map[string]string{"error": "Failed to fetch breeds: " + err.Error()}
	}

	c.ServeJSON()
}


func (c *MainController) FetchBreedImages() {
    breedID := c.GetString("breed_id")
    apiKey, _ := beego.AppConfig.String("catapi_key")
    url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?breed_ids=%s&limit=10", breedID)

    resultChan := make(chan []map[string]interface{})
    errorChan := make(chan error)

    go func() {
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            errorChan <- err
            return
        }
        req.Header.Add("x-api-key", apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errorChan <- err
            return
        }
        defer resp.Body.Close()

        var result []map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            errorChan <- err
            return
        }

        resultChan <- result
    }()

    select {
    case result := <-resultChan:
        c.Data["json"] = result
    case err := <-errorChan:
        c.Data["json"] = map[string]string{"error": err.Error()}
    }

    c.ServeJSON()
}



// VoteForCat handles voting (like/dislike) for a cat image
func (c *MainController) VoteForCat() {
	var vote models.Vote
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vote); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid input"}
		c.ServeJSON()
		return
	}

	apiKey, _ := beego.AppConfig.String("catapi_key")
	url := "https://api.thecatapi.com/v1/votes"

	payload, err := json.Marshal(vote)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to encode vote data"}
		c.ServeJSON()
		return
	}

	resultChan := make(chan string)
	errChan := make(chan error)

	go func() {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			errChan <- err
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			errChan <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			return
		}

		resultChan <- "Vote submitted successfully"
	}()

	select {
	case message := <-resultChan:
		c.Data["json"] = map[string]string{"message": message}
	case err := <-errChan:
		c.Data["json"] = map[string]string{"error": "Failed to submit vote: " + err.Error()}
	}

	c.ServeJSON()
}

// AddToFavorites adds a cat image to favorites
func (c *MainController) AddToFavorites() {
	var favorite models.Favorite
	body, err := io.ReadAll(c.Ctx.Request.Body)
	fmt.Println(string(body))
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Unable to read request body"}
		c.ServeJSON()
		return
	}

	if err := json.Unmarshal(body, &favorite); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON format. Ensure 'image_id' is included."}
		c.ServeJSON()
		return
	}

	if favorite.ImageID == "" {
		c.Data["json"] = map[string]string{"error": "Missing 'image_id' in the request body"}
		c.ServeJSON()
		return
	}

	apiKey, _ := beego.AppConfig.String("catapi_key")
	url := "https://api.thecatapi.com/v1/favourites"

	payload := map[string]string{"image_id": favorite.ImageID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to encode favorite data"}
		c.ServeJSON()
		return
	}

	resultChan := make(chan string)
	errChan := make(chan error)

	go func() {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			errChan <- err
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			errChan <- fmt.Errorf("external API error: %s", string(bodyBytes))
			return
		}

		resultChan <- "Added to favorites successfully"
	}()

	select {
	case message := <-resultChan:
		c.Data["json"] = map[string]string{"message": message}
	case err := <-errChan:
		c.Data["json"] = map[string]string{"error": "Failed to add to favorites: " + err.Error()}
	}

	c.ServeJSON()
}

// FetchFavorites fetches favorite cat images
func (c *MainController) FetchFavorites() {
	apiKey, _ := beego.AppConfig.String("catapi_key")
	url := "https://api.thecatapi.com/v1/favourites"

	resultChan := make(chan []map[string]interface{})
	errChan := make(chan error)

	go func() {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		var result []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			errChan <- err
			return
		}

		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		c.Data["json"] = result
	case err := <-errChan:
		c.Data["json"] = map[string]string{"error": "Failed to fetch favorites: " + err.Error()}
	}

	c.ServeJSON()
}

func (c *MainController) FetchNewCatImage() {
    apiKey, _ := beego.AppConfig.String("catapi_key")
    url := "https://api.thecatapi.com/v1/images/search"

    // Channel to handle responses
    resultChan := make(chan map[string]interface{})
    errorChan := make(chan error)

    // Goroutine for API request
    go func() {
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            errorChan <- err
            return
        }
        req.Header.Add("x-api-key", apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errorChan <- err
            return
        }
        defer resp.Body.Close()

        var result []map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            errorChan <- err
            return
        }

        if len(result) > 0 {
            resultChan <- result[0]
        } else {
            errorChan <- fmt.Errorf("no image found")
        }
    }()

    // Handle responses via select
    select {
    case result := <-resultChan:
        c.Data["json"] = result
    case err := <-errorChan:
        c.Data["json"] = map[string]string{"error": err.Error()}
    }

    c.ServeJSON()
}



// package controllers

// import (
// 	"bytes"
// 	"catapigo/models"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"

// 	beego "github.com/beego/beego/v2/server/web"
// )

// type MainController struct {
// 	beego.Controller
// }

// // InitialData represents the structure for all initial data
// type InitialData struct {
// 	Breeds    []map[string]interface{} `json:"breeds"`
// 	Favorites []map[string]interface{} `json:"favorites"`
// 	NewCat    map[string]interface{}   `json:"newCat"`
// 	Error     map[string]string        `json:"error,omitempty"`
// }

// func (c *MainController) Prepare() {
// 	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
// 	c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

// 	if c.Ctx.Request.Method == "OPTIONS" {
// 		c.Ctx.Output.SetStatus(200)
// 		c.StopRun()
// 	}
// }

// func (c *MainController) ServeFrontend() {
// 	c.TplName = "index.tpl"
// }

// // FetchInitialData fetches all required data concurrently
// func (c *MainController) FetchInitialData() {
// 	apiKey, _ := beego.AppConfig.String("catapi_key")

// 	var wg sync.WaitGroup
// 	var data InitialData
// 	var mu sync.Mutex

// 	errorChan := make(chan error, 3)

// 	// Fetch breeds
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		breeds, err := fetchBreeds(apiKey)
// 		if err != nil {
// 			errorChan <- fmt.Errorf("breeds error: %v", err)
// 			return
// 		}
// 		mu.Lock()
// 		data.Breeds = breeds
// 		mu.Unlock()
// 	}()

// 	// Fetch favorites
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		favorites, err := fetchFavorites(apiKey)
// 		if err != nil {
// 			errorChan <- fmt.Errorf("favorites error: %v", err)
// 			return
// 		}
// 		mu.Lock()
// 		data.Favorites = favorites
// 		mu.Unlock()
// 	}()

// 	// Fetch new cat
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		newCat, err := fetchNewCat(apiKey)
// 		if err != nil {
// 			errorChan <- fmt.Errorf("new cat error: %v", err)
// 			return
// 		}
// 		mu.Lock()
// 		data.NewCat = newCat
// 		mu.Unlock()
// 	}()

// 	// Wait for all goroutines to complete
// 	go func() {
// 		wg.Wait()
// 		close(errorChan)
// 	}()

// 	// Collect any errors
// 	var errors []string
// 	for err := range errorChan {
// 		errors = append(errors, err.Error())
// 	}

// 	if len(errors) > 0 {
// 		data.Error = map[string]string{"message": fmt.Sprintf("Errors occurred: %v", errors)}
// 	}

// 	c.Data["json"] = data
// 	c.ServeJSON()
// }

// // FetchCatBreeds fetches all available cat breeds
// func (c *MainController) FetchCatBreeds() {
// 	breeds, err := fetchBreeds(beego.AppConfig.DefaultString("catapi_key", ""))
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": "Failed to fetch breeds: " + err.Error()}
// 	} else {
// 		c.Data["json"] = breeds
// 	}
// 	c.ServeJSON()
// }

// func (c *MainController) FetchBreedImages() {
// 	breedID := c.GetString("breed_id")
// 	apiKey, _ := beego.AppConfig.String("catapi_key")
// 	url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?breed_ids=%s&limit=5", breedID)

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.ServeJSON()
// 		return
// 	}
// 	req.Header.Add("x-api-key", apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.ServeJSON()
// 		return
// 	}
// 	defer resp.Body.Close()

// 	var result []map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 	} else {
// 		c.Data["json"] = result
// 	}
// 	c.ServeJSON()
// }

// // VoteForCat handles voting for a cat image
// func (c *MainController) VoteForCat() {
// 	var vote models.Vote
// 	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vote); err != nil {
// 		c.Data["json"] = map[string]string{"error": "Invalid input"}
// 		c.ServeJSON()
// 		return
// 	}

// 	apiKey, _ := beego.AppConfig.String("catapi_key")
// 	url := "https://api.thecatapi.com/v1/votes"

// 	payload, err := json.Marshal(vote)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": "Failed to encode vote data"}
// 		c.ServeJSON()
// 		return
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.ServeJSON()
// 		return
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("x-api-key", apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.ServeJSON()
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
// 		c.Data["json"] = map[string]string{"error": "Failed to submit vote"}
// 	} else {
// 		c.Data["json"] = map[string]string{"message": "Vote submitted successfully"}
// 	}
// 	c.ServeJSON()
// }

// // AddToFavorites adds a cat image to favorites
// func (c *MainController) AddToFavorites() {
// 	var favorite models.Favorite
// 	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&favorite); err != nil {
// 		c.Data["json"] = map[string]string{"error": "Invalid input"}
// 		c.ServeJSON()
// 		return
// 	}

// 	apiKey, _ := beego.AppConfig.String("catapi_key")
// 	url := "https://api.thecatapi.com/v1/favourites"

// 	payload, err := json.Marshal(map[string]string{"image_id": favorite.ImageID})
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": "Failed to encode favorite data"}
// 		c.ServeJSON()
// 		return
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.ServeJSON()
// 		return
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("x-api-key", apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.ServeJSON()
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
// 		c.Data["json"] = map[string]string{"error": "Failed to add to favorites"}
// 	} else {
// 		c.Data["json"] = map[string]string{"message": "Added to favorites successfully"}
// 	}
// 	c.ServeJSON()
// }

// // Helper functions for API calls
// func fetchBreeds(apiKey string) ([]map[string]interface{}, error) {
// 	url := "https://api.thecatapi.com/v1/breeds"
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("x-api-key", apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var result []map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func fetchFavorites(apiKey string) ([]map[string]interface{}, error) {
// 	url := "https://api.thecatapi.com/v1/favourites"
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("x-api-key", apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var result []map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func fetchNewCat(apiKey string) (map[string]interface{}, error) {
// 	url := "https://api.thecatapi.com/v1/images/search"
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("x-api-key", apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var result []map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}

// 	if len(result) == 0 {
// 		return nil, fmt.Errorf("no image found")
// 	}
// 	return result[0], nil
// }

// func (c *MainController) FetchFavorites() {
//     favorites, err := fetchFavorites(beego.AppConfig.DefaultString("catapi_key", ""))
//     if err != nil {
//         c.Data["json"] = map[string]string{"error": err.Error()}
//     } else {
//         c.Data["json"] = favorites
//     }
//     c.ServeJSON()
// }

// func (c *MainController) FetchNewCatImage() {
// 	newCat, err := fetchNewCat(beego.AppConfig.DefaultString("catapi_key", ""))
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 	} else {
// 		c.Data["json"] = newCat
// 	}
// 	c.ServeJSON()
// }