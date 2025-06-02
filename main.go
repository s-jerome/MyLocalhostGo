package main

import (
	"mylocalhost/config"
	"mylocalhost/logger"
	youtube_ratedVideos "mylocalhost/sites/Youtube/ratedvideos"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if chdirErr := setChdir(); chdirErr != nil {
		os.Exit(1)
		return
	}

	if readConfigError := config.Read(); readConfigError != nil {
		logger.WriteError("Error reading the config\n%v", readConfigError)
		os.Exit(1)
	}

	var server = http.NewServeMux()
	server.HandleFunc("/youtube/rating/get-rated-videos", youtube_ratedVideos.GetRatedVideosRequestHandler)
	server.HandleFunc("/youtube/rating/set-video-rating", youtube_ratedVideos.SetVideoRatingRequestHandler)
	var serverPort = config.Get("server.port")
	var err = http.ListenAndServe(":"+serverPort, server)

	youtube_ratedVideos.CloseDatabaseConnection()

	if err != nil {
		logger.WriteError("Error listening at port "+serverPort+"\n%v", err)
		os.Exit(1)
	}
}

// Set the current working directory to the one where the current executable is.
func setChdir() error {
	var executableFilePath, executableErr = os.Executable()
	if executableErr != nil {
		return executableErr
	}
	var executableDirectoryPath = filepath.Dir(executableFilePath)
	var chdirErr = os.Chdir(executableDirectoryPath)
	if chdirErr != nil {
		return chdirErr
	}
	return nil
}
