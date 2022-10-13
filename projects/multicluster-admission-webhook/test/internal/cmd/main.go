package main

import (
	"context"
	"log"

	rbacconfig "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/config"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/test/internal/pkg/webhook"
)

func main() {
	rbacConfig := rbacconfig.NewConfig()
	err := webhook.Start(context.Background(), rbacConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("exiting...")
}
