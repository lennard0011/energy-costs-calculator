package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsagePassSanityCheck(t *testing.T) {
	// Test cases: usage should pass sanity check
	validUsages := []float64{50, 75, 100}
	for _, usage := range validUsages {
		assert.True(t, usageIsWithinBounds(usage), "Expected usage %f to pass sanity check, but it failed")
	}

	// Test cases: usage should fail sanity check
	invalidUsages := []float64{-10, 110}
	for _, usage := range invalidUsages {
		assert.False(t, usageIsWithinBounds(usage), "Expected usage %f to fail sanity check, but it succeeded")
	}
}
