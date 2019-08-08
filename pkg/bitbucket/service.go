package bitbucket

type Repositories struct {
	Next    string              `json:"next"`
	Page    int64               `json:"page"`
	Pagelen int64               `json:"pagelen"`
	Size    int64               `json:"size"`
	Values  []Repositories_sub8 `json:"values"`
}

type Repositories_sub3 struct {
	Avatar       Repositories_sub1   `json:"avatar"`
	Branches     Repositories_sub1   `json:"branches"`
	Clone        []Repositories_sub2 `json:"clone"`
	Commits      Repositories_sub1   `json:"commits"`
	Downloads    Repositories_sub1   `json:"downloads"`
	Forks        Repositories_sub1   `json:"forks"`
	Hooks        Repositories_sub1   `json:"hooks"`
	HTML         Repositories_sub1   `json:"html"`
	Issues       Repositories_sub1   `json:"issues"`
	Pullrequests Repositories_sub1   `json:"pullrequests"`
	Self         Repositories_sub1   `json:"self"`
	Source       Repositories_sub1   `json:"source"`
	Tags         Repositories_sub1   `json:"tags"`
	Watchers     Repositories_sub1   `json:"watchers"`
}

type Repositories_sub5 struct {
	Avatar Repositories_sub1 `json:"avatar"`
	HTML   Repositories_sub1 `json:"html"`
	Self   Repositories_sub1 `json:"self"`
}

type Repositories_sub8 struct {
	CreatedOn   string            `json:"created_on"`
	Description string            `json:"description"`
	ForkPolicy  string            `json:"fork_policy"`
	FullName    string            `json:"full_name"`
	HasIssues   bool              `json:"has_issues"`
	HasWiki     bool              `json:"has_wiki"`
	IsPrivate   bool              `json:"is_private"`
	Language    string            `json:"language"`
	Links       Repositories_sub3 `json:"links"`
	Mainbranch  Repositories_sub4 `json:"mainbranch"`
	Name        string            `json:"name"`
	Owner       Repositories_sub6 `json:"owner"`
	Project     Repositories_sub7 `json:"project"`
	Scm         string            `json:"scm"`
	Size        int64             `json:"size"`
	Slug        string            `json:"slug"`
	Type        string            `json:"type"`
	UpdatedOn   string            `json:"updated_on"`
	UUID        string            `json:"uuid"`
	Website     string            `json:"website"`
}

type Repositories_sub6 struct {
	DisplayName string            `json:"display_name"`
	Links       Repositories_sub5 `json:"links"`
	Type        string            `json:"type"`
	Username    string            `json:"username"`
	UUID        string            `json:"uuid"`
}

type Repositories_sub2 struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type Repositories_sub1 struct {
	Href string `json:"href"`
}

type Repositories_sub7 struct {
	Key   string            `json:"key"`
	Links Repositories_sub5 `json:"links"`
	Name  string            `json:"name"`
	Type  string            `json:"type"`
	UUID  string            `json:"uuid"`
}

type Repositories_sub4 struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

