package entity

// CodeBlock represents a single draggable line or block of code in a quest.
type CodeBlock struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

// Quest represents an algorithmic Drag & Drop challenge bound 1:1 to a locked zone.
type Quest struct {
	ID               string
	ZoneID           string
	Title            string
	ProblemStatement string
	CodeBlocks       []CodeBlock
	CorrectOrder     []string
}

// VerifySequence checks if the submitted slice of block IDs strictly matches CorrectOrder.
func (q *Quest) VerifySequence(submittedOrder []string) bool {
	if len(q.CorrectOrder) == 0 || len(submittedOrder) != len(q.CorrectOrder) {
		return false
	}
	for i := range q.CorrectOrder {
		if submittedOrder[i] != q.CorrectOrder[i] {
			return false
		}
	}
	return true
}
