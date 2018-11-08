package codegen

import (
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/errors"
)

const (
	// magic comments
	serviceDeclaration      = "@solo-kit:xds-service="
	noReferencesDeclaration = "@solo-kit:resource.no_references"
	nameFieldDeclaration    = "@solo-kit:resource.name"
	xdsEnabledDeclaration   = "@solo-kit:resource.xds-enabled"
)

type xdsService struct {
	name            string
	messageTypeName string
	groupName       string
}

type xdsMessage struct {
	name            string
	serviceTypeName string
	nameField       string
	noReferences    bool
	groupName       string
	protoPackage    string
}

func getXdsResources(project *Project, messages []ProtoMessageWrapper, services []*protokit.ServiceDescriptor) ([]*XDSResource, error) {
	var msgs []*xdsMessage
	var svcs []*xdsService

	for _, msg := range messages {
		msg, err := describeXdsResource(msg.Message)
		if err != nil {
			return nil, err
		}
		if msg == nil {
			// message is not a resource
			continue
		}
		if msg.groupName != project.GroupName {
			continue
		}
		msgs = append(msgs, msg)
	}
	for _, svc := range services {
		service, err := describeXdsService(svc)
		if err != nil {
			return nil, err
		}
		if service == nil {
			// message is not a resource
			continue
		}
		if service.groupName != project.GroupName {
			continue
		}
		svcs = append(svcs, service)
	}

	// match time!
	// for every service, match it with a config message.
	return processMessagesAndServices(project, msgs, svcs)
}

func processMessagesAndServices(project *Project, msgs []*xdsMessage, svcs []*xdsService) ([]*XDSResource, error) {
	var resources []*XDSResource
	for _, svc := range svcs {
		var message *xdsMessage
		for i, msg := range msgs {
			if msg.serviceTypeName == svc.name {
				message = msg
				msgs = append(msgs[:i], msgs[i+1:]...)
				break
			}
		}
		if message == nil {
			return nil, errors.Errorf("no message defined for service %v", svc.name)
		}

		resources = append(resources, &XDSResource{
			MessageType:  message.name,
			Name:         svc.name,
			NameField:    message.nameField,
			NoReferences: message.noReferences,
			GroupName:    message.groupName,
			Package:      message.protoPackage,
			Project:      project,
		})
	}

	if len(msgs) != 0 {
		var msgnames []string
		for _, msg := range msgs {
			msgnames = append(msgnames, msg.name)
		}
		return nil, errors.Errorf("orphaned messages: %s", strings.Join(msgnames, ","))
	}

	return resources, nil
}

func describeXdsResource(msg *protokit.Descriptor) (*xdsMessage, error) {
	commentsString := msg.GetComments().Leading
	comments := strings.Split(commentsString, "\n")
	service, ok := getCommentValue(comments, serviceDeclaration)
	if !ok {
		// no service definition - the object doesn't belong to us
		return nil, nil
	}

	// not a solo kit resource, or you messed up!
	name := ""

	for _, field := range msg.Fields {
		if strings.Contains(field.GetComments().Leading, nameFieldDeclaration) {
			if name != "" {
				return nil, errors.Errorf("can only have one name annotation")
			}
			if field.GetType() != descriptor.FieldDescriptorProto_TYPE_STRING {
				return nil, errors.Errorf("name type should be a string")
			}
			name = field.GetName()
		}
	}

	if name == "" && hasPrimitiveField(msg, "name", descriptor.FieldDescriptorProto_TYPE_STRING) {
		name = "name"
	}

	if name == "" {
		return nil, errors.Errorf("no name field found. please use " + nameFieldDeclaration + " to designate a field as a name")
	}

	var noRefs bool
	if strings.Contains(commentsString, noReferencesDeclaration) {
		noRefs = true
	}

	return &xdsMessage{
		name:            msg.GetName(),
		serviceTypeName: service,
		nameField:       name,
		noReferences:    noRefs,
		groupName:       msg.GetPackage(),
		protoPackage:    msg.GetPackage(),
	}, nil
}

func describeXdsService(service *protokit.ServiceDescriptor) (*xdsService, error) {
	comments := service.Comments.Leading
	if !strings.Contains(comments, xdsEnabledDeclaration) {
		return nil, nil
	}

	msgConfig := ""

	streamPrefix := "Stream"
	fetchPrefix := "Fetch"

	for _, method := range service.Methods {
		if strings.HasPrefix(method.GetName(), streamPrefix) {
			if msgConfig != "" {
				return nil, errors.Errorf("Only one stream method is allowed")
			}
			msgConfig = strings.TrimPrefix(method.GetName(), streamPrefix)
		}
	}

	if msgConfig == "" {
		return nil, errors.Errorf("No stream method found")
	}

	if service.GetNamedMethod(fetchPrefix+msgConfig) == nil {
		return nil, errors.Errorf("No fetch method found")
	}

	return &xdsService{
		name:            service.GetName(),
		messageTypeName: msgConfig,
		groupName:       service.GetPackage(),
	}, nil
}
