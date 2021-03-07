package main

import "encoding/xml"

type BuildersInner struct {
	BuilderId   int      `json:"builderid"`
	Description string   `json:"description"`
	MasterIds   []int    `json:"masterids"`
	Name        string   `json:"name"`
	Tags        []string `json:"tags"`
}

type Builders struct {
	Inner []BuildersInner `json:"builders"`
	Meta  Meta            `json:"meta"`
}

type BuildsInner struct {
	BuilderId      int  `json:"builderid"`
	BuildId        int  `json:"buildid"`
	BuildRequestId int  `json:"buildrequestid"`
	Complete       bool `json:"complete"`
	CompleteAt     *int64 `json:"complete_at"`
	MasterId       int  `json:"master_id"`
	Number         int  `json:"number"`
	//Properties []string `json:"properties"` XXX type
	Results     *int   `json:"results"`
	StartedAt   int64    `json:"started_at"`
	StateString string `json:"state_string"`
	WorkerId    int    `json:"workerid"`
}

type Builds struct {
	Inner []BuildsInner `json:"builds"`
	Meta  Meta          `json:"meta"`
}

type StepsInner struct {
	BuildId     int      `json:"buildid"`
	Complete    bool     `json:"complete"`
	CompleteAt  *int64     `json:"complete_at"`
	Hidden      bool     `json:"hidden"`
	Name        string   `json:"name"`
	Number      int      `json:""`
	Results     *int     `json:"results"`
	StartedAt   int64      `json:"started_at"`
	StateString string   `json:"state_string"`
	StepId      int      `json:"stepid"`
	Urls        []string `json:"urls"`
}

type Steps struct {
	Inner []StepsInner `json:"steps"`
	Meta  Meta         `json:"meta"`
}

type LogsInner struct {
	Complete bool   `json:"complete"`
	LogId    int    `json:"logid"`
	Name     string `json:"name"`
	NumLines int    `json:"num_lines"`
	Slug     string `json:"slug"`
	StepId   int    `json:"stepid"`
	Type     string `json:"type"`
}

type Logs struct {
	Inner []LogsInner `json:"logs"`
	Meta  Meta        `json:"meta"`
}

type BuildRequestsInner struct {
	BuilderId          int  `json:"builderid"`
	BuildRequestId     int  `json:"buildrequestid"`
	BuildSetId         int  `json:"buildsetid"`
	Claimed            bool `json:"claimed"`
	ClaimedAt          int64  `json:"claimed_at"`
	ClaimedByMasterId  int  `json:"claimed_by_masterid"`
	Complete           bool `json:"complete"`
	CompleteAt         int64  `json:"complete_at"`
	CompleteByMasterId int  `json:"complete_by_masterid"`
	Priority           int  `json:"priority"`
	//Properties []string `json:"properties"` XXX type
	Results     *int `json:"results"`
	SubmittedAt int64  `json:"submitted_at"`
	WaitedFor   bool `json:"waited_for"`
}

type BuildRequests struct {
	Inner []BuildRequestsInner `json:"buildrequests"`
	Meta  Meta                `json:"meta"`
}

type Sourcestamps struct {
	Branch     string `json:"branch"`
	Codebase   string `json:"codebase"`
	CreatedAt  int64    `json:"created_at"`
	Repository string `json:"repository"`
	Revision   string `json:"revision"`
	SSID       int    `json:"ssid"`
}

type BuildSetsInner struct {
	// XXX more fields that aren't very interesting
	Sourcestamps []Sourcestamps `json:"sourcestamps"`
}

type BuildSets struct {
	Inner []BuildSetsInner `json:"buildsets"`
	Meta  Meta             `json:"meta"`
}

type Meta struct {
	Total int `json:"total"`
}

type TestResults struct {
	XMLName xml.Name `xml:"tests-results"`
	TestPlans []struct {
		XMLName xml.Name `xml:"tp"`
		ID string `xml:"id,attr"`
		TestCases []struct{
			XMLName xml.Name `xml:"tc"`
			ID string `xml:"id,attr"`
			Failed string `xml:"failed"`
		} `xml:"tc"`
	} `xml:"tp"`
}
