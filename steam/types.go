package steam

type CommentResponse struct {
	Success      bool   `json:"success"`
	Start        int    `json:"start"`
	TotalCount   int    `json:"total_count"`
	CommentsHTML string `json:"comments_html"`
	Timelastpost int    `json:"timelastpost"`
}
