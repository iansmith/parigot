package eng

import (
	"bufio"
	"context"
	"io"

	pcontext "github.com/iansmith/parigot/context"

	"github.com/tetratelabs/wazero/experimental/logging"
)

type rawLineReader struct {
	ctx  context.Context
	src  pcontext.Source
	out  *io.PipeWriter
	in   *bufio.Scanner
	cont pcontext.LogContainer
}

func (w *rawLineReader) read(ctx context.Context) {
	for {
		err := w.in.Scan()
		if err == true && w.in.Err() != nil {
			pcontext.Errorf(ctx, "internal error with pipe inside wazerowriter: %v", w.in.Err())
			continue
		}
		t := w.in.Text() + "\n"
		ll := pcontext.NewLogLine(ctx, w.src, pcontext.Debug, w.src.String(), true, t)
		w.addLineToOrigContainer(ll)
		pcontext.Dump(ctx)
	}
}
func (w *rawLineReader) addLineToOrigContainer(ll pcontext.LogLine) {
	cont := pcontext.GetContainer(w.ctx)
	cont.AddLogLine(w.ctx, ll)
	cont.Dump(w.ctx)
}

func (w *rawLineReader) Write(p []byte) (int, error) {
	return w.out.Write(p)
}

func (w *rawLineReader) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func newRawLineReader(ctx context.Context, src pcontext.Source) logging.Writer {
	rd, wr := io.Pipe()
	buffered := bufio.NewScanner(rd)

	nCtx := wazeroContext(ctx)
	rlr := &rawLineReader{
		ctx:  nCtx,
		src:  src,
		out:  wr,
		in:   buffered,
		cont: pcontext.GetContainer(nCtx),
	}
	go rlr.read(rlr.ctx)
	return rlr
}
