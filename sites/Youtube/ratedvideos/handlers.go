package youtube

import (
	"bytes"
	"encoding/json"
	"io"
	responses "mylocalhost/utils/responses"
	"net/http"
	"strconv"
)

// My Chrome extension wants to get all the videos and their rating from database.
func GetRatedVideosRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var videos, videosErr = GetRatedVideos()
	if videosErr != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, videosErr, "Getting the rated videos from database")
		return
	}

	var buffer bytes.Buffer
	if encodeErr := json.NewEncoder(&buffer).Encode(videos); encodeErr == nil {
		buffer.WriteTo(w)
	} else {
		responses.SendErrorResponse(w, http.StatusInternalServerError, encodeErr, "Encoding the rated videos in JSON")
	}
}

// I just rated a Youtube video. My Chrome extension intercepted that and sent some video data to save them in database.
func SetVideoRatingRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The request must be POST")
		return
	}

	var requestBody, requestBodyErr = io.ReadAll(r.Body)
	if requestBodyErr != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, requestBodyErr, "Reading POST data")
		return
	}
	if len(requestBody) == 0 {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The POST data is empty")
		return
	}

	var postData map[string]interface{}
	if parseErr := json.Unmarshal(requestBody, &postData); parseErr != nil {
		responses.SendErrorResponse(w, http.StatusBadRequest, parseErr, "Parsing the POST data to JSON")
		return
	}

	videoIdValue, keyExists := postData["videoId"]
	if keyExists == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "No videoId given")
		return
	}
	videoId, typeOk := videoIdValue.(string)
	if typeOk == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The videoId is not a string")
		return
	}
	if videoId == "" {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The videoId is empty")
		return
	}

	ratingValue, keyExists := postData["rating"]
	if keyExists == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "No rating given")
		return
	}
	rating, typeOk := ratingValue.(string)
	if typeOk == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The rating is not a string")
		return
	}
	if rating == "" {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The rating is empty")
		return
	}
	if rating != "like" && rating != "dislike" && rating != "none" {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The rating is invalid (should be either like/dislike/none)")
		return
	}

	channelNameValue, keyExists := postData["channelName"]
	if keyExists == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "No channelName given")
		return
	}
	channelName, typeOk := channelNameValue.(string)
	if typeOk == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The channelName is not a string")
		return
	}
	if channelName == "" {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The channelName is empty")
		return
	}

	videoTitleValue, keyExists := postData["videoTitle"]
	if keyExists == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "No videoTitle given")
		return
	}
	videoTitle, typeOk := videoTitleValue.(string)
	if typeOk == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "the videoTitle is not a string")
		return
	}
	if videoTitle == "" {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The videoTitle is empty")
		return
	}

	channelIdValue, keyExists := postData["channelId"]
	if keyExists == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "No channelId given")
		return
	}
	channelId, typeOk := channelIdValue.(string)
	if typeOk == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The channelId is not a string")
		return
	}
	if channelId == "" {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The channelId is empty")
		return
	}

	videoDescriptionValue, keyExists := postData["videoDescription"]
	if keyExists == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "No videoDescription given")
		return
	}
	videoDescription, typeOk := videoDescriptionValue.(string)
	if typeOk == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The videoDescription is not a string")
		return
	}

	videoDurationSecondsValue, keyExists := postData["videoDurationSeconds"]
	if keyExists == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "No videoDurationSeconds given")
		return
	}
	//.. The duration of the video should be sent as a string, because it's stored as a string by Youtube.
	videoDurationSecondsString, typeOk := videoDurationSecondsValue.(string)
	if typeOk == false {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The videoDurationSeconds is not a string")
		return
	}
	videoDurationSeconds, convErr := strconv.ParseInt(videoDurationSecondsString, 10, 64)
	if convErr != nil {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The videoDurationSeconds is not a integer")
		return
	}

	var sqlError = SetVideoRating(videoId, rating, channelName, videoTitle, channelId, videoDescription, videoDurationSeconds)
	if sqlError != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, sqlError, "Saving the rating in database")
	}
}
