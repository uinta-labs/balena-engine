// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package batch

import (
	"github.com/aws/aws-sdk-go/private/protocol"
)

const (

	// ErrCodeClientException for service response error code
	// "ClientException".
	//
	// These errors are usually caused by a client action, such as using an action
	// or resource on behalf of a user that doesn't have permissions to use the
	// action or resource, or specifying an identifier that is not valid.
	ErrCodeClientException = "ClientException"

	// ErrCodeServerException for service response error code
	// "ServerException".
	//
	// These errors are usually caused by a server issue.
	ErrCodeServerException = "ServerException"
)

var exceptionFromCode = map[string]func(protocol.ResponseMetadata) error{
	"ClientException": newErrorClientException,
	"ServerException": newErrorServerException,
}
