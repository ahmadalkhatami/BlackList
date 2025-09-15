package utils

import (
	"fmt"
	"strings"
)

type FilterBuilder struct {
	conds []string
	args  []any
}

func (f *FilterBuilder) Add(field string, value any) {
	if value == nil {
		return
	}
	param := fmt.Sprintf("@p%d", len(f.args)+1)
	f.conds = append(f.conds, fmt.Sprintf("[%s] = %s", field, param))
	f.args = append(f.args, value)
}

func (f *FilterBuilder) Where(condition string, value any) {
	if value == nil {
		return
	}
	param := fmt.Sprintf("@p%d", len(f.args)+1)
	cond := strings.Replace(condition, "?", param, 1)
	f.conds = append(f.conds, cond)
	f.args = append(f.args, value)
}

func (f *FilterBuilder) Build(baseQuery string) (string, []any) {
	if len(f.conds) > 0 {
		baseQuery += " WHERE " + strings.Join(f.conds, " AND ")
	}
	return baseQuery, f.args
}

func NewQueryBuilder() *FilterBuilder {
	return &FilterBuilder{}
}

/*
-- Example With Filter --

fb := utils.NewFilterBuilder()

fb.Add("IsActive", filter.IsActive)
fb.Add("IsIndividual", filter.IsIndividual)
fb.Add("WatchlistSource", filter.WatchlistSource)

query := `
    SELECT [Id],[Nama],[Alias1],[Alias2],[Alias3],[Alias4],
           [Type],[TempatLahir],[TanggalLahir],[KTP],[NPWP],
           [NoPaspor],[CreatedAt],[IsActive]
    FROM [dbo].[MASTER_LOCAL_BLACKLIST]
`
query, args := fb.Build(query)

rows, err := r.DB.QueryContext(ctx, query, args...)
if err != nil {
    return nil, err
}
defer rows.Close()


-- With No Filter --
fb := utils.NewFilterBuilder()

// Tidak ada Add() dipanggil
query := `
    SELECT [Id],[Nama],[Alias1],[Alias2],[Alias3],[Alias4],
           [Type],[TempatLahir],[TanggalLahir],[KTP],[NPWP],
           [NoPaspor],[CreatedAt],[IsActive]
    FROM [dbo].[MASTER_LOCAL_BLACKLIST]
`
query, args := fb.Build(query)

fmt.Println(query) // --> SELECT ... FROM [dbo].[MASTER_LOCAL_BLACKLIST]
fmt.Println(args)  // --> []

*/
