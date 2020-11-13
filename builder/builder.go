package builder

import (
	//"bytes"
	"context"
	"fmt"
	// "bufio"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type BuildConfig struct {
	Image string `hcl:"image"`
	Env map[string]string `hcl:"env,optional"`
	Push bool `hcl:"push,optional"`
}

type Builder struct {
	config BuildConfig
}

// Implement Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// Implement ConfigurableNotify
func (b *Builder) ConfigSet(config interface{}) error {
	c, ok := config.(*BuildConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit
		return fmt.Errorf("Expected *BuildConfig as parameter")
	}

	// validate the config
	if c.Image == "" {
		return fmt.Errorf("Image must be set to a valid image reference")
	}

	return nil
}

// Implement Builder
func (b *Builder) BuildFunc() interface{} {
	// return a function which will be called by Waypoint
	return b.build
}

func write(f *os.File, str string) error {
	f.WriteString(fmt.Sprintf("%s\n", str))
	err := f.Sync()

	return err
}

// A BuildFunc does not have a strict signature, you can define the parameters
// you need based on the Available parameters that the Waypoint SDK provides.
// Waypoint will automatically inject parameters as specified
// in the signature at run time.
//
// Available input parameters:
// - context.Context
// - *component.Source
// - *component.JobInfo
// - *component.DeploymentConfig
// - *datadir.Project
// - *datadir.App
// - *datadir.Component
// - hclog.Logger
// - terminal.UI
// - *component.LabelSet
//
// The output parameters for BuildFunc must be a Struct which can
// be serialzied to Protocol Buffers binary format and an error.
// This Output Value will be made available for other functions
// as an input parameter.
// If an error is returned, Waypoint stops the execution flow and
// returns an error to the user.
func (b *Builder) build(ctx context.Context, ui terminal.UI, log hclog.Logger) (*Binary, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	// If we have a step set, abort it on exit
	var s terminal.Step
	defer func() {
		if s != nil {
			s.Abort()
		}
	}()

	s = sg.Add("Initializing dobi build context")
	var command string
	if (b.config.Push) {
		command = "push"
	} else {
		command = "build"
	}

	subcommand := fmt.Sprintf("%s:%s", b.config.Image, command)

	cmd := exec.Command(
		"dobi",
		subcommand,
	)

	cmd.Env = os.Environ()
	for key, value := range b.config.Env {
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("%s=%s", key, value),
		)
    }
	s.Done()


	s = sg.Add(fmt.Sprintf("Executing command: dobi %s", subcommand))

	cmd.Stdout = s.TermOutput()
	cmd.Stderr = cmd.Stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	s.Done()

	return &Binary{}, nil
}
