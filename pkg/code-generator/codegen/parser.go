package codegen

import (
	"encoding/json"
	"io/ioutil"
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

	// magic comments
	shortNameDeclaration      = "@solo-kit:resource.short_name="
	pluralNameDeclaration     = "@solo-kit:resource.plural_name="
	resourceGroupsDeclaration = "@solo-kit:resource.resource_groups="
)

func ParseRequest(req *plugin_go.CodeGeneratorRequest) (*Project, error) {
	log.Printf("parsing request %v", req.FileToGenerate, req.GetParameter())
	params := req.GetParameter()
	if params == "" {
		return nil, errors.Errorf("must provide path to project.json file")
	}

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

	var groupName string
	for _, desc := range descriptors {
		if groupName == "" {
			groupName = desc.GetPackage()
		}
		if groupName != desc.GetPackage() {
			return nil, errors.Errorf("package conflict: %v must match %v", groupName, desc.GetPackage())
		}
	}

	project := &Project{
		ProjectConfig: projectConfig,
		GroupName:     groupName,
	}
	resources, resourceGroups, err := getResources(project, messages)
	if err != nil {
		return nil, err
	}

	project.Resources = resources
	project.ResourceGroups = resourceGroups

	return project, nil
}

func dataTypeForMessage(groupName, messageName string) string {
	return "." + groupName + "." + messageName + ".DataEntry"
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
		resource, groups, err := describeResource(project.GroupName, msg)
		if err != nil {
			return nil, nil, err
		}
		resource.Project = project
		for _, group := range groups {
			resourcesByGroup[group] = append(resourcesByGroup[group], resource)
		}
		resources = append(resources, resource)
	}

	var resourceGroups []*ResourceGroup

	for group, resources := range resourcesByGroup {
		log.Printf("group: %v", group)
		rg := &ResourceGroup{
			Name:             group,
			BelongsToProject: project,
			Resources:        resources,
		}
		for _, res := range resources {
			res.Project = project
			res.ResourceGroups = append(res.ResourceGroups, rg)
		}
		resourceGroups = append(resourceGroups, rg)
	}
	for _, res := range resources {
		// sort for stability
		sort.SliceStable(res.ResourceGroups, func(i, j int) bool {
			return res.ResourceGroups[i].Name < res.ResourceGroups[j].Name
		})
	}
	return resources, resourceGroups, nil
}

func describeResource(groupName string, msg *protokit.Descriptor) (*Resource, []string, error) {
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
	resourceGroups := strings.Split(joinedResourceGroups, ",")
	if resourceGroups[0] == "" {
		resourceGroups = nil
	}

	hasStatus := hasField(msg, "status", statusTypeName)
	dataTypeName := dataTypeForMessage(groupName, msg.GetName())
	log.Printf("%v", dataTypeName)
	hasData := hasField(msg, "data", dataTypeName)

	fields := collectFields(msg)

	return &Resource{
		Name:       name,
		ShortName:  shortName,
		PluralName: pluralName,
		HasStatus:  hasStatus,
		HasData:    hasData,
		Fields:     fields,
	}, resourceGroups, nil
}

func collectFields(msg *protokit.Descriptor) []*Field {
	var fields []*Field
	for _, f := range msg.GetField() {
		fields = append(fields, &Field{
			Name:     f.GetName(),
			TypeName: f.GetTypeName(),
		})
	}
	log.Printf("%v", fields)
	return fields
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
