package main

import (
	"fmt"
)

func main() {
	builders := getBuilders()
	for _, builder := range builders {
		builds := getBuilds(builder.BuilderId)
		for _, build := range builds {
			steps := getSteps(builder.BuilderId, build.Number)
			for _, step := range steps {
				logs := getLogs(builder.BuilderId, build.Number, step.Number)
				logRaw := getLogRaw(logs)
				if len(logs) == 1 && logs[0].LogId == 207 {
					fmt.Printf("%s", logRaw)
				}
			}

		}
	}
}
