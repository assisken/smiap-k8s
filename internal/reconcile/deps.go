package reconcile

import (
	"context"
	"smiap-k8s/internal/differ"
)

type ActionRunner interface {
	Create(ctx context.Context, item differ.DiffItem) error
	Delete(ctx context.Context, item differ.DiffItem) error
	Update(ctx context.Context, item differ.DiffItem) error
}
