package vault

import (
    "context"
    "errors"
    "os"
    "time"

    "github.com/avast/retry-go"
    vault "github.com/hashicorp/vault/api"
    awsauth "github.com/hashicorp/vault/api/auth/aws"
    "github.com/rotisserie/eris"
    "github.com/solo-io/gloo/pkg/utils"
    v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
    "github.com/solo-io/go-utils/contextutils"
)

// In an ideal world, we would re-use the mocks provided by an external library.
// Since the vault.AuthMethod interface does not have corresponding mocks, we have to define our own.
//go:generate mockgen -destination mocks/mock_auth.go -package mocks github.com/hashicorp/vault/api AuthMethod

var _ vault.AuthMethod = &staticTokenAuthMethod{}
var _ vault.AuthMethod = &retryableAuthMethod{}

var (
    mLastLoginSuccess = utils.MakeLastValueCounter("gloo.solo.io/vault/aws/last_login_success", "Timestamp of last successful authentication of vault with AWS IAM")
    mLastLoginFailure = utils.MakeLastValueCounter("gloo.solo.io/vault/aws/last_login_failure", "Timestamp of last failed authentication of vault with AWS IAM")
    mLoginSuccesses   = utils.MakeSumCounter("gloo.solo.io/vault/aws/login_successes", "Number of successful authentications of vault with AWS IAM")
    mLoginFailures    = utils.MakeSumCounter("gloo.solo.io/vault/aws/login_failures", "Number of failed authentications of vault with AWS IAM")
)

var (
    ErrVaultAuthentication = errors.New("unable to authenticate to Vault")
    ErrNoAuthInfo          = errors.New("no auth info was returned after login")
)

// newAuthMethodForSettings returns a vault auth method based on the provided settings.
func newAuthMethodForSettings(ctx context.Context, vaultSettings *v1.Settings_VaultSecrets) (vault.AuthMethod, error) {
    switch tlsCfg := vaultSettings.GetAuthMethod().(type) {
    case *v1.Settings_VaultSecrets_AccessToken:
        return newStaticTokenAuthMethod(tlsCfg.AccessToken), nil

    case *v1.Settings_VaultSecrets_Aws:
        awsAuth, err := newAwsAuthMethod(ctx, tlsCfg.Aws)
        if err != nil {
            return nil, err
        }

        return newRetryableAuthMethod(awsAuth), nil

    default:
        // We don't have one of the defined auth methods, so try to fall back to the
        // deprecated token field before erroring
        return newStaticTokenAuthMethod(vaultSettings.GetToken()), nil
    }
}

func newStaticTokenAuthMethod(token string) vault.AuthMethod {
    return &staticTokenAuthMethod{
        token: token,
    }
}

type staticTokenAuthMethod struct {
    token string
}

func (s *staticTokenAuthMethod) Login(ctx context.Context, _ *vault.Client) (*vault.Secret, error) {
    if s.token == "" {
        utils.Measure(ctx, mLastLoginFailure, time.Now().Unix())
        utils.MeasureOne(ctx, mLoginFailures)
        return nil, eris.Errorf("unable to determine vault authentication method. check Settings configuration")
    }

    contextutils.LoggerFrom(ctx).Debugf("Successfully authenticated to Vault with static token")
    utils.Measure(ctx, mLastLoginSuccess, time.Now().Unix())
    utils.MeasureOne(ctx, mLoginSuccesses)
    return &vault.Secret{
        Auth: &vault.SecretAuth{
            ClientToken: s.token,
        },
    }, nil
}

func newRetryableAuthMethod(authMethod vault.AuthMethod, retryOptions ...retry.Option) vault.AuthMethod {
    // Standard retry options, which can be overridden by the retryOptions parameter
    defaultRetryOptions := []retry.Option{
        retry.Delay(1 * time.Second),
        retry.DelayType(retry.BackOffDelay),
        retry.Attempts(10),
    }
    options := append(defaultRetryOptions, retryOptions...)

    return &retryableAuthMethod{
        authMethod:   authMethod,
        retryOptions: options,
    }
}

type retryableAuthMethod struct {
    authMethod   vault.AuthMethod
    retryOptions []retry.Option
}

func (r *retryableAuthMethod) Login(ctx context.Context, client *vault.Client) (*vault.Secret, error) {
    var (
        loginResponse *vault.Secret
        loginErr      error
    )

    loginOnce := func() error {
        loginResponse, loginErr = r.authMethod.Login(ctx, client)
        if loginErr != nil {
            contextutils.LoggerFrom(ctx).Errorf("unable to authenticate to Vault: %v", loginErr)
            utils.Measure(ctx, mLastLoginFailure, time.Now().Unix())
            utils.MeasureOne(ctx, mLoginFailures)
            return ErrVaultAuthentication
        }

        if loginResponse == nil {
            contextutils.LoggerFrom(ctx).Errorf("no auth info was returned after login")
            utils.Measure(ctx, mLastLoginFailure, time.Now().Unix())
            utils.MeasureOne(ctx, mLoginFailures)
            return ErrNoAuthInfo
        }

        utils.Measure(ctx, mLastLoginSuccess, time.Now().Unix())
        utils.MeasureOne(ctx, mLoginSuccesses)
        contextutils.LoggerFrom(ctx).Debugf("Successfully authenticated to Vault %v", loginResponse)
        return nil
    }

    loginErr = retry.Do(loginOnce, r.retryOptions...)
    return loginResponse, loginErr
}

func newAwsAuthMethod(_ context.Context, aws *v1.Settings_VaultAwsAuth) (*awsauth.AWSAuth, error) {
    // The AccessKeyID and SecretAccessKey are not required in the case of using temporary credentials from assumed roles with AWS STS or IRSA.
    // STS: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html
    // IRSA: https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html
    var possibleErrStrings []string
    if accessKeyId := aws.GetAccessKeyId(); accessKeyId != "" {
        os.Setenv("AWS_ACCESS_KEY_ID", accessKeyId)
    } else {
        possibleErrStrings = append(possibleErrStrings, "access key id must be defined for AWS IAM auth")
    }

    if secretAccessKey := aws.GetSecretAccessKey(); secretAccessKey != "" {
        os.Setenv("AWS_SECRET_ACCESS_KEY", secretAccessKey)
    } else {
        possibleErrStrings = append(possibleErrStrings, "secret access key must be defined for AWS IAM auth")
    }

    // if we have only partial configuration set
    if len(possibleErrStrings) == 1 {
        return nil, errors.New("only partial credentials were provided for AWS IAM auth: " + possibleErrStrings[0])
    }

    // At this point, we either have full auth configuration set, or are in an ec2 environment, where vault will infer the credentials.
    loginOptions := []awsauth.LoginOption{awsauth.WithIAMAuth()}

    if role := aws.GetVaultRole(); role != "" {
        loginOptions = append(loginOptions, awsauth.WithRole(role))
    }

    if region := aws.GetRegion(); region != "" {
        loginOptions = append(loginOptions, awsauth.WithRegion(region))
    }

    if iamServerIdHeader := aws.GetIamServerIdHeader(); iamServerIdHeader != "" {
        loginOptions = append(loginOptions, awsauth.WithIAMServerIDHeader(iamServerIdHeader))
    }

    if mountPath := aws.GetMountPath(); mountPath != "" {
        loginOptions = append(loginOptions, awsauth.WithMountPath(mountPath))
    }

    if sessionToken := aws.GetSessionToken(); sessionToken != "" {
        os.Setenv("AWS_SESSION_TOKEN", sessionToken)
    }

    return awsauth.NewAWSAuth(loginOptions...)
}
