package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type StepContext struct {
	Builder BuildersInner
	Build   BuildsInner
	Step    StepsInner
}

type TestOutput struct {
	Architecture string
	FailedTests  []string
}

func main() {
	builders := getBuilders()
	for _, builder := range builders {
		firstTestResult := true
		var prevFailedTests []TestOutput
		builds := getBuilds(builder.BuilderId)
		for _, build := range builds {
			sourcestamps := getSourcestamps(build.BuildRequestId)
			steps := getSteps(builder.BuilderId, build.Number)
			var failedStepNames []string
			var failedTests []TestOutput
			var inProgress bool
			for _, step := range steps {
				sc := StepContext{
					Builder: builder,
					Build:   build,
					Step:    step,
				}
				sc.dumpLog()
				if sc.IsFailed() {
					failedStepNames = append(failedStepNames, step.Name)
				}

				if sc.IsInProgress() {
					inProgress = true
				}
				sc.dumpTestRawOutput()

				if sc.IsXML() {
					failedTests = append(failedTests, sc.getTestFailures())
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

			fmt.Printf("----------------------------------------\n")
			fmt.Printf("in progress: %v\n", inProgress)
			fmt.Printf("failed steps: %v\n", failedStepNames)
			fmt.Printf("commit time: %v\n commit hash: %v\n", time.Unix(sourcestamps[1].CreatedAt, 0), sourcestamps[1].Revision)
			fmt.Printf("failed steps: %v\n", failedStepNames)
			fmt.Printf("failed test cases: %v\n", failedTests)
			if !firstTestResult {
				newFail, newSuccess := compareTests(prevFailedTests, failedTests)
				fmt.Printf("New fail: %v\nNew success: %v\n", newFail, newSuccess)
			}

			firstTestResult = false
			prevFailedTests = failedTests

			//startedAt := build.StartedAt
		}
		//builderName := builder.Name
	}
}

func compareTests(a, b []TestOutput) ([]TestOutput, []TestOutput) {
	added := []TestOutput{}
	removed := []TestOutput{}
	for _, newTestOutput := range a {
		for _, prevTestOutput := range b {
			if newTestOutput.Architecture == prevTestOutput.Architecture {
				addedTestCases := difference(prevTestOutput.FailedTests, newTestOutput.FailedTests)
				removedTestCases := difference(newTestOutput.FailedTests, prevTestOutput.FailedTests)
				added = append(added, TestOutput{
					Architecture: newTestOutput.Architecture,
					FailedTests:  addedTestCases,
				})
				removed = append(removed, TestOutput{
					Architecture: newTestOutput.Architecture,
					FailedTests:  removedTestCases,
				})
			}
		}
	}

	return added, removed
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var removed []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			removed = append(removed, x)
		}
	}
	return removed
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
	html, err := exec.Command("xsltproc", "--nonet", "--novalid", dirName+sc.GetTargetName()+"-tests-results.xsl", dirName+sc.GetTargetName()+"-test.xml").Output()
	if err != nil {
		panic(err)
	}

	f, err := os.Create(dirName + sc.GetTargetName() + "-tests.html")
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
	return fmt.Sprintf("_out/%d/%s/", sc.Build.StartedAt, sc.Builder.Name)
}

func (sc StepContext) getTestFailures() TestOutput {
	pathName := sc.getOutputDir() + sc.GetTargetName() + "-test.xml"
	return TestOutput{
		Architecture: sc.GetTargetName(),
		FailedTests:  getTestFailuresPath(pathName),
	}
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

	f, err := os.Create(dirName + filename)
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
	sc.dump(fmt.Sprintf("%d.log", sc.Step.Number), false)
}

func (sc StepContext) dumpXML() {
	sc.dump(sc.GetTargetName()+"-test.xml", true)
}

func (sc StepContext) dumpXSL() {
	sc.dump(sc.GetTargetName()+"-tests-results.xsl", true)
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
	return nameWords[len(nameWords)-1]
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
