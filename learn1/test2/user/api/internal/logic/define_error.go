package logic

import "errors"

var (
	ErrPasswordNotEqual = errors.New("password not equal")
	ErrQueryFailed      = errors.New("query db failed")
	ErrUserExist        = errors.New("user already exist")
	ErrWrongCreateToken = errors.New("create token failed")
)
