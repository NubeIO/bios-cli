package commander

import (
	"time"
)

func (bt *BuildTool) time(params interface{}) (interface{}, error) {
	return time.Now(), nil
}
