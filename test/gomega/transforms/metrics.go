package transforms

import (
	errors "github.com/rotisserie/eris"
	"go.opencensus.io/stats/view"
)

func WithLastValueTransform() func(rows []*view.Row) (int, error) {
	return func(rows []*view.Row) (int, error) {
		if len(rows) == 0 {
			// In the future, it would be ideal to distinguish between "no rows" and "a row with a value of 0"
			// For now, we return 0 to represent both cases
			return 0, nil
		}

		if len(rows) > 1 {
			return 0, errors.Errorf("expected 1 row, found %d", len(rows))
		}

		return int(rows[0].Data.(*view.LastValueData).Value), nil
	}
}

func WithSumValueTransform() func(rows []*view.Row) (int, error) {
	return func(rows []*view.Row) (int, error) {
		if len(rows) == 0 {
			// In the future, it would be ideal to distinguish between "no rows" and "a row with a value of 0"
			// For now, we return 0 to represent both cases
			return 0, nil
		}

		if len(rows) > 1 {
			return 0, errors.Errorf("expected 1 row, found %d", len(rows))
		}

		return int(rows[0].Data.(*view.SumData).Value), nil
	}
}
