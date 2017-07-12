package main

type versionResponse struct {
	Version string `json:"version"`
}

type metadataResponse struct {
	Instance *Instance `json:"instance"`
}

type healthResponse struct {
	Status int `json:"status"`
}

type homeResponse struct {
	Document string `json:"document"`
}

type errorRequest struct {
	Message string `json:"message"`
}
