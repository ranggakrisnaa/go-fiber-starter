package script

import (
	"errors"

	"gorm.io/gorm"
)

func Script(scriptName string, db *gorm.DB) error {
	switch scriptName {
	case "example_script":
		exampleScript := NewExampleScript(db)
		return exampleScript.Run()
	case "kafka_setup":
		kafkaScript := NewKafkaScript()
		return kafkaScript.Run()
	default:
		return errors.New("script not found")
	}
}
