package main

import (
	"encoding/json"
	"fmt"
	commander "github.com/NubeIO/bios-cli/cmd"
	"log"
	"os"
	"strings"
)

func main() {
	bt := commander.NewBuildTool()
	var resp []*response
	command := os.Args[2]
	if command == "listCommands" {
		_, err := bt.ExecuteStep(commander.BuildStep{Name: "listCommands", Cmd: "listCommands", Params: nil})
		if err != nil {
			log.Fatalf("Error executing command listCommands: %v", err)
		}
		return
	}

	buildYAML, err := bt.LoadBuildYAML(command)
	if err != nil {
		out := &response{
			File:  command,
			Error: err.Error(),
		}
		resp = append(resp, out)
		dump(resp)
		return
	}

	args := parseArgs(os.Args[3:])

	for i, step := range buildYAML.Steps {
		out := &response{
			Name:      step.Name,
			Cmd:       step.Cmd,
			StepCount: i,
		}
		params := replaceParams(step.Params, buildYAML.Vars, args)
		ret, err := bt.ExecuteStep(commander.BuildStep{Name: step.Name, Cmd: step.Cmd, Params: params})
		if err != nil {
			out.Error = err.Error()
		} else {
			out.Response = ret
			resp = append(resp, out)
		}
	}
	dump(resp)
}

func dump(resp any) {
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}

type response struct {
	File      string      `json:"file,omitempty"`
	Name      string      `json:"name,omitempty"`
	Cmd       string      `json:"cmd,omitempty"`
	StepCount int         `json:"stepCount,omitempty"`
	Response  interface{} `json:"response,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func parseArgs(rawArgs []string) map[string]string {
	args := make(map[string]string)
	for _, arg := range rawArgs {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			args[parts[0]] = parts[1]
		}
	}
	return args
}

func replaceParams(params interface{}, vars []commander.Variable, args map[string]string) interface{} {
	switch p := params.(type) {
	case string:
		return replaceVarsAndArgs(p, vars, args)
	case []interface{}:
		var replacedParams []interface{}
		for _, param := range p {
			replacedParams = append(replacedParams, replaceVarsAndArgs(fmt.Sprintf("%v", param), vars, args))
		}
		return replacedParams
	case map[string]interface{}:
		replacedParams := make(map[string]interface{})
		for key, param := range p {
			replacedParams[key] = replaceVarsAndArgs(fmt.Sprintf("%v", param), vars, args)
		}
		return replacedParams
	}
	return params
}

func replaceVarsAndArgs(param string, vars []commander.Variable, args map[string]string) string {
	result := param
	for key, val := range args {
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", key), val)
	}
	for _, v := range vars {
		if val, ok := args[v.Name]; ok {
			result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", v.Name), val)
		} else {
			result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", v.Name), v.Value)
		}
	}
	return result
}
