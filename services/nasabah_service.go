package services

import (
	"BlackListWorker/models"
	"database/sql"
)

func LoadCIF(db *sql.DB) ([]models.CIF, error) {
	rows, err := db.Query(`
		SELECT cif_number, nama_nasabah, tempat_lahir, tanggal_lahir, ktp, npwp, no_paspor 
		FROM CIF
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.CIF
	for rows.Next() {
		var c models.CIF
		err := rows.Scan(&c.CIFNumber, &c.NamaNasabah, &c.TempatLahir, &c.TanggalLahir, &c.KTP, &c.NPWP, &c.NoPaspor)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}
