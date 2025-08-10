package services

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

func LoadCIF(db *sql.DB) ([]models.MasterNasabah, error) {
	rows, err := db.Query(`
		SELECT cif_number, nama_nasabah, tempat_lahir, tanggal_lahir, ktp, npwp, no_paspor 
		FROM MASTER_NASABAH
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.MasterNasabah
	for rows.Next() {
		var c models.MasterNasabah
		err := rows.Scan(&c.CIFNumber, &c.NamaNasabah, &c.TempatLahir, &c.TanggalLahir, &c.KTP, &c.NPWP, &c.NoPaspor)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}
