// Package errors contains variables with error explainations.
package errors

import "errors"

var (
	ErrBannerNotFound        = errors.New("banner not found")
	ErrBannerInCacheNotFound = errors.New("banner in cache not found")
	ErrBannerExpired         = errors.New("banner content expired")
	ErrBannerNotAllowed      = errors.New("not allowed for user")
)
