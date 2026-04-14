# personal-platform

A personal IDP CLI. Currently scaffolds GitHub Pages sites; more to come.

## Setup

```bash
go mod tidy
go build -o personal-platform .
```

Or install directly to your GOPATH:
```bash
go install .
```

## Usage

### new github-pages

Creates a public GitHub repo with a static landing page and a GitHub Actions
workflow that auto-deploys to GitHub Pages on every push to `main`.

```bash
# Required environment variables
export GITHUB_TOKEN=ghp_your_pat_here   # needs repo + pages scopes
export GITHUB_OWNER=your-github-username

personal-platform new github-pages my-cool-site
personal-platform new github-pages my-cool-site --description "My personal site"
personal-platform new github-pages my-cool-site --owner other-org
```

Your PAT needs the following scopes:
- `repo` (full control of repositories)
- (Pages is controlled via the repo scope)

The site will be live at `https://<owner>.github.io/<repo-name>` after the
first Actions run completes (~1 min).

## Project structure

```
personal-platform/
├── main.go
├── cmd/
│   ├── root.go               # cobra root command
│   ├── new.go                # 'new' parent command
│   └── new_github_pages.go   # 'new github-pages' subcommand + embed
├── internal/
│   ├── github/
│   │   └── client.go         # GitHub API wrapper
│   └── secret/
│       └── store.go          # SecretStore interface + env implementation
└── templates/
    └── github-pages/
        ├── index.html                     # Go template — {{.RepoName}} interpolated
        └── .github/workflows/pages.yml   # GitHub Actions deploy workflow
```

## Roadmap

- [ ] `personal-platform new github-pages` (static) ← **here now**
- [ ] `personal-platform new github-pages --type angular` (Angular + Node build step)
- [ ] `personal-platform token status` (JFrog PAT lifecycle in AWS Secrets Manager)
- [ ] AWS Secrets Manager `SecretStore` implementation
- [ ] `personal-platform new` interactive picker (fzf-style)
