package git

import (
	"context"
	"os"
	"strings"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/constants"

	"github.com/aokoli/goutils"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/api/v1"
)

type RemoteSyncer struct {
	CsClient *v1.ChangeSetClient
}

// TODO: for now everything here is synchronous, consider using goroutines
// Checks whether any of the changesets in the current snapshot are pending a commit status. If so, the changes are
// pushed to the remote repository and the changeset is updated (either with the new commit hash or with an error message)
func (s *RemoteSyncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {

	// Iterate over all available changesets
	for _, cs := range snap.Changesets[defaults.GlooSystem] {

		var err error

		// If the pending action is COMMIT, push the changeset to the git remote
		if cs.PendingAction == v1.Action_COMMIT {
			err = s.pushChanges(ctx, cs)
		}

		// If the pending action is CHECK_OUT, populate the changeset with data from the remote
		if cs.PendingAction == v1.Action_CHECK_OUT {
			err = s.checkout(ctx, cs)
		}

		// If an error occurred, set the error message on the changeset and set its pending_action to NONE
		if err != nil {
			contextutils.LoggerFrom(ctx).Error(err)
			err = s.markChangesetAsFailed(ctx, cs, err)
			if err != nil {
				// Panic if we can't mark the changeset as failed to avoid getting stuck in a loop
				contextutils.LoggerFrom(ctx).Panicf("Could not mark changeset [%v] as failed. Root cause: [%v]",
					cs.Metadata.Name, err)
			}
		}
	}
	return nil
}

// Pushes the state of the repository to the remote
func (s *RemoteSyncer) pushChanges(ctx context.Context, cs *v1.ChangeSet) error {
	contextutils.LoggerFrom(ctx).Infof("Preparing to commit changeset [%v] ", cs.Metadata.Name)

	err := validateForPush(cs)
	if err != nil {
		return err
	}

	// Create a new repository and delete it once the method returns
	repo, err := cloneRemote(os.Getenv(constants.RemoteUriEnvVariableName))
	defer repo.Delete()

	// TODO
	if err != nil {
		return err
	}

	// Check whether the given root commit exists
	if !repo.CommitExists(cs.RootCommit.GetValue()) {
		return errors.Errorf("Could not find a commit with hash [%v]", cs.RootCommit.GetValue())
	}

	// Get all the existing branches
	branches, err := repo.ListBranches(true)
	if err != nil {
		return err
	}

	// Check if the changeset branch exists
	if exists, name := branchExists(branches, cs.Branch.GetValue()); exists {
		contextutils.LoggerFrom(ctx).Infof("Found branch [%v] matching changeset branch [%v].", name, cs.Branch.GetValue())

		// If the branch exists, check if the HEAD of the branch matches the changeset root commit. If the commit that
		// HEAD points to is different than the root commit, it means the working copy is stale. Return an error.

		// Switch to the branch specified in the changeset
		contextutils.LoggerFrom(ctx).Infof("Checking out branch [%v].", cs.Branch.GetValue())
		err = repo.CheckoutBranch(cs.Branch.GetValue())
		if err != nil {
			return err
		}

		// Get the HEAD and compare it to the root commit hash specified in the changeset
		headRef, err := repo.Head()
		if err != nil {
			return err
		}
		if headRef != cs.RootCommit.GetValue() {
			return errors.Errorf("Changeset has root commit [%v] but the target branch [%v] is at commit [%v]. "+
				"Your working copy is stale. To solve this, either specify a different branch name or check out the branch again.",
				cs.RootCommit.GetValue(), cs.Branch.GetValue(), headRef)

		}
		contextutils.LoggerFrom(ctx).Infof("Current HEAD: [%v] ", headRef)

	} else {

		// If the branch does not exists, create a new branch starting at the given hash

		contextutils.LoggerFrom(ctx).Infof("No branch matching changeset branch [%v] found. "+
			"Creating new branch starting at [%v].", name, cs.Branch.GetValue(), cs.RootCommit.GetValue())
		err = repo.NewBranchFromHash(cs.Branch.GetValue(), cs.RootCommit.GetValue())
		if err != nil {
			return err
		}
	}

	// At this point we are on the desired branch

	// Write the content of the changeset to the local repository
	err = repo.Import(cs)
	if err != nil {
		return err
	}

	// Stage all files in the repository
	err = repo.Add(repo.root)
	if err != nil {
		return err
	}

	hash, err := repo.Commit(cs.Description.Value)
	if err != nil {
		return err
	}
	contextutils.LoggerFrom(ctx).Infof("Created new commit [%v] on branch [%v].", hash, cs.Branch.GetValue())

	contextutils.LoggerFrom(ctx).Infof("Pushing branch to remote [%v]", os.Getenv(constants.RemoteUriEnvVariableName))
	err = repo.Push(os.Getenv(constants.RemoteUriEnvVariableName))
	if err != nil {
		return err
	}

	cs.PendingAction = v1.Action_NONE
	cs.EditCount = types.UInt32Value{Value: 0}
	cs.RootCommit.Value = hash
	cs.RootDescription.Value = cs.Description.Value
	_, err = (*s.CsClient).Write(cs, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})

	return err
}

