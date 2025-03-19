package tests

import (
	"testing"
	"vectorDb/client/mocks"

	"github.com/stretchr/testify/assert"
)

func TestEmbed(t *testing.T) {
    // Test cases
    testCases := []struct {
        name           string
        inputKey       string
        mockResponse   []float32
        mockError      error
        expectedResult []float32
        expectedError  error
    }{
        {
            name:           "successful embedding",
            inputKey:       "test phrase",
            mockResponse:   []float32{0.1, 0.2, 0.3},
            mockError:      nil,
            expectedResult: []float32{0.1, 0.2, 0.3},
            expectedError:  nil,
        },
        // Add more test cases as needed
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create mock client
            mockClient := new(mocks.MockClient)
            
            // Set expectations
            mockClient.On("Embed", tc.inputKey).Return(tc.mockResponse, tc.mockError)
            
            // Call the function
            result, err := mockClient.Embed(tc.inputKey)
            
            // Assert results
            assert.Equal(t, tc.expectedResult, result)
            assert.Equal(t, tc.expectedError, err)
            
            // Verify expectations were met
            mockClient.AssertExpectations(t)
        })
    }
}