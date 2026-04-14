package appfr

import (
	"time"

	"github.com/robfig/cron/v3"
)

// AddCronJob registers a named cron job with the given schedule spec and function.
//
// Spec format supports either 5-field or 6-field (with optional seconds):
//
//	┌───────────── [second]   (0-59)  ← optional
//	│ ┌─────────── minute     (0-59)
//	│ │ ┌───────── hour       (0-23)
//	│ │ │ ┌─────── day        (1-31)
//	│ │ │ │ ┌───── month      (1-12)
//	│ │ │ │ │ ┌─── weekday    (0-6, Sunday=0)
//	│ │ │ │ │ │
//	* * * * * *   ← 6-field (with seconds)
//	  * * * * *   ← 5-field (without seconds)
//
// Special characters:
//   - *  every unit
//   - ,  list of values   (e.g. "1,15" = 1st and 15th)
//   - -  range            (e.g. "1-5"  = Monday to Friday)
//   - /  step             (e.g. "*/10" = every 10 units)
//
// Examples:
//
//	"* * * * * *"      every second          (6-field)
//	"0 * * * * *"      every minute          (6-field)
//	"* * * * *"        every minute          (5-field)
//	"0 0 * * * *"      every hour            (6-field)
//	"0 0 8 * * *"      every day at 08:00    (6-field)
//	"0 8 * * *"        every day at 08:00    (5-field)
//	"0 0 0 * * 1"      every Monday midnight (6-field)
//	"0 */30 * * * *"   every 30 minutes      (6-field)
func (a *App) AddCronJob(spec, jobName string, job func()) {
	if a.cron == nil {
		a.cron = newCron()
	}

	fn := func() {
		start := time.Now()
		a.logger.Infof("Starting cron job: %s", jobName)
		job()
		a.logger.Infof("Finished cron job: %s in %v", jobName, time.Since(start))
	}

	if _, err := a.cron.AddFunc(spec, fn); err != nil {
		a.logger.Errorf("[cron] failed to register job %s: %v", jobName, err)
	}
}

func newCron() *cron.Cron {
	return cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
}
