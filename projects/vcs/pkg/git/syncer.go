package git

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/projects/vcs/pkg"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/file"
	"go.uber.org/zap"
)

const (
	// TODO: temporary, get these from env variables
	remoteUrl      = "https://github.com/solo-io/gitbot-test"
	token          = "TOKEN HERE"
	deploymentType = "kube"
)

var logger *zap.SugaredLogger

type RemoteSyncer struct {
	CsClient *v1.ChangeSetClient
}

// Checks whether any of the changesets in the current snapshot are pending a commit status. If so, the changes are
// pushed to the remote repository and the changeset is updated (either with the new commit hash or with an error message)
func (s *RemoteSyncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {

	logger = contextutils.LoggerFrom(ctx)

	// Iterate over all available changesets
	for _, cs := range snap.Changesets[defaults.GlooSystem] {

		// If this commit is pending, push it to the git remote
		if cs.PendingAction == v1.Action_COMMIT {

			s.pushChanges(cs, ctx)

		}
	}

	return nil
}

func (s *RemoteSyncer) pushChanges(cs *v1.ChangeSet, ctx context.Context) error {
	logger.Infof("Preparing to commit changeset [%v] ", cs.Metadata.Name)

	// create temp dir
	tempDir, err := ioutil.TempDir("", pkg.AppName)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// TODO: WATCH OUT! pending flag has to be set to false at some point to not get stuck in an infinite loop
	if cs.EditCount.Value == 0 {
		logger.Error("Changeset edit count is zero. Nothing to commit.")
		return errors.Errorf("Changeset %v does not contain any edits", cs.Metadata.Name)
	}

	// TODO: improve error handling
	if cs.Branch.GetValue() == "master" {
		logger.Error("Direct changes to master are not allowed.")
		return errors.Errorf("Changeset %v contains edits to the master branch", cs.Metadata.Name)
	}

	repo := NewRepo(tempDir).WithTokenAuth(token)
	err = repo.Clone(remoteUrl)
	if err != nil {
		return err
	}

	// Switch to master branch if somehow the remote HEAD is pointing to another branch
	err = repo.CheckoutBranch("master", true)
	if err != nil {
		return err
	}

	branches, err := repo.ListBranches(true)
	if err != nil {
		return err
	}

	if exists, name := branchExists(branches, cs.Branch.GetValue()); exists {
		logger.Infof("Found branch [%v] matching changeset branch [%v].", name, cs.Branch.GetValue())
	} else {
		logger.Infof("No branch matching changeset branch [%v] found. Creating new branch.", name, cs.Branch.GetValue())
		err = repo.NewBranch(cs.Branch.GetValue())
		if err != nil {
			return err
		}
	}

	// Switch to the branch specified in the changeset
	logger.Infof("Checking out branch [%v].", cs.Branch.GetValue())
	err = repo.CheckoutBranch(cs.Branch.GetValue(), false)
	if err != nil {
		return err
	}

	// Get the HEAD and compare it to the root commit hash specified in the changeset
	headRef, err := repo.Head()
	if err != nil {
		return err
	}
	if headRef != cs.RootCommit.GetValue() {
		logger.Errorf("Changeset root hash [%v] does not match the HEAD [%v] of branch [%v]", cs.RootCommit.GetValue(), headRef, cs.Branch.GetValue())
	}
	logger.Infof("Current HEAD: [%v] ", headRef)

	// Write the content of the changeset to the local repository
	err = writeResourcesToFile(cs, tempDir, ctx)
	if err != nil {
		return err
	}

	err = repo.Add(tempDir)
	if err != nil {
		return err
	}

	hash, err := repo.Commit(cs.Description.Value)
	if err != nil {
		return err
	}

	err = repo.Push(remoteUrl)
	if err != nil {
		return err
	}

	// TODO maybe get message from commit object itself
	if err == nil {
		cs.RootCommit.Value = hash
		cs.RootDescription.Value = cs.Description.Value
		(*s.CsClient).Write(cs, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
	}

	// TODO: set error element if error occurs
	return nil
}

func (s *RemoteSyncer) markChangesetAsFailed(ctx context.Context, cs *v1.ChangeSet, msg string) error {

	cs.ErrorMsg = &types.StringValue{Value: msg}
	_, err := (*s.CsClient).Write(cs, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})

	return err
}

// Write the changeset data to the file system
func writeResourcesToFile(cs *v1.ChangeSet, dir string, ctx context.Context) error {

	// Create file client
	dc, err := file.NewDualClient(deploymentType, dir)
	if err != nil {
		return err
	}

	for _, vs := range cs.Data.VirtualServices {
		_, err = dc.File.VirtualServiceClient.Write(vs, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, gateway := range cs.Data.Gateways {
		_, err = dc.File.GatewayClient.Write(gateway, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, proxy := range cs.Data.Proxies {
		_, err = dc.File.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, resolverMap := range cs.Data.ResolverMaps {
		_, err = dc.File.ResolverMapClient.Write(resolverMap, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, schema := range cs.Data.Schemas {
		_, err = dc.File.SchemaClient.Write(schema, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	for _, setting := range cs.Data.Settings {
		_, err = dc.File.SettingsClient.Write(setting, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		if err != nil {
			return err
		}
	}

	// TODO: do we do upstreams?
	for _, upstream := range cs.Data.Upstreams {
		_, err = dc.File.UpstreamClient.Write(upstream, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		if err != nil {
			return err
		}
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
