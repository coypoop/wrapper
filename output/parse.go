package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"sort"
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
	TestsURL     template.URL
}

type StepOutput struct {
	Name string
	URL  template.URL
}

type PageData struct {
	Builds []Build
}

type Build struct {
	InProgress          bool
	Buildtype           string
	CommitUnix          int64
	CommitDate          string
	CommitRevision      template.URL
	FailedStepsMetric   string
	FailedTestsTotal    string
	HasFailedSteps      bool
	FailedSteps         []StepOutput
	HasNewTestFailures  bool
	NewTestFailures     []TestOutput
	HasNewTestSuccesses bool
	NewTestSuccesses    []TestOutput
	HasTestFailures     bool
	TotalTestResults    []TestOutput
}

func main() {
	buildsParsed := []Build{}
	builders := getBuilders()
	for _, builder := range builders {
		firstTestResult := true
		var prevFailedTests []TestOutput
		builds := getBuilds(builder.BuilderId)
		for _, build := range builds {
			sourcestamp := getSrcSourcestamp(build.BuildRequestId)
			commitDate := time.Unix(sourcestamp.CreatedAt, 0)

			// Don't show builds that are too old
			if commitDate.Add(7 * 24 * 60 * time.Minute).Before(time.Now()) {
				break
			}

			steps := getSteps(builder.BuilderId, build.Number)
			var failedSteps []StepOutput
			var failedTests, newTestFailures, newTestSuccesses []TestOutput
			var inProgress bool
			for _, step := range steps {
				sc := StepContext{
					Builder: builder,
					Build:   build,
					Step:    step,
				}
				sc.dumpLog()
				if sc.IsFailed() {
					failedSteps = append(failedSteps, StepOutput{
						Name: step.Name,
						URL:  template.URL(fmt.Sprintf("%s%d.log", sc.getExternalDir(), sc.Step.Number)),
					})
				}

				if sc.IsInProgress() {
					inProgress = true
				}
				sc.dumpTestRawOutput()

				if sc.IsXML() {
					failures := sc.getTestFailures()
					if len(failures.FailedTests) > 0 {
						failedTests = append(failedTests, sc.getTestFailures())
					}
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

			if !firstTestResult {
				newTestFailures, newTestSuccesses = compareTests(prevFailedTests, failedTests)
			}

			firstTestResult = false
			prevFailedTests = failedTests

			buildsParsed = append(buildsParsed, Build{
				InProgress:          inProgress,
				Buildtype:           builder.Name,
				CommitUnix:          commitDate.Unix(),
				CommitDate:          commitDate.String(),
				CommitRevision:      template.URL(sourcestamp.Revision),
				FailedStepsMetric:   fmt.Sprintf("%d/%d", len(failedSteps), len(steps)),
				FailedTestsTotal:    fmt.Sprintf("%d", len(failedTests)),
				HasFailedSteps:      len(failedSteps) > 0,
				FailedSteps:         failedSteps,
				HasNewTestFailures:  len(newTestFailures) > 0,
				NewTestFailures:     newTestFailures,
				HasNewTestSuccesses: len(newTestSuccesses) > 0,
				NewTestSuccesses:    newTestSuccesses,
				HasTestFailures:     len(failedTests) > 0,
				TotalTestResults:    failedTests,
			})
			/*
				fmt.Printf("----------------------------------------\n")
				fmt.Printf("in progress: %v\n", inProgress)
				fmt.Printf("failed steps: %v\n", failedSteps)
				fmt.Printf("commit time: %v\n commit hash: %v\n", time.Unix(sourcestamps[1].CreatedAt, 0), sourcestamps[1].Revision)
				fmt.Printf("failed steps: %v\n", failedSteps)
				fmt.Printf("failed test cases: %v\n", failedTests)
			*/
			//startedAt := build.StartedAt
		}
		//builderName := builder.Name
	}
	var tpl bytes.Buffer
	tmpl := template.Must(template.ParseFiles("index.html"))

	data := PageData{
		Builds: sortBuilds(buildsParsed),
	}
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("_out/index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(tpl.String())
	if err != nil {
		panic(err)
	}
}

func sortBuilds(input []Build) []Build {
	sort.Slice(input, func(i, j int) bool {
		if input[i].Buildtype == input[j].Buildtype {
			return input[i].CommitUnix > input[j].CommitUnix
		}
		return input[i].Buildtype < input[j].Buildtype
	})

	return input
}

func compareTests(a, b []TestOutput) ([]TestOutput, []TestOutput) {
	added := []TestOutput{}
	removed := []TestOutput{}
	for _, newTestOutput := range a {
		for _, prevTestOutput := range b {
			if newTestOutput.Architecture == prevTestOutput.Architecture {
				addedTestCases := difference(prevTestOutput.FailedTests, newTestOutput.FailedTests)
				removedTestCases := difference(newTestOutput.FailedTests, prevTestOutput.FailedTests)
				if len(addedTestCases) > 0 {
					added = append(added, TestOutput{
						Architecture: newTestOutput.Architecture,
						FailedTests:  addedTestCases,
						TestsURL:     newTestOutput.TestsURL,
					})
				}
				if len(removedTestCases) > 0 {
					removed = append(removed, TestOutput{
						Architecture: newTestOutput.Architecture,
						FailedTests:  removedTestCases,
						TestsURL:     newTestOutput.TestsURL,
					})
				}
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

func (sc StepContext) getExternalDir() string {
	return fmt.Sprintf("%d/%s/", sc.Build.StartedAt, sc.Builder.Name)
}

func (sc StepContext) getTestFailures() TestOutput {
	fileName := sc.GetTargetName() + "-test.xml"
	outFilename := sc.GetTargetName() + "-tests.html"
	return TestOutput{
		TestsURL:     template.URL(sc.getExternalDir() + outFilename),
		Architecture: sc.GetTargetName(),
		FailedTests:  getTestFailuresPath(sc.getOutputDir() + fileName),
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
