package mocks

import "context"

type MockEventEmitter struct {
	EmitFunc func(ctx context.Context, name string, payload any) error
}

func (m *MockEventEmitter) Emit(ctx context.Context, name string, payload any) error {
	if m.EmitFunc != nil {
		return m.EmitFunc(ctx, name, payload)
	}
	return nil
}
