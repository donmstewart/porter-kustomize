package kustomize

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

// Validate that the output directory for Kubernetes Manifest `yaml` files
//   exists (if not make it) and trim any trailing '/' characters
func (m *Mixin) manifestHandling(step interface{}) error {
	var manifests string

	switch s := step.(type) {
	case InstallStep:
		manifests = s.Manifests
	case UpgradeStep:
		manifests = s.Manifests
	case UninstallStep:
		manifests = s.Manifests
	default:
		return errors.New("Unsupported Step type " + reflect.TypeOf(step).String())
	}

	// Do we have a trailing '/' if so remove it
	if strings.HasSuffix(manifests, "/") {
		manifests = strings.TrimSuffix(manifests, "/")
	}
	// Check if the manifest directory exists if not create it
	if _, err := m.FileSystem.DirExists(manifests); err != nil {
		if m.Debug {
			fmt.Println("DEBUG: Manifests directory missing creating...")
		}
		if err := m.FileSystem.MkdirAll(manifests, os.ModePerm); err != nil {
			// We failed to create the manifest output directory so return an error
			return errors.Wrap(err, "couldn't make output directory")
		}
	}
	return nil
}

// Loops round the list of entries in the `kustomization_input` field in `porter.yaml` and
//   build a command to execute `kustomize build` on that directory. Once we have iterated over the list
//   and built the commands then we execute the list of commands to generate the Kubernetes
//   `yaml` files.
func (m *Mixin) buildAndExecuteKustomizeCmds(step interface{}, commands []*exec.Cmd) error {
	var kustomization []string
	var manifests string
	var reorder = "legacy"

	switch s := step.(type) {
	case InstallStep:
		kustomization = s.Kustomization
		manifests = s.Manifests
		if s.Reorder != "" {
			reorder = s.Reorder
		}
	case UpgradeStep:
		kustomization = s.Kustomization
		manifests = s.Manifests
		if s.Reorder != "" {
			reorder = s.Reorder
		}
	case UninstallStep:
		kustomization = s.Kustomization
		manifests = s.Manifests
		if s.Reorder != "" {
			reorder = s.Reorder
		}
	default:
		return errors.New("Unsupported Step type")
	}

	if m.Debug {
		fmt.Println("DEBUG: Reorder: " + reorder)
	}

	// Loop around the list of kustomization directories specified in the `porter.yaml`
	for _, kustomizationFile := range kustomization {
		// The path to write out the generatred Kubernetes Manifests
		pathSegments := strings.Split(kustomizationFile, string(os.PathSeparator))

		if strings.HasSuffix(manifests, string(os.PathSeparator)) == false {
			manifests = manifests + string(os.PathSeparator)
		}

		// Build the kustomize command string and pipe it to the output file in the manifests directory
		cmd := m.NewCommand("kustomize", "build", kustomizationFile, "--reorder", reorder, "-o", manifests +
			//pathSegments[len(pathSegments)-1]+".yaml", "--reorder", reorder)
			pathSegments[len(pathSegments)-1]+".yaml")

		commands = append(commands, cmd)
	}
	// Loop and execute the list of kustomization commands
	for _, cmd := range commands {
		//buf := new(bytes.Buffer)
		//cmd.Stdout = buf
		cmd.Stdout = m.Out
		cmd.Stderr = m.Err

		prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
		if m.Debug {
			fmt.Println("DEBUG: " + prettyCmd)
		}
		err := cmd.Start()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("couldn't run command %s", prettyCmd))
		}
		err = cmd.Wait()
		if err != nil {
			prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
			return errors.Wrap(err, fmt.Sprintf("error running command %s", prettyCmd))
		}
		//output := buf.String()
		sensitiveFields := []string{"kustomizeBaseGHToken"}
		m.Context.SetSensitiveValues(sensitiveFields)
		//_, err = m.Out.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

// If using git::https:// to access Kustomize assets in remote repos then a GITHUB TOKEN is used
//   When the token is specified then git config is run in the container to update the git
//   configuration to use the token on all git requests.
func (m *Mixin) configureGithubToken(ghToken string) error {
	if ghToken != "" {
		var gitArgs = strings.Builder{}
		if _, err := gitArgs.WriteString("url.https://"); err != nil {
			return err
		}
		if _, err := gitArgs.WriteString(ghToken); err != nil {
			return err
		}
		if _, err := gitArgs.WriteString(":@github.com/.insteadOf"); err != nil {
			return err
		}

		gitCmd := m.NewCommand("git", "config", "--global", gitArgs.String(), "https://github.com/")
		if m.Debug {
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
		} else {
			gitCmd.Stdout = m.Out
			gitCmd.Stderr = m.Err
		}

		err := gitCmd.Start()
		if err != nil {
			return err
		}
		err = gitCmd.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}
