package main

import (
	"os"
	"strings"

	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-kit/pkg/code-generator/cmd"
	"github.com/solo-io/solo-kit/pkg/code-generator/docgen/options"
	"github.com/solo-io/solo-kit/pkg/code-generator/schemagen"
	"github.com/solo-io/solo-kit/pkg/code-generator/sk_anyvendor"
)

//go:generate go run generate.go

func main() {
	log.Printf("starting generate for gloo")

	// Explicitly specify the directories to be built (i.e. do not build gateway2 since
	// it causes compilation errors in solo-kit, and also because gateway2 protos are not
	// needed for gloo edge classic). See `projects/gateway2/api/README.md` for more info.
	protoImports := sk_anyvendor.CreateDefaultMatchOptions(
		[]string{
			"projects/gloo/**/*.proto",
			"projects/gateway/**/*.proto",
			"projects/ingress/**/*.proto",
			sk_anyvendor.SoloKitMatchPattern,
		},
	)
	protoImports.External["github.com/solo-io/solo-apis"] = []string{
		"api/rate-limiter/**/*.proto", // Import rate limit API
		"api/gloo-fed/fed/**/*.proto", // Import gloo fed gloo instance API
	}
	// Import gloo instance API dependencies
	protoImports.External["github.com/solo-io/skv2"] = []string{
		"api/**/**/*.proto",
	}

	generateOptions := cmd.GenerateOptions{
		SkipGenMocks: true,
		CustomCompileProtos: []string{
			"github.com/solo-io/gloo/projects/gloo/api/grpc",
		},
		SkipGeneratedTests: true,
		// helps to cut down on time spent searching for imports, not strictly necessary
		SkipDirs: []string{
			"docs",
			"test",
			"projects/gloo/api/grpc",
		},
		RelativeRoot:  ".",
		CompileProtos: true,
		GenDocs: &cmd.DocsOptions{
			Output: options.Hugo,
			HugoOptions: &options.HugoOptions{
				DataDir: "/docs/data",
				ApiDir:  "reference/api",
			},
			RenderOptions: &options.RenderOptions{
				SkipLinksForPathPrefixes: []string{
					"github.com/solo-io/gloo/projects/gloo/api/external",
				},
			},
		},
		ExternalImports: protoImports,
		ValidationSchemaOptions: &schemagen.ValidationSchemaOptions{
			CrdDirectory:                 "install/helm/gloo/crds",
			JsonSchemaTool:               "protoc",
			RemoveDescriptionsFromSchema: true,
			EnumAsIntOrString:            true,
			MessagesWithEmptySchema: []string{
				// These messages are recursive, and will cause codegen to enter an infinite loop
				// If you're adding a new message to this list, consider addressing the issue
				// solo-io/protoc-gen-openapi or switching to another implementation
				"core.solo.io.Status",
				"ratelimit.api.solo.io.Descriptor",
				"als.options.gloo.solo.io.AndFilter",
				"als.options.gloo.solo.io.OrFilter",
				"opentelemetry.proto.common.v1.AnyValue",

				// These messages are part of our internal API, and therefore aren't required
				// Also they are quite large and can cause the Proxy CRD to become too large,
				// resulting in: https://github.com/solo-io/gloo/issues/4789
				"gloo.solo.io.HttpListener",
				"gloo.solo.io.TcpListener",
				"gloo.solo.io.HybridListener",
				"gloo.solo.io.AggregateListener",
			},
			DisableKubeMarkers: true,
		},
	}
	if err := cmd.Generate(generateOptions); err != nil {
		// check the error for "generateMessageSchema", if we see it many times
		// it's likely a cycle in the the proto files
		count := strings.Count(err.Error(), "generateMessageSchema")
		// this is arbitrary, raise if you're seeing false positives
		if count > 10 {
			log.Fatalf("generate failed!: %v. The error message indicates a cycle the protos! See MessagesWithEmptySchema in generate.go.", err)
		}

		log.Fatalf("generate failed!: %v", err)
	}

	err := removeExternalApiDocs()
	if err != nil {
		log.Fatalf("failed to remove external api docs: %v", err)
	}

	err = fixExtauthCrossReferences()
	if err != nil {
		log.Fatalf("failed to fix extauth cross-references: %v", err)
	}

	log.Printf("finished generating code for gloo")
}

func removeExternalApiDocs() error {
	const externalApiDocsPath = "docs/content/reference/api/github.com/solo-io/gloo/projects/gloo/api/external"
	return os.RemoveAll(externalApiDocsPath)
}

// fixExtauthCrossReferences fixes incorrect cross-references in the generated documentation
// where extauth-internal.proto.sk is referenced instead of extauth.proto.sk
func fixExtauthCrossReferences() error {
	const extauthDocPath = "docs/content/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth.proto.sk.md"
	const extauthInternalDocPath = "docs/content/reference/api/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth-internal.proto.sk.md"

	// Change extauth-internal.proto.sk/#config to extauth.proto.sk/#config
	content, err := os.ReadFile(extauthDocPath)
	if err != nil {
		return err
	}
	content = []byte(strings.ReplaceAll(string(content), "extauth-internal.proto.sk/#config", "extauth.proto.sk/#config"))
	err = os.WriteFile(extauthDocPath, content, 0644)
	if err != nil {
		return err
	}

	// Change extauth-internal.proto.sk/#config to extauth-internal.proto.sk/#config-1
	content, err = os.ReadFile(extauthInternalDocPath)
	if err != nil {
		return err
	}
	content = []byte(strings.ReplaceAll(string(content), "extauth-internal.proto.sk/#config", "extauth-internal.proto.sk/#config-1"))
	err = os.WriteFile(extauthInternalDocPath, content, 0644)
	if err != nil {
		return err
	}

	log.Printf("fixed extauth cross-references in %s", extauthDocPath)
	return nil
}
