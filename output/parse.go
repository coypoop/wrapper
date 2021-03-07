package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"os"
	"os/exec"
)

type StepContext struct {
	Builder BuildersInner
	Build BuildsInner
	Step StepsInner
}

func main() {
	builders := getBuilders()
	for _, builder := range builders {
		builds := getBuilds(builder.BuilderId)
		for _, build := range builds {
			sourcestamps := getSourcestamps(build.BuildRequestId)
			steps := getSteps(builder.BuilderId, build.Number)
			var failedSteps []StepsInner
			var failedTests []string
			var inProgress bool
			for _, step := range steps {
				sc := StepContext{
					Builder: builder,
					Build: build,
					Step: step,
				}
				if sc.IsFailed() {
					failedSteps = append(failedSteps, step)
				}

				if sc.IsInProgress() {
					inProgress = true
				}
				sc.dumpTestRawOutput()

				if sc.IsXML() {
					failedTests = sc.getTestFailures()
					sc.dumpTestHTML()
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
			fmt.Printf("commit time: %v\n commit hash: %v\n", time.Unix(sourcestamps[1].CreatedAt, 0), sourcestamps[1].Revision)
			fmt.Printf("failed steps: %v\n", failedSteps)
			fmt.Printf("failed test cases: %v\n", failedTests)
			//startedAt := build.StartedAt
		}
		//builderName := builder.Name
	}
}

func (sc StepContext) dumpTestRawOutput() {
	if sc.IsXSL() {
		sc.dumpXSL()
	}

	if sc.IsXML() {
		sc.dumpXML()
	}

	if sc.IsCSS() {
		sc.dumpCSS()
	}
}

func (sc StepContext) dumpTestHTML() {
	dirName := sc.getOutputDir()
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

func (sc StepContext) getOutputDir() string {
	return fmt.Sprintf("_out/%d/%s/%s/", sc.Build.StartedAt, sc.Builder.Name, sc.GetTargetName())
}

func (sc StepContext) getTestFailures() []string {
	pathName := sc.getOutputDir() + "test.xml"
	return getTestFailuresPath(pathName)
}

func (sc StepContext) dump(filename string, stripDebug bool) {
	dirName := sc.getOutputDir()
	err := os.MkdirAll(dirName, 0744)
	if err != nil {
		panic(err)
	}

	logs := getLogs(sc.Builder.BuilderId, sc.Build.Number, sc.Step.Number)
	logRaw := getLogRaw(logs)
	if stripDebug {
		logRaw = stripBuildbotDebug(logRaw)
	}

	f, err := os.Create(dirName + "/" + filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(logRaw)
	if err != nil {
		panic(err)
	}
}

func (sc StepContext) dumpLog() {
	sc.dump(string(sc.Step.Number) + ".log", false)
}

func (sc StepContext) dumpXML() {
	sc.dump("test.xml", true)
}

func (sc StepContext) dumpXSL() {
	sc.dump("tests-results.xsl", true)
}

func (sc StepContext) dumpCSS() {
	sc.dump("tests-results.css", true)
}

// Buildbot adds some information about the command being executed.
// (Exit status, command executed) - strip it
func stripBuildbotDebug(log []byte) []byte {
	split := bytes.Split(log, []byte("\n"))
	return bytes.Join(split[5:len(split)-2], []byte("\n"))
}

func (sc StepContext) GetTargetName() string {
	nameWords := strings.Split(sc.Step.Name, " ")
	return nameWords[len(nameWords) - 1]
}

func (sc StepContext) IsCSS() bool {
	return strings.Contains(sc.Step.Name, "CSS")
}

func (sc StepContext) IsXSL() bool {
	return strings.Contains(sc.Step.Name, "XSL")
}

func (sc StepContext) IsXML() bool {
	return strings.Contains(sc.Step.Name, "XML")
}

func (sc StepContext) IsFailed() bool {
	return (sc.Step.Results != nil &&
		*sc.Step.Results != 0)
}

func (sc StepContext) IsInProgress() bool {
	return sc.Step.Results == nil
}
