package errors

import "errors"

var (
	ErrUnknownTypeOfStorage = errors.New("unknown type of storage")
	ErrCommentNotFound      = errors.New("comment not found")
	ErrPostNotFound         = errors.New("post not found")
	ErrCanNotCreate         = errors.New("can not create object")
)
