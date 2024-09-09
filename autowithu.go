package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Schedules struct {
		Path struct {
			BatchFile string `yaml:"batch_file"`
			FxServer  string `yaml:"fxserver"`
			Server    string `yaml:"server"`
		} `yaml:"path"`
		Schedule []struct {
			Time string `yaml:"time"`
		} `yaml:"schedule"`
	} `yaml:"schedules"`
}

func loadConfig() (Config, error) {
	var config Config
	configPath := "m308.yaml"
	file, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("config file not found at %s", configPath)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %v", err)
	}
	return config, nil
}

func killExistingProcesses(processName string) {
	query := fmt.Sprintf("SELECT * FROM Win32_Process WHERE Name='%s'", processName)
	var processes []struct {
		ProcessID uint32
	}
	err := wmi.Query(query, &processes)
	if err != nil {
		fmt.Println("Error querying processes:", err)
		return
	}

	for _, process := range processes {
		fmt.Printf("Terminating existing process: %d\n", process.ProcessID)
		cmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprint(process.ProcessID))
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to terminate process %d: %v\n", process.ProcessID, err)
		}
	}
}

func closeOldCmdWindows() {
	query := "SELECT * FROM Win32_Process WHERE Name='cmd.exe'"
	var processes []struct {
		ProcessID uint32
	}
	err := wmi.Query(query, &processes)
	if err != nil {
		fmt.Println("Error querying cmd processes:", err)
		return
	}

	for _, process := range processes {
		fmt.Printf("Closing old cmd window: %d\n", process.ProcessID)
		cmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprint(process.ProcessID))
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to close cmd window %d: %v\n", process.ProcessID, err)
		}
	}
}

func runBatchFile(config Config) {
	closeOldCmdWindows()
	killExistingProcesses("fxserver.exe")
	killExistingProcesses("server.exe")
	time.Sleep(5 * time.Second)

	batchFilePath := config.Schedules.Path.BatchFile
	fmt.Printf("Attempting to run batch file at: %s\n", time.Now().Format(time.RFC3339))
	cmd := exec.Command("cmd", "/C", "start", batchFilePath)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running batch file: %v\n", err)
	}
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	c := cron.New()
	for _, task := range config.Schedules.Schedule {
		_, err := c.AddFunc(task.Time, func() { runBatchFile(config) })
		if err != nil {
			fmt.Printf("Error scheduling task: %v\n", err)
		}
	}

	c.Start()
	fmt.Println("Scheduler started. Waiting for the scheduled time...")

	select {} // Block forever
}
