package textutil

import "database/sql"

type StringCombine interface {
	CombineAliases(a1, a2, a3, a4 sql.NullString) []string
	WMDAliases(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 sql.NullString) []string
}

func CombineAliases(a1, a2, a3, a4 sql.NullString) []string {
	aliases := []string{}
	for _, a := range []sql.NullString{a1, a2, a3, a4} {
		if a.Valid && a.String != "" {
			aliases = append(aliases, a.String)
		}
	}
	return aliases
}

func WMDAliases(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 sql.NullString) []string {
	aliases := []string{}
	for _, a := range []sql.NullString{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10} {
		if a.Valid && a.String != "" {
			aliases = append(aliases, a.String)
		}
	}
	return aliases
}
