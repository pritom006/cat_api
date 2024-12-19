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
	body, err := io.ReadAll(c.Ctx.Request.Body)
	fmt.Println(body)
	fmt.Println(err)
	// err := json.Unmarshal(c.Ctx.Input.RequestBody, &favorite)
	if err != nil {
		// Log the error and payload for debugging
		fmt.Printf("Failed to parse payload: %s\n", string(c.Ctx.Input.RequestBody))
		c.Data["json"] = map[string]string{"error": "Invalid input. Ensure 'image_id' is included in the request body."}
		c.ServeJSON()
		return
	}

	apiKey, _ := beego.AppConfig.String("catapi_key")
	subId := c.GetString("sub_id")
	fmt.Println(subId)
	url := "https://api.thecatapi.com/v1/favourites?sub_id="+ subId

	// Create the payload for TheCatAPI
	payload := map[string]string{
		"image_id": favorite.ImageID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to encode favorite data"}
		c.ServeJSON()
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
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
		c.Data["json"] = map[string]string{"error": "Failed to submit favorite"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Data["json"] = map[string]string{"error": "Failed to add to favorites. External API error."}
		c.ServeJSON()
		return
	}

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
