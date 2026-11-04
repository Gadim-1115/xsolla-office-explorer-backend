package shuffler

import (
	"crypto/rand"
	"math/big"

	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

func ShuffleCodeBlocks(blocks []entity.CodeBlock, correctOrder []string) []entity.CodeBlock {
	n := len(blocks)
	shuffled := make([]entity.CodeBlock, n)
	copy(shuffled, blocks)

	if n <= 1 {
		return shuffled
	}

	for i := n - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			j := i - 1
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
			continue
		}
		j := int(jBig.Int64())
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	if matchesOrder(shuffled, correctOrder) {
		shuffled[0], shuffled[1] = shuffled[1], shuffled[0]
	}

	return shuffled
}

func matchesOrder(blocks []entity.CodeBlock, correctOrder []string) bool {
	if len(blocks) != len(correctOrder) {
		return false
	}
	for i := range blocks {
		if blocks[i].ID != correctOrder[i] {
			return false
		}
	}
	return true
}
