package logger

import (
	"fmt"
	"log"
	"mylocalhost/utils"
	"os"
	"path/filepath"
	"strings"
)

func WriteError(err string, v ...interface{}) {
	appendToFile("logs/errors.log", err, v...)
}

func WriteLog(text string, v ...interface{}) {
	appendToFile("logs/log.log", text, v...)
}

func appendToFile(filePath string, text string, v ...interface{}) {
	var textToWrite string
	if v == nil {
		textToWrite = text
	} else {
		textToWrite = fmt.Sprintf(text, v...)
	}
	var newLine = ""
	if strings.HasSuffix(textToWrite, "\n") == false {
		newLine = "\n"
	}
	textToWrite = fmt.Sprintf("%s -- %s%s", utils.NowToString(), textToWrite, newLine)

	if mkdirErr := mkdirAll(filePath); mkdirErr != nil {
		log.Printf("[logger.appendToFile] Error for mkdir the file \"%s\":\n%v\ntext to write:\n%s", filePath, mkdirErr, textToWrite)
		return
	}

	var file, fileErr = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if fileErr != nil {
		log.Printf("[logger.appendToFile] Error opening the file \"%s\":\n%v\ntext to write:\n%s", filePath, fileErr, textToWrite)
		return
	}
	defer file.Close()

	var dataToWrite = []byte(textToWrite)
	if _, writeErr := file.Write(dataToWrite); writeErr != nil {
		log.Printf("[logger.appendToFile] Error writing to the file \"%s\":\n%v\ntext to write:\n%s", filePath, writeErr, textToWrite)
	}
}

func mkdirAll(filePath string) error {
	var directoryPath = filepath.Dir(filePath)
	var mkdirErr = os.MkdirAll(directoryPath, os.ModePerm)
	return mkdirErr
}
