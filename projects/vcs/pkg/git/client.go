package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/solo-io/solo-projects/projects/vcs/pkg/constants"
	"gopkg.in/src-d/go-git.v4/config"

	"github.com/solo-io/solo-kit/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

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

// TODO: enhance errors, include more info
var (
	AbsPathNotInRepo    = errors.Errorf("The given absolute path does not point to a file in the repository")
	CloneInExistingRepo = errors.Errorf("Cannot clone into an existing repository")
	InvalidBranchName   = errors.Errorf("The branch name can contain only alphanumeric characters, hyphens, and underscores")

	branchNameRegExp = regexp.MustCompile(constants.BranchRegExp)
)

type Repository struct {
	root string
	auth transport.AuthMethod
}

// Creates a repository in the default temporary files directory
func NewTempRepo() (*Repository, error) {
	tempDir, err := ioutil.TempDir("", constants.AppName)
	if err != nil {
		return &Repository{}, err
	}
	return NewRepo(tempDir)
}

// Creates a repository in the given directory
func NewRepo(root string) (*Repository, error) {
	_, err := os.Stat(root)
	if err != nil {
		return &Repository{}, err
	}
	return &Repository{root: root}, nil
}

// Deletes the directory that contains the repository
func (r *Repository) Delete() error {
	return os.RemoveAll(r.root)
}

// Configure client for token authentication with remote
func (r *Repository) WithTokenAuth(token string) *Repository {
	r.auth = &http.BasicAuth{Username: "gitbot", Password: token} // username just has to be a non-empty string
	return r
}

// Configure client for basic authentication with remote
func (r *Repository) WithBasicAuth(username, password string) *Repository {
	r.auth = &http.BasicAuth{Username: username, Password: password}
	return r
}

func (r *Repository) Root() string {
	return r.root
}

// Returns true if a non-bare git repository already exists
func (r *Repository) IsRepo() (bool, error) {
	_, err := os.Stat(path.Join(r.root, gitDirName))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Initialize the repository
// Creates a first commit containing a dummy README file on the master branch
func (r *Repository) Init() (string, error) {

	// For an explanation of bare vs non-bare, see here:
	// http://www.saintsjd.com/2011/01/what-is-a-bare-git-repository/
	repo, err := goGit.PlainInit(r.root, false)
	if err != nil {
		return "", err
	}

	w, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	// create a file for first commit
	filePath, err := createReadMeFile(r)
	if err != nil {
		return "", err
	}

	// add it to the index
	if err = r.Add(filePath); err != nil {
		return "", err
	}

	hash, err := w.Commit("First commit", &goGit.CommitOptions{Author: signature()})

	return hash.String(), err
}

// List all branches.
// If the parameter is 'true', this will include remote-tracking branches that do not have local ref. Clients should
// never have to use this feature. It is mainly present for tests.
//
// Note: this is a slightly modified version of go-git.Repository.Branches()
func (r *Repository) ListBranches(includeRemotes bool) ([]string, error) {
	repo, err := goGit.PlainOpen(r.root)
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

// Create a new local branch starting from the current HEAD
func (r *Repository) NewBranch(name string) error {
	if !branchNameRegExp.MatchString(name) {
		return InvalidBranchName
	}

	repo, err := goGit.PlainOpen(r.root)
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

// Returns the reference where HEAD is pointing to
func (r *Repository) Head() (string, error) {
	repo, err := goGit.PlainOpen(r.root)
	if err != nil {
		return "", err
	}

	headRef, err := repo.Head()
	if err != nil {
		return "", err
	}

	if headRef.Type() != plumbing.HashReference {
		return "", errors.Errorf("Unsupported reference type [%v] for reference [%v]", headRef.Type(), headRef.Hash())
	}

	return headRef.Hash().String(), nil
}

// Returns the hash and message for the last commit
func (r *Repository) LastCommit() (hash, message string, e error) {
	repo, err := goGit.PlainOpen(r.root)
	if err != nil {
		return "", "", err
	}
	commitIter, err := repo.Log(&goGit.LogOptions{})
	if err != nil {
		return "", "", err
	}
	commit, err := commitIter.Next()
	if err != nil {
		return "", "", err
	}

	return commit.Hash.String(), commit.Message, nil
}

// Checkout a reference by name
// Name must be a short refname (without the refs/... prefix)
func (r *Repository) CheckoutBranch(name string) error {
	workTree, _, err := r.getWorkTree()
	if err != nil {
		return err
	}

	if !branchNameRegExp.MatchString(name) {
		return InvalidBranchName
	}

	return workTree.Checkout(&goGit.CheckoutOptions{
		Branch: plumbing.ReferenceName(LocalBranchPrefix + name),
	})
}

// Checkout a commit by its hash. HEAD will be in detached mode.
// The hash parameter must be the full 40-byte hexadecimal commit object name.
func (r *Repository) CheckoutCommit(hash string) error {
	workTree, _, err := r.getWorkTree()
	if err != nil {
		return err
	}
	return workTree.Checkout(&goGit.CheckoutOptions{
		Hash: plumbing.NewHash(hash),
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

// Clone the repository at the given URL
func (r *Repository) Clone(remoteUrl string) error {

	if isRepo, _ := r.IsRepo(); isRepo {
		return CloneInExistingRepo
	}

	repo, err := goGit.PlainClone(r.root, false, &goGit.CloneOptions{URL: remoteUrl, Auth: r.auth})
	if err != nil {
		return err
	}

	// The ref spec creates local references for all the remote references.
	//
	// If e.g. the remote contains
	// [
	// 		refs/heads/master,
	// 		refs/heads/branch_1
	// ]
	// after this method is called the local repo will contain
	// [
	// 		refs/heads/master,
	// 		refs/heads/branch_1,
	// 		refs/remotes/origin/master,
	// 		refs/remotes/origin/branch_1
	// ]
	// This allows us to avoid distinguishing between local and remote-tracking references during checkout.
	err = repo.Fetch(&goGit.FetchOptions{
		RefSpecs: []config.RefSpec{"+refs/heads/*:refs/heads/*"}, Auth: r.auth,
	})

	return err
}

func (r *Repository) Push(remoteUrl string) error {
	repo, err := goGit.PlainOpen(r.root)
	if err != nil {
		return err
	}
	return repo.Push(&goGit.PushOptions{Auth: r.auth})
}

// Create a simple README file at the root of the given repository
func createReadMeFile(r *Repository) (string, error) {
	filename := path.Join(r.root, readmeFile)
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
// Returns a path relative to the root of the repository.
func validatePath(filePath string, r *Repository) (string, error) {

	// path is absolute
	if path.IsAbs(filePath) {
		if !strings.Contains(filePath, r.root) {
			return "", AbsPathNotInRepo
		}

		filePath = strings.Replace(filePath, r.root, "", -1)

		// Remove leading '/'
		if len(filePath) > 0 && filePath[0] == '/' {
			filePath = filePath[1:]
		}
	}

	// At this point the path is relative
	_, err := os.Stat(path.Join(r.root, filePath))
	if err == nil {
		return filePath, nil
	} else {
		return "", err
	}
}

func (r *Repository) getWorkTree() (*goGit.Worktree, *goGit.Repository, error) {
	repo, err := goGit.PlainOpen(r.root)
	if err != nil {
		return nil, nil, err
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return nil, nil, err
	}

	return workTree, repo, nil
}
