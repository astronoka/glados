package github

import "time"

// GitUser is user object
type GitUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"url"`
	HTMLURL   string `json:"html_url"`
}

// GitPullRequest is pull request object
type GitPullRequest struct {
	ID             int64      `json:"id"`
	URL            string     `json:"url"`
	HTMLURL        string     `json:"html_url"`
	Number         int64      `json:"number"`
	State          string     `json:"state"`
	Title          string     `json:"title"`
	Body           string     `json:"body"`
	User           GitUser    `json:"user"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	MergedAt       *time.Time `json:"merged_at"`
	MergeCommitSha string     `json:"merge_commit_sha"`
}

// GitRepository is repository object
type GitRepository struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	Owner     GitUser   `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GitIssue is issue object
type GitIssue struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	HTMLURL   string    `json:"html_url"`
	Number    int64     `json:"number"`
	State     string    `json:"state"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	User      GitUser   `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GitComment is comment object
type GitComment struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	HTMLURL   string    `json:"html_url"`
	Body      string    `json:"body"`
	User      GitUser   `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
