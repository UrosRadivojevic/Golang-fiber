package requests

type CreateMovieRequest struct {
	Movie    string `json:"movie"`
	Watched  bool   `json:"watched"`
	Year     int    `json:"year"`
	LeadRole string `json:"leadrole"`
}
