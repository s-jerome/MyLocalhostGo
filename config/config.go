package config

import (
	"os"
	"strings"
)

var configs = map[string]string{
	"server.port": "8801",
}

func Read() error {
	var fileData, readErr = os.ReadFile("config.txt")
	if readErr != nil {
		if os.IsNotExist(readErr) == false {
			return readErr
		} else {
			return nil
		}
	}

	var fileContent = string(fileData)
	var lines = strings.Split(fileContent, "\n")
	for i := 0; i < len(lines); i++ {
		var line = strings.TrimSpace(lines[i])
		if len(line) == 0 || line == "" || strings.Index(line, "#") == 0 || strings.Index(line, "//") == 0 {
			continue
		}
		var key, value, found = strings.Cut(line, "=")
		if found == false {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		configs[key] = value
	}

	return nil
}

func Get(key string) string {
	//.. If the key is not present, the value is an empty string.
	var value = configs[key]
	return value
}
