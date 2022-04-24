package resolver

import (
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/res"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"github.com/saime-0/messenger-for-employee/pkg/kit"
	"time"
)

// да, я знаю что это криндж, не душите

func (r *Resolver) RegularSchedule(interval int64) (err error) {
	ready := make(chan int8)

	regFn := func() {
		end := kit.After(interval)

		err = r.prepareScheduleRefreshSessions(end)
		if err != nil {

			r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
			panic(err)
		}

		select {
		case ready <- 1:
			//  сигнал о готовности был услышан
		default:
			// сигнал о готовности никто не услышал
		}
	}

	regFn()

	go func() {
		for {
			runAt := kit.After(interval)
			_, err = r.Services.Scheduler.AddTask(
				regFn,
				runAt,
			)
			if err != nil {
				panic(err)
			}
			r.Services.Cache.Set(res.CacheNextRunRegularScheduleAt, runAt)

			// начато прослушивание канала
			<-ready
			// получен сигнал о готовности
		}
	}()

	return nil

}

func (r *Resolver) prepareScheduleRefreshSessions(before int64) (err error) {
	sessions, err := r.Services.Repos.Prepares.ScheduleRefreshSessions(before)
	if err != nil {
		return err
	}
	for _, rs := range sessions {

		if rs.Exp <= time.Now().Unix() {
			// удаляю сессию, тк она уже истекла
			err := r.Services.Repos.Employees.DeleteRefreshSession(rs.ID)
			if err != nil {

				r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
			}
			continue
		}
		_, err = r.Services.Scheduler.AddTask(
			func() {
				// спланированное удаление
				err := r.Services.Repos.Employees.DeleteRefreshSession(rs.ID)
				if err != nil {

					r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()+""))
				}
			},
			rs.Exp,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
