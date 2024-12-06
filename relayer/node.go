package relayer

import "context"

// TODO: bootstrap nodes
type Relayer struct{}

func New() (*Relayer, error) {
	return &Relayer{}, nil
}

func (r *Relayer) Run(ctx context.Context) error {
	return nil
}
