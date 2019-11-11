package kustomize

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type UpgradeTest struct {
	expectedCommand string
	upgradeStep     UpgradeStep
}

func TestMixin_UnmarshalUpgradeStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/upgrade-robotshop-input.yaml")
	require.NoError(t, err)

	var action UpgradeAction
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	require.Len(t, action.Steps, 2)
	step := action.Steps[0]

	assert.Equal(t, "porter-robotshop-cart", step.Name)
	assert.Equal(t, "Generate the Kubernetes deployment file for the Shopping Cart Microservice", step.Description)
	assert.Contains(t, step.Kustomization, "kustomize/robotshop/overlays/local/cart")
	assert.Equal(t, "manifests/", step.Manifests)

	assert.Equal(t, map[string]string{"kustomizeBaseGHToken": "{{ bundle.parameters.gh_token }}"}, step.Set)
}

func TestMixin_Upgrade(t *testing.T) {
	microService := "cart"
	name := "porter-robotshop-" + microService
	kustomization := []string{"kustomize/robotshop/overlays/local/" + microService}
	manifests := "manifests"
	setArgs := map[string]string{
		"kustomizeBaseGHToken": "{{ bundle.parameters.gh_token }}",
	}

	expectedCmd := fmt.Sprintf("kustomize build %s -o %s/%s.yaml", kustomization[0], manifests, microService)
	expectedGitCmd := fmt.Sprintf("git config --global url.https://{{ bundle.parameters.gh_token }}:@github.com/.insteadOf https://github.com/")
	upgradeTests := []UpgradeTest{
		{
			expectedCommand: expectedGitCmd + "\n" + expectedCmd,
			upgradeStep: UpgradeStep {
				UpgradeArguments: UpgradeArguments{
					Step:          Step{Description: "Upgrade Robotshop"},
					Name:          name,
					Kustomization: kustomization,
					Manifests:     manifests,
					Set:           setArgs,
				},
			},
		},
	}

	defer os.Unsetenv(test.ExpectedCommandEnv)
	for _, upgradeTest := range upgradeTests {
		t.Run(upgradeTest.expectedCommand, func(t *testing.T) {
			err := os.Setenv(test.ExpectedCommandEnv, upgradeTest.expectedCommand)

			if err != nil {
				os.Exit(-1)
			}
			action := UpgradeAction{Steps: []UpgradeStep{upgradeTest.upgradeStep}}
			b, _ := yaml.Marshal(action)

			h := NewTestMixin(t)
			h.In = bytes.NewReader(b)

			err = h.Upgrade()

			require.NoError(t, err)
		})
	}
}
