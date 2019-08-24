package bitbucket

type RepositoriesResponse struct {
	Next    string								`json:"next"`
	Page    int64								`json:"page"`
	Pagelen int64								`json:"pagelen"`
	Size    int64								`json:"size"`
	Values  []Repository						`json:"values"`
}

type RepositoryApiLinks struct {
	Avatar       LinkData						`json:"avatar"`
	Branches     LinkData						`json:"branches"`
	Clone        []RepositoryCloneLink			`json:"clone"`
	Commits      LinkData						`json:"commits"`
	Downloads    LinkData						`json:"downloads"`
	Forks        LinkData						`json:"forks"`
	Hooks        LinkData						`json:"hooks"`
	HTML         LinkData						`json:"html"`
	Issues       LinkData						`json:"issues"`
	Pullrequests LinkData						`json:"pullrequests"`
	Self         LinkData						`json:"self"`
	Source       LinkData						`json:"source"`
	Tags         LinkData						`json:"tags"`
	Watchers     LinkData						`json:"watchers"`
}

type UserApiLinks struct {
	Avatar LinkData								`json:"avatar"`
	HTML   LinkData								`json:"html"`
	Self   LinkData								`json:"self"`
}

type Repository struct {
	CreatedOn   string							`json:"created_on"`
	Description string							`json:"description"`
	ForkPolicy  string							`json:"fork_policy"`
	FullName    string							`json:"full_name"`
	HasIssues   bool							`json:"has_issues"`
	HasWiki     bool							`json:"has_wiki"`
	IsPrivate   bool							`json:"is_private"`
	Language    string							`json:"language"`
	Links       RepositoryApiLinks				`json:"links"`
	Mainbranch  RepositoryBranch				`json:"mainbranch"`
	Name        string							`json:"name"`
	Owner       RepositoryUser					`json:"owner"`
	Project     RepositoryProject				`json:"project"`
	Scm         string							`json:"scm"`
	Size        int64							`json:"size"`
	Slug        string							`json:"slug"`
	Type        string							`json:"type"`
	UpdatedOn   string							`json:"updated_on"`
	UUID        string							`json:"uuid"`
	Website     string							`json:"website"`
}

type RepositoryUser struct {
	DisplayName string							`json:"display_name"`
	Links       UserApiLinks					`json:"links"`
	Type        string							`json:"type"`
	Username    string							`json:"username"`
	UUID        string							`json:"uuid"`
}

type RepositoryCloneLink struct {
	Href string									`json:"href"`
	Name string									`json:"name"`
}

type LinkData struct {
	Href string									`json:"href"`
}

type RepositoryProject struct {
	Key   string								`json:"key"`
	Links UserApiLinks							`json:"links"`
	Name  string								`json:"name"`
	Type  string								`json:"type"`
	UUID  string								`json:"uuid"`
}

type RepositoryBranch struct {
	Name string									`json:"name"`
	Type string									`json:"type"`
}
