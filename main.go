package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/docker/docker/api/types"
// 	"github.com/docker/docker/api/types/container"
// 	"github.com/docker/docker/client"
// 	"github.com/docker/docker/pkg/stdcopy"
// )

// var images = map[string]string{
// 	"cpp":    "gcr.io/kfr-ci/kfr-cpp",
// 	"go":     "gcr.io/kfr-ci/kfr-go",
// 	"java":   "gcr.io/kfr-ci/kfr-java",
// 	"python": "gcr.io/kfr-ci/kfr-python",
// }

// type pipeline struct {
// 	RepositoryID int64  `json:"repository_id,omitempty"`
// 	URL          string `json:"url,omitempty"`
// 	Repository   string `json:"repository,omitempty"`
// 	Branch       string `json:"branch,omitempty"`
// 	LogFileName  string `json:"log_file_name,omitempty"`
// 	Language     string `json:"language,omitempty"`
// }

// func (p *pipeline) runContainer() {
// 	ctx := context.Background()
// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		panic(err)
// 	}
// 	cli.NegotiateAPIVersion(ctx)
// 	if !supportedLanguage(p.Language) {
// 		log.Println("Language currently not supported")
// 	}
// 	resp, err := cli.ContainerCreate(ctx, &container.Config{
// 		Image: getImage(p.Language),
// 		Env:   p.loadEnvVars(),
// 	}, nil, nil, "")
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
// 		panic(err)
// 	}
// 	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
// 	select {
// 	case err := <-errCh:
// 		if err != nil {
// 			panic(err)
// 		}
// 	case <-statusCh:
// 	}
// 	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
// 	if err != nil {
// 		panic(err)
// 	}
// 	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
// }

// func supportedLanguage(lang string) bool {
// 	_, ok := images[lang]
// 	return ok
// }

// func getImage(lang string) string {
// 	img, _ := images[lang]
// 	return img
// }

// func (p *pipeline) loadEnvVars() []string {
// 	return []string{
// 		keyValueEnv("REPO_URL", p.URL),
// 		keyValueEnv("REPO_NAME", p.Repository),
// 		keyValueEnv("REPO_BRANCH", p.Branch),
// 	}
// }

// func keyValueEnv(key, value string) string {
// 	return fmt.Sprintf("%s=%s", key, value)
// }

// func main() {
// 	p := pipeline{
// 		RepositoryID: int64(3242398998436593423),
// 		URL:          "https://github.com/krosf-university/POO.git",
// 		Repository:   "poo",
// 		Branch:       "P4",
// 		LogFileName:  "3213468762138765168765138765135768424",
// 		Language:     "cpp",
// 	}
// 	p.runContainer()
// }
type status struct {
	Owner      string `json:"owner,omitempty"`
	RepoName   string `json:"repo_name,omitempty"`
	Status     string `json:"status,omitempty"`
	TelegramID string `json:"telegram_id,omitempty"`
}

type updateStatus struct {
	Message struct {
		Attributes struct {
			Key string `json:"key,omitempty"`
		} `json:"attributes,omitempty"`
		Data      status `json:"data,omitempty"`
		MessageID string `json:"messageId,omitempty"`
	} `json:"message,omitempty"`
	Subscription string `json:"subscription,omitempty"`
}

func (u *updateStatus) UnmarshalJSON(data []byte) error {
	type Alias updateStatus
	aux := &struct {
		Message struct {
			Data string `json:"data"`
		} `json:"message,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(aux.Message.Data)
	if err = json.Unmarshal(data, &u.Message.Data); err != nil {
		return err
	}
	return nil
}

func main() {
	data := "{\"message\":{\"attributes\":{\"key\":\"value\"},\"data\":\"eyJvd25lciI6Imtyb3NmIiwicmVwb19uYW1lIjoicG9vIiwic3RhdHVzIjoiYnVpbGRpbmciLCJ0ZWxlZ3JhbV9pZCI6IjIzNDQ1MzI0NCJ9\",\"messageId\":\"136969346945\"},\"subscription\":\"projects/myproject/subscriptions/mysubscription\"}"
	var s updateStatus
	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		fmt.Println(err, "err")
	}
	fmt.Println(s.Message.Data)
}
