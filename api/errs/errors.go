package errs

import (
	"errors"
)

// General errors
var (
	ErrInternal    = errors.New("internal_server_error")
	ErrInvalidJSON = errors.New("invalid_json")
)

// Token authentication errors
var (
	ErrMissingHeaders     = errors.New("missing_headers")
	ErrInvalidUserToken   = errors.New("invalid_user_token")
	ErrAccountNotVerified = errors.New("account_not_verified")
)

// Access errors
var (
	ErrDuplicateEmail           = errors.New("duplicate_email")
	ErrUnsupportedAccountType   = errors.New("unsupported_account_type")
	ErrInvalidLogin             = errors.New("invalid_login")
	ErrInvalidVerificationToken = errors.New("invalid_verification_token")
	ErrExpiredVerificationToken = errors.New("expired_verification_token")
)

// GCMError wraps an error from Google regarding GCM registration
type GCMError string

func (e GCMError) Error() string {
	return string(e)
}

// Device registration errors
var (
	ErrInvalidRegistrationToken = errors.New("invalid_registration_token")
	ErrDuplicateDeviceToken     = errors.New("duplicate_device_token")
	ErrUnableToRegisterDevice   = errors.New("unable_to_register_device")
	ErrGCMServiceUnavailable    = GCMError("gcm_service_unavailable")
)

// Errors from Google Login
var (
	ErrInvalidGoogleIDToken     = errors.New("invalid_google_id_token")
	ErrGoogleAccountNotVerified = errors.New("google_account_not_verified")
	ErrGoogleOAuthUnavailable   = errors.New("google_oauth_unavailable")
)
