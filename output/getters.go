package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getLogRaw(logs []LogsInner) []byte {
	var out []byte

	for _, log := range logs {
		response, err := http.Get(fmt.Sprintf("http://localhost:8010/api/v2/logs/%d/raw", log.LogId))
		if err != nil {
			panic("Failed to interact with buildbot REST API")
		}
		defer response.Body.Close()

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic("Failed to interact with buildbot REST API")
		}
		out = append(out, responseData...)
	}

	return out
}

func getLogs(builder, build, step int) []LogsInner {
	var logs Logs

	response, err := http.Get(fmt.Sprintf("http://localhost:8010/api/v2/builders/%d/builds/%d/steps/%d/logs", builder, build, step))
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}

	err = json.Unmarshal(responseData, &logs)
	if err != nil {
		panic("Failed to unmarshal buildbot REST API - logs")
	}

	return logs.Inner
}

func getSteps(builder, build int) []StepsInner {
	var steps Steps

	response, err := http.Get(fmt.Sprintf("http://localhost:8010/api/v2/builders/%d/builds/%d/steps", builder, build))
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}

	err = json.Unmarshal(responseData, &steps)
	if err != nil {
		panic("Failed to unmarshal buildbot REST API - steps")
	}

	return steps.Inner
}

const MaxBuilds = 10

func getBuilds(builder int) []BuildsInner {
	var builds Builds

	response, err := http.Get(fmt.Sprintf("http://localhost:8010/api/v2/builders/%d/builds?order=-started_at&limit=%d", builder, MaxBuilds))
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}

	err = json.Unmarshal(responseData, &builds)
	if err != nil {
		panic("Failed to unmarshal buildbot REST API - builds")
	}

	return builds.Inner
}

func getBuilders() []BuildersInner {
	var builders Builders

	response, err := http.Get("http://localhost:8010/api/v2/builders")
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}

	err = json.Unmarshal(responseData, &builders)
	if err != nil {
		panic("Failed to unmarshal buildbot REST API - builders")
	}

	return builders.Inner
}