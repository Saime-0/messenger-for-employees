package piper

import (
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/cdl"
	"github.com/saime-0/messenger-for-employee/internal/clog"
	"github.com/saime-0/messenger-for-employee/internal/config"
	"github.com/saime-0/messenger-for-employee/internal/healer"
	"github.com/saime-0/messenger-for-employee/internal/repository"
	"github.com/saime-0/messenger-for-employee/internal/resp"
	"github.com/saime-0/messenger-for-employee/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Rows []interface{}

type Request struct {
	Timestamp time.Time
	ID        string
	Status    int
	Method    string
	Path      string
	Duration  string
	Body      *Rows
}

type Method struct {
	Method   string
	Vars     *bson.M
	startAt  time.Time
	Duration string
	Body     *Rows
}

var _ clog.Logger = (*Node)(nil)

type Node struct {
	repos      *repository.Repositories
	Healer     *healer.Healer
	Dataloader *cdl.Dataloader
	err        **model.AdvancedError

	ID            *string
	RootContainer interface{}
	scope         *Rows
	ScopeMethod   *Method

	cfg *config.Config2

	Request *Request
}

func (p *Pipeline) CreateNode(id string) (*Node, *Request) {
	scope := &Rows{}

	request := &Request{
		Timestamp: time.Now(),
		ID:        kit.RandomSecret(6),
		Body:      scope,
	}

	n := &Node{
		repos:      p.repos,
		Healer:     p.healer,
		Dataloader: p.dataloader,
		err:        new(*model.AdvancedError),

		ID: &id,
		RootContainer: bson.M{
			"Request": request,
		},
		scope: scope,

		cfg: p.cfg,

		Request: request,
	}
	p.Nodes[id] = n
	return n, request
}

func (p *Pipeline) DeleteNode(id string) {
	delete(p.Nodes, id)
}

func (n *Node) Execute() {
	if len(*n.Request.Body) != 0 {
		n.Healer.Log(n.RootContainer)
	}
}

func (n *Node) SwitchMethod(name string, vars *bson.M) {
	meth := &Method{
		Method:  name,
		Vars:    vars,
		Body:    &Rows{},
		startAt: time.Now(),
	}
	*n.scope = append(*n.scope, meth)
	n.scope = meth.Body

	n.ScopeMethod = meth
}
func (n *Node) MethodTiming() {
	if n.ScopeMethod != nil {
		n.ScopeMethod.Duration = time.Since(n.ScopeMethod.startAt).String()
	}
}

func (n Node) SetError(code resp.ErrCode, msg string) {
	*n.err = resp.Error(code, msg)
}
func (n Node) GetError() *model.AdvancedError {
	return *n.err
}
