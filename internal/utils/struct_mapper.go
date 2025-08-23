package utils

func MapSlice[S any, D any](src []S, mapper func(S) D) []D {
	dst := make([]D, 0, len(src))
	for _, s := range src {
		dst = append(dst, mapper(s))
	}
	return dst
}

/**
Example Use:
dtoList := utils.MapSlice(dataA, func(m MatchingResult) MatchingResultDTO {
    return MatchingResultDTO{
        ID: m.Id,
        // ...
    }
})
*/
