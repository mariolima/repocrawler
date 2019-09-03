package bitbucket

type CommitsResponse struct {
	Next    string   `json:"next"`
	Pagelen int64    `json:"pagelen"`
	Values  []Commit `json:"values"`
}

type CommitsResponse_sub5 struct {
	Approve  CommitsResponse_sub1 `json:"approve"`
	Comments CommitsResponse_sub1 `json:"comments"`
	Diff     CommitsResponse_sub1 `json:"diff"`
	HTML     CommitsResponse_sub1 `json:"html"`
	Patch    CommitsResponse_sub1 `json:"patch"`
	Self     CommitsResponse_sub1 `json:"self"`
	Statuses CommitsResponse_sub1 `json:"statuses"`
}

type Commit struct {
	Author     CommitsResponse_sub4   `json:"author"`
	Date       string                 `json:"date"`
	Hash       string                 `json:"hash"`
	Links      CommitsResponse_sub5   `json:"links"`
	Message    string                 `json:"message"`
	Parents    []CommitsResponse_sub7 `json:"parents"`
	Rendered   CommitsResponse_sub9   `json:"rendered"`
	Repository CommitsResponse_sub10  `json:"repository"`
	Summary    CommitsResponse_sub8   `json:"summary"`
	Type       string                 `json:"type"`
}

type CommitsResponse_sub2 struct {
	Avatar CommitsResponse_sub1 `json:"avatar"`
	HTML   CommitsResponse_sub1 `json:"html"`
	Self   CommitsResponse_sub1 `json:"self"`
}

type CommitsResponse_sub10 struct {
	FullName string               `json:"full_name"`
	Links    CommitsResponse_sub2 `json:"links"`
	Name     string               `json:"name"`
	Type     string               `json:"type"`
	UUID     string               `json:"uuid"`
}

type CommitsResponse_sub6 struct {
	HTML CommitsResponse_sub1 `json:"html"`
	Self CommitsResponse_sub1 `json:"self"`
}

type CommitsResponse_sub8 struct {
	HTML   string `json:"html"`
	Markup string `json:"markup"`
	Raw    string `json:"raw"`
	Type   string `json:"type"`
}

type CommitsResponse_sub7 struct {
	Hash  string               `json:"hash"`
	Links CommitsResponse_sub6 `json:"links"`
	Type  string               `json:"type"`
}

type CommitsResponse_sub1 struct {
	Href string `json:"href"`
}

type CommitsResponse_sub9 struct {
	Message CommitsResponse_sub8 `json:"message"`
}

type CommitsResponse_sub4 struct {
	Raw  string `json:"raw"`
	Type string `json:"type"`
	User User   `json:"user"`
}
