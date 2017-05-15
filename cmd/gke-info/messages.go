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
