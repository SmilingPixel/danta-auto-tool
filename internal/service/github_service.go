package service

import (
	"dantaautotool/internal/entity"
	"dantaautotool/pkg/utils/http"
	"os"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/rs/zerolog/log"
)

// GithubServiceIntf defines the interface for GithubService.
type GithubServiceIntf interface {

	// GetFileContent retrieves the content of a file given its path.
	// It returns the content as a string and an error if any occurs.
	GetFileContent(owner, repo, path string) (string, string, error)
	
	// CreateOrUpdateFileContent creates or updates the content of a file given its path.
    // It returns an error if any occurs.
    CreateOrUpdateFileContent(owner, repo, path, message, content, sha, branch string, committer *entity.Committer) error
}

// GithubService provides methods to interact with Github.
type GithubService struct {
	// client is used to interact with Github
	client *http.HTTPClient

	// config for github authentication
	authHeaders map[string]string
}

// NewGithubService creates a new instance of GithubService.
func NewGithubService() *GithubService {
	cli := http.NewHTTPClient(
		"https://api.github.com",
		[]string{},
		http.EmptyHTTPClientMiddlewareSlice(),
	)

	pat := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	if pat == "" {
		log.Error().Msg("[NewGithubService] GITHUB_PERSONAL_ACCESS_TOKEN is empty")
		return nil
	}

	return &GithubService{
		client:        cli,
		authHeaders:   map[string]string{
			"Authorization": "token " + pat,
			"X-GitHub-Api-Version": "2022-11-28",
		},
	}
}

// GetFileContent retrieves the content of a file given its path.
// It returns the content and sha as strings and an error if any occurs.
func (s *GithubService) GetFileContent(owner, repo, path string) (string, string, error) {
	headers := make(map[string]string)
	for k, v := range s.authHeaders {
		headers[k] = v
	}
	pathParams := map[string]string{
		"owner": owner,
		"repo":  repo,
		"path":  path,
	}
	queryParams := map[string]string{}

	_, _, respBytes, err := s.client.PerformGet("/repos/{owner}/{repo}/contents/{path}", headers, pathParams, queryParams)
	if err != nil {
		log.Error().Err(err).Msg("[GetFileContent] Failed to get file content")
		return "", "", err
	}

	var getFileContentResp entity.GetRepoContentResponse
	err = sonic.Unmarshal(respBytes, &getFileContentResp)
	if err != nil {
		log.Error().Err(err).Msg("[GetFileContent] Failed to unmarshal response")
		return "", "", err
	}

	fileContent := getFileContentResp.Content
	sha := getFileContentResp.SHA
	return fileContent, sha, nil
}

// CreateOrUpdateFileContent creates or updates the content of a file given its path.
// It returns an error if any occurs.
func (s *GithubService) CreateOrUpdateFileContent(owner, repo, path, message, content, sha, branch string, committer *entity.Committer) error {
    headers := make(map[string]string)
	for k, v := range s.authHeaders {
		headers[k] = v
	}
    pathParams := map[string]string{
        "owner": owner,
        "repo":  repo,
        "path":  path,
    }
    queryParams := map[string]string{}

    body := entity.CreateOrUpdateFileContentRequest{
        Message:   message,
        Content:   content,
        SHA:       sha,
        Branch:    branch,
        Committer: committer,
    }

    bodyBytes, err := sonic.Marshal(body)
    if err != nil {
        log.Error().Err(err).Msg("[CreateOrUpdateFileContent] Failed to marshal request body")
        return err
    }

    _, _, _, err = s.client.PerformRequest(consts.MethodPut, "/repos/{owner}/{repo}/contents/{path}", headers, pathParams, queryParams, bodyBytes)
    if err != nil {
        log.Error().Err(err).Msg("[CreateOrUpdateFileContent] Failed to create or update file content")
        return err
    }

    return nil
}