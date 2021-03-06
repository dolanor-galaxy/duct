package duct

import (
	"context"
	"log"
	"os"

	dc "github.com/fsouza/go-dockerclient"
)

// Build is a set of instructions for building a container image. All paths are
// relative to the working directory of the test.
type Build struct {
	// Dockerfile is the path to the dockerfile.
	Dockerfile string
	// Context is the directory to use.
	Context string
}

// Builder is a named collection of builds.
type Builder map[string]Build

// Run runs the builds. It logs them to stderr similarly to `docker build`.
func (bc Builder) Run(ctx context.Context) error {
	client, err := dc.NewClientFromEnv()
	if err != nil {
		return err
	}

	for name, build := range bc {
		dir := build.Context
		if dir == "" {
			dir = "."
		}

		log.Printf("Building image: [%s]", name)
		err := client.BuildImage(dc.BuildImageOptions{
			Context:      ctx,
			Name:         name,
			ContextDir:   dir,
			Dockerfile:   build.Dockerfile,
			OutputStream: os.Stderr,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
