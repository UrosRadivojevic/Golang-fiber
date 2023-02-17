package requests

type CreateMovieRequest struct {
	Movie    string `json:"movie" validate:"required"`
	Watched  bool   `json:"watched" validate:"required"`
	Year     int    `json:"year" validate:"required"`
	LeadRole string `json:"leadrole" validate:"required"`
}
