package piper

import (
	"context"
	"github.com/saime-0/messenger-for-employee/internal/cdl"
	"github.com/saime-0/messenger-for-employee/internal/config"
	"github.com/saime-0/messenger-for-employee/internal/healer"
	"github.com/saime-0/messenger-for-employee/internal/repository"
	"github.com/saime-0/messenger-for-employee/internal/res"
	"sync"
)

type Pipeline struct {
	Nodes      map[string]*Node
	mu         *sync.Mutex
	repos      *repository.Repositories
	healer     *healer.Healer
	dataloader *cdl.Dataloader

	cfg *config.Config2
}

func NewPipeline(cfg *config.Config2, repos *repository.Repositories, healer *healer.Healer, dataloader *cdl.Dataloader) *Pipeline {
	return &Pipeline{
		Nodes:      map[string]*Node{},
		mu:         new(sync.Mutex),
		repos:      repos,
		healer:     healer,
		dataloader: dataloader,
		cfg:        cfg,
	}
}

func (p *Pipeline) NodeFromContext(ctx context.Context) *Node {
	return ctx.Value(res.CtxNode).(*Node)
}
