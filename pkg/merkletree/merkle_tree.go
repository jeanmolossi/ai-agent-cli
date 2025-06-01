package merkletree

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeChunkHash calculate the SHA-256 hash from a content and returns the hex.
func ComputeChunkHash(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])
}

// BuildMerkleRoot receives a slice of hashes, ordered deterministic then
// returns the Merkle Root as a hex string. If hashes is empty, back with
// empty string
func BuildMerkleRoot(hashes []string) string {
	n := len(hashes)
	if n == 0 {
		return ""
	}

	// If has single hash it's the root
	if n == 1 {
		return hashes[0]
	}

	var nextLevel []string
	// build next level of tree
	for i := 0; i < n; i += 2 {
		var combined string

		// when have a pair of hashes
		if i+1 < n {
			combined = hashes[i] + hashes[i+1]
		} else {
			// when is odd number we duplicate the last one
			combined = hashes[i] + hashes[i]
		}

		sum := sha256.Sum256([]byte(combined))
		nextLevel = append(nextLevel, hex.EncodeToString(sum[:]))
	}

	return BuildMerkleRoot(nextLevel)
}
