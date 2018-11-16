package git_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/git"
	goGit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var _ = Describe("Git Client", func() {

	var (
		repoClient       *git.Repository
		masterCommitHash string
		err              error
	)

	Context("working with a local repository", func() {

		// Initialize git repository in a random temp directory
		BeforeEach(func() {
			root, err := ioutil.TempDir("", "go_git_client_test")
			Expect(err).To(BeNil())
			Expect(root).To(BeADirectory())

			repoClient, err = git.NewRepo(root)
			Expect(err).To(BeNil())
			Expect(repoClient).NotTo(BeNil())

			masterCommitHash, err = repoClient.Init()
		})

		// Clean up
		AfterEach(func() {
			os.RemoveAll(repoClient.Root())
		})

		Describe("initializing a repository", func() {

			It("does not generate an error", func() {
				Expect(err).To(BeNil())
			})

			It("generated a commit hash", func() {
				Expect(masterCommitHash).To(Not(BeEmpty()))
			})

			It("creates a non-bare repository", func() {
				isRepo, err := repoClient.IsRepo()
				Expect(err).To(BeNil())
				Expect(isRepo).To(BeTrue())
			})

			It("creates a master branch", func() {
				repo, err := goGit.PlainOpen(repoClient.Root())
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
				_, err = os.Stat(path.Join(repoClient.Root(), "README.md"))
				Expect(err).To(BeNil())
			})
		})

		Describe("listing branches", func() {

			It("lists the existing local branches", func() {
				branches, err := repoClient.ListBranches(false)
				Expect(err).To(BeNil())
				Expect(branches).To(ConsistOf(git.LocalBranchPrefix + "master"))
			})
		})

		Describe("committing a file", func() {

			var (
				fileName = "test.txt"
				fileDir  string
			)

			BeforeEach(func() {
				fileDir = repoClient.Root()
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
					repo, err := goGit.PlainOpen(repoClient.Root())
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
					fileDir, _ = ioutil.TempDir("", "go_git_client_test")
				})

				It("generates an error", func() {
					Expect(err).To(HaveOccurred())
				})

				AfterEach(func() {
					os.RemoveAll(fileDir)
				})
			})

			Describe("when the file is committed", func() {

				var (
					commitMsg = "Added file: " + fileName
					hash      string
				)

				// We need a JustBeforeEach instead of a BeforeEach to execute these preconditions in the correct order
				// given the JustBeforeEach in the enclosing function
				JustBeforeEach(func() {
					hash, err = repoClient.Commit(commitMsg)
				})

				It("does not generate an error", func() {
					Expect(err).To(BeNil())
				})

				It("does return a non-empty hash", func() {
					Expect(hash).To(Not(BeEmpty()))
				})

				It("correctly stores the commit", func() {
					repo, err := goGit.PlainOpen(repoClient.Root())
					Expect(err).To(BeNil())

					commit, err := repo.CommitObject(plumbing.NewHash(hash))
					Expect(err).To(BeNil())
					Expect(commit).To(Not(BeNil()))
					Expect(commit.Message).To(BeEquivalentTo(commitMsg))
					Expect(commit.Author.Name).To(BeEquivalentTo(git.Author))
					Expect(commit.Author.Email).To(BeEquivalentTo(git.AuthorEmail))
				})

				It("the working tree is clean", func() {
					repo, err := goGit.PlainOpen(repoClient.Root())
					Expect(err).To(BeNil())

					workTree, err := repo.Worktree()
					Expect(err).To(BeNil())

					status, err := workTree.Status()
					Expect(err).To(BeNil())
					Expect(status.IsClean()).To(BeTrue())
				})

				Describe("when we check the current HEAD", func() {

					var headRef string

					// We need a JustBeforeEach instead of a BeforeEach to execute these preconditions in the correct order
					// given the JustBeforeEach in the enclosing function
					JustBeforeEach(func() {
						headRef, err = repoClient.Head()
					})

					It("does not generate an error", func() {
						Expect(err).To(BeNil())
					})

					It("points to the right commit", func() {
						Expect(headRef).To(BeEquivalentTo(hash))
					})
				})

				Describe("when we check the last commit", func() {

					var lastHash, lastMsg string

					// We need a JustBeforeEach instead of a BeforeEach to execute these preconditions in the correct order
					// given the JustBeforeEach in the enclosing function
					JustBeforeEach(func() {
						lastHash, lastMsg, err = repoClient.LastCommit()
					})

					It("does not generate an error", func() {
						Expect(err).To(BeNil())
					})

					It("points to the right commit", func() {
						Expect(lastHash).To(BeEquivalentTo(hash))
					})

					It("returns the right commit message", func() {
						Expect(lastMsg).To(BeEquivalentTo(commitMsg))
					})
				})

				Describe("when we checkout the previous commit by its hash", func() {

					// We need a JustBeforeEach instead of a BeforeEach to execute these preconditions in the correct order
					// given the JustBeforeEach in the enclosing function
					JustBeforeEach(func() {
						err = repoClient.CheckoutCommit(masterCommitHash)
					})

					It("does not generate an error", func() {
						Expect(err).To(BeNil())
					})

					It("the new HEAD points to the right commit", func() {
						newHead, err := repoClient.Head()
						Expect(err).To(BeNil())
						Expect(newHead).To(BeEquivalentTo(masterCommitHash))
					})

					It("the file we just committed does not exists", func() {
						_, err = os.Stat(path.Join(fileDir, fileName))
						Expect(os.IsNotExist(err)).To(BeTrue())
					})
				})
			})
		})

		Describe("committing a directory structure", func() {

			var file1, file2, file3, hash string

			BeforeEach(func() {
				err = os.MkdirAll(repoClient.Root()+"/dir1/dir2", os.ModePerm)
				Expect(err).To(BeNil())

				file1 = path.Join(repoClient.Root(), "dir1", "file_1_1")
				file2 = path.Join(repoClient.Root(), "dir1", "file_1_2")
				file3 = path.Join(repoClient.Root(), "dir1", "dir2", "file_2_1")

				err = ioutil.WriteFile(file1, []byte(fmt.Sprint("Some random content...\n")), 0644)
				Expect(err).To(BeNil())
				err = ioutil.WriteFile(file2, []byte(fmt.Sprint("Some random content...\n")), 0644)
				Expect(err).To(BeNil())
				err = ioutil.WriteFile(file3, []byte(fmt.Sprint("Some random content...\n")), 0644)
				Expect(err).To(BeNil())

				err = repoClient.Add(repoClient.Root())
				Expect(err).To(BeNil())
			})

			It("all files have been added correctly", func() {
				repo, err := goGit.PlainOpen(repoClient.Root())
				Expect(err).To(BeNil())

				workTree, err := repo.Worktree()
				Expect(err).To(BeNil())

				status, err := workTree.Status()
				Expect(err).To(BeNil())
				Expect(status.IsClean()).To(BeFalse())

				Expect(status.File(strings.Replace(file1, repoClient.Root(), "", -1)[1:]).Staging).To(BeEquivalentTo(goGit.Added))
				Expect(status.File(strings.Replace(file2, repoClient.Root(), "", -1)[1:]).Staging).To(BeEquivalentTo(goGit.Added))
				Expect(status.File(strings.Replace(file3, repoClient.Root(), "", -1)[1:]).Staging).To(BeEquivalentTo(goGit.Added))
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
					repo, err := goGit.PlainOpen(repoClient.Root())
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
				repo, _ := goGit.PlainOpen(repoClient.Root())
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
					err = repoClient.CheckoutBranch(branchName)
				})

				It("does not generate an error", func() {
					Expect(err).To(BeNil())
				})

				It("the correct branch has been checked out", func() {
					repo, err := goGit.PlainOpen(repoClient.Root())
					Expect(err).To(BeNil())

					ref, err := repo.Head()
					Expect(err).To(BeNil())

					Expect(ref.Name()).To(BeEquivalentTo(git.LocalBranchPrefix + branchName))
				})
			})

			Describe("checking out the new branch", func() {

				BeforeEach(func() {
					err = repoClient.CheckoutBranch(branchName)
				})

				It("does not generate an error", func() {
					Expect(err).To(BeNil())
				})
			})

			Describe("new file committed to new branch", func() {

				const (
					fileName  = "only_on_new_branch.txt"
					commitMsg = "Another file: " + fileName
				)
				var filePath, hash string

				BeforeEach(func() {
					_ = repoClient.CheckoutBranch(branchName)

					filePath = path.Join(repoClient.Root(), fileName)
					err = ioutil.WriteFile(filePath, []byte(fmt.Sprint("I am just using up memory...\n")), 0644)
					Expect(err).To(BeNil())

					repoClient.Add(filePath)
					hash, _ = repoClient.Commit(commitMsg)
				})

				It("commit has been created", func() {
					repo, err := goGit.PlainOpen(repoClient.Root())
					Expect(err).To(BeNil())

					commit, err := repo.CommitObject(plumbing.NewHash(hash))
					Expect(err).To(BeNil())
					Expect(commit).To(Not(BeNil()))
				})

				It("file is not present on master branch", func() {
					repoClient.CheckoutBranch("master")
					_, err = os.Stat(path.Join(repoClient.Root(), fileName))
					Expect(os.IsNotExist(err)).To(BeTrue())
				})
			})
		})
	})

	Context("working with remote repositories", func() {

		var remoteClient *git.Repository

		BeforeEach(func() {
			repoRoot, _ := ioutil.TempDir("", "go_git_client_remotes_test")
			remoteRoot, _ := ioutil.TempDir("", "go_git_client_remotes_test")

			// Initialize git repository in a random temp dir
			repoClient, err = git.NewRepo(repoRoot)
			Expect(err).To(BeNil())
			Expect(repoClient.IsRepo()).To(BeFalse())

			// Initialize another local repository that will act as remote
			remoteClient, err = git.NewRepo(remoteRoot)
			Expect(err).To(BeNil())

			_, err = remoteClient.Init()
			Expect(err).To(BeNil())

			// Create another branch on the remote
			remoteClient.NewBranch("branch_1")
			remoteClient.CheckoutBranch("branch_1")
			Expect(ioutil.WriteFile(path.Join(remoteClient.Root(), "test_file"), []byte("Lorem ipsum..."), 0644)).To(BeNil())
			Expect(remoteClient.Add(path.Join(remoteClient.Root(), "test_file"))).To(BeNil())
			_, err := remoteClient.Commit("commit on branch_1")
			Expect(err).To(BeNil())
			remoteClient.CheckoutBranch("master")
		})

		Describe("cloning a remote repository (with two branches)", func() {

			BeforeEach(func() {
				err = repoClient.Clone(remoteClient.Root())
			})

			It("does not generate an error", func() {
				Expect(err).To(BeNil())
			})

			It("creates a repository in the given directory", func() {
				Expect(repoClient.IsRepo()).To(BeTrue())
			})

			It("the local repository points to the correct remote", func() {
				repo, _ := goGit.PlainOpen(repoClient.Root())
				remotes, err := repo.Remotes()

				Expect(err).To(BeNil())
				Expect(len(remotes)).To(BeEquivalentTo(1))
				Expect(remotes[0].Config().URLs[0]).To(BeEquivalentTo(remoteClient.Root()))
			})

			It("creates both local and remote-tracking refs for the remote refs", func() {
				branches, err := repoClient.ListBranches(true)

				Expect(err).To(BeNil())
				Expect(len(branches)).To(BeEquivalentTo(4))
				Expect(branches).To(ContainElement("refs/heads/master"))
				Expect(branches).To(ContainElement("refs/heads/branch_1"))
				Expect(branches).To(ContainElement("refs/remotes/origin/master"))
				Expect(branches).To(ContainElement("refs/remotes/origin/branch_1"))
			})

			It("is on the master branch", func() {
				repo, _ := goGit.PlainOpen(repoClient.Root())
				ref, err := repo.Head()

				Expect(err).To(BeNil())
				Expect(ref.Name()).To(BeEquivalentTo(git.LocalBranchPrefix + "master"))
			})

			It("can check out the other remote branch", func() {
				err = repoClient.CheckoutBranch("branch_1")
				Expect(err).To(BeNil())
			})
		})

		Describe("pushing a commit to a remote repository", func() {

			var hash string

			BeforeEach(func() {
				Expect(repoClient.Clone(remoteClient.Root())).To(BeNil())
				Expect(repoClient.NewBranch("to_be_pushed")).To(BeNil())
				Expect(repoClient.CheckoutBranch("to_be_pushed")).To(BeNil())
				Expect(ioutil.WriteFile(path.Join(repoClient.Root(), "test_push"), []byte("to be pushed..."), 0777)).To(BeNil())
				Expect(repoClient.Add(path.Join(repoClient.Root(), "test_push"))).To(BeNil())
				hash, err = repoClient.Commit("Commit on branch to be pushed")
				Expect(err).To(BeNil())
				Expect(hash).To(Not(BeNil()))

				err = repoClient.Push(remoteClient.Root())
			})

			It("does not generate an error", func() {
				Expect(err).To(BeNil())
			})

			It("correctly pushed the branch to the remote", func() {
				branches, err := repoClient.ListBranches(true)
				Expect(err).To(BeNil())
				Expect(branches).To(ContainElement("refs/remotes/origin/to_be_pushed"))
			})

			It("the remote branch contains the commit", func() {

				// Clone the remote into a new location and check that the file exists
				newRepoRoot, _ := ioutil.TempDir("", "go_git_client_remotes_test")
				newRepoClient, err := git.NewRepo(newRepoRoot)
				Expect(err).To(BeNil())
				Expect(newRepoClient.Clone(remoteClient.Root())).To(BeNil())
				Expect(newRepoClient.CheckoutBranch("to_be_pushed")).To(BeNil())

				fileInfo, err := os.Stat(path.Join(newRepoRoot, "test_push"))
				Expect(err).To(BeNil())
				Expect(fileInfo.Name()).To(BeEquivalentTo("test_push"))
			})
		})
	})
})
