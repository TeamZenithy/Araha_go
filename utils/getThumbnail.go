package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/TeamZenithy/Araha/logger"
	"github.com/rojters/opengraph"
)

//GetYTThumbnail returns youtube video's thumbnail
func GetYTThumbnail(urlVideo string) (string, error) {
	const (
		vi     = "https://i.ytimg.com/vi/"
		resMax = "/maxresdefault.jpg"
	)

	equalIndex := strings.Index(urlVideo, "=")
	ampIndex := strings.Index(urlVideo, "&")
	slash := strings.LastIndex(urlVideo, "/")
	questionIndex := strings.Index(urlVideo, "?")
	var id string

	if equalIndex != -1 {

		if ampIndex != -1 {
			id = urlVideo[equalIndex+1 : ampIndex]
		} else if questionIndex != -1 && strings.Contains(urlVideo, "?t=") {
			id = urlVideo[slash+1 : questionIndex]
		} else {
			id = urlVideo[equalIndex+1:]
		}

	} else {
		id = urlVideo[slash+1:]
	}

	if strings.ContainsAny(id, "?&/<%=") {
		return "", errors.New("invalid characters in video id")
	}
	if len(id) < 10 {
		return "", errors.New("the video id must be at least 10 characters long")
	}

	return vi + id + resMax, nil
}

//GetSCThumbnail returns soundcloud song's thumbnail
func GetSCThumbnail(urlSong string) (string, error) {
	metaData := fetchMetaData(urlSong)
	artwork := extractMetaData(metaData)
	if artwork == "" {
		logger.Warn(fmt.Sprintf("Failed to parse soundcloud thumbnail: %s", urlSong))
		return "", nil
	}
	return artwork, nil
}
func fetchMetaData(url string) []opengraph.MetaData {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	md, err := opengraph.Extract(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return md
}

func extractMetaData(md []opengraph.MetaData) string {
	result := ""
	for i := range md {
		switch md[i].Property {
		case "image":
			result = md[i].Content
		}
	}

	return result
}