// Copies the state of the repository into the changeset
func (s *RemoteSyncer) checkout(ctx context.Context, cs *v1.ChangeSet) error {

	// Either branch name OR commit hash will be non-empty
	branchName, commitHash, err := validateForCheckout(cs)
	if err != nil {
		return err
	}

	// Create a new repository and delete it once the method returns
	repo, err := cloneRemote(os.Getenv(constants.RemoteUriEnvVariableName))
	defer repo.Delete()

	if !goutils.IsEmpty(branchName) {
		contextutils.LoggerFrom(ctx).Infof("Checking out branch [%v]", branchName)
		err = repo.CheckoutBranch(branchName)
		if err != nil {
			if err == plumbing.ErrReferenceNotFound {
				err = errors.Errorf("Could not find any branch with name [%v]", branchName)
			}
			return err
		}
	} else { // checkout by commit hash
		contextutils.LoggerFrom(ctx).Infof("Checking out commit [%v]", commitHash)
		err = repo.CheckoutCommit(commitHash)
		if err != nil {
			if err == plumbing.ErrObjectNotFound {
				err = errors.Errorf("Could not find any commit with hash [%v]", commitHash)
			}
			return err
		}
	}

	lastCommitHash, lastCommitMessage, err := repo.LastCommit()
	if err != nil {
		return err
	}

	// Update changeset
	if data, err := repo.ToChangeSetData(); err != nil {
		return err
	} else {
		cs.Data = *data
	}
	cs.PendingAction = v1.Action_NONE
	cs.EditCount = types.UInt32Value{Value: 0}
	cs.RootCommit = types.StringValue{Value: lastCommitHash}
	cs.RootDescription = types.StringValue{Value: lastCommitMessage}

	_, err = (*s.CsClient).Write(cs, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})

	return err
}

func (s *RemoteSyncer) markChangesetAsFailed(ctx context.Context, cs *v1.ChangeSet, e error) error {

	cs.ErrorMsg = &types.StringValue{Value: e.Error()}
	cs.PendingAction = v1.Action_NONE

	_, err := (*s.CsClient).Write(cs, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})

	return err
}

// Verifies that only one between branch and hash is present
func validateForCheckout(cs *v1.ChangeSet) (branch, hash string, err error) {
	branchName, commitHash := cs.Branch.GetValue(), cs.RootCommit.GetValue()

	if !goutils.IsEmpty(branchName) && !goutils.IsEmpty(commitHash) {
		return "", "", errors.Errorf("Both branch name and commit hash are present. Use only one of them. Changeset [%v]", cs.Metadata.Name)
	}

	if goutils.IsEmpty(branchName) && goutils.IsEmpty(commitHash) {
		return "", "", errors.Errorf("Specify either branch name or commit hash. Changeset [%v]", cs.Metadata.Name)
	}

	return branchName, commitHash, nil
}

func validateForPush(cs *v1.ChangeSet) error {

	if cs.EditCount.Value == 0 {
		return errors.Errorf("Changeset %v does not contain any edits", cs.Metadata.Name)
	}

	if cs.Branch.GetValue() == "master" {
		return errors.Errorf("Changeset %v contains edits to the master branch", cs.Metadata.Name)
	}

	if cs.Description.GetValue() == "" {
		return errors.Errorf("A commit message must be provided", cs.Metadata.Name)
	}

	if cs.Branch.GetValue() == "" {
		return errors.Errorf("No branch name specified", cs.Metadata.Name)
	}

	if cs.RootCommit.GetValue() == "" {
		return errors.Errorf("Root commit hash missing", cs.Metadata.Name)
	}

	return nil
}

func branchExists(branches []string, name string) (bool, string) {
	for _, branchName := range branches {
		if strings.Contains(branchName, name) {
			return true, branchName
		}
	}
	return false, ""
}

func cloneRemote(url string) (*Repository, error) {
	// Create a new repository and delete it once the method returns
	repo, err := NewTempRepo()
	if err != nil {
		return nil, err
	}

	// Set token for authentication with remote
	repo = repo.WithTokenAuth(os.Getenv(constants.AuthTokenEnvVariableName))

	err = repo.Clone(url)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
