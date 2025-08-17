package services

import (
	"BlackListWorker/models"
	"database/sql"
)

func LoadCIF(db *sql.DB) ([]models.MASTER_NASABAH, error) {
	rows, err := db.Query(`
		SELECT CifNumber, NamaNasabah, TempatLahir, TanggalLahir, Ktp, Npwp, NoPaspor 
		FROM MASTER_NASABAH
		WHERE StatusNasabah = 'ACTIVE'
		AND NamaNasabah = 'Prince Eddie Reinger'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.MASTER_NASABAH
	for rows.Next() {
		var c models.MASTER_NASABAH
		err := rows.Scan(&c.CIFNumber, &c.NamaNasabah, &c.TempatLahir, &c.TanggalLahir, &c.KTP, &c.NPWP, &c.NoPaspor)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}
