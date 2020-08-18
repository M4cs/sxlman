package main

// Mods Struct for Mod Get Req
type Mods struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	LastUpdate int64  `json:"date_updated"`
	Modfile    struct {
		Version      string `json:"version"`
		DownloadInfo struct {
			URL string `json:"binary_url"`
		} `json:"download"`
	} `json:"modfile"`
}

// Results Struct for Mods RESULTS
type Results struct {
	Data []Mods `json:"data"`
}

// Config struct to define environment variables
type Config struct {
	APIKey          string `json:"apikey"`
	TrackedPackages []Mods `json:"tracking"`
	DownloadFolder  string `json:"downloadFolder"`
	AutoUpdate      bool   `json:"autoUpdate"`
	GamePath        string `json:"gamePath"`
}
