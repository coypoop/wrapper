package main

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
	CompleteAt     *int `json:"complete_at"`
	MasterId       int  `json:"master_id"`
	Number         int  `json:"number"`
	//Properties []string `json:"properties"` XXX type
	Results     *int   `json:"results"`
	StartedAt   int    `json:"started_at"`
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
	CompleteAt  *int     `json:"complete_at"`
	Hidden      bool     `json:"hidden"`
	Name        string   `json:"name"`
	Number      int      `json:""`
	Results     *int     `json:"results"`
	StartedAt   int      `json:"started_at"`
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

type Meta struct {
	Total int `json:"total"`
}
