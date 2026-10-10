//go:build !synthetic

package transcript

import (
	"context"
	"time"
)

func BeforeYield(context.Context) error         { return nil }
func StreamFailure() (<-chan time.Time, func()) { return nil, func() {} }

func BeforeAccept(context.Context) error { return nil }
