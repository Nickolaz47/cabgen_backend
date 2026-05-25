package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func (m *MockEventRepo) MarkFailed(ctx context.Context, eventID uuid.UUID, err string) error {
	args := m.Called(ctx, eventID, err)
	return args.Error(0)
}

func (m *MockEventRepo) MarkDone(ctx context.Context, eventID uuid.UUID) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *MockEventRepo) MarkProcessing(ctx context.Context, eventID uuid.UUID) error {
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
	fn, _ := args.Get(0).(HandlerFunc)
	return fn, args.Bool(1)
}

func TestNewDispatcher(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	reg := NewRegistry()
	interval := 5 * time.Minute
	nWorkers := 3

	result := NewDispatcher(eventRepo, reg, interval, nWorkers)

	assert.NotNil(t, result)
}

func TestPoll(t *testing.T) {
	mockRepo := new(MockEventRepo)
	jobCh := make(chan models.Event, 5)

	d := &dispatcher{
		eventRepo: mockRepo,
		interval:  10 * time.Millisecond,
		jobCh:     jobCh,
	}

	id1 := uuid.New()
	id2 := uuid.New()
	eventsData := []models.Event{
		{ID: id1},
		{ID: id2},
	}

	mockRepo.On("GetEvents", mock.Anything, 10).Return(eventsData, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go d.poll(ctx)

	recv := func(label string) models.Event {
		t.Helper()
		select {
		case e := <-jobCh:
			return e
		case <-time.After(time.Second):
			t.Fatalf("Timeout waiting for %s", label)
			return models.Event{}
		}
	}

	e1 := recv("event 1")
	assert.Equal(t, id1, e1.ID)

	e2 := recv("event 2")
	assert.Equal(t, id2, e2.ID)

	mockRepo.AssertExpectations(t)
}

func TestWorker(t *testing.T) {
	handlerOK := HandlerFunc(func(ctx context.Context, payload []byte) error {
		return nil
	})
	handlerErr := HandlerFunc(func(ctx context.Context, payload []byte) error {
		return errors.New("handler error")
	})
	handlerPanic := HandlerFunc(func(ctx context.Context, payload []byte) error {
		panic("boom")
	})

	run := func(t *testing.T, evt models.Event, reg *MockRegistry) result {
		t.Helper()
		jobCh := make(chan models.Event, 1)
		resultCh := make(chan result, 1)

		d := &dispatcher{
			registry: reg,
			jobCh:    jobCh,
			resultCh: resultCh,
		}

		jobCh <- evt
		close(jobCh)
		d.worker(context.Background())

		select {
		case res := <-resultCh:
			return res
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for worker result")
			return result{}
		}
	}

	t.Run("Success", func(t *testing.T) {
		mockReg := new(MockRegistry)
		evt := models.Event{ID: uuid.New(), Name: "test.ok"}
		mockReg.On("Get", "test.ok").Return(handlerOK, true)

		res := run(t, evt, mockReg)

		assert.Equal(t, evt.ID, res.eventID)
		assert.NoError(t, res.err)
		mockReg.AssertExpectations(t)
	})

	t.Run("HandlerError", func(t *testing.T) {
		mockReg := new(MockRegistry)
		evt := models.Event{ID: uuid.New(), Name: "test.err"}
		mockReg.On("Get", "test.err").Return(handlerErr, true)

		res := run(t, evt, mockReg)

		assert.Equal(t, evt.ID, res.eventID)
		assert.EqualError(t, res.err, "handler error")
		mockReg.AssertExpectations(t)
	})

	t.Run("HandlerPanic", func(t *testing.T) {
		mockReg := new(MockRegistry)
		evt := models.Event{ID: uuid.New(), Name: "test.panic"}
		mockReg.On("Get", "test.panic").Return(handlerPanic, true)

		res := run(t, evt, mockReg)

		assert.Equal(t, evt.ID, res.eventID)
		assert.Error(t, res.err)
		assert.Contains(t, res.err.Error(), "panic in handler")
		mockReg.AssertExpectations(t)
	})

	t.Run("HandlerNotFound", func(t *testing.T) {
		mockReg := new(MockRegistry)
		evt := models.Event{ID: uuid.New(), Name: "test.missing"}
		mockReg.On("Get", "test.missing").Return(nil, false)

		res := run(t, evt, mockReg)

		assert.Equal(t, evt.ID, res.eventID)
		assert.EqualError(t, res.err, "handler not found")
		mockReg.AssertExpectations(t)
	})
}

func TestListenResults(t *testing.T) {
	t.Run("MarkDone", func(t *testing.T) {
		mockRepo := new(MockEventRepo)
		resultCh := make(chan result, 1)

		d := &dispatcher{
			eventRepo: mockRepo,
			resultCh:  resultCh,
		}

		id := uuid.New()
		mockRepo.On("MarkDone", mock.Anything, id).Return(nil)

		resultCh <- result{eventID: id, err: nil}
		close(resultCh)

		d.listenResults(context.Background())
		mockRepo.AssertExpectations(t)
	})

	t.Run("MarkFailed", func(t *testing.T) {
		mockRepo := new(MockEventRepo)
		resultCh := make(chan result, 1)

		d := &dispatcher{
			eventRepo: mockRepo,
			resultCh:  resultCh,
		}

		id := uuid.New()
		mockRepo.On("MarkFailed", mock.Anything, id, "fail").Return(nil)

		resultCh <- result{eventID: id, err: errors.New("fail")}
		close(resultCh)

		d.listenResults(context.Background())
		mockRepo.AssertExpectations(t)
	})
}

func TestRun(t *testing.T) {
	mockRepo := new(MockEventRepo)
	mockReg := new(MockRegistry)

	evt := models.Event{
		ID:      uuid.New(),
		Name:    "test.event",
		Payload: []byte("{}"),
		Status:  models.EventPending,
	}

	handlerCalled := make(chan struct{}, 1)
	handler := HandlerFunc(func(ctx context.Context, p []byte) error {
		handlerCalled <- struct{}{}
		return nil
	})

	mockRepo.On("GetEvents", mock.Anything, 10).
		Return([]models.Event{evt}, nil).Once()
	mockRepo.On("GetEvents", mock.Anything, 10).
		Return([]models.Event{}, nil)

	mockReg.On("Get", "test.event").Return(handler, true)

	dbFinished := make(chan struct{}, 1)
	mockRepo.On("MarkDone", mock.Anything, evt.ID).
		Run(func(args mock.Arguments) { dbFinished <- struct{}{} }).
		Return(nil)

	d := NewDispatcher(mockRepo, mockReg, 10*time.Millisecond, 2)

	ctx, cancel := context.WithCancel(context.Background())

	runFinished := make(chan struct{})
	go func() {
		d.Run(ctx)
		close(runFinished)
	}()

	select {
	case <-handlerCalled:
	case <-time.After(time.Second):
		t.Fatal("Timeout: handler never called")
	}

	select {
	case <-dbFinished:
	case <-time.After(time.Second):
		t.Fatal("Timeout: status never updated in database")
	}

	cancel()

	select {
	case <-runFinished:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: dispatcher did not shut down correctly (possible deadlock)")
	}

	mockRepo.AssertExpectations(t)
	mockReg.AssertExpectations(t)
}
