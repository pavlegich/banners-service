// Package utils contains additional methods for server.
package utils

import (
	"context"
	"fmt"
)

type contextKey int

const (
	// List of const variables contains variables for
	// put values into and get values from the context.
	ContextRoleKey contextKey = iota
)

// GetUserRoleFromContext finds and returns user role from the context.
func GetUserRoleFromContext(ctx context.Context) (string, error) {
	ctxValue := ctx.Value(ContextRoleKey)
	if ctxValue == nil {
		return "", fmt.Errorf("GetUserRoleFromContext: get context value failed")
	}
	userRole, ok := ctxValue.(string)
	if !ok {
		return "", fmt.Errorf("GetUserRoleFromContext: convert context value into string failed")
	}
	return userRole, nil
}
