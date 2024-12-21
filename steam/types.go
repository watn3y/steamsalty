package steam

type CommentsPage struct {
	Success      bool   `json:"success"`
	Start        int    `json:"start"`
	TotalCount   int    `json:"total_count"`
	CommentsHTML string `json:"comments_html"`
	TimeLastPost int64  `json:"timelastpost"`
}

type Comment struct {
	ID               uint64
	Timestamp        int64
	Author           string
	AuthorProfileURL string
	Text             string
}
