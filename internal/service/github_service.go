package service

import (
	"dantaautotool/internal/entity"
	"dantaautotool/pkg/utils/http"
	"encoding/base64"
	"fmt"
	"os"

	"maps"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/rs/zerolog/log"
)

// GithubServiceIntf defines the interface for GithubService.
type GithubServiceIntf interface {

	// GetFileContent retrieves the content of a file given its path.
	// It returns the content and an error if any occurs.
	GetFileContent(owner, repo, path string) (*entity.RepoContent, error)
	
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
			"Accept": "application/vnd.github+json",
			"Authorization": "token " + pat,
			"X-GitHub-Api-Version": "2022-11-28",
		},
	}
}

// GetFileContent retrieves the content of a file given its path.
// It returns the content and an error if any occurs.
func (s *GithubService) GetFileContent(owner, repo, path string) (*entity.RepoContent, error) {
	headers := make(map[string]string)
	maps.Copy(headers, s.authHeaders)
	pathParams := map[string]string{
		"owner": owner,
		"repo":  repo,
		// "path":  path,
	}
	queryParams := map[string]string{}

	// To avoid '/' in path being encoded, we need to put it in path in advance
	_, _, respBytes, err := s.client.PerformGet(fmt.Sprintf("/repos/{owner}/{repo}/contents/%s", path), headers, pathParams, queryParams)
	if err != nil {
		log.Err(err).Msg("[GetFileContent] Failed to get file content")
		return nil, err
	}

	var getFileContentResp entity.GetRepoContentResponse
	err = sonic.Unmarshal(respBytes, &getFileContentResp)
	if err != nil {
		log.Err(err).Msg("[GetFileContent] Failed to unmarshal response")
		return nil, err
	}

	// Decode the content from base64
	repoContent := entity.RepoContent{}
	repoContent.GetRepoContentResponse = getFileContentResp
	if getFileContentResp.Encoding == "base64" {
		decodedContent, err := base64.StdEncoding.DecodeString(getFileContentResp.Content)
		if err != nil {
			log.Err(err).Msg("[GetFileContent] Failed to decode base64 content")
			return nil, err
		}
		repoContent.DecodedContent = string(decodedContent)
	} else {
		log.Error().Msg("[GetFileContent] Content encoding is not base64")
		return nil, fmt.Errorf("content encoding is not base64")
	}

	return &repoContent, nil
}

// CreateOrUpdateFileContent creates or updates the content of a file given its path.
// It returns an error if any occurs.
func (s *GithubService) CreateOrUpdateFileContent(owner, repo, path, message, content, sha, branch string, committer *entity.Committer) error {
    headers := make(map[string]string)
	maps.Copy(headers, s.authHeaders)
    pathParams := map[string]string{
        "owner": owner,
        "repo":  repo,
        // "path":  path,
    }
    queryParams := map[string]string{}

	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))

    body := entity.CreateOrUpdateFileContentRequest{
        Message:   message,
        Content:   encodedContent,
        SHA:       sha,
        Branch:    branch,
        Committer: committer,
    }

    bodyBytes, err := sonic.Marshal(body)
    if err != nil {
        log.Error().Err(err).Msg("[CreateOrUpdateFileContent] Failed to marshal request body")
        return err
    }
	
    statusCode, _, respBodyBytes, err := s.client.PerformRequest(fmt.Sprintf("/repos/{owner}/{repo}/contents/%s", path), consts.MethodPut, headers, pathParams, queryParams, bodyBytes)
    if err != nil {
        log.Error().Err(err).Msg("[CreateOrUpdateFileContent] Failed to create or update file content")
        return err
    }
	if statusCode != consts.StatusOK {
		respBody := string(respBodyBytes)
		log.Error().Err(err).Msgf("[CreateOrUpdateFileContent] Failed to create or update file content, status code: %d, response: %s", statusCode, respBody)
		return fmt.Errorf("failed to create or update file content, status code: %d", statusCode)
	}

    return nil
}