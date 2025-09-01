package pelucio

import "github.com/devmalloni/pelucio/x/xtime"

type PelucionOpt func(p *Pelucio)

func WithReadWriter(rw ReadWriter) PelucionOpt {
	return func(p *Pelucio) {
		p.readWriter = rw
	}
}

func WithClock(clock xtime.Clock) PelucionOpt {
	return func(p *Pelucio) {
		p.clock = clock
	}
}

func WithSyncer(syncer Syncer) PelucionOpt {
	return func(p *Pelucio) {
		p.syncer = syncer
	}
}
