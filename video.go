package main

type VideoData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Id          string `json:"id"`
}

type YoutubePlaylistDataResponse struct {
	Name   string      `json:"name"`
	Videos []VideoData `json:"videos"`
}
