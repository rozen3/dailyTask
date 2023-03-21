# API


Creates a new daily task that runs at the specified time with a random delay.  
parameter ignoreHoliday specifies whether to skip holidays.  
parameter handler specifies the task handler function.  

```
import "github.com/rozen3/dailyTask"

f := func() {
    println("hello world!)
}

t, err := NewDailyTask("task1", 19, 30, 0, 1800, true, f, context.TODO())

// Starts the task.
t.Start()

// Stops the task.
t.Close()


```