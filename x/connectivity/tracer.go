package connectivity

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type contextKey struct{}

var activeSpanKey = contextKey{}

type GenericSpan interface {
	End()
	RecordError(err error)
}

type GenericTracer interface {
	Start(ctx context.Context, spanName string) (context.Context, GenericSpan)
}

// customTracer is an implementation of GenericTracer.
type customTracer struct{}

func NewCustomTracer() GenericTracer {
	return &customTracer{}
}

func (t *customTracer) Start(ctx context.Context, spanName string) (context.Context, GenericSpan) {
	parentSpan, _ := ctx.Value(activeSpanKey).(GenericSpan)

	span := &customSpan{
		name:       spanName,
		startTime:  time.Now(),
		parentSpan: parentSpan,
	}
	fmt.Printf("Span: %s, Start Time: %v, Parent: %v\n", span.name, span.startTime, span.parentSpan)
	// Store the new span in the context
	newCtx := context.WithValue(ctx, activeSpanKey, span)
	return newCtx, span
}

type customSpan struct {
	name       string
	startTime  time.Time
	endTime    time.Time
	err        error
	parentSpan GenericSpan
	mutex      sync.Mutex
}

func (s *customSpan) End() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.endTime = time.Now()
	fmt.Printf("Span: %s, Duration: %v\n", s.name, s.endTime.Sub(s.startTime))
}

func (s *customSpan) RecordError(err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.err = err
	fmt.Println(err)
}
