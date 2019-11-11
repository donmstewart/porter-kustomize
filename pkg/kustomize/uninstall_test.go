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

type UninstallTest struct {
	expectedCommand string
	uninstallStep   UninstallStep
}

func TestMixin_UnmarshalUninstallStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/uninstall-robotshop-input.yaml")
	require.NoError(t, err)

	var action UninstallAction
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

func TestMixin_Uninstall(t *testing.T) {
	microService := "cart"
	name := "porter-robotshop-" + microService
	kustomization := []string{"kustomize/robotshop/overlays/local/" + microService}
	manifests := "manifests"
	setArgs := map[string]string{
		"kustomizeBaseGHToken": "{{ bundle.parameters.gh_token }}",
	}

	expectedCmd := fmt.Sprintf("kustomize build %s -o %s/%s.yaml", kustomization[0], manifests, microService)
	expectedGitCmd := fmt.Sprintf("git config --global url.https://{{ bundle.parameters.gh_token }}:@github.com/.insteadOf https://github.com/")
	uninstallTests := []UninstallTest{
		{
			expectedCommand: expectedGitCmd + "\n" + expectedCmd,
			uninstallStep: UninstallStep {
				UninstallArguments: UninstallArguments{
					Step:          Step{Description: "Unnstall Robotshop"},
					Name:          name,
					Kustomization: kustomization,
					Manifests:     manifests,
					Set:           setArgs,
				},
			},
		},
	}

	defer os.Unsetenv(test.ExpectedCommandEnv)
	for _, uninstallTest := range uninstallTests {
		t.Run(uninstallTest.expectedCommand, func(t *testing.T) {
			err := os.Setenv(test.ExpectedCommandEnv, uninstallTest.expectedCommand)

			if err != nil {
				os.Exit(-1)
			}
			action := UninstallAction{Steps: []UninstallStep{uninstallTest.uninstallStep}}
			b, _ := yaml.Marshal(action)

			h := NewTestMixin(t)
			h.In = bytes.NewReader(b)

			err = h.Uninstall()

			require.NoError(t, err)
		})
	}
}
