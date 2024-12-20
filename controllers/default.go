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

// Render the main page
func (c *MainController) Get() {
	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

// FetchCatBreeds fetches all available cat breeds from TheCatAPI
func (c *MainController) FetchCatBreeds() {
	apiKey, _ := beego.AppConfig.String("catapi_key")
	url := "https://api.thecatapi.com/v1/breeds"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch breeds"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	c.Data["json"] = result
	c.ServeJSON()
}

// VoteForCat handles voting (like/dislike) for a cat image
func (c *MainController) VoteForCat() {
	var vote models.Vote
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &vote)
	if err != nil {
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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to submit vote"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	c.Data["json"] = map[string]string{"message": "Vote submitted successfully"}
	c.ServeJSON()
}

func (c *MainController) FetchFavorites() {
	apiKey, _ := beego.AppConfig.String("catapi_key")
	url := "https://api.thecatapi.com/v1/favourites"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("Cache-Control", "no-cache") // Prevent caching
	req.Header.Add("Pragma", "no-cache")       // Ensure fresh data

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch favorites"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to decode API response"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = result
	c.ServeJSON()
}




func (c *MainController) AddToFavorites() {
	// Parse the incoming JSON payload
	var favorite models.Favorite
	body, err := io.ReadAll(c.Ctx.Request.Body) // Read the request body
	if err != nil {
		fmt.Printf("Failed to read request body: %v\n", err)
		c.Data["json"] = map[string]string{"error": "Unable to read request body"}
		c.ServeJSON()
		return
	}

	fmt.Printf("Raw Request Body: %s\n", string(body)) // Debugging the raw payload

	// Unmarshal the JSON into the Favorite struct
	err = json.Unmarshal(body, &favorite)
	if err != nil {
		fmt.Printf("Failed to parse JSON payload: %v\n", err)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format. Ensure 'image_id' is included."}
		c.ServeJSON()
		return
	}

	fmt.Printf("Parsed Favorite Struct: %+v\n", favorite) // Debugging the parsed struct

	// Validate the parsed input
	if favorite.ImageID == "" {
		c.Data["json"] = map[string]string{"error": "Missing 'image_id' in the request body"}
		c.ServeJSON()
		return
	}

	// Retrieve API key and sub_id
	apiKey, _ := beego.AppConfig.String("catapi_key")
	subID := c.GetString("sub_id") // Optional sub_id parameter
	fmt.Printf("Sub ID: %s\n", subID) // Debugging sub_id

	// Construct the API URL
	url := "https://api.thecatapi.com/v1/favourites"
	if subID != "" {
		url += "?sub_id=" + subID
	}

	// Create the payload for TheCatAPI
	payload := map[string]string{"image_id": favorite.ImageID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal payload: %v\n", err)
		c.Data["json"] = map[string]string{"error": "Failed to encode favorite data"}
		c.ServeJSON()
		return
	}

	fmt.Printf("Payload Sent to TheCatAPI: %s\n", string(payloadBytes)) // Debugging payload

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %v\n", err)
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		c.Data["json"] = map[string]string{"error": "Failed to submit favorite"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read error response body for debugging
		fmt.Printf("TheCatAPI Error Response: %s\n", string(bodyBytes))
		c.Data["json"] = map[string]string{"error": "Failed to add to favorites. External API error."}
		c.ServeJSON()
		return
	}

	// Successfully added to favorites
	c.Data["json"] = map[string]string{"message": "Added to favorites successfully"}
	c.ServeJSON()
}

// FetchNewCatImage fetches a random cat image
func (c *MainController) FetchNewCatImage() {
	apiKey, _ := beego.AppConfig.String("catapi_key")
	url := "https://api.thecatapi.com/v1/images/search"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch new image"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result) > 0 {
		c.Data["json"] = result[0]
	} else {
		c.Data["json"] = map[string]string{"error": "No image found"}
	}
	c.ServeJSON()
}
