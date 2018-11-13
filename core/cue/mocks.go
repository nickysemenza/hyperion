package cue

import (
	"context"
	"sync"

	"github.com/stretchr/testify/mock"
)

//MockMaster is a mocked Master
type MockMaster struct {
	mock.Mock
}

//AddIDsRecursively is a mock implementation.
func (m *MockMaster) AddIDsRecursively(c *Cue) {
	m.Called()
}

//ProcessStack is a mock implementation.
func (m *MockMaster) ProcessStack(ctx context.Context, cs *Stack) {
	m.Called(ctx, cs)
}

//ProcessCue is a mock implementation.
func (m *MockMaster) ProcessCue(ctx context.Context, c *Cue, wg *sync.WaitGroup) {
	m.Called(ctx, c, wg)
}

//ProcessFrame is a mock implementation.
func (m *MockMaster) ProcessFrame(ctx context.Context, cf *Frame, wg *sync.WaitGroup) {
	m.Called(ctx, cf, wg)
}

//ProcessFrameAction is a mock implementation.
func (m *MockMaster) ProcessFrameAction(ctx context.Context, cfa *FrameAction, wg *sync.WaitGroup) {
	m.Called(ctx, cfa, wg)
}

//EnQueueCue is a mock implementation.
func (m *MockMaster) EnQueueCue(c Cue, cs *Stack) *Cue {
	args := m.Called(c, cs)
	return args.Get(0).(*Cue)
}

//GetDefaultCueStack is a mock implementation.
func (m *MockMaster) GetDefaultCueStack() *Stack {
	args := m.Called()
	return args.Get(0).(*Stack)
}
