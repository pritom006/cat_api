package models

type Favorite struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id,omitempty"` 
}

type FavoriteResponse struct {
    ID        int    `json:"id"`
    ImageID   string `json:"image_id"`
    SubID     string `json:"sub_id"`
    CreatedAt string `json:"created_at"`
    Image     struct {
        ID  string `json:"id"`
        URL string `json:"url"`
    } `json:"image"`
}