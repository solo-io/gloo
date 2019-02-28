package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/solo-io/go-utils/contextutils"

	"cloud.google.com/go/storage"
	"github.com/ghodss/yaml"
)

const (
	distBucket         = "gloo-ee-distribution"
	distBucketReleases = "releases"
	distBucketTags     = "tags"
	indexFile          = "index.yaml"

	objectDNE = "storage: object doesn't exist"
)

type distributionBucketClient struct {
	ctx    context.Context
	client *storage.Client
}

type distributionVersion struct {
	Version string `json:"version"`
	Id      string `json:"id"`
}

type distributionVersions struct {
	Versions []distributionVersion `json:"versions"`
}

func newDistributionBucketClient(ctx context.Context) (*distributionBucketClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	distBucketCli := &distributionBucketClient{
		client: client,
		ctx:    ctx,
	}
	return distBucketCli, nil
}

func syncDataToBucket(db *distributionBucketClient) error {
	if err := db.saveZipToBucket(); err != nil {
		return err
	}
	if err := db.uploadDistributionFolder(); err != nil {
		return err
	}
	if err := db.syncIndexFile(); err != nil {
		return err
	}
	return nil
}

func (db *distributionBucketClient) saveZipToBucket() error {
	logger := contextutils.LoggerFrom(db.ctx)
	zipFileName := fmt.Sprintf("%s%s%s", glooe, version, zipExt)

	bkt := db.client.Bucket(distBucket)
	obj := bkt.Object(path.Join(distBucketReleases, id.String(), zipFileName))
	logger.Infof("writing object: %v", obj.ObjectName())

	wr := obj.NewWriter(db.ctx)
	if err := writeDistributionArchive(db.ctx, wr, zipFileName); err != nil {
		return err
	}
	// Explicit non-deferred close to ensure object exists
	if err := wr.Close(); err != nil {
		return err
	}
	// Make zip file accessible to everyone with the link/public internet
	if err := obj.ACL().Set(db.ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	return nil
}

func (db *distributionBucketClient) uploadDistributionFolder() error {
	logger := contextutils.LoggerFrom(db.ctx)
	files, err := ioutil.ReadDir(outputDistributionDir)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	bkt := db.client.Bucket(distBucket)
	errch := make(chan error)

	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			go func(file os.FileInfo, wg *sync.WaitGroup) {
				defer wg.Done()
				logger.Infof("syncing file: (%s)", file.Name())
				r, err := os.Open(filepath.Join(outputDistributionDir, file.Name()))
				if err != nil {
					errch <- err
					return
				}
				obj := bkt.Object(path.Join(distBucketTags, version, file.Name()))
				wr := obj.NewWriter(db.ctx)
				defer wr.Close()
				_, err = io.Copy(wr, r)
				if err != nil {
					errch <- err
					return
				}
				logger.Infof("finished syncing file: (%s)", file.Name())
			}(file, wg)
		}
	}

	go func(wg *sync.WaitGroup, errch chan error) {
		wg.Wait()
		close(errch)
	}(wg, errch)

	select {
	case err := <-errch:
		if err != nil {
			return err
		}
	case <-time.After(10 * time.Minute):
		return fmt.Errorf("go routine timed out after 10 minutes")
	}
	return nil
}

func (db *distributionBucketClient) syncIndexFile() error {
	logger := contextutils.LoggerFrom(db.ctx)
	logger.Info("reading index.yaml file with version map")
	bkt := db.client.Bucket(distBucket)
	obj := bkt.Object(indexFile)

	_, err := obj.Attrs(db.ctx)
	if err != nil {
		if err.Error() == objectDNE {
			if err := db.createIndexFile(); err != nil {
				return err
			}
		}
	} else {
		if err := db.updateIndexFile(obj); err != nil {
			return err
		}
	}
	return nil
}

func (db *distributionBucketClient) createIndexFile() error {
	logger := contextutils.LoggerFrom(db.ctx)
	logger.Info("creating index.yaml file")
	bkt := db.client.Bucket(distBucket)
	obj := bkt.Object(indexFile)
	versions := distributionVersions{
		Versions: []distributionVersion{},
	}
	w := obj.NewWriter(db.ctx)
	if err := db.writeIndexFile(versions, w); err != nil {
		return err
	}
	logger.Info("finished creating index.yaml file")
	return nil
}

func (db *distributionBucketClient) writeIndexFile(versions distributionVersions, w *storage.Writer) error {
	newVersion := distributionVersion{
		Id:      id.String(),
		Version: version,
	}
	versions.Versions = append(versions.Versions, newVersion)

	outByt, err := yaml.Marshal(versions)
	if err != nil {
		return err
	}

	defer w.Close()
	if _, err = io.WriteString(w, string(outByt)); err != nil {
		return err
	}
	return nil
}

func (db *distributionBucketClient) updateIndexFile(obj *storage.ObjectHandle) error {
	logger := contextutils.LoggerFrom(db.ctx)
	r, err := obj.NewReader(db.ctx)
	if err != nil {
		return err
	}
	byt, err := ioutil.ReadAll(r)

	var versions distributionVersions
	err = yaml.Unmarshal(byt, &versions)
	if err != nil {
		return err
	}

	w := obj.NewWriter(db.ctx)
	if err := db.writeIndexFile(versions, w); err != nil {
		return err
	}
	logger.Infof("finished updating index.yaml file with new version: (%s)", version)
	return nil
}
