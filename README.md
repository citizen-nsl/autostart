# AutoStart

This project is a Go application that manages processes and schedules tasks based on a configuration file. The application is designed to handle process termination, window management, and batch file execution on Windows systems.

## Features

- Terminate existing processes by name.
- Close old command prompt windows.
- Execute a batch file based on the schedule defined in a YAML configuration file.
- Schedule tasks using a cron-like scheduler.

## Getting Started

### Prerequisites

- Go 1.17 or higher
- Access to a Windows system for process management and window handling

### Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/yourusername/my-go-project.git
    cd my-go-project
    ```

2. **Initialize the Go module:**

    If not already initialized, create a `go.mod` file by running:

    ```sh
    go mod init my-go-project
    ```

3. **Install dependencies:**

    Install the required Go packages using:

    ```sh
    go get github.com/StackExchange/wmi
    go get github.com/robfig/cron/v3
    ```

4. **Prepare configuration file:**

    Create a file named `m308.yaml` in the root of the project directory with the following format:

    ```yaml
    schedules:
      path:
        batch_file: "path/to/your/batchfile.bat"
        fxserver: "path/to/fxserver.exe"
        server: "path/to/server.exe"
      schedule:
        - time: "12:00"
        - time: "18:00"
    ```

    Replace the paths and times with your actual values.

### Running the Application

To run the application, use:

```sh
go run main.go
```

```sh
go build
```

```sh
./autostart
```
