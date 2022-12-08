package repositories

import (
	"context"
	"fmt"
)

func (s *SQLStore) DeleteSlideTx(ctx context.Context, id string) error {
	return s.ExecTx(ctx, func(q *Queries) error {
		var err error
		err = q.DeleteAnswersBySlide(ctx, id)
		if err != nil {
			return fmt.Errorf("delete answers: %w", err)
		}

		err = q.DeleteQuestionsBySlide(ctx, id)
		if err != nil {
			return fmt.Errorf("delete questions: %w", err)
		}

		err = q.DeleteSlide(ctx, id)
		if err != nil {
			return fmt.Errorf("delete slide: %w", err)
		}

		return nil
	})
}

func (s *SQLStore) DeleteQuestionTx(ctx context.Context, id string) error {
	return s.ExecTx(ctx, func(q *Queries) error {
		var err error
		err = q.DeleteAnswersByQuestion(ctx, id)
		if err != nil {
			return fmt.Errorf("delete answers: %w", err)
		}

		err = q.DeleteQuestion(ctx, id)
		if err != nil {
			return fmt.Errorf("delete question: %w", err)
		}

		return nil
	})
}
