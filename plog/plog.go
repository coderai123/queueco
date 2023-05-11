package plog

import (
	"encoding/json"
	"log"
)

type Fields map[string]interface{}

func LogInfo(location, message string, fields map[string]interface{}) {
	LogFields("INFO", location, message, fields)
}

func LogError(location, message string, fields map[string]interface{}) {
	LogFields("ERROR", location, message, fields)
}

func LogFields(level, location, message string, fields map[string]interface{}) {
	logData := map[string]interface{}{
		"level":    level,
		"location": location,
		"message":  message,
	}
	for key, value := range fields {
		logData[key] = value
	}
	logBytes, _ := json.Marshal(logData)
	log.Println(string(logBytes))
}
