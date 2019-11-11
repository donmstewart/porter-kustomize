package kustomize

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os/exec"
)

type UpgradeAction struct {
	Steps []UpgradeStep `yaml:"upgrade"`
}

// UpgradeStep represents the structure of an Upgrade step
type UpgradeStep struct {
	UpgradeArguments `yaml:"kustomize"`
}

// UpgradeArguments represent the arguments available to the Upgrade step
type UpgradeArguments struct {
	Step `yaml:",inline"`

	Name          string            `yaml:"name"`
	Kustomization []string          `yaml:"kustomization_input"`
	Manifests     string            `yaml:"kubernetes_manifest_output"`
	Reorder       string            `yaml:"reorder"`
	Set           map[string]string `yaml:"set"`
}

// Upgrade issues a kustomize upgrade command for a release using the provided UpgradeArguments
func (m *Mixin) Upgrade() error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var action UpgradeAction
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
