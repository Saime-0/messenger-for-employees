package subix

import (
	"github.com/saime-0/messenger-for-employee/internal/repository"
	"github.com/saime-0/messenger-for-employee/pkg/scheduler"
)

type Subix struct {
	rooms     Rooms
	employees Employees
	clients   Clients
	repo      *repository.Repositories
	sched     *scheduler.Scheduler
}

func NewSubix(repo *repository.Repositories, sched *scheduler.Scheduler) *Subix {
	return &Subix{
		rooms:     Rooms{},
		employees: Employees{},
		clients:   Clients{},
		repo:      repo,
		sched:     sched,
	}
}
