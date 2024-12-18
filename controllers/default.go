// package controllers

// import (
// 	beego "github.com/beego/beego/v2/server/web"
// )

// type MainController struct {
// 	beego.Controller
// }

// func (c *MainController) Get() {
// 	c.Data["Website"] = "beego.vip"
// 	c.Data["Email"] = "astaxie@gmail.com"
// 	c.TplName = "index.tpl"
// }

package controllers

import (
	"bytes"
	"catapigo/models"
	"encoding/json"
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

	// Marshal the Vote object into JSON
	payload, err := json.Marshal(vote)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to encode vote data"}
		c.ServeJSON()
		return
	}

	// Create the request with the payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to submit vote"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Send success response
	c.Data["json"] = map[string]string{"message": "Vote submitted successfully"}
	c.ServeJSON()
}

// FetchFavorites retrieves the list of favorite cat images
func (c *MainController) FetchFavorites() {
	apiKey, _ := beego.AppConfig.String("catapi_key")
	url := "https://api.thecatapi.com/v1/favourites"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch favorites"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	c.Data["json"] = result
	c.ServeJSON()
}
