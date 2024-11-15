package web

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"pc-beantragung/internal/database"
	. "pc-beantragung/internal/domain/signon"
	"strconv"
)

type Filter struct {
	State  ProcessingState
	Active bool
}

func (self Filter) has(state ProcessingState) bool {
	return self.Active && self.State == state
}

func ListSignonsHandler(db database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var signons []SignOn
		var filter Filter

		if !r.URL.Query().Has("filter[state]") {
			signons, _ = db.SignOnRepo().GetAll()
			filter = Filter{}
		} else {
			state := ProcessingStateFromString(r.URL.Query().Get("filter[state]"))
			filter = Filter{State: state, Active: true}
			signons, _ = db.SignOnRepo().GetByState(state)
		}

		SignOnList(filter, signons).Render(r.Context(), w)
	}
}

func ToggleSidebarHandler(db database.Service, toogleOn bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		signon, _ := db.SignOnRepo().GetById(id)

		if !toogleOn {
			return
		}

		Sidebar(signon).Render(r.Context(), w)
	}
}

func UpdateSignonHandler(db database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(w, "Error wrong id", http.StatusInternalServerError)
			return
		}

		signon, err := db.SignOnRepo().GetById(id)

		if err != nil {
			http.Error(w, "Error retrieven the Sigon", http.StatusInternalServerError)
			return
		}

		signon.MyComment = r.FormValue("comment")

		newState := ProcessingStateFromString(r.FormValue("state"))
		if newState != signon.MyState {
			signon.MyState = newState
			RemovedTr(signon.Id).Render(r.Context(), w)
		}

		db.SignOnRepo().UpdateContext(signon)
	}
}

func UploadFileHandler(db database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("signons")

		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
			return
		}

		defer file.Close()

		signons, _ := parseCSV(file)

		db.SignOnRepo().SaveAll(signons)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func parseCSV(file io.Reader) ([]SignOn, error) {
	var signons []SignOn

	// Initialize CSV reader.
	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Read the header row first.
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	// Read each row and convert it to a SignOn struct.
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %v", err)
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID: %v", err)
		}

		user := SignOn{
			Id:                   2,
			IdPc:                 id,
			Company:              sql.NullString{String: record[4], Valid: true},
			Firstname:            sql.NullString{String: record[5], Valid: true},
			Lastname:             sql.NullString{String: record[6], Valid: true},
			Zip:                  sql.NullString{String: record[7], Valid: true},
			City:                 sql.NullString{String: record[8], Valid: true},
			Street:               sql.NullString{String: record[9], Valid: true},
			HouseNo:              sql.NullString{String: record[10], Valid: true},
			PCState:              sql.NullString{String: record[11], Valid: true},
			DesiredDeliveryStart: sql.NullString{String: record[12], Valid: true},
			MeterNo:              sql.NullString{String: record[14], Valid: true},
			Malo:                 sql.NullString{String: record[15], Valid: true},
			Melo:                 sql.NullString{String: record[16], Valid: true},
			ConfigId:             sql.NullString{String: record[17], Valid: true},
		}

		signons = append(signons, user)
	}

	return signons, nil
}
