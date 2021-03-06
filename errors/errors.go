package errors

import (
	"net/http"

	"github.com/docker/distribution/registry/api/errcode"
)

// The notary API is on version 1, but URLs start with /v2/ to be consistent
// with the registry API
const errGroup = "notary.api.v1"

var (
	// ErrNoStorage lint comment
	ErrNoStorage = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "NO_STORAGE",
		Message:        "The server is misconfigured and has no storage.",
		Description:    "No storage backend has been configured for the server.",
		HTTPStatusCode: http.StatusInternalServerError,
	})
	// ErrNoFilename lint comment
	ErrNoFilename = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "NO_FILENAME",
		Message:        "No file/role name provided.",
		Description:    "No file/role name is provided to associate an update with.",
		HTTPStatusCode: http.StatusBadRequest,
	})
	// ErrInvalidRole lint comment
	ErrInvalidRole = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "INVALID_ROLE",
		Message:        "The role you are attempting to operate on is invalid.",
		Description:    "The user attempted to operate on a role that is not deemed valid.",
		HTTPStatusCode: http.StatusBadRequest,
	})
	// ErrMalformedJSON lint comment
	ErrMalformedJSON = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "MALFORMED_JSON",
		Message:        "JSON sent by the client could not be parsed by the server",
		Description:    "The client sent malformed JSON.",
		HTTPStatusCode: http.StatusBadRequest,
	})
	// ErrUpdating lint comment
	ErrUpdating = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "UPDATING",
		Message:        "An error has occurred while updating the TUF repository.",
		Description:    "An error occurred when attempting to apply an update at the storage layer.",
		HTTPStatusCode: http.StatusInternalServerError,
	})
	// ErrMetadataNotFound lint comment
	ErrMetadataNotFound = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "METADATA_NOT_FOUND",
		Message:        "You have requested metadata that does not exist.",
		Description:    "The user requested metadata that is not known to the server.",
		HTTPStatusCode: http.StatusNotFound,
	})
	// ErrMalformedUpload lint comment
	ErrMalformedUpload = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "MALFORMED_UPLOAD",
		Message:        "The body of your request is malformed.",
		Description:    "The user uploaded new TUF data and the server was unable to parse it as multipart/form-data.",
		HTTPStatusCode: http.StatusBadRequest,
	})
	// ErrGenericNotFound lint comment
	ErrGenericNotFound = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "GENERIC_NOT_FOUND",
		Message:        "You have requested a resource that does not exist.",
		Description:    "The user requested a non-specific resource that is not known to the server.",
		HTTPStatusCode: http.StatusNotFound,
	})
	// ErrNoCryptoService lint comment
	ErrNoCryptoService = errcode.Register(errGroup, errcode.ErrorDescriptor{
		Value:          "NO_CRYPTOSERVICE",
		Message:        "The server does not have a signing service configured.",
		Description:    "No signing service has been configured for the server and it has been asked to perform an operation that requires either signing, or key generation.",
		HTTPStatusCode: http.StatusInternalServerError,
	})
	// ErrUnknown is the generic internal server error
	ErrUnknown = errcode.ErrorCodeUnknown
)
