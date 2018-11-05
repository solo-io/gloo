package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	goGit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

const (
	Author            = "gitbot"
	AuthorEmail       = "gitbot@solo.io"
	LocalBranchPrefix = "refs/heads/"

	gitDirName = ".git"
	repoTitle  = "Gloo state repository"
	readmeFile = "README.md"
)

var AbsPathNotInRepo = errors.New("the given absolute path does not point to a file in the repository")

type Repository struct {
	Root string
}

// Returns true if a non-bare git repository already exists
func (r *Repository) IsRepo() (bool, error) {
	_, err := os.Stat(path.Join(r.Root, gitDirName))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Create a working directory
func (r *Repository) Init() error {

	// For an explanation of bare vs non-bare, see here:
	// http://www.saintsjd.com/2011/01/what-is-a-bare-git-repository/
	repo, err := goGit.PlainInit(r.Root, false)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// create a file for first commit
	filePath, err := createReadMeFile(r)
	if err != nil {
		return err
	}

	// add it to the index
	if err = r.Add(filePath); err != nil {
		return err
	}

	_, err = w.Commit("First commit", &goGit.CommitOptions{Author: signature()})

	return err
}

// List all branches
func (r *Repository) ListBranches(includeRemotes bool) ([]string, error) {
	repo, err := goGit.PlainOpen(r.Root)
	if err != nil {
		return nil, err
	}

	refIterator, err := repo.Storer.IterReferences()
	if err != nil {
		return nil, err
	}

	// filter references
	refIterator = storer.NewReferenceFilteredIter(
		func(r *plumbing.Reference) bool {
			return r.Name().IsBranch() || (includeRemotes && r.Name().IsRemote())
		}, refIterator)

	// collect all branch names
	branches := make([]string, 0)
	refIterator.ForEach(func(b *plumbing.Reference) error {
		branches = append(branches, b.Name().String())
		return nil
	})

	return branches, nil
}

// Create a new branch
func (r *Repository) NewBranch(name string) error {
	repo, err := goGit.PlainOpen(r.Root)
	if err != nil {
		return err
	}

	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	// Reference name has to be a complete reference name (i.e. with the ref/heads/ prefix)
	ref := plumbing.NewHashReference(
		plumbing.ReferenceName(fmt.Sprint(LocalBranchPrefix, name)),
		headRef.Hash())

	return repo.Storer.SetReference(ref)
}

// Checkout a branch by name
func (r *Repository) Checkout(name string) error {
	workTree, _, err := r.getWorkTree()
	if err != nil {
		return err
	}

	return workTree.Checkout(&goGit.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprint(LocalBranchPrefix, name)),
	})
}

// Add a file or the content of a directory to the index.
func (r *Repository) Add(filePath string) error {
	filePath, err := validatePath(filePath, r)
	if err != nil {
		return err
	}

	workTree, _, err := r.getWorkTree()
	if err != nil {
		return err
	}

	_, err = workTree.Add(filePath)

	return err
}

// Commit index to repository.
//
// Returns a string representing the commit hash.
func (r *Repository) Commit(msg string) (string, error) {

	workTree, _, err := r.getWorkTree()
	if err != nil {
		return "", err
	}

	hash, err := workTree.Commit(msg, &goGit.CommitOptions{Author: signature()})
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

// Create a simple README file at the Root of the given repository
func createReadMeFile(r *Repository) (string, error) {
	filename := path.Join(r.Root, readmeFile)
	err := ioutil.WriteFile(filename, []byte(fmt.Sprint("# ", repoTitle, "\n")), 0644)
	return readmeFile, err
}

// Used to sign commits
func signature() *object.Signature {
	return &object.Signature{
		Name:  Author,
		Email: AuthorEmail,
		When:  time.Now(),
	}
}

// Checks whether the given path (absolute or relative) correctly points to a file/directory in the repository.
//
// Returns a path relative to the Root of the repository.
func validatePath(filePath string, r *Repository) (string, error) {

	// path is absolute
	if path.IsAbs(filePath) {
		if !strings.Contains(filePath, r.Root) {
			return "", AbsPathNotInRepo
		}

		filePath = strings.Replace(filePath, r.Root, "", -1)

		// Remove leading '/'
		if len(filePath) > 0 && filePath[0] == '/' {
			filePath = filePath[1:]
		}
	}

	// At this point the path is relative
	_, err := os.Stat(path.Join(r.Root, filePath))
	if err == nil {
		return filePath, nil
	} else {
		return "", err
	}
}

func (r *Repository) getWorkTree() (*goGit.Worktree, *goGit.Repository, error) {
	repo, err := goGit.PlainOpen(r.Root)
	if err != nil {
		return nil, nil, err
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return nil, nil, err
	}

	return workTree, repo, nil
}
