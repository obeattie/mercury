package client

import (
	"fmt"

	terrors "github.com/mondough/typhon/errors"
)

const (
	errUidField      = "Client-Uid"
	errServiceField  = "Client-Service"
	errEndpointField = "Client-Endpoint"
)

type ErrorSet []*terrors.Error

// Copy returns a new ErrorSet containing the same errors as the receiver
func (es ErrorSet) Copy() ErrorSet {
	result := make(ErrorSet, len(es))
	copy(result, es)
	return result
}

// ForUid returns the error for a given request uid (or nil)
func (es ErrorSet) ForUid(uid string) *terrors.Error {
	for _, e := range es {
		if euid, ok := e.Params[errUidField]; ok && euid == uid {
			return e
		}
	}
	return nil
}

// Any returns whether there are any contained errors
func (es ErrorSet) Any() bool {
	return len(es) > 0
}

// Errors returns a map of request uids to their error, for requests which had errors
func (es ErrorSet) Errors() map[string]*terrors.Error {
	result := make(map[string]*terrors.Error, len(es)) // Never return nil; with a map it's just fraught
	for _, err := range es {
		result[err.Params[errUidField]] = err
	}
	return result
}

// IgnoreCode returns a new ErrorSet without errors of the given codes
func (es ErrorSet) IgnoreCode(codes ...string) ErrorSet {
	if len(codes) == 0 {
		return es
	}
	codesMap := make(map[string]struct{}, len(codes))
	for _, c := range codes {
		codesMap[c] = struct{}{}
	}

	result := make(ErrorSet, 0, len(es)-len(codes))
	for _, err := range es {
		if _, excluded := codesMap[err.Code]; !excluded {
			result = append(result, err)
		}
	}
	return result
}

// IgnoreEndpoint returns a new ErrorSet without errors from the given service endpoint
func (es ErrorSet) IgnoreEndpoint(service, endpoint string) ErrorSet {
	result := make(ErrorSet, 0, len(es)-1)
	for _, err := range es {
		if !(err.Params[errServiceField] == service && err.Params[errEndpointField] == endpoint) {
			result = append(result, err)
		}
	}
	return result
}

// IgnoreService returns a new ErrorSet without errors from the given service(s)
func (es ErrorSet) IgnoreService(services ...string) ErrorSet {
	if len(services) == 0 {
		return es
	}
	servicesMap := stringsMap(services...)
	result := make(ErrorSet, 0, len(es)-len(services))
	for _, err := range es {
		if _, excluded := servicesMap[err.Params[errServiceField]]; !excluded {
			result = append(result, err)
		}
	}
	return result
}

// IgnoreUid returns a new ErrorSet without errors from the given request uid(s)
func (es ErrorSet) IgnoreUid(uids ...string) ErrorSet {
	if len(uids) == 0 {
		return es
	}
	uidsMap := stringsMap(uids...)
	result := make(ErrorSet, 0, len(es)-len(uids))
	for _, err := range es {
		if _, excluded := uidsMap[err.Params[errUidField]]; !excluded {
			result = append(result, err)
		}
	}
	return result
}

// sanitiseContext takes an error context and removes client-specific things from it (in-place)
func (es ErrorSet) sanitiseContext(ctx map[string]string) {
	delete(ctx, errUidField)
	delete(ctx, errServiceField)
	delete(ctx, errEndpointField)
}

// Combined returns a combined error from the set. If there is only one error, it is returned unmolested. If there are
// more, they are all "flattened" into a single error. Where codes differ, they are normalised to that with the lowest
// index.
func (es ErrorSet) Combined() error {
	switch len(es) {
	case 0:
		return nil

	case 1:
		return es[0]

	default:
		msg := fmt.Sprintf("%s, and %d more errors", es[0].Message, len(es)-1)
		result := terrors.New(es[0].Code, msg, nil)

		params := []map[string]string{}
		for _, err := range es {
			// TODO how do we replicate this flattening logic with string error codes
			// if err.Code < result.Code {
			result.Code = err.Code
			// }
			params = append(params, err.Params)
		}

		result.Params = mergeContexts(params...)
		es.sanitiseContext(result.Params)
		return result
	}
}

// Error satisfies Go's Error interface
func (es ErrorSet) Error() string {
	if err := es.Combined(); err != nil {
		return err.Error()
	}
	return ""
}
