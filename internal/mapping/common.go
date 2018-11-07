package mapping

import (
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
)

func OrderBy(in []*charonrpc.Order) []model.RowOrder {
	out := make([]model.RowOrder, 0, len(in))
	for _, i := range in {
		out = append(out, model.RowOrder{
			Name:       i.GetName(),
			Descending: i.GetDescending(),
		})
	}
	return out
}
