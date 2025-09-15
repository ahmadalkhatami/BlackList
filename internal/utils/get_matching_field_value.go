package utils

import (
	"BlackListWorker/internal/domain/models"
	"strings"
)

func GetCIFValueByField(cif models.MasterNasabah, field string) string {
	switch field {
	case "nama", "namanasabah":
		return cif.NamaNasabah
	case "tempatlahir":
		return *cif.TempatLahir
	case "tanggallahir":
		if !cif.TanggalLahir.IsZero() {
			return cif.TanggalLahir.Format("2006-01-02")
		}
		return ""
	case "ktp":
		return *cif.KTP
	case "npwp":
		return *cif.NPWP
	case "nopaspor":
		return *cif.NoPaspor
	default:
		return ""
	}
}

func GetWatchlistValuesByField(wl models.MasterWatchlist, field string) []string {
	switch field {
	case "nama":
		values := []string{}
		if strings.TrimSpace(wl.Nama) != "" {
			values = append(values, wl.Nama)
		}
		for _, alias := range wl.Aliases {
			if strings.TrimSpace(alias) != "" {
				values = append(values, alias)
			}
		}
		return values
	case "tempatlahir":
		if wl.TempatLahir != nil {
			return []string{*wl.TempatLahir}
		}
	case "tanggallahir":
		if !wl.TanggalLahir.IsZero() {
			return []string{wl.TanggalLahir.Format("2006-01-02")}
		}
	case "ktp":
		if wl.KTP != nil {
			return []string{*wl.KTP}
		}
	case "npwp":
		if wl.NPWP != nil {
			return []string{*wl.NPWP}
		}
	case "nopaspor":
		if wl.NoPaspor != nil {
			return []string{*wl.NoPaspor}
		}
	}
	return []string{}
}
