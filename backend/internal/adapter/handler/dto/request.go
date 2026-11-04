package dto

import "github.com/xsolla-school/office-explorer/backend/internal/domain"

type VerifyQuestRequest struct {
	SubmittedOrder []string `json:"submitted_order"`
}

func (r *VerifyQuestRequest) Validate() error {
	if len(r.SubmittedOrder) == 0 {
		return domain.ErrInvalidInput
	}
	for _, id := range r.SubmittedOrder {
		if id == "" {
			return domain.ErrInvalidInput
		}
	}
	return nil
}
