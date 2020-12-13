package covid

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// type board struct {
// 	grafik_pasien graph_pasien
// }

type graph_pasien struct {
	Grafik_odp       []graph_status_kasus `json:"grafik_odp"`
	Grafik_pdp       []graph_status_kasus `json:"grafik_pdp"`
	Grafik_positif   []graph_status_kasus `json:"grafik_positif"`
	Grafik_sembuh    []graph_status_kasus `json:"grafik_sembuh"`
	Grafik_meninggal []graph_status_kasus `json:"grafik_meninggal"`
	Grafik_kasus     []graph_status_kasus `json:"grafik_kasus"`
	Grafik_kelamin   []graph_kelamin      `json:"grafik_kelamin"`
	Grafik_usia      []graph_usia         `json:"grafik_usia"`
	Gedung_peta      []Peta_mapper        `json:"gedung_peta"`
}

type graph_status_kasus struct {
	Tgl_data      time.Time `json:"time"`
	Jumlah_pasien int       `json:"jumlah_pasien"`
}

type graph_kelamin struct {
	Kelamin       string `json:"kelamin"`
	Jumlah_pasien int    `json:"jumlah_pasien"`
}

type graph_usia struct {
	Usia          string `json:"usia"`
	Jumlah_pasien int    `json:"jumlah_pasien"`
}

type Data_pasien struct {
	tgl_data           time.Time
	tgl_batch          time.Time
	nama_jenis_pasien  string
	jenis_kelamin      string
	kd_gedung_peta     string
	umur               int
	nama_status_pasien string
	jumlah_pasien      int
}

type data_rs struct {
	Kd_rs int `json:"kd_rs"   binding:"required"`
}

func Get_data_board(c *gin.Context) {
	var grafik_pasien graph_pasien
	// var pasien Data_pasien
	var data data_rs

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
		// return "", jwt.ErrMissingLoginValues
	}

	// session := sessions.Default(c)
	// kd_rs := session.Get("kd_rs").(int)
	kd_rs := data.Kd_rs
	fmt.Printf("data %+v", data)
	pasien := data_pasien_covid(kd_rs)

	grafik_pasien.Grafik_odp = status_kasus(pasien, "ODP")
	grafik_pasien.Grafik_pdp = status_kasus(pasien, "PDP")
	grafik_pasien.Grafik_positif = status_kasus(pasien, "Terkonfirmasi")
	grafik_pasien.Grafik_sembuh = status_kasus(pasien, "Sembuh")
	grafik_pasien.Grafik_meninggal = status_kasus(pasien, "Meninggal")
	grafik_pasien.Grafik_kasus = status_kasus(pasien, "all")
	grafik_pasien.Grafik_kelamin = graf_kelamin(pasien)
	grafik_pasien.Grafik_usia = graf_usia(pasien)
	// grafik_pasien.gedung_peta = []Peta_mapper{}
	// graf_kelamin(pasien)
	// fmt.Printf("cek data : %+v\n", grafik_pasien)

	// data_json, err := json.Marshal(grafik_odp)
	// fmt.Printf("cek data_json : %+v\n", data_json)

	// if err != nil {
	// 	fmt.Printf("error json : %v", err)
	// }

	c.JSON(200, gin.H{"message": grafik_pasien})
}

func data_pasien_covid(kd_rs int) []Data_pasien {

	rows, _ := connect.Query(context.Background(), `
		select pasien_covid_batch.tgl_data,
			   pasien_covid_batch.tgl_batch,
			   case pasien_covid_batch.kd_jenis_pasien
				   when '1' then 'Pasien non pegawai RS'
				   when '2' then 'Pasien pegawai RS'
				   when '3' then 'Pegawai RS tidak di rawat' end nama_jenis_pasien,
			   jenis_kelamin,
			   kd_gedung_peta,
			   umur,
			   case pasien_covid_batch.kd_status_pasien
				   when '1' then 'ODP'
				   when '2' then 'PDP'
				   when '3' then 'Terkonfirmasi'
				   when '4' then 'Sembuh'
				   when '5' then 'Meninggal' end nama_status_pasien,
			   pasien_covid_batch.jumlah_pasien
		  from pasien_covid_batch
		 where pasien_covid_batch.tgl_data >= cast(date_trunc('month', current_date) as date)
		   and pasien_covid_batch.hospital_id = $1
	  order by tgl_data`, kd_rs)

	defer rows.Close()

	data := []Data_pasien{}
	for rows.Next() {
		var pasien_row Data_pasien

		rows.Scan(&pasien_row.tgl_data, &pasien_row.tgl_batch, &pasien_row.nama_jenis_pasien,
			&pasien_row.jenis_kelamin, &pasien_row.kd_gedung_peta, &pasien_row.umur,
			&pasien_row.nama_status_pasien, &pasien_row.jumlah_pasien)

		data = append(data, pasien_row)
	}

	return data
}

