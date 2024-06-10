package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_htons(t *testing.T) {
	testcases := []struct {
		name   string
		input  uint16
		expect uint16
	}{
		{
			name:   "simple",
			input:  0x1234,
			expect: 0x3412,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := htons(tc.input)
			assert.Equal(t, actual, tc.expect)
		})
	}
}
