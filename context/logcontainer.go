package context

import (
	"bufio"
	"context"
	"runtime"
	"strings"
	"sync"
)

var UseBlack = true

type logContainer struct {
	lock        *sync.Mutex
	front, back int
	line        [MaxContainerSize]*logLine
}

func newLogContainer() *logContainer {
	return &logContainer{lock: new(sync.Mutex)}
}

func (c *logContainer) StackTrace(ctx context.Context, detailPrefix, header, footer, funcName string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	b := make([]byte, 4096) // adjust buffer size to be larger than expected stack
	n := runtime.Stack(b, false)
	s := string(b[:n])

	rd := strings.NewReader(s)
	c.addLogLineNoLock(ctx, LogLineFromPrintf(ctx, StackTraceInternal, Debug, funcName, header))
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		curr := ">   " + scanner.Text()
		c.addLogLineNoLock(ctx, LogLineFromPrintf(ctx, StackTraceInternal, Debug, funcName, curr))
	}

	c.addLogLineNoLock(ctx, LogLineFromPrintf(ctx, StackTraceInternal, Debug, funcName, header))
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

func (c *logContainer) Dump() {
	i := c.back
	for i < c.front {
		// put this line's data in the buffer
		l := c.line[i]
		l.Print()
		i = (i + 1) % MaxContainerSize
	}

}

func GetContainer(ctx context.Context) LogContainer {
	cont := ctx.Value(ParigotLogContainer)
	if cont == nil {
		return nil
	}
	container := cont.(*logContainer)
	return container

}
