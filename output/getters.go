package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

func getBuildRequests(buildrequestid int) []BuildRequestsInner {
	var buildrequests BuildRequests

	response, err := http.Get(fmt.Sprintf("http://localhost:8010/api/v2/buildrequests/%d", buildrequestid))
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}

	err = json.Unmarshal(responseData, &buildrequests)
	if err != nil {
		panic("Failed to unmarshal buildbot REST API - buildrequestid")
	}

	return buildrequests.Inner
}

func getBuildSets(buildsetid int) []BuildSetsInner {
	var buildsets BuildSets

	response, err := http.Get(fmt.Sprintf("http://localhost:8010/api/v2/buildsets/%d", buildsetid))
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Failed to interact with buildbot REST API")
	}

	err = json.Unmarshal(responseData, &buildsets)
	if err != nil {
		panic("Failed to unmarshal buildbot REST API - buildsetid")
	}

	return buildsets.Inner
}

func getSrcSourcestamp(buildrequestid int) Sourcestamps {
	for _, sourcestamp := range getSourcestamps(buildrequestid) {
		if !strings.Contains(sourcestamp.Codebase, "xsrc") {
			return sourcestamp
		}
	}
	panic("Didn't find src sourcestamp")
}

func getSourcestamps(buildrequestid int) []Sourcestamps {
	buildrequests := getBuildRequests(buildrequestid)
	var sourcestamps []Sourcestamps

	for _, buildrequest := range buildrequests {
		buildsets := getBuildSets(buildrequest.BuildSetId)
		for _, buildset := range buildsets {
			sourcestamps = append(sourcestamps, buildset.Sourcestamps...)
		}
	}

	return sourcestamps
}

func getTestFailuresPath(xmlPath string) []string {
	var testFailures []string

	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		panic(err)
	}
	defer xmlFile.Close()

	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}
	testResults := TestResults{}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&testResults)

	for _, testPlan := range testResults.TestPlans {
		for _, testCase := range testPlan.TestCases {
			if testCase.Failed != "" {
				testFailures = append(testFailures, testPlan.ID+":"+testCase.ID)
			}
		}
	}

	return testFailures
}
