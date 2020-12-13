package covid

import (
	"context"
	"covid/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type Peta_mapper struct {
	Kd_gedung_peta   string
	Nama_gedung_peta string
	Bentuk           string
	Koordinat        string
	Prefill_color    string
	Fill_color       string
	Data_pasien      []peta_pasien
}

type peta_pasien struct {
	tgl_batch      time.Time
	status         string
	odp            string
	pdp            int
	terkontaminasi int
	sembuh         int
	meninggal      int
}

var connect *pgx.Conn

func init() {
	var section string
	section = "database_covid"

	connect = config.Connect(section)

}

func data_gedung_peta(hospital_id int, data_pasien []interface{}) []Peta_mapper {

	rows, _ := connect.Query(context.Background(),
		`select gedung_peta.kd_gedung_peta,
				gedung_peta.nama_gedung_peta,
				gedung_mapper.bentuk,
				gedung_mapper.koordinat,
				gedung_mapper.prefill_color,
				gedung_mapper.fill_color
		   from gedung_peta, gedung_mapper
		  where gedung_peta.kd_gedung_peta = gedung_mapper.kd_gedung_peta
			and gedung_peta.status_aktif = '1'
			and gedung_mapper.status_aktif = '1'
			and gedung_peta.hospital_id = $1`, hospital_id)

	defer rows.Close()

	data_peta := []Peta_mapper{}
	for rows.Next() {
		var data_row Peta_mapper

		err := rows.Scan(&data_row.Kd_gedung_peta, &data_row.Nama_gedung_peta, &data_row.Bentuk,
			&data_row.Koordinat, &data_row.Prefill_color, &data_row.Fill_color)
		if err != nil {
			fmt.Printf("Query failed : %v", err)
		}

		// data_row.Data_pasien = filter(data_pasien, func() {})

		data_peta = append(data_peta, data_row)

	}

	return data_peta
}

// func data_pasien_peta(kd_gedung_peta string) []Peta_pasien {

// }
