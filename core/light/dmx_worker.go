package light

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nickysemenza/hyperion/util/metrics"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

//DMXState holds the DMX512 values for each channel
type DMXState struct {
	universes map[int][]byte
	m         sync.Mutex
}

func (s *DMXState) getValue(universe, channel int) int {
	return int(s.universes[universe][channel-1])
}

func (s *DMXState) set(ctx context.Context, ops ...dmxOperation) error {
	ctx, span := trace.StartSpan(ctx, "setDMXValues")
	//span.SetTag("service.name", "dmx")
	defer span.End()
	//span.SetTag("operations", ops)
	s.m.Lock()
	defer s.m.Unlock()
	for _, op := range ops {
		channel := op.channel
		universe := op.universe
		value := op.value
		if channel < 1 || channel > 255 {
			return fmt.Errorf("dmx channel (%d) not in range, op=%v", channel, op)
		}

		s.initializeUniverse(universe)
		s.universes[universe][channel-1] = byte(value)
	}

	return nil
}

func (s *DMXState) initializeUniverse(universe int) {
	if s.universes[universe] == nil {
		chans := make([]byte, 255)
		s.universes[universe] = chans
	}
}

//GetDMXState returns the current dmx state
func (m *StateManager) GetDMXState() *DMXState {
	return &m.dmxState
}

//SetDMXState updates the dmxstate
func (m *StateManager) SetDMXState(ctx context.Context, ops ...dmxOperation) error {
	return m.dmxState.set(ctx, ops...)
}

//OLAClient is the interface for communicating with ola
type OLAClient interface {
	SendDmx(universe int, values []byte) (status bool, err error)
	Close()
}

//SendDMXWorker sends OLA the current dmxState across all universes
func SendDMXWorker(ctx context.Context, client OLAClient, tick time.Duration, manager Manager, wg *sync.WaitGroup) error {
	defer wg.Done()
	defer client.Close()

	t := time.NewTimer(tick)
	defer t.Stop()
	log.Printf("timer started at %v", time.Now())

	for {
		select {
		case <-ctx.Done():
			log.Println("SendDMXWorker shutdown")
			return ctx.Err()
		case <-t.C:
			for k, v := range manager.GetDMXState().universes {
				timer := prometheus.NewTimer(metrics.ExternalResponseTime.WithLabelValues("ola"))
				client.SendDmx(k, v)
				timer.ObserveDuration()
			}
			t.Reset(tick)
		}
	}
}
