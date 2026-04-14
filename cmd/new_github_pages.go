package cmd

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"strings"
	"time"

	"github.com/spf13/cobra"

	ghclient "personal-platform/internal/github"
	"personal-platform/internal/secret"
	"personal-platform/templates"
)

// templateData is passed into index.html when rendering.
type templateData struct {
	RepoName string
}

var (
	flagDescription string
	flagOwner       string
)

var githubPagesCmd = &cobra.Command{
	Use:   "github-pages <repo-name>",
	Short: "Create a GitHub repo with automatic Pages deployment",
	Long: `Scaffolds a new public GitHub repository with:
  - A static index.html landing page
  - A GitHub Actions workflow that deploys to GitHub Pages on every push to main

Example:
  personal-platform new github-pages my-cool-site`,
	Args: cobra.ExactArgs(1),
	RunE: runGithubPages,
}

func init() {
	githubPagesCmd.Flags().StringVarP(
		&flagDescription, "description", "d", "",
		"Repository description",
	)
	githubPagesCmd.Flags().StringVarP(
		&flagOwner, "owner", "o", "",
		"GitHub username (overrides GITHUB_OWNER env var)",
	)
	newCmd.AddCommand(githubPagesCmd)
}

func runGithubPages(cmd *cobra.Command, args []string) error {
	repoName := args[0]
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// --- Resolve credentials ---
	store := secret.NewEnvSecretStore()

	token, err := store.Get("GITHUB_TOKEN")
	if err != nil {
		return fmt.Errorf("GitHub PAT not found: set the GITHUB_TOKEN environment variable\n%w", err)
	}

	owner := flagOwner
	if owner == "" {
		owner, err = store.Get("GITHUB_OWNER")
		if err != nil {
			return fmt.Errorf("GitHub owner not found: set --owner flag or GITHUB_OWNER environment variable\n%w", err)
		}
	}

	client := ghclient.NewClient(ctx, token, owner)

	// --- Step 1: Create the repo ---
	fmt.Printf("\u2192 Creating repository %s/%s...\n", owner, repoName)
	_, err = client.CreateRepo(ctx, repoName, flagDescription)
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}
	fmt.Printf("  \u2713 Repository created\n")

	// --- Step 2: Render and push files ---
	fmt.Printf("\u2192 Pushing template files...\n")
	if err := pushTemplateFiles(ctx, client, repoName); err != nil {
		return err
	}
	fmt.Printf("  \u2713 Files pushed\n")

	// --- Step 3: Enable GitHub Pages ---
	fmt.Printf("\u2192 Enabling GitHub Pages...\n")
	_, err = client.EnablePages(ctx, repoName)
	if err != nil {
		// Pages API can be flaky right after repo creation; warn but don't fatal.
		fmt.Printf("  \u26A0 Could not auto-enable Pages (you may need to enable it manually): %v\n", err)
	} else {
		fmt.Printf("  \u2713 GitHub Pages enabled\n")
	}

	// --- Done ---
	fmt.Printf("\n\u2728 Done! Your site will be live at:\n")
	fmt.Printf("   https://%s.github.io/%s\n\n", owner, repoName)
	fmt.Printf("   It may take a minute for the first Actions run to complete.\n")
	fmt.Printf("   Track it at: https://github.com/%s/%s/actions\n", owner, repoName)

	return nil
}

// pushTemplateFiles walks the embedded template FS, renders any .html files
// as Go templates, and pushes each file to the new repo.
func pushTemplateFiles(ctx context.Context, client *ghclient.Client, repoName string) error {
	data := templateData{RepoName: repoName}

	// The embed root is the github-pages folder inside package templates — strip it
	// so the files land at the repo root.
	const embedRoot = "github-pages"

	return fs.WalkDir(templates.GithubPages, embedRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		raw, err := templates.GithubPages.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read embedded file %q: %w", path, err)
		}

		// Render Go template substitutions for .html files only
		content := raw
		if strings.HasSuffix(path, ".html") {
			tmpl, err := template.New("").Parse(string(raw))
			if err != nil {
				return fmt.Errorf("parse template %q: %w", path, err)
			}
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				return fmt.Errorf("render template %q: %w", path, err)
			}
			content = buf.Bytes()
		}

		// Strip the embedRoot prefix to get the in-repo path
		repoPath := strings.TrimPrefix(path, embedRoot+"/")

		commitMsg := fmt.Sprintf("chore: scaffold %s via personal-platform CLI", repoPath)
		if err := client.PushFile(ctx, repoName, repoPath, commitMsg, content); err != nil {
			return err
		}
		fmt.Printf("    pushed: %s\n", repoPath)
		return nil
	})
}
