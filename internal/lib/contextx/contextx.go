//go:build !modern

package contextx

import "context"

func AfterFunc(ctx context.Context, f func()) {
	go func() {
		<-ctx.Done()
		f()
	}()
}
