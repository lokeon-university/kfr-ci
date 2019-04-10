package ci

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//LogExt log file extension
const (
	LogExt  = ".log"
	Error   = "error"
	Pending = "pending"
	Success = "success"
	Failure = "failure"
)

var (
	// ciDIR is the absolute path to the CI directory
	ciDIR = filepath.Join(os.Getenv("CI_ROOT"), "ci")
	// LogDIR is the absolute path to the CI log directory
	LogDIR = filepath.Join(ciDIR, "logs")
	//availableImages list of available of supported languages
	availableImages = map[string]string{
		"go":         "docker/go",
		"javascript": "docker/node",
		"c++":        "docker/c++",
		"python":     "docker/python",
	}
)

func supportedLanguage(lang string) (ok bool) {
	_, ok = availableImages[lang]
	return
}

// Pipeline contains required information to run pipeline for a given project.
type Pipeline struct {
	// LogFileName name to be used for pipeline output.
	// It should always be the commit hash.
	LogFileName string
	logFilePath string
	// Repository name of pipeline project.
	Repository string
	// Branch name of pipeline project.
	Branch string
	// URL of repository
	URL string
	// Language of the pipeline project.
	Language string
	// UpdateStatus callback function to update status of pipeline.
	UpdateStatus func(string)
}

func checkError(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
	return
}

func createDirFor(filePath string) error {
	dir, file := filepath.Split(filePath)
	log.Printf("Making dir: %s for file: %s\n", dir, file)
	return os.MkdirAll(dir, 0755)
}

//ActivePipeline return true if pipeline is active
func ActivePipeline(filepath string) bool {
	cmd := exec.Command("lsof", filepath)
	return cmd.Run() == nil
}

// TriggerPipeline it builds the absolute path to the job log file, creating necessary parent directories.
// It terminates if a routine is currently active for the given pipeline.
// Otherwise, sets up a new routine for the pipeline.
func TriggerPipeline(pipe *Pipeline) {
	pipe.logFilePath = filepath.Join(LogDIR, fmt.Sprintf("%s%s", pipe.LogFileName, LogExt))
	err := createDirFor(pipe.logFilePath)
	if err != nil {
		log.Println("Couldn't create directory for job: ", err)
		return
	}
	if ActivePipeline(pipe.logFilePath) {
		log.Println("A pipeline is currently in progress: ", pipe.logFilePath)
		return
	}
	if err := exec.Command("bash", "-c", "> ", pipe.logFilePath).Run(); err != nil {
		log.Printf("Error: %s occurred while trying to clear logfile %s\n", err, pipe.logFilePath)
		return
	}
	if !supportedLanguage(pipe.Language) {
		log.Println("Project Language is currently not supported.")
		return
	}
	log.Printf("Running pipeline: %v\n", pipe)
	go RunPipeline(pipe)
}

//RunPipeline run current pipeline
func RunPipeline(pipe *Pipeline) {
	logFile, err := os.OpenFile(pipe.logFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("Error %s occurred while opening log file: %s\n", err, pipe.logFilePath)
		pipe.UpdateStatus(Error)
		return
	}
	defer logFile.Close()
	pipe.UpdateStatus(Pending)
	image := availableImages[pipe.Language]
	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s '%s' %s", filepath.Join(ciDIR, "run.sh"), "ENV", image))
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Run()
	msg := "Pipeline completed successfully"
	status := Success
	log.Println("Exit code: ", err)
	if err != nil {
		msg = fmt.Sprintf("Test failed with exit code: %s", err)
		status = Failure
	}
	pipe.UpdateStatus(status)
	logFile.WriteString(fmt.Sprintf("<h4>%s</h4>", msg))
}
