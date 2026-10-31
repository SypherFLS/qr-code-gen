package dto

type CreateQRRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type CreateQRResponse struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
