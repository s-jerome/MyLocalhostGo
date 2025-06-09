package youtube

import (
	"database/sql"
	"fmt"
	"mylocalhost/config"
	database "mylocalhost/utils/database"
	dates "mylocalhost/utils/dates"

	_ "github.com/mattn/go-sqlite3"
)

type RatedVideo struct {
	Rowid   int64  `json:"-"`
	VideoId string `json:"videoId"`
	Rating  string `json:"rating"`
}

var _connection *sql.DB

// Cache of the videos I just rated. Not very usefull I know.
var _videosByVideoId = make(map[string]*RatedVideo)
var _channelIdsByName = make(map[string]int64)

func openConnection() error {
	if _connection != nil {
		return nil
	}

	var dbFilePath = config.Get("Youtube.ratedVideos.databaseFilePath")

	var connection, dbFileExists, connectionErr = database.OpenSQLiteConnection(dbFilePath)
	if connectionErr != nil {
		return connectionErr
	}
	_connection = connection

	if dbFileExists == false {
		var _, execErr = connection.Exec(`
		CREATE TABLE "channels" ("id" INTEGER, "name" TEXT NOT NULL CHECK("name" != '') UNIQUE, "channel_id" TEXT NOT NULL CHECK("channel_id" != ''),
		PRIMARY KEY("id"));
		
		CREATE TABLE "videos" ("video_id" TEXT NOT NULL CHECK("video_id" != '') UNIQUE, "rating" TEXT NOT NULL CHECK("rating" != ''), 
		"channel_id" INTEGER NOT NULL, "title" TEXT NOT NULL CHECK("title" != ''),
		"description" TEXT NOT NULL DEFAULT '', "comment" TEXT NOT NULL DEFAULT '',
		"created_at" TEXT NOT NULL CHECK("created_at" != ''), "updated_at" TEXT NOT NULL DEFAULT '',
		"downloaded_at" TEXT NOT NULL DEFAULT '', "deleted_at" TEXT NOT NULL DEFAULT '',
		FOREIGN KEY("channel_id") REFERENCES "channels"("id") ON DELETE RESTRICT ON UPDATE CASCADE);
		
		CREATE INDEX "idx_channels_name" ON "channels" ("name");
		CREATE INDEX "idx_videos_video_id" ON "videos" ("video_id");`)
		if execErr != nil {
			return execErr
		}
	}

	return nil
}

func GetRatedVideos() ([]RatedVideo, error) {
	if openErr := openConnection(); openErr != nil {
		return nil, openErr
	}

	var stmt, stmtErr = _connection.Prepare("SELECT video_id, rating FROM videos;")
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer stmt.Close()

	var rows, queryErr = stmt.Query()
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	var ratedVideos []RatedVideo
	for rows.Next() {
		var ratedVideo = RatedVideo{}
		if scanErr := rows.Scan(&ratedVideo.VideoId, &ratedVideo.Rating); scanErr != nil {
			return nil, scanErr
		}
		ratedVideos = append(ratedVideos, ratedVideo)
	}

	return ratedVideos, nil
}

// Insert or update the rating for a video.
func SetVideoRating(videoId string, rating string, channelName string, videoTitle string, channelId string, videoDescription string, videoDurationSeconds int64) error {
	if openErr := openConnection(); openErr != nil {
		return openErr
	}

	var ratedVideo, err = getVideoFromVideoId(videoId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = insertVideo(videoId, rating, channelName, videoTitle, channelId, videoDescription, videoDurationSeconds)
		}
	} else if ratedVideo.Rating != rating {
		err = updateRating(ratedVideo, rating)
	}
	return err
}

func getVideoFromVideoId(videoid string) (*RatedVideo, error) {
	var video, keyExists = _videosByVideoId[videoid]
	if keyExists {
		return video, nil
	}

	var stmt, stmtErr = _connection.Prepare("SELECT rowid, rating FROM videos WHERE video_id = ?")
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer stmt.Close()

	var ratedVideo = &RatedVideo{}
	var scanErr = stmt.QueryRow(videoid).Scan(&ratedVideo.Rowid, &ratedVideo.Rating)
	if scanErr == nil {
		_videosByVideoId[videoid] = ratedVideo
	}
	return ratedVideo, scanErr
}

func insertVideo(videoid string, rating string, channelName string, videoTitle string, channelId string, videoDescription string, videoDurationSeconds int64) error {
	var channelRowid, err = getChannelIdByName(channelName)
	if err != nil {
		if err == sql.ErrNoRows {
			channelRowid, err = insertChannel(channelName, channelId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	var stmt, stmtErr = _connection.Prepare("INSERT INTO videos(video_id, rating, channel_id, title, description, duration_seconds, created_at) VALUES(?, ?, ?, ?, ?, ?, ?);")
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	var result, execErr = stmt.Exec(videoid, rating, channelRowid, videoTitle, videoDescription, videoDurationSeconds, dates.NowToString())
	_ = result
	return execErr
}

func getChannelIdByName(name string) (int64, error) {
	var channelId, keyExists = _channelIdsByName[name]
	if keyExists {
		return channelId, nil
	}

	var stmt, stmtErr = _connection.Prepare("SELECT id FROM channels WHERE name = ?")
	if stmtErr != nil {
		return 0, stmtErr
	}
	defer stmt.Close()

	var scanErr = stmt.QueryRow(name).Scan(&channelId)
	return channelId, scanErr
}

func insertChannel(name string, id string) (int64, error) {
	var stmt, stmtErr = _connection.Prepare("INSERT INTO channels(name, channel_id) VALUES(?, ?);")
	if stmtErr != nil {
		return 0, stmtErr
	}
	defer stmt.Close()

	var result, execErr = stmt.Exec(name, id)
	if execErr != nil {
		return 0, execErr
	}
	var lastInsertId, _ = result.LastInsertId()
	_channelIdsByName[name] = lastInsertId
	return lastInsertId, nil
}

func updateRating(ratedVideo *RatedVideo, rating string) error {
	var stmt, stmtErr = _connection.Prepare("UPDATE videos SET rating = ?, updated_at = ? WHERE rowid = ?;")
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	var result, execErr = stmt.Exec(rating, dates.NowToString(), ratedVideo.Rowid)
	if execErr != nil {
		return execErr
	}
	var rowsAffected, _ = result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("The update of the rating \"%s\" for the rowid %d has affected no rows", rating, ratedVideo.Rowid)
	} else if rowsAffected > 1 {
		return fmt.Errorf("The update of the rating \"%s\" for the rowid %d has affected %d rows", rating, ratedVideo.Rowid, rowsAffected)
	}
	ratedVideo.Rating = rating
	return nil
}

func CloseDatabaseConnection() {
	if _connection != nil {
		_connection.Close()
	}
}
