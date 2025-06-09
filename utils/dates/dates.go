package utils

import (
	"fmt"
	"time"
)

const millisecond int64 = 1000000

// Return the current locale time formatted like: 1988-09-26 02:10:00.123Z
func NowToString() string {
	var now = time.Now()

	//.. I don't use something like: time.Now().Format(time.RFC3339Nano)
	//.. because the nanoseconds are not aligned. Usually they have 7 digits, but sometimes they have 6
	//.. therefore the lines in the log are not well aligned, and I don't like that.
	//.. So I made this code based on "formatHeader" function in "log.go" file in golang source code.
	var year, month, day = now.Date()
	var hour, min, sec = now.Clock()

	var result = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%03d", year, month, day, hour, min, sec, now.Nanosecond()/int(millisecond))
	return result
}
