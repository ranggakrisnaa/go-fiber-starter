package script

import (
	"fmt"
	"os"
	"os/exec"
)

type (
	KafkaScript struct{}
)

func NewKafkaScript() *KafkaScript {
	return &KafkaScript{}
}

func (s *KafkaScript) Run() error {
	fmt.Println("Running Kafka setup script...")

	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	// Path to the kafka-setup.sh script
	scriptPath := dir + "/kafka-setup.sh"

	// Check if the script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("kafka setup script not found at: %s", scriptPath)
	}

	// Execute the script
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Executing Kafka setup script...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing kafka setup script: %v", err)
	}

	fmt.Println("Kafka setup completed successfully!")
	return nil
}
