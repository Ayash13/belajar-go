package dto

type PersonCreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PersonResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
