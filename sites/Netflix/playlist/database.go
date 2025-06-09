package netflix

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mylocalhost/config"
	utils "mylocalhost/utils/database"
	dates "mylocalhost/utils/dates"

	_ "github.com/mattn/go-sqlite3"
)

type videoData struct {
	Rowid   int64
	VideoId int64
	Type    string
	Title   string

	Status string

	Casting   string
	Creators  string
	Directors string
	Writers   string

	Genres string
	Mood   string
	Tags   string

	AgeAdvised       int
	AgeAdvisedReason string

	Synopsis string

	SeasonCount    int
	NumSeasonLabel string
	EpisodeCount   int

	DurationSec int64

	AvailabilityStartTime string

	DataFrom_ string `json:"_dataFrom"`

	CreatedAt string
	UpdatedAt string
}

type saveVideoToPlaylistResult struct {
	Rowid int64  `json:"rowid"`
	Error string `json:"error"`
	// The kind of SQL request made (INSERT/UPDATE/NONE)
	Query string `json:"query"`

	UpdatedColumns []string `json:"updatedColumns"`
	OldValues      []any    `json:"oldValues"`
	NewValues      []any    `json:"newValues"`
}

var _connection *sql.DB

func openConnection() error {
	if _connection != nil {
		return nil
	}

	var dbFilePath = config.Get("Netflix.databaseFilePath")

	var connection, dbFileExists, connectionErr = utils.OpenSQLiteConnection(dbFilePath)
	if connectionErr != nil {
		return connectionErr
	}
	_connection = connection

	if dbFileExists == false {
		var _, execErr = connection.Exec(`
		CREATE TABLE "playlist" (
			"video_id"	INTEGER NOT NULL CHECK("video_id" > 0) UNIQUE,
			"type"	TEXT NOT NULL DEFAULT '',
			"title"	TEXT NOT NULL CHECK("title" != ''),
			"status"	TEXT NOT NULL CHECK("status" != ''),
			"casting"	TEXT NOT NULL DEFAULT '',
			"creators"	TEXT NOT NULL DEFAULT '',
			"directors"	TEXT NOT NULL DEFAULT '',
			"writers"	TEXT NOT NULL DEFAULT '',
			"genres"	TEXT NOT NULL DEFAULT '',
			"mood"	TEXT NOT NULL DEFAULT '',
			"tags"	TEXT NOT NULL DEFAULT '',
			"age_advised"	INTEGER NOT NULL DEFAULT 0,
			"age_advised_reason"	TEXT NOT NULL DEFAULT '',
			"synopsis"	TEXT NOT NULL DEFAULT '',
			"season_count"	INTEGER NOT NULL DEFAULT 0,
			"num_season_label"	TEXT NOT NULL DEFAULT '',
			"episode_count"	INTEGER NOT NULL DEFAULT 0,
			"duration_sec"	INTEGER NOT NULL DEFAULT 0,
			"availability_starttime"	TEXT NOT NULL DEFAULT '',
			"_data_from"	TEXT NOT NULL DEFAULT '',
			"created_at"	TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now', 'localtime')),
			"updated_at"	TEXT NOT NULL DEFAULT ''
		);
		
		CREATE TABLE "playlist_updates" (
			"video_id"	INTEGER NOT NULL,
			"updated_at"	TEXT NOT NULL,
			"updates"	TEXT NOT NULL
		);
		
		CREATE INDEX "idx_playlist_video_id" ON "playlist" ("video_id");`)
		if execErr != nil {
			return execErr
		}
	}

	return nil
}

