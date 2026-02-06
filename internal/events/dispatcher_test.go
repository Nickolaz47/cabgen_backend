package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDispatcher(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	reg := NewRegistry()
	interval := 5 * time.Minute
	nWorkers := 3

	result := NewDispatcher(
		eventRepo, reg, interval, nWorkers,
	)

	assert.NotNil(t, result)
}

type MockEventRepo struct {
	mock.Mock
}

func (m *MockEventRepo) GetEvents(ctx context.Context, limit int) ([]models.Event, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.Event), args.Error(1)
}

func (m *MockEventRepo) CreateEvent(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepo) MarkFailed(ctx context.Context, eventID uint, err string) error {
	args := m.Called(ctx, eventID, err)
	return args.Error(0)
}

func (m *MockEventRepo) MarkDone(ctx context.Context, eventID uint) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockEventRepo) MarkProcessing(ctx context.Context, eventID uint) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

type MockRegistry struct {
	mock.Mock
}

func (m *MockRegistry) Register(name string, handler HandlerFunc) error {
	args := m.Called(name, handler)
	return args.Error(0)
}

func (m *MockRegistry) Get(name string) (HandlerFunc, bool) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(HandlerFunc), args.Bool(1)
}

func TestPoll(t *testing.T) {
	mockRepo := new(MockEventRepo)
	jobCh := make(chan models.Event, 5)

	d := &dispatcher{
		eventRepo: mockRepo,
		interval:  time.Millisecond * 10,
		jobCh:     jobCh,
	}

	eventsData := []models.Event{{ID: 1}, {ID: 2}}

	mockRepo.On("GetEvents", mock.Anything, 10).Return(eventsData, nil)

	ctx, cancel := context.WithCancel(context.Background())
	go d.poll(ctx)

	select {
	case e := <-jobCh:
		assert.Equal(t, uint(1), e.ID)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting event 1")
	}

	select {
	case e := <-jobCh:
		assert.Equal(t, uint(2), e.ID)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting event 2")
	}

	cancel()
	mockRepo.AssertExpectations(t)
}

func TestWorker(t *testing.T) {
	mockReg := new(MockRegistry)

	handlerMock := HandlerFunc(
		func(ctx context.Context, payload []byte) error {
			return nil
		})
	handlerPanic := HandlerFunc(
		func(ctx context.Context, payload []byte) error {
			panic("boom")
		})

	t.Run("Success", func(t *testing.T) {
		jobCh := make(chan models.Event, 1)
		resultCh := make(chan result, 1)

		d := &dispatcher{
			registry: mockReg,
			jobCh:    jobCh,
			resultCh: resultCh,
		}

		evt := models.Event{ID: 10, Name: "test.event"}
		mockReg.On("Get", "test.event").Return(handlerMock, true)

		jobCh <- evt
		close(jobCh)

		d.worker(context.Background())

		res := <-resultCh
		assert.Equal(t, evt.ID, res.eventID)
		assert.NoError(t, res.err)
	})

	t.Run("Error", func(t *testing.T) {
		jobCh := make(chan models.Event, 1)
		resultCh := make(chan result, 1)

		d := &dispatcher{
			registry: mockReg,
			jobCh:    jobCh,
			resultCh: resultCh,
		}

		evt := models.Event{ID: 666, Name: "test.panic"}

		mockReg.On("Get", "test.panic").Return(handlerPanic, true)

		jobCh <- evt
		close(jobCh)

		d.worker(context.Background())

		res := <-resultCh
		assert.Equal(t, evt.ID, res.eventID)
		assert.Error(t, res.err)
		assert.Contains(t, res.err.Error(), "panic in handler")
	})
}

func TestListenResults(t *testing.T) {
	mockRepo := new(MockEventRepo)

	t.Run("MarkDone", func(t *testing.T) {
		resultCh := make(chan result, 2)

		d := &dispatcher{
			eventRepo: mockRepo,
			resultCh:  resultCh,
		}

		resOk := result{eventID: 100, err: nil}
		mockRepo.On("MarkDone", mock.Anything, uint(100)).Return(nil)

		resultCh <- resOk
		close(resultCh)

		d.listenResults(context.Background())
		mockRepo.AssertExpectations(t)
	})

	t.Run("MarkFailed", func(t *testing.T) {
		resultCh := make(chan result, 2)

		d := &dispatcher{
			eventRepo: mockRepo,
			resultCh:  resultCh,
		}

		resFail := result{eventID: 100, err: errors.New("fail")}
		mockRepo.On("MarkFailed", mock.Anything, uint(100), "fail").Return(nil)

		resultCh <- resFail
		close(resultCh)

		d.listenResults(context.Background())
		mockRepo.AssertExpectations(t)
	})
}

func TestRun(t *testing.T) {
	mockRepo := new(MockEventRepo)
	mockReg := new(MockRegistry)

	evt := models.Event{ID: 1, Name: "test.event", Payload: []byte("{}")}
	handlerCalled := make(chan bool, 1)

	handler := HandlerFunc(func(ctx context.Context, p []byte) error {
		handlerCalled <- true
		return nil
	})

	mockRepo.On("GetEvents", mock.Anything, 10).Return([]models.Event{evt}, nil).Once()
	mockRepo.On("GetEvents", mock.Anything, 10).Return([]models.Event{}, nil)

	mockReg.On("Get", "test.event").Return(handler, true)

	dbFinished := make(chan bool, 1)
	mockRepo.On("MarkDone", mock.Anything, evt.ID).Run(func(args mock.Arguments) {
		dbFinished <- true
	}).Return(nil)

	d := NewDispatcher(mockRepo, mockReg, 10*time.Millisecond, 2)

	ctx, cancel := context.WithCancel(context.Background())

	runFinished := make(chan bool)
	go func() {
		d.Run(ctx)
		close(runFinished)
	}()

	select {
	case <-handlerCalled:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout: Handler never called")
	}

	select {
	case <-dbFinished:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout: Status never updated on database")
	}

	cancel()

	select {
	case <-runFinished:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: O Dispatcher não desligou corretamente (possível deadlock)")
	}

	mockRepo.AssertExpectations(t)
	mockReg.AssertExpectations(t)
}
