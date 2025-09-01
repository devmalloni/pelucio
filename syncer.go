package pelucio

import (
	"context"
	"sync"
)

// Add syncrhonization to operations if pelucio is used concurrently.
// For simple usecases, it can be used with NewMutextSyncer().
// Default is NoOpSyncer, which does nothing.
//
// For horizontal scaling, it is recommended to implement a distributed lock
// and use it as Syncer.
type Syncer interface {
	Lock(ctx context.Context)
	Unlock(ctx context.Context)
}

type NoOpSyncer struct{}

func (s *NoOpSyncer) Lock(ctx context.Context)   {}
func (s *NoOpSyncer) Unlock(ctx context.Context) {}

type MutexSyncer struct {
	l *sync.Mutex
}

func NewMutexSyncer() *MutexSyncer {
	return &MutexSyncer{
		l: &sync.Mutex{},
	}
}
func (p *MutexSyncer) Lock(ctx context.Context) {
	p.l.Lock()
}
func (p *MutexSyncer) Unlock(ctx context.Context) {
	p.l.Unlock()
}
