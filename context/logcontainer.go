package context

import (
	"bufio"
	"context"
	"runtime"
	"strings"
	"sync"

	"github.com/iansmith/parigot/apishared"
)

var UseBlack = true

type logContainer struct {
	lock        *sync.Mutex
	front, back int
	line        [MaxContainerSize]*logLine
	origin      string
}

func newLogContainer(orig string) *logContainer {
	return &logContainer{lock: new(sync.Mutex), origin: orig}
}

func (c *logContainer) StackTrace(ctx context.Context) {
	c.lock.Lock()
	defer c.lock.Unlock()

	b := make([]byte, 4096) // adjust buffer size to be larger than expected stack
	n := runtime.Stack(b, false)
	s := string(b[:n])

	fn := pullFunc(ctx, "")
	src := PullSource(ctx, StackTraceInternal)

	c.addLogLineNoLock(ctx, NewLogLine(ctx, src, Debug,
		fn, true, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"))

	rd := strings.NewReader(s)
	scanner := bufio.NewScanner(rd)

	total := 0
	for scanner.Scan() {
		curr := "> " + scanner.Text()
		if total+len(curr) > apishared.ExpectedStackDumpSize {
			break
		}
		c.addLogLineNoLock(ctx, NewLogLine(ctx, src, Debug, fn, true, curr))
	}

	c.addLogLineNoLock(ctx, NewLogLine(ctx, src, Debug, fn, true,
		"<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"))
}

func (c *logContainer) AddLogLine(ctx context.Context, l LogLine) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.addLogLineNoLock(ctx, l.(*logLine))
}

func (c *logContainer) addLogLineNoLock(ctx context.Context, l *logLine) {
	c.line[c.front] = l
	c.front = (c.front + 1) % MaxContainerSize
	c.front %= MaxContainerSize
	if c.front == c.back {
		c.back = (c.back + 1) % MaxContainerSize
	}
}

func (c *logContainer) Dump(ctx context.Context) {
	if c.front == c.back {
		return
	}
	i := c.back
	for i < c.front {
		// put this line's data in the buffer
		l := c.line[i]
		l.Print()
		i = (i + 1) % MaxContainerSize
	}
	c.back = i
}

func GetContainer(ctx context.Context) LogContainer {
	cont := ctx.Value(ParigotLogContainer)
	if cont == nil {
		return nil
	}
	container := cont.(*logContainer)
	return container

}
