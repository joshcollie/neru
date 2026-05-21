//go:build !darwin

package accessibility

import (
	"context"

	derrors "github.com/y3owk1n/neru/internal/core/errors"
)

// StreamClickableNodes returns an empty closed channel on non-Darwin platforms.
// This satisfies the AXClient interface contract without requiring the darwin-only
// tree traversal implementation.
func (c *InfraAXClient) StreamClickableNodes(
	_ context.Context,
	root AXElement,
	_ []string,
	_ int,
) (<-chan AXNode, error) {
	if root == nil {
		return nil, derrors.New(derrors.CodeInvalidInput, "element is nil")
	}

	ch := make(chan AXNode)
	close(ch)

	return ch, nil
}