func status_kasus(pasien []Data_pasien, status string) []graph_status_kasus {
	var old_tgl time.Time
	var jumlah int

	data_graph := []graph_status_kasus{}
	for i := range pasien {
		var data_row graph_status_kasus
		if pasien[i].nama_status_pasien == status {
			if (pasien[i].tgl_data).Equal(old_tgl) {
				// fmt.Printf("cek tgl_data : %v\n", pasien[i].tgl_data)
				// fmt.Printf("cek old_tgl : %v\n", old_tgl)
				jumlah = jumlah + pasien[i].jumlah_pasien
				old_tgl = pasien[i].tgl_data
			} else {
				if old_tgl.Year() > 1900 {
					data_row.Tgl_data = old_tgl
					data_row.Jumlah_pasien = jumlah
					data_graph = append(data_graph, data_row)
					// fmt.Printf("cek assign old_tgl tgl_data : %v\n", data_row.tgl_data)
					fmt.Printf("cek assign old_tgl : %v\n", data_row)
				}
				old_tgl = pasien[i].tgl_data
				jumlah = pasien[i].jumlah_pasien
			}
		} else if status == "all" {
			if (pasien[i].tgl_data).Equal(old_tgl) {
				data_row.Jumlah_pasien = data_row.Jumlah_pasien + pasien[i].jumlah_pasien
			} else {
				if old_tgl.Year() > 1900 {
					data_row.Tgl_data = old_tgl
					data_row.Jumlah_pasien = jumlah
					data_graph = append(data_graph, data_row)
				}
				old_tgl = pasien[i].tgl_data
				jumlah = pasien[i].jumlah_pasien
			}

		}

		if (len(pasien) - 1) == i {
			data_row.Tgl_data = old_tgl
			data_row.Jumlah_pasien = jumlah
			data_graph = append(data_graph, data_row)
		}
	}

	return data_graph
}

func graf_kelamin(pasien []Data_pasien) []graph_kelamin {

	var jumlah_laki, jumlah_perempuan int = 0, 0
	var data_graph []graph_kelamin
	fmt.Printf("graf_kelamin %v", len(pasien))
	tgl_data := pasien[len(pasien)-1].tgl_data

	for i := range pasien {
		if (pasien[i].tgl_data).Equal(tgl_data) {
			if pasien[i].jenis_kelamin == "L" {
				jumlah_laki = jumlah_laki + pasien[i].jumlah_pasien
			} else {
				jumlah_perempuan = jumlah_perempuan + pasien[i].jumlah_pasien
			}
		}
	}

	data_graph = []graph_kelamin{
		graph_kelamin{
			Kelamin:       "L",
			Jumlah_pasien: jumlah_laki,
		},
		graph_kelamin{
			Kelamin:       "P",
			Jumlah_pasien: jumlah_perempuan,
		},
	}

	return data_graph
}

func graf_usia(pasien []Data_pasien) []graph_usia {

	var jumlah_025, jumlah_2550, jumlah_51 int = 0, 0, 0
	var data_graph []graph_usia

	tgl_data := pasien[len(pasien)-1].tgl_data

	for i := range pasien {
		if (pasien[i].tgl_data).Equal(tgl_data) {
			if pasien[i].umur <= 25 {
				jumlah_025 = jumlah_025 + pasien[i].jumlah_pasien
			} else if 25 < pasien[i].umur && pasien[i].umur <= 50 {
				jumlah_2550 = jumlah_2550 + pasien[i].jumlah_pasien
			} else {
				jumlah_51 = jumlah_51 + pasien[i].jumlah_pasien
			}
		}
	}

	data_graph = []graph_usia{
		graph_usia{
			Usia:          "< 25",
			Jumlah_pasien: jumlah_025,
		},
		graph_usia{
			Usia:          "25 <= 50",
			Jumlah_pasien: jumlah_2550,
		},
		graph_usia{
			Usia:          " 50 <",
			Jumlah_pasien: jumlah_51,
		},
	}

	return data_graph
}
