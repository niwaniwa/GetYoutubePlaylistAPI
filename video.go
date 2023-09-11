package main

type VideoData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type YoutubePlaylistDataResponse struct {
	Name   string      `json:"name"`
	Videos []VideoData `json:"videos"`
}
