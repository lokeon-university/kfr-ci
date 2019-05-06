package main

import (
	"context"
	"io"
	"log"

	"cloud.google.com/go/storage"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type agent struct {
	docker  *client.Client
	storage *storage.Client
	ctx     context.Context
}

func newAgent() *agent {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	stgcli, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	cli.NegotiateAPIVersion(ctx)
	return &agent{
		cli,
		stgcli,
		ctx,
	}
}

func (a *agent) buildPipeline(p *pipeline) {
	if !p.supportedLanguage() {
		p.status(":earth_americas: Language is currently not supported.")
		log.Printf("%s was requested\n", p.Language)
		return
	}
	contr, err := a.docker.ContainerCreate(a.ctx, &container.Config{
		Image: p.getImage(),
		Env:   p.envVars(),
	}, nil, nil, "")
	if err != nil {
		p.status(":construction: Failed to create pipeline.")
		log.Printf("Unable to create container of %s\n", p.getImage())
		return
	}
	if err := a.docker.ContainerStart(a.ctx, contr.ID, types.ContainerStartOptions{}); err != nil {
		p.status(":rotating_light: Failed to build pipeline")
		log.Printf("Unable to create container of %s\n", p.getImage())
		return
	}
	p.status(":tools: Building Pipeline :tools:")
	correct := true
	statusCh, errCh := a.docker.ContainerWait(a.ctx, contr.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Println(err)
			panic(err)
		}
	case sts := <-statusCh:
		switch sts.StatusCode {
		case int64(2):
			p.status(":x: File .kfr-ci.json not found")
			break
		case int64(4):
			p.status(":x: key stepts cannot be empty")
			break
		case int64(0):
			p.status(":tada: pipeline finished.")
			break
		}
	}
	logfile, err := a.docker.ContainerLogs(a.ctx, contr.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		correct = false
		log.Printf("Failed get log file: %v, %v", contr.ID, p.LogFileName)
	}
	file, err := a.savePipelineLog(p, logfile)
	if err != nil {
		correct = false
		log.Printf("Failed to close writer: %v, %v", contr.ID, p.LogFileName)
	}
	err = a.docker.ContainerRemove(a.ctx, contr.ID, types.ContainerRemoveOptions{})
	if err != nil {
		log.Printf("Unable to remove container %v\n", contr.ID)
	}
	if correct {
		p.status(file)
	} else {
		p.status(":bomb: failed to create logs please retry.")
	}
}

func (a *agent) savePipelineLog(p *pipeline, logfile io.ReadCloser) (string, error) {
	file := p.getLogFileName()
	wc := a.storage.Bucket("kfr-ci-pipelines").Object(file).NewWriter(a.ctx)
	if _, err := io.Copy(wc, logfile); err != nil {
		p.status(":page_with_curl: Failed saving log")
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	return file, nil
}