// Insert or update the given video.
func saveVideoToPlaylist(videoToAdd *videoData) saveVideoToPlaylistResult {
	var result = saveVideoToPlaylistResult{}
	if openErr := openConnection(); openErr != nil {
		result.Error = openErr.Error()
		return result
	}

	var savedVideo, getVideoErr = getVideoFromVideoId(videoToAdd.VideoId)
	if getVideoErr != nil {
		if getVideoErr == sql.ErrNoRows {
			result.Query = "INSERT"
			var insertErr = insertVideo(videoToAdd)
			if insertErr == nil {
				result.Rowid = videoToAdd.Rowid
			} else {
				result.Error = insertErr.Error()
			}
		} else {
			result.Error = getVideoErr.Error()
		}
	} else {
		result.Rowid = savedVideo.Rowid
		videoToAdd.Rowid = savedVideo.Rowid

		var columnsToUpdate []string
		var oldValues []any
		var newValues []any

		if savedVideo.AgeAdvised != videoToAdd.AgeAdvised {
			columnsToUpdate = append(columnsToUpdate, "age_advised")
			oldValues = append(oldValues, savedVideo.AgeAdvised)
			newValues = append(newValues, videoToAdd.AgeAdvised)
		}
		if savedVideo.AgeAdvisedReason != videoToAdd.AgeAdvisedReason {
			columnsToUpdate = append(columnsToUpdate, "age_advised_reason")
			oldValues = append(oldValues, savedVideo.AgeAdvisedReason)
			newValues = append(newValues, videoToAdd.AgeAdvisedReason)
		}
		if savedVideo.AvailabilityStartTime != videoToAdd.AvailabilityStartTime {
			columnsToUpdate = append(columnsToUpdate, "availability_starttime")
			oldValues = append(oldValues, savedVideo.AvailabilityStartTime)
			newValues = append(newValues, videoToAdd.AvailabilityStartTime)
		}
		if savedVideo.Casting != videoToAdd.Casting {
			columnsToUpdate = append(columnsToUpdate, "casting")
			oldValues = append(oldValues, savedVideo.Casting)
			newValues = append(newValues, videoToAdd.Casting)
		}
		if savedVideo.Creators != videoToAdd.Creators {
			columnsToUpdate = append(columnsToUpdate, "creators")
			oldValues = append(oldValues, savedVideo.Creators)
			newValues = append(newValues, videoToAdd.Creators)
		}
		if savedVideo.Directors != videoToAdd.Directors {
			columnsToUpdate = append(columnsToUpdate, "directors")
			oldValues = append(oldValues, savedVideo.Directors)
			newValues = append(newValues, videoToAdd.Directors)
		}
		if savedVideo.DurationSec != videoToAdd.DurationSec {
			columnsToUpdate = append(columnsToUpdate, "duration_sec")
			oldValues = append(oldValues, savedVideo.DurationSec)
			newValues = append(newValues, videoToAdd.DurationSec)
		}
		if savedVideo.EpisodeCount != videoToAdd.EpisodeCount {
			columnsToUpdate = append(columnsToUpdate, "episode_count")
			oldValues = append(oldValues, savedVideo.EpisodeCount)
			newValues = append(newValues, videoToAdd.EpisodeCount)
		}
		if savedVideo.Genres != videoToAdd.Genres {
			columnsToUpdate = append(columnsToUpdate, "genres")
			oldValues = append(oldValues, savedVideo.Genres)
			newValues = append(newValues, videoToAdd.Genres)
		}
		if savedVideo.Mood != videoToAdd.Mood {
			columnsToUpdate = append(columnsToUpdate, "mood")
			oldValues = append(oldValues, savedVideo.Mood)
			newValues = append(newValues, videoToAdd.Mood)
		}
		if savedVideo.NumSeasonLabel != videoToAdd.NumSeasonLabel {
			columnsToUpdate = append(columnsToUpdate, "num_season_label")
			oldValues = append(oldValues, savedVideo.NumSeasonLabel)
			newValues = append(newValues, videoToAdd.NumSeasonLabel)
		}
		if savedVideo.SeasonCount != videoToAdd.SeasonCount {
			columnsToUpdate = append(columnsToUpdate, "season_count")
			oldValues = append(oldValues, savedVideo.SeasonCount)
			newValues = append(newValues, videoToAdd.SeasonCount)
		}
		if savedVideo.Synopsis != videoToAdd.Synopsis {
			columnsToUpdate = append(columnsToUpdate, "synopsis")
			oldValues = append(oldValues, savedVideo.Synopsis)
			newValues = append(newValues, videoToAdd.Synopsis)
		}
		if savedVideo.Tags != videoToAdd.Tags && videoToAdd.Tags != "" {
			//.. If there are tags saved in database, but the tags from the webpage are empty,
			//.. it might be because the data were retrieved in the variable "netflix.falcorCache",
			//.. wich doesn't contain the tags.
			columnsToUpdate = append(columnsToUpdate, "tags")
			oldValues = append(oldValues, savedVideo.Tags)
			newValues = append(newValues, videoToAdd.Tags)
		}
		if savedVideo.Title != videoToAdd.Title {
			columnsToUpdate = append(columnsToUpdate, "title")
			oldValues = append(oldValues, savedVideo.Title)
			newValues = append(newValues, videoToAdd.Title)
		}
		if savedVideo.Type != videoToAdd.Type {
			columnsToUpdate = append(columnsToUpdate, "type")
			oldValues = append(oldValues, savedVideo.Type)
			newValues = append(newValues, videoToAdd.Type)
		}
		if savedVideo.Writers != videoToAdd.Writers {
			columnsToUpdate = append(columnsToUpdate, "writers")
			oldValues = append(oldValues, savedVideo.Writers)
			newValues = append(newValues, videoToAdd.Writers)
		}

		var numberColumnsToUpdate = len(columnsToUpdate)
		if numberColumnsToUpdate > 0 {
			if savedVideo.DataFrom_ != videoToAdd.DataFrom_ {
				columnsToUpdate = append(columnsToUpdate, "_data_from")
				oldValues = append(oldValues, savedVideo.DataFrom_)
				newValues = append(newValues, videoToAdd.DataFrom_)
			}

			videoToAdd.UpdatedAt = dates.NowToString()

			result.Query = "UPDATE"
			result.UpdatedColumns = columnsToUpdate
			result.OldValues = oldValues
			result.NewValues = newValues
		} else {
			result.Query = "NONE"
		}

		var transaction, transactionErr = _connection.Begin()
		if transactionErr != nil {
			result.Error = "TransactionErr: " + transactionErr.Error()
			return result
		}

		var newStatus = savedVideo.Status + "\r\n\r\n" + videoToAdd.Status
		videoToAdd.Status = newStatus
		var finalColumnsToUpdate = append(columnsToUpdate, "status")
		newValues = append(newValues, videoToAdd.Status)
		if updateErr := update(transaction, videoToAdd, finalColumnsToUpdate, newValues); updateErr != nil {
			result.Error = "UpdateErr: " + updateErr.Error()
			return result
		}

		if numberColumnsToUpdate > 0 {
			//.. The video data has changed. I keep a historic of the changes.
			if insertUpdatesErr := insertPlaylistUpdates(transaction, videoToAdd, columnsToUpdate, newValues, oldValues); insertUpdatesErr != nil {
				result.Error = "InsertUpdatesErr: " + insertUpdatesErr.Error()
				if rollbackErr := transaction.Rollback(); rollbackErr != nil {
					result.Error += "\nRollbackErr: " + rollbackErr.Error()
				}
				return result
			}
		}

		if commitErr := transaction.Commit(); commitErr != nil {
			result.Error = "CommitErr: " + commitErr.Error()
		}
	}

	return result
}

