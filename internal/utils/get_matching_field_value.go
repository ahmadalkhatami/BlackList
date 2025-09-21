package utils

import (
	"BlackListWorker/internal/domain/models"
	"strings"
)

func GetCIFValueByField(cif models.MasterNasabah, field string) string {
	switch strings.ToLower(field) {
	case "nama", "namanasabah":
		return cif.NamaNasabah
	case "tempatlahir":
		if cif.TempatLahir != nil {
			return *cif.TempatLahir
		}
		return ""
	case "tanggallahir":
		if !cif.TanggalLahir.IsZero() {
			return cif.TanggalLahir.Format("2006-01-02")
		}
		return ""
	case "ktp":
		if cif.KTP != nil {
			return *cif.KTP
		}
		return ""
	case "npwp":
		if cif.NPWP != nil {
			return *cif.NPWP
		}
		return ""
	case "nopaspor":
		if cif.NoPaspor != nil {
			return *cif.NoPaspor
		}
		return ""
	default:
		return ""
	}
}

func GetWatchlistValuesByField(wl models.MasterWatchlist, field string) []string {
	field = strings.ToLower(field)

	if field == "nama" {
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
	}

	if strings.HasPrefix(field, "alias") {
		values := []string{}
		for _, alias := range wl.Aliases {
			if strings.TrimSpace(alias) != "" {
				values = append(values, alias)
			}
		}
		return values
	}

	switch field {
	case "tempatlahir":
		if wl.TempatLahir != nil {
			return []string{*wl.TempatLahir}
		}
	case "tanggallahir":
		if wl.TanggalLahir != nil && !wl.TanggalLahir.IsZero() {
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
