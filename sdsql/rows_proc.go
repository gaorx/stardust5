package sdsql

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
)

type RowsProc[ROW any] interface {
	ProcRows(rows []ROW) ([]ROW, error)
}

var (
	_ RowsProc[struct{}] = (RowsProcFunc[struct{}])(nil)
	_ RowsProc[struct{}] = (Completer[struct{}])(nil)
	_ RowsProc[struct{}] = (InplaceCompleter[struct{}])(nil)
	_ RowsProc[struct{}] = (Filter[struct{}])(nil)
	_ RowsProc[struct{}] = Aggregator[struct{}, string, struct{}]{}
)

func ProcRows[ROW any](rows []ROW, procs ...RowsProc[ROW]) ([]ROW, error) {
	if len(procs) <= 0 {
		return rows, nil
	}
	for _, proc := range procs {
		if proc == nil {
			continue
		}
		var err error
		rows, err = proc.ProcRows(rows)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func ProcRow[ROW any](row ROW, procs ...RowsProc[ROW]) (ROW, error) {
	if len(procs) <= 0 {
		return row, nil
	}
	rows, err := ProcRows([]ROW{row}, procs...)
	if err != nil {
		var zero ROW
		return zero, err
	}
	return rows[0], nil
}

type RowsProcFunc[ROW any] func(rows []ROW) ([]ROW, error)

func (f RowsProcFunc[ROW]) ProcRows(rows []ROW) ([]ROW, error) {
	if f == nil {
		return rows, nil
	}
	return f(rows)
}

type Completer[ROW any] func(ROW) (ROW, error)

func (c Completer[ROW]) ProcRows(rows []ROW) ([]ROW, error) {
	if c == nil {
		return rows, nil
	}
	if len(rows) <= 0 {
		return rows, nil
	}
	newRows := make([]ROW, len(rows))
	for _, row := range rows {
		newRow, err := c(row)
		if err != nil {
			return nil, err
		}
		newRows = append(newRows, newRow)
	}
	return newRows, nil
}

type InplaceCompleter[ROW any] func(ROW)

func (c InplaceCompleter[ROW]) ProcRows(rows []ROW) ([]ROW, error) {
	if c == nil {
		return rows, nil
	}
	if len(rows) <= 0 {
		return rows, nil
	}
	for _, row := range rows {
		c(row)
	}
	return rows, nil
}

type Filter[ROW any] func(ROW) bool

func (f Filter[ROW]) ProcRows(rows []ROW) ([]ROW, error) {
	if f == nil {
		return rows, nil
	}
	if len(rows) <= 0 {
		return rows, nil
	}
	newRows := make([]ROW, 0)
	for _, row := range rows {
		if f(row) {
			newRows = append(newRows, row)
		}
	}
	return newRows, nil
}

type Aggregator[ROW any, ON comparable, COMPLEMENT any] struct {
	disabled        bool
	Collect         func(ROW) []ON
	Fetch           func([]ON) (map[ON]COMPLEMENT, error)
	Complete        func(ROW, map[ON]COMPLEMENT) (ROW, error)
	CompleteInplace func(ROW, map[ON]COMPLEMENT)
}

func (a Aggregator[ROW, ON, COMPLEMENT]) If(enabled bool) Aggregator[ROW, ON, COMPLEMENT] {
	a1 := a
	a1.disabled = !enabled
	return a1
}

func (a Aggregator[ROW, ON, COMPLEMENT]) ProcRows(rows []ROW) ([]ROW, error) {
	if a.disabled {
		return rows, nil
	}

	if a.Collect == nil || a.Fetch == nil || (a.Complete == nil && a.CompleteInplace == nil) {
		return nil, sderr.New("invalid aggregator")
	}

	if len(rows) <= 0 {
		return rows, nil
	}
	collected := make(map[ON]struct{})
	for _, row := range rows {
		ons := a.Collect(row)
		if len(ons) <= 0 {
			// do nothing
		} else if len(ons) == 1 {
			collected[ons[0]] = struct{}{}
		} else {
			for _, on := range ons {
				collected[on] = struct{}{}
			}
		}
	}
	complements, err := a.Fetch(lo.Keys(collected))
	if err != nil {
		return nil, err
	}
	if a.Complete != nil {
		newRows := make([]ROW, len(rows))
		for _, row := range rows {
			newRow, err := a.Complete(row, complements)
			if err != nil {
				return nil, err
			}
			newRows = append(newRows, newRow)
		}
		return newRows, nil
	} else if a.CompleteInplace != nil {
		for _, row := range rows {
			a.CompleteInplace(row, complements)
		}
		return rows, nil
	} else {
		return nil, sderr.New("aggregator.complete is nil")
	}
}
