package errs

import "errors"

var (
	ErrInvalidLoginOrPassword = errors.New("Invalid login or password")
	ErrInvalidOrderNumber     = errors.New("invalid otder number")
	ErrOrderExist             = errors.New("order has already been added")
	ErrUnreachableOrder       = errors.New("order has already been added by another user")
	ErrNoReward               = errors.New("user does not have enough rewards")
	ErrLoginAlreadyExist      = errors.New("login already exist")
)
