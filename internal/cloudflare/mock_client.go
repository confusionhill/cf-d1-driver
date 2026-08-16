package cloudflare

import "context"

// MockClient is a test-friendly mock implementation of the Client interface.
type MockClient struct {
	QueryFunc func(ctx context.Context, req QueryRequest) ([]QueryResult, error)
	BatchFunc func(ctx context.Context, req BatchRequest) ([]QueryResult, error)
}

func (m *MockClient) Query(ctx context.Context, req QueryRequest) ([]QueryResult, error) {
	if m == nil || m.QueryFunc == nil {
		return nil, nil
	}
	return m.QueryFunc(ctx, req)
}

func (m *MockClient) Batch(ctx context.Context, req BatchRequest) ([]QueryResult, error) {
	if m == nil || m.BatchFunc == nil {
		return nil, nil
	}
	return m.BatchFunc(ctx, req)
}
