package main

import (
	"net/http"
	"os"

	"github.com/ilackarms/go-github-webhook-server/github"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/example/source_events_from_github/base"
	"github.com/solo-io/gloo/pkg/log"
)

func main() {
	log.Fatalf("err:", run())
}

func run() error {
	endpoint := "minio-service.default.svc.cluster.local:9000"
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false // Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return errors.Wrap(err, "setting up minio client")
	}

	// Make a new bucket called images
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	location := os.Getenv("MINIO_REGION")
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil || !exists {
		err = minioClient.MakeBucket(bucketName, location)
		if err != nil {
			return errors.Wrap(err, "creating bucket failed")
		}
	}
	log.Printf("Successfully created bucket %s\n", bucketName)

	opts := base.Opts{
		ClientID:  os.Getenv("HOSTNAME"),
		ClusterID: "test-cluster",
		NatsURL:   "nats://nats-streaming.default.svc.cluster.local:4222",
		Subject:   "github-webhooks",
		Handler:   handleWatch(minioClient, bucketName),
	}
	base.Run(opts)
	log.Printf("terminated")
	return nil
}

func handleWatch(minioClient *minio.Client, bucketName string) func(watch *github.WatchEvent) error {
	return func(watch *github.WatchEvent) error {
		imgUrl := watch.Sender.AvatarURL
		res, err := http.Get(imgUrl)
		if err != nil {
			return errors.Wrap(err, "downloading image from url "+imgUrl)
			// upload the
		}
		// Upload the image
		objectName := watch.Sender.Login + ".png"
		contentType := "image/png"

		// Upload the zip file with FPutObject
		n, err := minioClient.PutObject(bucketName,
			objectName,
			res.Body,
			res.ContentLength,
			minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			return errors.Wrap(err, "uploading object")
		}
		log.Printf("uploaded %v/%v: %v bytes", bucketName, objectName, n)

		return nil
	}
}
