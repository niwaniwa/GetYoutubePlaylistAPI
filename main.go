package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	service  *youtube.Service
	interval = 2
	burst    = 2
	limit    = rate.NewLimiter(rate.Every(time.Duration(interval)*time.Second), burst)
)

func main() {
	err := godotenv.Load("value.env")

	apikey := os.Getenv("apikey")

	if err != nil {
		panic(err)
	}

	service, err = youtube.NewService(context.Background(), option.WithAPIKey(apikey))
	if err != nil {
		panic(err)
	}

	engine := gin.Default()
	engine.Use(rateLimiter())
	engine.GET("/playlist/:id", func(c *gin.Context) {
		// パラメータからIDを取得
		id := c.Param("id")

		c.JSON(200, GetPlaylist(id))
	})

	engine.Run()
}

func GetPlaylist(id string) YoutubePlaylistDataResponse {
	var videoDataList []VideoData

	request := service.Playlists.List([]string{"snippet"}).Id(id)

	do, err := request.Do()
	if err != nil {
		panic(err)
		return YoutubePlaylistDataResponse{}
	}

	videoDataList = GetVideo(id, "", nil)

	if videoDataList == nil {
		videoDataList = []VideoData{}
	}

	var name string
	if len(do.Items) == 0 {
		return YoutubePlaylistDataResponse{}
	} else if do.Items[0].Snippet.Title == "" {
		name = id
	} else {
		name = do.Items[0].Snippet.Title
	}

	data := YoutubePlaylistDataResponse{
		Name:   name,
		Videos: videoDataList,
	}

	return data
}

func GetVideo(id string, nextPageToken string, videoDataListSource []VideoData) []VideoData {
	newRequest := service.PlaylistItems.List([]string{"snippet"}).PlaylistId(id).MaxResults(50).PageToken(nextPageToken)

	var videoDataList []VideoData

	if videoDataListSource != nil {
		videoDataList = videoDataListSource
	}

	response, err := newRequest.Do()
	if err != nil {
		return videoDataList
	}

	log.Print(len(response.Items))

	for _, item := range response.Items {
		if item.Snippet.Title == "Private video" || item.Snippet.Title == "Deleted video" {
			continue
		}

		videoData := VideoData{
			Id:    item.Snippet.ResourceId.VideoId,
			Title: item.Snippet.Title,
		}
		log.Print(videoData.Title)
		videoDataList = append(videoDataList, videoData)
	}

	if response.NextPageToken == "" {
		return videoDataList
	}

	return GetVideo(id, response.NextPageToken, videoDataList)
}

func rateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if limit.Allow() == false {
			c.JSON(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
			c.Abort()
		}
	}
}