// Get the row id and the status for the given video id.
func getVideoFromVideoId(videoId int64) (*videoData, error) {
	var stmt, stmtErr = _connection.Prepare("SELECT rowid, type, title, status, casting, creators, directors, writers, genres, mood, tags, age_advised, age_advised_reason, synopsis, season_count, num_season_label, episode_count, duration_sec, availability_starttime, _data_from FROM playlist WHERE video_id = ?")
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer stmt.Close()

	var rowid int64
	var ttype string
	var title string
	var status string
	var casting string
	var creators string
	var directors string
	var writers string
	var genres string
	var mood string
	var tags string
	var ageAdvised int
	var ageAdvisedReason string
	var synopsis string
	var seasonCount int
	var numSeasonLabel string
	var episodeCount int
	var durationSec int64
	var availabilityStartTime string
	var dataFrom_ string
	var scanErr = stmt.QueryRow(videoId).Scan(&rowid, &ttype, &title, &status, &casting, &creators, &directors, &writers, &genres, &mood, &tags, &ageAdvised, &ageAdvisedReason, &synopsis, &seasonCount, &numSeasonLabel, &episodeCount, &durationSec, &availabilityStartTime, &dataFrom_)
	if scanErr != nil {
		return nil, scanErr
	}
	var video = &videoData{Rowid: rowid, VideoId: videoId, Type: ttype, Title: title, Status: status, Casting: casting, Creators: creators, Directors: directors, Writers: writers, Genres: genres, Mood: mood, Tags: tags, AgeAdvised: ageAdvised, AgeAdvisedReason: ageAdvisedReason, Synopsis: synopsis, SeasonCount: seasonCount, NumSeasonLabel: numSeasonLabel, EpisodeCount: episodeCount, DurationSec: durationSec, AvailabilityStartTime: availabilityStartTime, DataFrom_: dataFrom_}
	return video, nil
}

