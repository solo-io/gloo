package codegen

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/log"
)

const (
	// solo-kit types
	// required fields
	metadataTypeName = ".core.solo.io.Metadata"
	statusTypeName   = ".core.solo.io.Status"
	dataTypeName     = ".core.solo.io.Data"

	// magic comments
	shortNameDeclaration      = "@solo-kit:resource.short_name="
	pluralNameDeclaration     = "@solo-kit:resource.plural_name="
	resourceGroupsDeclaration = "@solo-kit:resource.resource_groups="
)

func parseRequest(req *plugin_go.CodeGeneratorRequest) (*Project, error) {
	if req.Parameter == nil {
		return nil, errors.Errorf("must provide path to project.json file")
	}

	params := *req.Parameter
	log.DefaultOut = os.Stderr
	log.Printf("got cli param from protoc invoke: %v", params)
	projectConfig, err := loadProjectConfig(params)
	if err != nil {
		return nil, err
	}

	descriptors := protokit.ParseCodeGenRequest(req)
	var messages []*protokit.Descriptor
	for _, file := range descriptors {
		messages = append(messages, file.GetMessages()...)
	}

	project := &Project{
		ProjectConfig: projectConfig,
	}
	resources, resourceGroups, err := getResources(project, messages)
	if err != nil {
		return nil, err
	}

	project.Resources = resources
	project.ResourceGroups = resourceGroups

	return project, nil
}

func loadProjectConfig(path string) (ProjectConfig, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ProjectConfig{}, err
	}
	var pc ProjectConfig
	err = json.Unmarshal(b, &pc)
	return pc, err
}

func getResources(project *Project, messages []*protokit.Descriptor) ([]*Resource, []*ResourceGroup, error) {
	resourcesByGroup := make(map[string][]*Resource)
	var resources []*Resource
	for _, msg := range messages {
		resource, groups, err := describeResource(msg)
		if err != nil {
			return nil, nil, err
		}
		for _, group := range groups {
			resourcesByGroup[group] = append(resourcesByGroup[group], resource)
		}
		resources = append(resources, resource)
	}

	var resourceGroups []*ResourceGroup

	for group, resources := range resourcesByGroup {
		rg := &ResourceGroup{
			Name:             group,
			BelongsToProject: project,
			Resources:        resources,
		}
		for _, res := range resources {
			res.BelongsToProject = project
			res.BelongsToResourceGroups = append(res.BelongsToResourceGroups, rg)
		}
		resourceGroups = append(resourceGroups, rg)
	}
	for _, res := range resources {
		// sort for stability
		sort.SliceStable(res.BelongsToResourceGroups, func(i, j int) bool {
			return res.BelongsToResourceGroups[i].Name < res.BelongsToResourceGroups[j].Name
		})
	}
	return resources, resourceGroups, nil
}

func describeResource(msg *protokit.Descriptor) (*Resource, []string, error) {
	// not a solo kit resource, or you messed up!
	if !hasField(msg, "metadata", metadataTypeName) {
		return nil, nil, nil
	}

	comments := strings.Split(msg.GetComments().Leading, "\n")

	name := msg.GetName()
	// required flags
	shortName, ok := getCommentValue(comments, shortNameDeclaration)
	if !ok {
		return nil, nil, errors.Errorf("must provide %s", shortNameDeclaration)
	}
	pluralName, ok := getCommentValue(comments, pluralNameDeclaration)
	if !ok {
		return nil, nil, errors.Errorf("must provide %s", pluralNameDeclaration)
	}

	// optional flags
	joinedResourceGroups, _ := getCommentValue(comments, resourceGroupsDeclaration)

	hasStatus := hasField(msg, "status", statusTypeName)
	hasData := hasField(msg, "data", dataTypeName)

	return &Resource{
		Name:       name,
		ShortName:  shortName,
		PluralName: pluralName,
		HasStatus:  hasStatus,
		HasData:    hasData,
	}, strings.Split(joinedResourceGroups, ","), nil
}

func hasField(msg *protokit.Descriptor, fieldName, fieldType string) bool {
	for _, field := range msg.Fields {
		if field.GetName() == fieldName && field.GetTypeName() == fieldType {
			return true
		}
	}
	return false
}

func getCommentValue(comments []string, key string) (string, bool) {
	for _, c := range comments {
		if strings.HasPrefix(c, key) {
			return strings.TrimPrefix(c, key), true
		}
	}
	return "", false
}
