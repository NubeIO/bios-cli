package main

import (
	"fmt"
	commander "github.com/NubeIO/ros-bios/cmd"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type requestBody struct {
	File string            `json:"file"`
	Args map[string]string `json:"args"`
}

func main() {
	bt := commander.NewBuildTool()

	if len(os.Args) > 1 && os.Args[1] == "server" {
		// Start the Gin server
		r := gin.Default()

		// List all the YAML files
		r.GET("/api/files", func(c *gin.Context) {
			files, err := filepath.Glob("*.yaml")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"files": files})
		})

		// Run a YAML file
		r.POST("/api/run", func(c *gin.Context) {
			var request *requestBody
			if err := c.BindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			fmt.Println("requestBody", request)
			buildYAML, err := bt.LoadBuildYAML(request.File)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			resp := make(map[string]interface{})
			for i, step := range buildYAML.Steps {
				params := replaceParams(step.Params, buildYAML.Vars, request.Args)
				fmt.Println("params", params)
				if ret, err := bt.ExecuteStep(commander.BuildStep{Name: step.Name, Cmd: step.Cmd, Params: params}); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				} else {
					resp[fmt.Sprintf("%s_%d", step.Cmd, i)] = ret
				}
			}

			c.JSON(http.StatusOK, gin.H{"message": resp})
		})

		// Run the server
		if err := r.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}

	} else if len(os.Args) < 3 {
		fmt.Println("Usage: build-tool build build.yaml MESSAGE=hello-world")
		_, err := bt.ExecuteStep(commander.BuildStep{Name: "listCommands", Cmd: "listCommands", Params: nil})
		if err != nil {
			log.Fatalf("Error executing command listCommands: %v", err)

		}
		os.Exit(1)
	}

	command := os.Args[2]
	if command == "listCommands" {
		_, err := bt.ExecuteStep(commander.BuildStep{Name: "listCommands", Cmd: "listCommands", Params: nil})
		if err != nil {
			log.Fatalf("Error executing command listCommands: %v", err)
		}
		return
	}

	buildYAML, err := bt.LoadBuildYAML(os.Args[2])
	if err != nil {
		log.Fatalf("Error loading build YAML: %v", err)
	}

	args := parseArgs(os.Args[3:])

	fmt.Printf("Name: %s\n", buildYAML.Name)
	for _, step := range buildYAML.Steps {
		fmt.Printf("Step: %s\n", step.Name)

		params := replaceParams(step.Params, buildYAML.Vars, args)

		_, err := bt.ExecuteStep(commander.BuildStep{Name: step.Name, Cmd: step.Cmd, Params: params})
		if err != nil {
			log.Fatalf("Error executing step %s: %v", step.Name, err)
		}
	}
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
