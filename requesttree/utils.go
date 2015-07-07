package requesttree

import (
	"golang.org/x/net/context"
)

const parentIdCtxKey = parentIdHeader

// ParentRequestIdFor returns the parent request ID for the provided context (if any).
func ParentRequestIdFor(ctx context.Context) string {
	switch v := ctx.Value(parentIdCtxKey).(type) {
	case string:
		return v
	}
	return ""
}
