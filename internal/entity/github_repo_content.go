package entity

// GetRepoContentResponse represents the response structure for the GitHub API's "Get repository content" endpoint.
// This struct maps to the JSON response returned by the API.
//
// Reference: https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28
type GetRepoContentResponse struct {
	// Type specifies the type of the content (e.g., "file", "dir").
	Type string `json:"type"`

	// Encoding specifies the encoding of the content (e.g., "base64").
	Encoding string `json:"encoding"`

	// Size specifies the size of the content in bytes.
	Size int `json:"size"`

	// Name specifies the name of the file or directory.
	Name string `json:"name"`

	// Path specifies the path of the file or directory in the repository.
	Path string `json:"path"`

	// Content contains the actual content of the file, encoded in base64.
	// This field is only populated for files, not directories.
	Content string `json:"content"`

	// SHA is the Git blob SHA of the content.
	SHA string `json:"sha"`

	// URL is the API URL to fetch this content.
	URL string `json:"url"`

	// GitURL is the Git blob URL to fetch the raw content.
	GitURL string `json:"git_url"`

	// HTMLURL is the URL to view the content on GitHub's web interface.
	HTMLURL string `json:"html_url"`

	// DownloadURL is the URL to download the raw content.
	DownloadURL string `json:"download_url"`

	// Links contains hypermedia links related to the content.
	Links Links `json:"_links"`
}

// Links represents the hypermedia links included in the GitHub API response.
type Links struct {
	// Git is the URL to fetch the content as a Git blob.
	Git string `json:"git"`

	// Self is the API URL to fetch this content.
	Self string `json:"self"`

	// HTML is the URL to view the content on GitHub's web interface.
	HTML string `json:"html"`
}

// CreateOrUpdateFileContentRequest represents the request body for creating or updating file content.
type CreateOrUpdateFileContentRequest struct {
    Message   string    `json:"message"`
    Content   string    `json:"content"`
    SHA       string    `json:"sha,omitempty"`
    Branch    string    `json:"branch,omitempty"`
    Committer *Committer `json:"committer,omitempty"`
}

// Committer represents the person that committed the file.
type Committer struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
