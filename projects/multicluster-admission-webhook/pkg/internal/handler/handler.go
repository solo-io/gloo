package handler

import (
	"context"
	"fmt"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/validation"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	NotAllowed = func(username string) string {
		return fmt.Sprintf("User %s does not have the permissions necessary to perform this action.", username)
	}
	Allowed = func(username string) string {
		return fmt.Sprintf("User %s has the necessary permissions to perform this action.", username)
	}
	InternalError = func(err error) string {
		return eris.Wrap(err, "Error computing cluster permissions").Error()
	}
)

type AdmissionWebhookHandler admission.Handler

func NewAdmissionWebhookHandler(
	mcAdmissionValidator validation.MultiClusterAdmissionValidator,
) AdmissionWebhookHandler {
	return &admissionWebhookHandler{
		mcAdmissionValidator: mcAdmissionValidator,
	}
}

type admissionWebhookHandler struct {
	mcAdmissionValidator validation.MultiClusterAdmissionValidator
}

func (r *admissionWebhookHandler) Handle(ctx context.Context, request admission.Request) admission.Response {
	logger := contextutils.LoggerFrom(contextutils.WithLoggerValues(ctx, zap.String("UID", string(request.UID))))
	logger.Debug("handling event")
	// Get all MultiClusterRoleBindings which apply to the user making the request
	matchingRoleBindings, err := r.mcAdmissionValidator.GetMatchingMultiClusterRoleBindings(ctx, request.UserInfo)
	if err != nil {
		return admission.Denied(InternalError(err))
	}

	// Useful debugging info if no matching role bindings are found
	if len(matchingRoleBindings) == 0 {
		logger.Debug("found no matching bindings, rejecting")
	}

	for _, matchingBinding := range matchingRoleBindings {
		childCtx := contextutils.WithLoggerValues(ctx,
			zap.String("gvk", request.Kind.String()),
			zap.String("role", matchingBinding.Spec.GetRoleRef().GetName()),
		)
		allowed, err := r.mcAdmissionValidator.ActionIsAllowed(childCtx, matchingBinding, &request)
		if err != nil {
			return admission.Denied(InternalError(err))
		}
		if !allowed {
			continue
		}
		return admission.Allowed(Allowed(request.UserInfo.Username))
	}
	return admission.Denied(NotAllowed(request.UserInfo.Username))
}
