package kustomize

import (
	"os/exec"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

type UninstallAction struct {
	Steps []UninstallStep `yaml:"uninstall"`
}

// UninstallStep represents the structure of an Uninstall action
type UninstallStep struct {
	UninstallArguments `yaml:"kustomize"`
}

// UninstallArguments are the arguments available for the Uninstall action
type UninstallArguments struct {
	Step `yaml:",inline"`

	Name          string            `yaml:"name"`
	Kustomization []string          `yaml:"kustomization_input"`
	Manifests     string            `yaml:"kubernetes_manifest_output"`
	Reorder       string            `yaml:"reorder"`
	Set           map[string]string `yaml:"set"`
	Purge         bool              `yaml:"purge"`
}

// Uninstall deletes a provided set of Kustomize releases, supplying optional flags/params
func (m *Mixin) Uninstall() error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var action UninstallAction
	err = yaml.Unmarshal(payload, &action)
	if err != nil {
		return err
	}
	if len(action.Steps) != 1 {
		return errors.Errorf("expected a single step, but got %d", len(action.Steps))
	}
	step := action.Steps[0]

	var commands []*exec.Cmd

	ghToken := step.Set["kustomizeBaseGHToken"]

	err = m.configureGithubToken(ghToken)
	if err != nil {
		return err
	}

	err = m.manifestHandling(step)
	if err != nil {
		return err
	}

	err = m.buildAndExecuteKustomizeCmds(step, commands)
	if err != nil {
		return err
	}

	return nil
}
