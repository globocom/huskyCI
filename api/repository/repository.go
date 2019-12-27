package repository

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/src-d/enry/v2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Repository holds all information of a repository.
type Repository struct {
	GitInfo    *git.Repository
	URL        string `bson:"repositoryURL" json:"repositoryURL"`
	Branch     string `bson:"repositoryBranch" json:"repositoryBranch"`
	HeadCommit string
	Commits    []Commit
	Files      []File
	Languages  []string
}

// Commit holds all information of a commit.
type Commit struct {
	Hash        string
	Author      string
	Description string
	Titte       string
	Date        time.Time
}

// File holds all the information of a file.
type File struct {
	Name string
	Hash string
}

// Scan will process all information of a repositoy based on its URL and branch.
func (r *Repository) Scan() error {

	if err := r.clone(); err != nil {
		return err
	}

	if err := r.setCommits(); err != nil {
		return err
	}

	if err := r.setFiles(); err != nil {
		return err
	}

	if err := r.setLanguages(); err != nil {
		return err
	}

	return nil
}

// clone clones the repository by setting the GitInfo field from the Repository struct.
func (r *Repository) clone() error {

	repoInfo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           r.URL,
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", r.Branch)),
	})
	if err != nil {
		return err
	}

	r.GitInfo = repoInfo
	return nil
}

// setCommits sets the head commit and all others of a repository and its branch after clonning it.
func (r *Repository) setCommits() error {

	ref, err := r.GitInfo.Head()
	if err != nil {
		return err
	}

	headCommit, err := r.GitInfo.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	r.HeadCommit = headCommit.Hash.String()

	commitIter, err := r.GitInfo.Log(&git.LogOptions{From: headCommit.Hash})
	if err != nil {
		return err
	}

	err = commitIter.ForEach(func(c *object.Commit) error {
		commitFound := Commit{
			Hash:        c.Hash.String(),
			Description: c.Message,
			Author:      c.Author.Email,
			Date:        c.Author.When,
		}
		r.Commits = append(r.Commits, commitFound)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// setFiles set all files found of a repository and its branch after clonning it.
func (r *Repository) setFiles() error {

	ref, err := r.GitInfo.Head()
	if err != nil {
		return err
	}

	headCommit, err := r.GitInfo.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	r.HeadCommit = headCommit.Hash.String()

	tree, err := headCommit.Tree()
	if err != nil {
		return err
	}

	tree.Files().ForEach(func(f *object.File) error {
		fileFound := File{Name: f.Name, Hash: f.Hash.String()}
		r.Files = append(r.Files, fileFound)
		return nil
	})

	return nil
}

// setLanguages will check every file from the repository and set the languages found using enry.
func (r *Repository) setLanguages() error {

	for _, file := range r.Files {
		lang, _ := enry.GetLanguageByExtension(file.Name)
		r.Languages = appendIfMissing(r.Languages, lang)
	}

	return nil
}

// CheckInput checks if URL and branch are safe.
func (r *Repository) CheckInput() error {

	if err := checkMaliciousRepoURL(r.URL); err != nil {
		return err
	}

	if err := checkMaliciousRepoBranch(r.Branch); err != nil {
		return err
	}

	return nil
}

// AppendIfMissing will append a string inside a slice of string if it is unique
func appendIfMissing(slice []string, s string) []string {

	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}

	return append(slice, s)
}

func checkMaliciousRepoURL(repositoryURL string) error {

	regexpGit := `((git|ssh|http(s)?)|((git@|gitlab@)[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`

	valid, err := regexp.MatchString(regexpGit, repositoryURL)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid URL format")
	}

	return nil
}

func checkMaliciousRepoBranch(repositoryBranch string) error {

	regexpBranch := `^[a-zA-Z0-9_\/.-]*$`

	valid, err := regexp.MatchString(regexpBranch, repositoryBranch)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid branch format")
	}

	return nil
}
