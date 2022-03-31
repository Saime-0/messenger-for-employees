package cdl

import (
	"database/sql"
	"github.com/saime-0/messenger-for-employee/internal/healer"
	"time"
)

type (
	ID      = string
	chanPtr = string
	Any     = interface{}

	RequestsCount uint8
	//categories    map[CategoryName]*parentCategory
)

type Dataloader struct {
	wait             time.Duration
	capacityRequests RequestsCount
	categories       *Categories
	db               *sql.DB
	healer           *healer.Healer
}

func NewDataloader(wait time.Duration, maxBatch RequestsCount, db *sql.DB, hlr *healer.Healer) *Dataloader {
	d := &Dataloader{
		wait:             wait,
		capacityRequests: maxBatch,
		db:               db,
		healer:           hlr,
	}
	d.ConfigureDataloader()
	return d
}
