package presidioclient

import (
	"context"
)

type PresidioClient interface {
	ProcessText(ctx context.Context, text string) (string, error)
}
