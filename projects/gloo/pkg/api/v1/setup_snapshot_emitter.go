package v1

import (
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type SetupEmitter interface {
	Register() error
	Settings() SettingsClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *SetupSnapshot, <-chan error, error)
}

func NewSetupEmitter(settingsClient SettingsClient) SetupEmitter {
	return NewSetupEmitterWithEmit(settingsClient, make(chan struct{}))
}

func NewSetupEmitterWithEmit(settingsClient SettingsClient, emit <-chan struct{}) SetupEmitter {
	return &setupEmitter{
		settings:  settingsClient,
		forceEmit: emit,
	}
}

type setupEmitter struct {
	forceEmit <-chan struct{}
	settings  SettingsClient
}

func (c *setupEmitter) Register() error {
	if err := c.settings.Register(); err != nil {
		return err
	}
	return nil
}

func (c *setupEmitter) Settings() SettingsClient {
	return c.settings
}

func (c *setupEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *SetupSnapshot, <-chan error, error) {
	errs := make(chan error)
	var done sync.WaitGroup
	/* Create channel for Settings */
	type settingsListWithNamespace struct {
		list      SettingsList
		namespace string
	}
	settingsChan := make(chan settingsListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup watch for Settings */
		settingsNamespacesChan, settingsErrs, err := c.settings.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Settings watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, settingsErrs, namespace+"-settings")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case settingsList := <-settingsNamespacesChan:
					select {
					case <-opts.Ctx.Done():
						return
					case settingsChan <- settingsListWithNamespace{list: settingsList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *SetupSnapshot)
	go func() {
		currentSnapshot := SetupSnapshot{}
		sync := func(newSnapshot SetupSnapshot) {
			if currentSnapshot.Hash() == newSnapshot.Hash() {
				return
			}
			currentSnapshot = newSnapshot
			sentSnapshot := currentSnapshot.Clone()
			snapshots <- &sentSnapshot
		}
		for {
			select {
			case <-opts.Ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case settingsNamespacedList := <-settingsChan:
				namespace := settingsNamespacedList.namespace
				settingsList := settingsNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Settings.Clear(namespace)
				newSnapshot.Settings.Add(settingsList...)
				sync(newSnapshot)
			}
		}
	}()
	return snapshots, errs, nil
}
