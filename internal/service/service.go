package service

import (
	"github.com/saime-0/messenger-for-employee/internal/cache"
	"github.com/saime-0/messenger-for-employee/internal/repository"
	"github.com/saime-0/messenger-for-employee/pkg/scheduler"
)

type Services struct {
	Repos     *repository.Repositories
	Scheduler *scheduler.Scheduler
	Cache     *cache.Cache
}
