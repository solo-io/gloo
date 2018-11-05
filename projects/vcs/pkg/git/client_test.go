package git_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/git"
	goGit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var _ = Describe("Git Client", func() {

	var (
		repoRoot   string
		repoClient *git.Repository
		err        error
	)

	// Initialize git repository in os.TempDir
	BeforeEach(func() {
		repoRoot, err = ioutil.TempDir("", "go_git_client_test")
		Expect(err).To(BeNil())
		Expect(repoRoot).To(BeADirectory())

		repoClient = &git.Repository{Root: repoRoot}
		Expect(repoClient).NotTo(BeNil())

		err = repoClient.Init()
	})

	// Clean up
	AfterEach(func() {
		os.RemoveAll(repoRoot)
	})

	Describe("initializing the client", func() {

		It("does not generate an error", func() {
			Expect(err).To(BeNil())
		})

		It("creates a non-bare repository", func() {
			isRepo, err := repoClient.IsRepo()
			Expect(err).To(BeNil())
			Expect(isRepo).To(BeTrue())
		})

		It("creates a master branch", func() {
			repo, err := goGit.PlainOpen(repoRoot)
			Expect(err).To(BeNil())

			refIterator, err := repo.Branches()
			Expect(err).To(BeNil())

			branches := make([]string, 0)
			refIterator.ForEach(func(b *plumbing.Reference) error {
				branches = append(branches, b.Name().String())
				return nil
			})

			Expect(branches).To(ConsistOf(git.LocalBranchPrefix + "master"))
		})

		It("generated a README file", func() {
			_, err = os.Stat(path.Join(repoRoot, "README.md"))
			Expect(err).To(BeNil())
		})
	})

	// TODO: update when adding remotes
	Describe("listing branches", func() {

		It("lists the existing local branches", func() {
			branches, err := repoClient.ListBranches(false)
			Expect(err).To(BeNil())
			Expect(branches).To(ConsistOf(git.LocalBranchPrefix + "master"))
		})
	})

	Describe("committing a file", func() {

		var (
			fileName  = "test.txt"
			fileDir   string
			commitMsg = "Added file: " + fileName
			hash      string
		)

		BeforeEach(func() {
			fileDir = repoRoot
		})

		JustBeforeEach(func() {
			filePath := path.Join(fileDir, fileName)
			err = ioutil.WriteFile(filePath, []byte(fmt.Sprint("Some random content...\n")), 0644)
			Expect(err).To(BeNil())

			// Stage the file
			err = repoClient.Add(filePath)
		})

		Context("when the file is staged for commit", func() {

			It("does not generate an error", func() {
				Expect(err).To(BeNil())
			})

			It("file is present in the index", func() {

				repo, err := goGit.PlainOpen(repoRoot)
				Expect(err).To(BeNil())

				workTree, err := repo.Worktree()
				Expect(err).To(BeNil())

				status, err := workTree.Status()
				Expect(err).To(BeNil())
				Expect(status.IsClean()).To(BeFalse())
				Expect(status.File(fileName).Staging).To(BeEquivalentTo(goGit.Added))
			})
		})

		Context("when we try to stage a file that is not in the repository", func() {

			BeforeEach(func() {
				fileDir = os.TempDir()
			})

			It("generates an error", func() {
				Expect(err).To(HaveOccurred())
			})

			AfterEach(func() {
				os.RemoveAll(fileDir)
			})

		})

		Describe("when the file is committed", func() {
			BeforeEach(func() {
				hash, err = repoClient.Commit(commitMsg)
			})

			It("does not generate an error", func() {
				Expect(err).To(BeNil())
			})

			It("does return a non-empty hash", func() {
				Expect(hash).To(Not(BeEmpty()))
			})

			It("correctly stores the commit", func() {
				repo, err := goGit.PlainOpen(repoRoot)
				Expect(err).To(BeNil())

				commit, err := repo.CommitObject(plumbing.NewHash(hash))
				Expect(err).To(BeNil())
				Expect(commit).To(Not(BeNil()))
				Expect(commit.Message).To(BeEquivalentTo(commitMsg))
				Expect(commit.Author.Name).To(BeEquivalentTo(git.Author))
				Expect(commit.Author.Email).To(BeEquivalentTo(git.AuthorEmail))
			})
		})
	})

	Describe("committing a directory structure", func() {

		var file1, file2, file3, hash string

		BeforeEach(func() {
			err = os.MkdirAll(repoRoot+"/dir1/dir2", os.ModePerm)
			Expect(err).To(BeNil())

			file1 = path.Join(repoRoot, "dir1", "file_1_1")
			file2 = path.Join(repoRoot, "dir1", "file_1_2")
			file3 = path.Join(repoRoot, "dir1", "dir2", "file_2_1")

			err = ioutil.WriteFile(file1, []byte(fmt.Sprint("Some random content...\n")), 0644)
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(file2, []byte(fmt.Sprint("Some random content...\n")), 0644)
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(file3, []byte(fmt.Sprint("Some random content...\n")), 0644)
			Expect(err).To(BeNil())

			err = repoClient.Add(repoRoot)
			Expect(err).To(BeNil())
		})

		It("all files have been added correctly", func() {

			repo, err := goGit.PlainOpen(repoRoot)
			Expect(err).To(BeNil())

			workTree, err := repo.Worktree()
			Expect(err).To(BeNil())

			status, err := workTree.Status()
			Expect(err).To(BeNil())
			Expect(status.IsClean()).To(BeFalse())

			Expect(status.File(strings.Replace(file1, repoRoot, "", -1)[1:]).Staging).To(BeEquivalentTo(goGit.Added))
			Expect(status.File(strings.Replace(file2, repoRoot, "", -1)[1:]).Staging).To(BeEquivalentTo(goGit.Added))
			Expect(status.File(strings.Replace(file3, repoRoot, "", -1)[1:]).Staging).To(BeEquivalentTo(goGit.Added))
		})

		Describe("when the contents of the index are committed", func() {
			BeforeEach(func() {
				hash, err = repoClient.Commit("Multiple files")
			})

			It("does not generate an error", func() {
				Expect(err).To(BeNil())
			})

			It("does return a non-empty hash", func() {
				Expect(hash).To(Not(BeEmpty()))
			})

			It("generates a valid commit", func() {
				repo, err := goGit.PlainOpen(repoRoot)
				Expect(err).To(BeNil())

				commit, err := repo.CommitObject(plumbing.NewHash(hash))
				Expect(err).To(BeNil())
				Expect(commit).To(Not(BeNil()))
				Expect(commit.Message).To(BeEquivalentTo("Multiple files"))
				Expect(commit.Author.Name).To(BeEquivalentTo(git.Author))
				Expect(commit.Author.Email).To(BeEquivalentTo(git.AuthorEmail))
			})
		})

	})

	Describe("creating a new branch", func() {

		const branchName = "newBranch"

		BeforeEach(func() {
			err = repoClient.NewBranch(branchName)
		})

		It("does not generate an error", func() {
			Expect(err).To(BeNil())
		})

		It("the branch has been successfully created", func() {
			repo, _ := goGit.PlainOpen(repoRoot)
			refIterator, _ := repo.Branches()

			branches := make([]string, 0)
			refIterator.ForEach(func(b *plumbing.Reference) error {
				branches = append(branches, b.Name().String())
				return nil
			})

			Expect(err).To(BeNil())
			Expect(branches).To(Not(BeNil()))
			Expect(branches).To(ConsistOf(
				git.LocalBranchPrefix+"master",
				git.LocalBranchPrefix+branchName,
			))
		})

		Describe("checking out the new branch", func() {

			BeforeEach(func() {
				err = repoClient.Checkout(branchName)
			})

			It("does not generate an error", func() {
				Expect(err).To(BeNil())
			})

			It("the correct branch has been checked out", func() {
				repo, err := goGit.PlainOpen(repoRoot)
				Expect(err).To(BeNil())

				ref, err := repo.Head()
				Expect(err).To(BeNil())

				Expect(ref.Name()).To(BeEquivalentTo(git.LocalBranchPrefix + branchName))
			})
		})

		Describe("new file committed to new branch", func() {

			const (
				fileName  = "only_on_new_branch.txt"
				commitMsg = "Another file: " + fileName
			)
			var filePath, hash string

			BeforeEach(func() {

				_ = repoClient.Checkout(branchName)

				filePath = path.Join(repoRoot, fileName)
				err = ioutil.WriteFile(filePath, []byte(fmt.Sprint("I am just using up memory...\n")), 0644)
				Expect(err).To(BeNil())

				repoClient.Add(filePath)
				hash, _ = repoClient.Commit(commitMsg)
			})

			It("commit has been created", func() {
				repo, err := goGit.PlainOpen(repoRoot)
				Expect(err).To(BeNil())

				commit, err := repo.CommitObject(plumbing.NewHash(hash))
				Expect(err).To(BeNil())
				Expect(commit).To(Not(BeNil()))
			})

			It("file is not present on master branch", func() {
				repoClient.Checkout("master")
				_, err = os.Stat(path.Join(repoRoot, fileName))
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})
	})
})