// Insert a new video to the playlist.
func insertVideo(video *videoData) error {
	var stmt, stmtErr = _connection.Prepare("INSERT INTO playlist(video_id, type, title, status, casting, creators, directors, writers, genres, mood, tags, age_advised, age_advised_reason, synopsis, season_count, num_season_label, episode_count, duration_sec, availability_starttime, _data_from) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	var result, execErr = stmt.Exec(video.VideoId, video.Type, video.Title, video.Status, video.Casting, video.Creators, video.Directors, video.Writers, video.Genres, video.Mood, video.Tags, video.AgeAdvised, video.AgeAdvisedReason, video.Synopsis, video.SeasonCount, video.NumSeasonLabel, video.EpisodeCount, video.DurationSec, video.AvailabilityStartTime, video.DataFrom_)
	if execErr != nil {
		return execErr
	}
	var lastInsertId, _ = result.LastInsertId()
	video.Rowid = lastInsertId
	return nil
}

func update(transaction *sql.Tx, video *videoData, columnsToUpdate []string, newValues []any) error {
	var columns = ""
	for _, column := range columnsToUpdate {
		if columns != "" {
			columns += ", "
		}
		columns += column + " = ?"
	}

	if video.UpdatedAt != "" {
		//.. If the column "status" is the only one to be updated,
		//.. I don't change the value of the column "updated_at".
		if columns != "" {
			columns += ", "
		}
		columns += "updated_at = ?"
	}

	var stmt, stmtErr = transaction.Prepare("UPDATE playlist SET " + columns + " WHERE rowid = ?;")
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	var values []any
	for _, newValue := range newValues {
		values = append(values, newValue)
	}
	if video.UpdatedAt != "" {
		values = append(values, video.UpdatedAt)
	}
	values = append(values, video.Rowid)

	var result, execErr = stmt.Exec(values...)
	if execErr != nil {
		return execErr
	}
	var rowsAffected, _ = result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("The update of the video \"%s\" for the rowid %d has affected no rows", video.Title, video.Rowid)
	} else if rowsAffected > 1 {
		return fmt.Errorf("The update of the video \"%s\" for the rowid %d has affected %d rows", video.Title, video.Rowid, rowsAffected)
	}
	return nil
}

// When video data have changed, I keep a historic of the changes.
func insertPlaylistUpdates(transaction *sql.Tx, video *videoData, columnsToUpdate []string, newValues []any, oldValues []any) error {
	var stmt, stmtErr = transaction.Prepare("INSERT INTO playlist_updates(video_id, updated_at, updates) VALUES(?, ?, ?);")
	if stmtErr != nil {
		return stmtErr
	}
	defer stmt.Close()

	var updatesArray []map[string]any
	for i, v := range columnsToUpdate {
		var data = make(map[string]any)
		data["column"] = v
		data["oldValue"] = oldValues[i]
		data["newValue"] = newValues[i]
		updatesArray = append(updatesArray, data)
	}
	var updatesData, marshalErr = json.MarshalIndent(updatesArray, "", "\t")
	if marshalErr != nil {
		return marshalErr
	}
	var updates = string(updatesData)

	var result, execErr = stmt.Exec(video.VideoId, video.UpdatedAt, updates)
	if execErr != nil {
		return execErr
	}
	var rowsAffected, _ = result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("The update of the video \"%s\" for the rowid %d has affected no rows", video.Title, video.Rowid)
	} else if rowsAffected > 1 {
		return fmt.Errorf("The update of the video \"%s\" for the rowid %d has affected %d rows", video.Title, video.Rowid, rowsAffected)
	}
	return nil
}

func CloseDatabaseConnection() {
	if _connection != nil {
		_connection.Close()
	}
}
