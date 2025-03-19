package mocks

import (
    "github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the Client interface
type MockClient struct {
    mock.Mock
}

// Embed implements the Client interface
func (m *MockClient) Embed(key string) ([]float32, error) {
    args := m.Called(key)
    
    // Handle return values
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    
    return args.Get(0).([]float32), args.Error(1)
}