// package routers

// import (
// 	"catapigo/controllers"
// 	beego "github.com/beego/beego/v2/server/web"
// )

// func init() {
//     beego.Router("/", &controllers.MainController{})
// }





package routers

import (
	"catapigo/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Route for rendering the home page
	beego.Router("/", &controllers.MainController{})

	// Route for fetching cat breeds
	beego.Router("/fetch-breeds", &controllers.MainController{}, "get:FetchCatBreeds")

	// Route for submitting a vote (like/dislike)
	beego.Router("/vote", &controllers.MainController{}, "post:VoteForCat")

	// Route for fetching favorite cats
	beego.Router("/favorites", &controllers.MainController{}, "get:FetchFavorites")
}