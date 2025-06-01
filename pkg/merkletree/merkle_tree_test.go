package merkletree_test

import (
	"testing"

	"github.com/jeanmolossi/ai-agent-cli/pkg/merkletree"
	"github.com/stretchr/testify/assert"
)

func TestBuildMerkleRoot(t *testing.T) {
	content1 := merkletree.ComputeChunkHash("hello")
	content2 := merkletree.ComputeChunkHash("world")
	content3 := merkletree.ComputeChunkHash("foo")
	content4 := merkletree.ComputeChunkHash("bar")
	content5 := merkletree.ComputeChunkHash("baz")

	hashes := []string{
		content1,
		content2,
		content3,
		content4,
		content5,
	}

	merkleroot := merkletree.BuildMerkleRoot(hashes)

	wantedHash := "3e1a159911b0753eea72fa67e12398f51c9910abb2d00f54ffd62df280591bbe"

	assert.Equal(t, wantedHash, merkleroot, "expects the root hash match")

	hashes[2] = merkletree.ComputeChunkHash("foo1") // modify a hash from list

	merkleroot = merkletree.BuildMerkleRoot(hashes)

	assert.NotEqual(t, wantedHash, merkleroot, "the hash should be changed")
}
