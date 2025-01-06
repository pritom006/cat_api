package controllers

import (
	"bytes"
	"catapigo/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "sort"
	"time"

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
		//c.StopRun()
		c.Ctx.WriteString("")
		c.Ctx.ResponseWriter.Flush()
		return
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



func (c *MainController) FetchFavorites() {
    apiKey, _ := beego.AppConfig.String("catapi_key")
    // Add sub_id to the URL to filter favorites
    url := "https://api.thecatapi.com/v1/favourites?sub_id=test"

    resultChan := make(chan []models.FavoriteResponse)
    errChan := make(chan error)

    go func() {
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            errChan <- err
            return
        }
        
        req.Header.Add("x-api-key", apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errChan <- err
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            bodyBytes, _ := io.ReadAll(resp.Body)
            errChan <- fmt.Errorf("external API error (status %d): %s", resp.StatusCode, string(bodyBytes))
            return
        }

        var favorites []models.FavoriteResponse
        if err := json.NewDecoder(resp.Body).Decode(&favorites); err != nil {
            errChan <- fmt.Errorf("failed to decode response: %v", err)
            return
        }

        // Add additional validation
        if len(favorites) == 0 {
            resultChan <- []models.FavoriteResponse{}
            return
        }

        resultChan <- favorites
    }()

    select {
    case result := <-resultChan:
        c.Data["json"] = result
    case err := <-errChan:
        c.Data["json"] = map[string]string{"error": err.Error()}
    case <-time.After(10 * time.Second):
        c.Data["json"] = map[string]string{"error": "Request timed out"}
    }

    c.ServeJSON()
}

func (c *MainController) AddToFavorites() {
    var favorite models.Favorite
    body, err := io.ReadAll(c.Ctx.Request.Body)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Unable to read request body"}
        c.ServeJSON()
        return
    }

    if err := json.Unmarshal(body, &favorite); err != nil {
        c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
        c.ServeJSON()
        return
    }

    if favorite.ImageID == "" {
        c.Data["json"] = map[string]string{"error": "Missing image_id"}
        c.ServeJSON()
        return
    }

    apiKey, _ := beego.AppConfig.String("catapi_key")
    url := "https://api.thecatapi.com/v1/favourites"

    // Use your existing Favorite model
    payload := models.Favorite{
        ImageID: favorite.ImageID,
        SubID:   "test", // Set a default if not provided
    }
    if favorite.SubID != "" {
        payload.SubID = favorite.SubID
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to encode data"}
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

        body, _ := io.ReadAll(resp.Body)
        if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
            errChan <- fmt.Errorf("API error: %s", string(body))
            return
        }

        resultChan <- "Added to favorites successfully"
    }()

    select {
    case message := <-resultChan:
        c.Data["json"] = map[string]string{"message": message}
    case err := <-errChan:
        c.Data["json"] = map[string]string{"error": err.Error()}
    case <-time.After(5 * time.Second):
        c.Data["json"] = map[string]string{"error": "Request timed out"}
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
