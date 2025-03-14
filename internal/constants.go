package internal

type StatusResponse struct {
	Status string `json:"status"`
}

type InfoResponse struct {
	*StatusResponse
	Description string `json:"description"`
}
