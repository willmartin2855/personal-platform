package github

import (
	"context"
	"encoding/base64"
	"fmt"

	gh "github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API for the operations we care about.
type Client struct {
	client *gh.Client
	owner  string
}

func NewClient(ctx context.Context, token, owner string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &Client{
		client: gh.NewClient(tc),
		owner:  owner,
	}
}

// CreateRepo creates a new public repo under the authenticated user.
func (c *Client) CreateRepo(ctx context.Context, name, description string) (*gh.Repository, error) {
	repo := &gh.Repository{
		Name:        gh.String(name),
		Description: gh.String(description),
		Private:     gh.Bool(false),
		AutoInit:    gh.Bool(false), // we'll push our own initial commit
	}
	created, _, err := c.client.Repositories.Create(ctx, "", repo)
	if err != nil {
		return nil, fmt.Errorf("create repo: %w", err)
	}
	return created, nil
}

// PushFile creates or updates a single file in the repo via the Contents API.
func (c *Client) PushFile(ctx context.Context, repo, path, commitMsg string, content []byte) error {
	encoded := base64.StdEncoding.EncodeToString(content)
	opts := &gh.RepositoryContentFileOptions{
		Message: gh.String(commitMsg),
		Content: []byte(encoded),
		Branch:  gh.String("main"),
	}

	_, _, err := c.client.Repositories.CreateFile(ctx, c.owner, repo, path, opts)
	if err != nil {
		return fmt.Errorf("push file %q: %w", path, err)
	}
	return nil
}

// EnablePages turns on GitHub Pages with the GitHub Actions source.
func (c *Client) EnablePages(ctx context.Context, repo string) (*gh.Pages, error) {
	source := &gh.PagesSource{
		// Empty branch/path means "use GitHub Actions" deployment
	}
	pages, _, err := c.client.Repositories.EnablePages(ctx, c.owner, repo, &gh.Pages{
		Source: source,
		BuildType: gh.String("workflow"),
	})
	if err != nil {
		return nil, fmt.Errorf("enable pages: %w", err)
	}
	return pages, nil
}
