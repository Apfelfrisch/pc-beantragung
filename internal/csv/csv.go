package csv

import (
	"database/sql"
	"encoding/csv"
	"io"
	"pc-beantragung/internal"
	"pc-beantragung/internal/signon"

	"github.com/gocarina/gocsv"
)

type Csv = []*CsvRow

type CsvRow struct {
	ID            int64         `csv:"#"`
	CreatedAt     internal.Time `csv:"Erstellt"`
	EnergyType    string        `csv:"Energietyp"`
	Company       string        `csv:"Firma"`
	Firstname     string        `csv:"Vorname"`
	Lastname      string        `csv:"Vorname"`
	Zip           string        `csv:"PLZ"`
	City          string        `csv:"Ort"`
	Street        string        `csv:"Straße"`
	HouseNo       string        `csv:"Hnr."`
	PcState       string        `csv:"Status"`
	DeliveryStart string        `csv:"Lieferbeginn"`
	DeliveryEnd   string        `csv:"Lieferende"`
	MeterNo       string        `csv:"Zähler"`
	Malo          string        `csv:"Marktlokation"`
	Melo          string        `csv:"Messlokation"`
	ConfigID      string        `csv:"Config-Id"`
}

func (self *CsvRow) ToSignOn() signon.Signon {
	return signon.Signon{
		IDPc:                 self.ID,
		EnergyType:           sql.NullString{String: self.EnergyType, Valid: true},
		Company:              sql.NullString{String: self.Company, Valid: true},
		Firstname:            sql.NullString{String: self.Firstname, Valid: true},
		Lastname:             sql.NullString{String: self.Lastname, Valid: true},
		Zip:                  sql.NullString{String: self.Zip, Valid: true},
		City:                 sql.NullString{String: self.City, Valid: true},
		Street:               sql.NullString{String: self.Street, Valid: true},
		HouseNo:              sql.NullString{String: self.HouseNo, Valid: true},
		PcState:              sql.NullString{String: self.PcState, Valid: true},
		DesiredDeliveryStart: sql.NullString{String: self.DeliveryStart, Valid: true},
		MeterNo:              sql.NullString{String: self.MeterNo, Valid: true},
		Malo:                 sql.NullString{String: self.Malo, Valid: true},
		Melo:                 sql.NullString{String: self.Melo, Valid: true},
		ConfigID:             sql.NullString{String: self.ConfigID, Valid: true},
	}
}

func ParseCsv(file io.Reader) (Csv, error) {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma = ';'
		return r
	})

	csvContent := Csv{}

	if err := gocsv.Unmarshal(file, &csvContent); err != nil {
		return nil, err
	}

	return csvContent, nil
}
