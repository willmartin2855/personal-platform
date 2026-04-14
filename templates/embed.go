package templates

import "embed"

// GithubPages contains the static site and workflow scaffold for `new github-pages`.
// The "all:" prefix includes dotfiles (e.g. .github/workflows).
//
//go:embed all:github-pages
var GithubPages embed.FS
