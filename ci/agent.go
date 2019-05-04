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
		p.status("Language is Currently not supported")
		log.Printf("%s was requested\n", p.Language)
		return
	}
	contr, err := a.docker.ContainerCreate(a.ctx, &container.Config{
		Image: p.getImage(),
		Env:   p.envVars(),
	}, nil, nil, "")
	if err != nil {
		p.status("Failed to run container")
		log.Printf("Unable to create container of %s\n", p.getImage())
		return
	}
	if err := a.docker.ContainerStart(a.ctx, contr.ID, types.ContainerStartOptions{}); err != nil {
		p.status("Unable to run container")
		log.Printf("Unable to create container of %s\n", p.getImage())
		return
	}
	statusCh, errCh := a.docker.ContainerWait(a.ctx, contr.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}
	logfile, err := a.docker.ContainerLogs(a.ctx, contr.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Println("Failed get log File")
	}
	err = a.savePipelineLog(p, logfile)
	if err != nil {
		log.Println("failed to close writer")
	}
	err = a.docker.ContainerRemove(a.ctx, contr.ID, types.ContainerRemoveOptions{})
	if err != nil {
		log.Printf("Unable to remove container %v", contr.ID)
	}
}

func (a *agent) savePipelineLog(p *pipeline, logfile io.ReadCloser) error {
	wc := a.storage.Bucket("kfr-ci-pipelines").Object(p.LogFileName).NewWriter(a.ctx)
	if _, err := io.Copy(wc, logfile); err != nil {
		p.status("Unable to save log")
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}
