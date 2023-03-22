package dailyTask

import (
	"context"
	"fmt"
	"github.com/bastengao/chinese-holidays-go/holidays"
	"math/rand"
	"time"
)

func init() {
	LoggerInit("./dailyTask.log")
	SetLogLevel(ParseLogLevel("debug"))
}

type DailyTask struct {
	taskName              string
	hour                  int    // 执行任务的小时
	min                   int    // 执行任务的分钟
	sec                   int    // 执行任务的秒钟
	randomDelayMaxSeconds int    // 最大延迟时间
	handler               func() // 任务处理程序
	ignoreHoliday         bool
	queryer               holidays.Queryer
	ctx                   context.Context
	cancel                context.CancelFunc
	targetTime            time.Time
}

func NewDailyTask(name string, hour, min, sec, randomDelayMaxSeconds int, ignoreHoliday bool, handler func(), ctx context.Context) (*DailyTask, error) {
	if hour < 0 || min < 0 || sec < 0 || randomDelayMaxSeconds < 0 {
		return nil, fmt.Errorf("invalid args, must bigger than or euqal to 0")
	}

	if hour*24*60+min*60+sec+randomDelayMaxSeconds >= 86400 {

		return nil, fmt.Errorf("invalid args, time must in 00:00-23:59:59")
	}

	myCtx, MyCancel := context.WithCancel(ctx)
	t := &DailyTask{
		taskName:              name,
		hour:                  hour,
		min:                   min,
		sec:                   sec,
		randomDelayMaxSeconds: randomDelayMaxSeconds,
		ignoreHoliday:         ignoreHoliday,
		queryer:               nil,
		handler:               handler,
		ctx:                   myCtx,
		cancel:                MyCancel,
	}

	// 是否智能跳过节假日
	if ignoreHoliday {
		queryer, err := holidays.BundleQueryer()
		if err != nil {
			Err("BundleQueryer failed: %s", err)

			return nil, err
		}

		t.queryer = queryer
	}

	return t, nil
}

func (t *DailyTask) Close() error {
	if t.cancel != nil {
		t.cancel()
		t.cancel = nil
	}

	return nil
}

func (t *DailyTask) Start() {
	Info("%s start", t)
	defer func() {
		Info("%t finished loop", t)
	}()

	for {
		now := time.Now()
		// 计算今天的任务执行时间
		targetTime := time.Date(now.Year(), now.Month(), now.Day(), t.hour, t.min, t.sec, 0, time.Local)

		// 如果今天的执行时间已过，就移到明天
		if now.After(targetTime) {
			targetTime = targetTime.Add(24 * time.Hour)
		}

		// 等待到执行时间
		rndDelay := time.Duration(rand.Intn(t.randomDelayMaxSeconds)) * time.Second
		targetTime = targetTime.Add(rndDelay)

		t.targetTime = targetTime

		delay := t.TimeToExecute()

		Info("%s target time %s, time to execute %s", t, targetTime, delay)

		timer := time.NewTimer(delay)

		select {
		case <-timer.C:
			if t.ignoreHoliday {
				// 检查跳过节假日
				isHoliday, err := t.queryer.IsHoliday(time.Now())
				if err != nil {
					Err("%s check holiday error: %s", t.taskName, err)
				} else {
					if isHoliday { // 节假日跳过
						break
					}
				}
			}

			Info("%s execute handler at %s", t, time.Now())
			t.handler()

		case <-t.ctx.Done():
			Info("%s ctx done break", t)
			return
		}
	}
}

func (t *DailyTask) TargetTime() time.Time {
	return t.targetTime
}

func (t *DailyTask) TimeToExecute() time.Duration {
	return t.TargetTime().Sub(time.Now().Truncate(time.Second))
}

func (t *DailyTask) String() string {
	return fmt.Sprintf("[%s %02d-%02d-%02d:%04d %t]", t.taskName, t.hour, t.min, t.sec, t.randomDelayMaxSeconds, t.ignoreHoliday)
}
