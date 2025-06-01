package str_test

import (
	"testing"

	"github.com/jeanmolossi/ai-agent-cli/app/support/str"
	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	assert.Len(t, str.Random(10), 10)
	assert.Empty(t, str.Random(0))
	assert.Panics(t, func() {
		str.Random(-1)
	})
}
