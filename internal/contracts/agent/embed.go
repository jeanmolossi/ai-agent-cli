package contractsagent

import "context"

type EmbedProvider interface {
	Embed(ctx context.Context, content string) ([][]float32, error)
}
