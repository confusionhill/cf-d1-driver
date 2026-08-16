package cloudflare

import "context"

// MockClient is a test-friendly mock implementation of the Client interface.
type MockClient struct {
    QueryFunc func(ctx context.Context, req QueryRequest) error
    BatchFunc func(ctx context.Context, req BatchRequest) error
}

func (m *MockClient) Query(ctx context.Context, req QueryRequest) error {
    if m == nil || m.QueryFunc == nil {
        return nil
    }
    return m.QueryFunc(ctx, req)
}

func (m *MockClient) Batch(ctx context.Context, req BatchRequest) error {
    if m == nil || m.BatchFunc == nil {
        return nil
    }
    return m.BatchFunc(ctx, req)
}
