package github_own

type GithubCode struct {
	IncompleteResults bool              `json:"incomplete_results"`
	Items             []GithubCode_sub3 `json:"items"`
	TotalCount        int64             `json:"total_count"`
}

type GithubCode_sub2 struct {
	ArchiveURL       string          `json:"archive_url"`
	AssigneesURL     string          `json:"assignees_url"`
	BlobsURL         string          `json:"blobs_url"`
	BranchesURL      string          `json:"branches_url"`
	CollaboratorsURL string          `json:"collaborators_url"`
	CommentsURL      string          `json:"comments_url"`
	CommitsURL       string          `json:"commits_url"`
	CompareURL       string          `json:"compare_url"`
	ContentsURL      string          `json:"contents_url"`
	ContributorsURL  string          `json:"contributors_url"`
	DeploymentsURL   string          `json:"deployments_url"`
	Description      string          `json:"description"`
	DownloadsURL     string          `json:"downloads_url"`
	EventsURL        string          `json:"events_url"`
	Fork             bool            `json:"fork"`
	ForksURL         string          `json:"forks_url"`
	FullName         string          `json:"full_name"`
	GitCommitsURL    string          `json:"git_commits_url"`
	GitRefsURL       string          `json:"git_refs_url"`
	GitTagsURL       string          `json:"git_tags_url"`
	HooksURL         string          `json:"hooks_url"`
	HTMLURL          string          `json:"html_url"`
	ID               int64           `json:"id"`
	IssueCommentURL  string          `json:"issue_comment_url"`
	IssueEventsURL   string          `json:"issue_events_url"`
	IssuesURL        string          `json:"issues_url"`
	KeysURL          string          `json:"keys_url"`
	LabelsURL        string          `json:"labels_url"`
	LanguagesURL     string          `json:"languages_url"`
	MergesURL        string          `json:"merges_url"`
	MilestonesURL    string          `json:"milestones_url"`
	Name             string          `json:"name"`
	NodeID           string          `json:"node_id"`
	NotificationsURL string          `json:"notifications_url"`
	Owner            GithubCode_sub1 `json:"owner"`
	Private          bool            `json:"private"`
	PullsURL         string          `json:"pulls_url"`
	ReleasesURL      string          `json:"releases_url"`
	StargazersURL    string          `json:"stargazers_url"`
	StatusesURL      string          `json:"statuses_url"`
	SubscribersURL   string          `json:"subscribers_url"`
	SubscriptionURL  string          `json:"subscription_url"`
	TagsURL          string          `json:"tags_url"`
	TeamsURL         string          `json:"teams_url"`
	TreesURL         string          `json:"trees_url"`
	URL              string          `json:"url"`
}

type GithubCode_sub1 struct {
	AvatarURL         string `json:"avatar_url"`
	EventsURL         string `json:"events_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	GravatarID        string `json:"gravatar_id"`
	HTMLURL           string `json:"html_url"`
	ID                int64  `json:"id"`
	Login             string `json:"login"`
	NodeID            string `json:"node_id"`
	OrganizationsURL  string `json:"organizations_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	ReposURL          string `json:"repos_url"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	URL               string `json:"url"`
}

type GithubCode_sub3 struct {
	GitURL     string          `json:"git_url"`
	HTMLURL    string          `json:"html_url"`
	Name       string          `json:"name"`
	Path       string          `json:"path"`
	Repository GithubCode_sub2 `json:"repository"`
	Score      float64         `json:"score"`
	Sha        string          `json:"sha"`
	URL        string          `json:"url"`
}
