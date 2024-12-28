package manager

import (
	"context"
)

type Repository interface {
	UpdateConfig(ctx context.Context, value any) error
	GetConfig(ctx context.Context, obj, defaultValue any) error
}
