package routers

import (
	"catapigo/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:ServeFrontend")
	beego.Router("/fetch-breeds", &controllers.MainController{}, "get:FetchCatBreeds")
	beego.Router("/vote", &controllers.MainController{}, "post:VoteForCat")
	beego.Router("/favorites", &controllers.MainController{}, "get:FetchFavorites")
	beego.Router("/addToFavorites", &controllers.MainController{}, "post:AddToFavorites")
	beego.Router("/fetch-new-cat", &controllers.MainController{}, "get:FetchNewCatImage")
	beego.Router("/fetch-breed-images", &controllers.MainController{}, "get:FetchBreedImages")
}
