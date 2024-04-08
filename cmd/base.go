package commander

import (
	"fmt"
	"gopkg.in/yaml.v3"

	"os"
	"os/exec"
	"strings"
)

// CommandHandler defines the signature for functions that handle commands.
type CommandHandler func(params interface{}) error

// BuildTool represents a build tool instance.
type BuildTool struct {
	CommandMap map[string]CommandHandler
	Commands   map[string]Command
	buildYAML  BuildYAML
}

type Command struct {
	Func CommandHandler
	Name string
	Help string
}

// NewBuildTool creates a new BuildTool instance.
func NewBuildTool() *BuildTool {
	bt := &BuildTool{}
	bt.Commands = make(map[string]Command)
	bt.CommandMap = map[string]CommandHandler{
		"listCommands":    bt.handleListCommands,
		"systemctl":       bt.handleSystemctl,
		"bash":            bt.handleRunBash,
		"http":            bt.handleRestyHTTPRequest,
		"github-download": bt.handleGitHubDownload,
		"dirs":            bt.handleFiles,
		"systemctl-file":  bt.handleSystemctlFile,
	}
	bt.Commands["listCommands"] = Command{Func: bt.handleListCommands, Name: "listCommands", Help: "List all available commands"}
	bt.Commands["systemctl"] = Command{Func: bt.handleSystemctl, Name: "systemctl", Help: "Manage systemd services"}
	bt.Commands["bash"] = Command{Func: bt.handleRunBash, Name: "runBash", Help: "Execute a bash command"}
	bt.Commands["http"] = Command{Func: bt.handleRestyHTTPRequest, Name: "http", Help: "Make an HTTP request using Resty"}
	bt.Commands["github-download"] = Command{Func: bt.handleGitHubDownload, Name: "github-download", Help: "Download and unzip a GitHub release"}
	bt.Commands["dirs"] = Command{Func: bt.handleFiles, Name: "dirs", Help: "Add/Edit files and dirs"}
	bt.Commands["dirs"] = Command{Func: bt.handleFiles, Name: "dirs", Help: "Add/Edit files and dirs"}
	bt.Commands["systemctl-file"] = Command{Func: bt.handleFiles, Name: "dirs", Help: "Generates a systemctl file"}

	return bt
}

// handleListCommands lists all available commands and their descriptions.
func (bt *BuildTool) handleListCommands(_ interface{}) error {
	fmt.Println("Available commands:")
	for _, cmd := range bt.Commands {
		fmt.Printf("%s - %s\n", cmd.Name, cmd.Help)
	}
	return nil
}

// BuildStep represents a single step in the build process.
type BuildStep struct {
	Name   string      `yaml:"name"`
	Cmd    string      `yaml:"cmd"`
	Params interface{} `yaml:"params"`
}

// BuildYAML represents the structure of the build.yaml file.
type BuildYAML struct {
	Shell string      `yaml:"shell"`
	Name  string      `yaml:"name"`
	Args  []string    `yaml:"args"`
	Vars  []Variable  `yaml:"vars"`
	Steps []BuildStep `yaml:"steps"`
}

// Variable represents a variable in the YAML file.
type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// LoadBuildYAML loads the build.yaml file into a BuildYAML struct.
func (bt *BuildTool) LoadBuildYAML(filename string) (*BuildYAML, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var buildYAML BuildYAML
	err = yaml.Unmarshal(yamlFile, &buildYAML)
	if err != nil {
		return nil, err
	}
	bt.buildYAML = buildYAML
	return &buildYAML, nil
}

// UpdateVar updates a variable in the BuildYAML.
func (bt *BuildTool) UpdateVar(name string, value string) {
	for i, v := range bt.buildYAML.Vars {
		if v.Name == name {
			bt.buildYAML.Vars[i].Value = value
			return
		}
	}
	// If the variable doesn't exist, add it
	bt.buildYAML.Vars = append(bt.buildYAML.Vars, Variable{Name: name, Value: value})
}

// ExecuteStep executes a single step in the build process.
func (bt *BuildTool) ExecuteStep(step BuildStep) error {
	handler, ok := bt.CommandMap[step.Cmd]
	if !ok {
		return fmt.Errorf("unknown command: %s", step.Cmd)
	}

	return handler(step.Params)
}

func (bt *BuildTool) executeCommand(commandName string, params interface{}) error {
	var cmdString string
	// Check if params is a slice of strings
	switch p := params.(type) {
	case string:
		cmdString = p
	case []string:
		cmdString = strings.Join(p, " ")
	default:
		return fmt.Errorf("invalid params type for %s command", commandName)
	}

	fmt.Println("[", cmdString, "]")
	// Split the command string into command and arguments
	parts := strings.Fields(cmdString)
	if len(parts) < 1 {
		return fmt.Errorf("%s command requires at least one argument", commandName)
	}
	cmd := parts[0]
	args := parts[1:]
	execCmd := exec.Command(commandName, append([]string{cmd}, args...)...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	err := execCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run %s command: %v", commandName, err)
	}

	return nil
}