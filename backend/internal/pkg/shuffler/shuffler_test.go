package shuffler

import (
	"testing"

	"github.com/xsolla-school/office-explorer/backend/internal/domain/entity"
)

func TestShuffleCodeBlocks_PreservesElementsAndAvoidsSolutionOrder(t *testing.T) {
	original := []entity.CodeBlock{
		{ID: "b1", Code: "line 1"},
		{ID: "b2", Code: "line 2"},
		{ID: "b3", Code: "line 3"},
		{ID: "b4", Code: "line 4"},
	}
	correctOrder := []string{"b1", "b2", "b3", "b4"}

	for i := 0; i < 50; i++ {
		shuffled := ShuffleCodeBlocks(original, correctOrder)

		if len(shuffled) != len(original) {
			t.Fatalf("expected length %d, got %d", len(original), len(shuffled))
		}

		if original[0].ID != "b1" || original[1].ID != "b2" {
			t.Fatal("original slice was mutated")
		}

		sameOrder := true
		for idx := range correctOrder {
			if shuffled[idx].ID != correctOrder[idx] {
				sameOrder = false
				break
			}
		}

		if sameOrder {
			t.Fatal("shuffled blocks matched correctOrder")
		}
	}
}

func TestShuffleCodeBlocks_SingleElement(t *testing.T) {
	single := []entity.CodeBlock{{ID: "b1", Code: "only line"}}
	res := ShuffleCodeBlocks(single, []string{"b1"})

	if len(res) != 1 || res[0].ID != "b1" {
		t.Fatalf("unexpected result for single element: %+v", res)
	}
}
