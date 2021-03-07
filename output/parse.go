package main

import (
	"bytes"
	"fmt"
	"strings"
	"os"
	"os/exec"
)

func main() {
	builders := getBuilders()
	for _, builder := range builders {
		builds := getBuilds(builder.BuilderId)
		for _, build := range builds {
			steps := getSteps(builder.BuilderId, build.Number)
			var failedSteps []StepsInner
			var inProgress bool
			for _, step := range steps {
				if step.IsFailed() {
					failedSteps = append(failedSteps, step)
				}

				if step.IsInProgress() {
					inProgress = true
				}
				dumpTestRawOutput(builder, build, step)

				if step.IsXML() {
					//testFailures = parseTestOutput(step)
					dumpTestHTML(builder, build, step)
				}
/*
				logs := getLogs(builder.BuilderId, build.Number, step.Number)
				logRaw := getLogRaw(logs)
				if len(logs) == 1 && logs[0].LogId == 207 {
					fmt.Printf("%s", logRaw)
				}
*/
			}
			fmt.Printf("in progress: %v\n", inProgress)
			fmt.Printf("failed steps: %v\n", failedSteps)
			//startedAt := build.StartedAt
		}
		//builderName := builder.Name
	}
}

func dumpTestRawOutput(builder BuildersInner, build BuildsInner, step StepsInner) {
	if step.IsXSL() {
		dumpXSL(builder, build, step)
	}

	if step.IsXML() {
		dumpXML(builder, build, step)
	}

	if step.IsCSS() {
		dumpCSS(builder, build, step)
	}
}

func dumpTestHTML(builder BuildersInner, build BuildsInner, step StepsInner) {
	dirName := getOutputDir(builder, build, step)
	html, err := exec.Command("xsltproc", "--nonet", "--novalid", dirName + "tests-results.xsl", dirName + "test.xml").Output()
	if err != nil {
		panic(err)
	}

	f, err := os.Create(dirName + "/tests.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(html)
	if err != nil {
		panic(err)
	}
}

func getOutputDir(builder BuildersInner, build BuildsInner, step StepsInner) string {
	return fmt.Sprintf("_out/%d/%s/%s/", build.StartedAt, builder.Name, step.GetTargetName())
}

func dumpRaw(builder BuildersInner, build BuildsInner, step StepsInner, filename string) {
	dirName := getOutputDir(builder, build, step)
	err := os.MkdirAll(dirName, 0744)
	if err != nil {
		panic(err)
	}

	logs := getLogs(builder.BuilderId, build.Number, step.Number)
	logRaw := getLogRaw(logs)
	XML := stripBuildbotDebug(logRaw)

	f, err := os.Create(dirName + "/" + filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(XML)
	if err != nil {
		panic(err)
	}
}

func dumpXML(builder BuildersInner, build BuildsInner, step StepsInner) {
	dumpRaw(builder, build, step, "test.xml")
}

func dumpXSL(builder BuildersInner, build BuildsInner, step StepsInner) {
	dumpRaw(builder, build, step, "tests-results.xsl")
}

func dumpCSS(builder BuildersInner, build BuildsInner, step StepsInner) {
	dumpRaw(builder, build, step, "tests-results.css")
}

// Buildbot adds some information about the command being executed.
// (Exit status, command executed) - strip it
func stripBuildbotDebug(log []byte) []byte {
	split := bytes.Split(log, []byte("\n"))
	return bytes.Join(split[5:len(split)-2], []byte("\n"))
}

func (step StepsInner) GetTargetName() string {
	nameWords := strings.Split(step.Name, " ")
	return nameWords[len(nameWords) - 1]
}

func (step StepsInner) IsCSS() bool {
	return strings.Contains(step.Name, "CSS")
}

func (step StepsInner) IsXSL() bool {
	return strings.Contains(step.Name, "XSL")
}

func (step StepsInner) IsXML() bool {
	return strings.Contains(step.Name, "XML")
}

func (step StepsInner) IsFailed() bool {
	return (step.Results != nil &&
		*step.Results != 0)
}

func (step StepsInner) IsInProgress() bool {
	return step.Results == nil
}
