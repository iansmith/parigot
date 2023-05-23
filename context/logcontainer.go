package context

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
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
	c.addLogLineNoLock(ctx, LogLineFromString(ctx, header, UnknownS, stackTraceInternal))
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		curr := ">   " + scanner.Text()
		c.addLogLineNoLock(ctx, LogLineFromString(ctx, curr, UnknownS, stackTraceInternal))
	}

	c.addLogLineNoLock(ctx, LogLineFromString(ctx, footer, UnknownS, stackTraceInternal))
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
	buf := &bytes.Buffer{}
	for i < c.front {
		// put this line's data in the buffer
		l := c.line[i]
		lastIsCR := false
		for k := 0; k < MaxLineLen; k++ {
			if l.data[k] == 0 {
				break
			}
			buf.WriteByte(l.data[k])
			if l.data[k] == 10 {
				lastIsCR = true
			} else {
				lastIsCR = false
			}
		}
		if !lastIsCR {
			buf.WriteString("\n")
		}
		i = (i + 1) % MaxContainerSize

		fmt.Printf(l.color()(buf.String()))
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
