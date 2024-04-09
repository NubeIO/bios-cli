package commander

import (
	"fmt"
)

func (bt *BuildTool) handleSystemInfo(params interface{}) (interface{}, error) {
	ops, err := parseParams(params)
	if err != nil {
		return nil, fmt.Errorf("system information: %v", err)
	}
	return bt.system.ExecuteMethods(ops)
}

func parseParams(params interface{}) (operations []string, err error) {
	var paramList []string
	switch p := params.(type) {
	case []string:
		paramList = p
	case []interface{}:
		// Try to convert []interface{} to []string
		for _, v := range p {
			if str, ok := v.(string); ok {
				paramList = append(paramList, str)
			} else {
				return nil, fmt.Errorf("invalid element type in []interface{}")
			}
		}
	default:
		return nil, fmt.Errorf("invalid params type")
	}
	return paramList, nil
}
