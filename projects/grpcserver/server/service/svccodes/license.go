package svccodes

import (
	"context"

	"github.com/rotisserie/eris"

	"github.com/solo-io/go-utils/contextutils"

	"go.uber.org/zap"

	"github.com/solo-io/solo-projects/pkg/license"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LicenseStatusErrorType                     = "license"
	LicenseStatusErrorSubject_InvalidLicense   = "invalid license"
	PermissionsErrorType                       = "permissions"
	PermissionsErrorSubject_InvalidPermissions = "invalid permissions"
	invalidLicenseLoggedMutationErrorMsg       = "This feature requires an Enterprise Gloo license." +
		" Visit http://www.solo.io/gloo-trial to request a trial license."
	emptyLicenseLoggedMutationErrorMsg = "This feature requires admin permissions. Has your Gloo admin" +
		" removed the Gloo Enterprise license key from the Gloo UI apiserver?"
)

var (
	InvalidLicenseErrorDetails = func(err error) *errdetails.PreconditionFailure_Violation {
		return &errdetails.PreconditionFailure_Violation{
			Type:        LicenseStatusErrorType,
			Subject:     LicenseStatusErrorSubject_InvalidLicense,
			Description: err.Error(),
		}
	}
	InvalidPermissionsErrorDetails = func(err error) *errdetails.PreconditionFailure_Violation {
		return &errdetails.PreconditionFailure_Violation{
			Type:        PermissionsErrorType,
			Subject:     PermissionsErrorSubject_InvalidPermissions,
			Description: err.Error(),
		}
	}
)

func CheckLicenseForGlooUiMutations(ctx context.Context, licenseClient license.Client) error {
	if err := licenseClient.IsLicenseValid(); err != nil {

		var details *errdetails.PreconditionFailure
		if eris.Is(err, license.LicenseEmptyError) {
			contextutils.LoggerFrom(ctx).Warnw(emptyLicenseLoggedMutationErrorMsg, zap.Error(err))
			// attach details to error message to distinguish this from other PermissionDenied errors
			details = &errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					InvalidPermissionsErrorDetails(err),
				},
			}
		} else {
			contextutils.LoggerFrom(ctx).Warnw(invalidLicenseLoggedMutationErrorMsg, zap.Error(err))
			// attach details to error message to distinguish this from other PermissionDenied errors
			details = &errdetails.PreconditionFailure{
				Violations: []*errdetails.PreconditionFailure_Violation{
					InvalidLicenseErrorDetails(err),
				},
			}
		}

		eStatus := status.New(codes.PermissionDenied, err.Error())
		var appendErr error
		eStatus, appendErr = eStatus.WithDetails(details)
		if appendErr != nil {
			// this should not happen
			return appendErr
		}
		return eStatus.Err()
	}
	return nil
}
