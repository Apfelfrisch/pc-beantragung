package web

import (
	"context"
	"net/http"
	"pc-beantragung/internal/csv"
	"pc-beantragung/internal/database"
	so "pc-beantragung/internal/signon"
	"strconv"

	"github.com/samber/lo"
)

type Filter struct {
	State  string
	Active bool
}

func (self Filter) has(state string) bool {
	return self.Active && self.State == state
}

func ListSignonsHandler(db database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var signons []so.Signon
		var filter Filter

		ctx := context.Background()
		if !r.URL.Query().Has("filter[state]") {
			signons, _ = db.SignonRepo().ListAll(ctx)
			filter = Filter{}
		} else {
			state := r.URL.Query().Get("filter[state]")
			filter = Filter{State: state, Active: true}
			signons, _ = db.SignonRepo().ListForState(ctx, r.URL.Query().Get("filter[state]"))
		}

		SignOnList(filter, signons).Render(r.Context(), w)
	}
}

func ToggleSidebarHandler(db database.Service, toogleOn bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !toogleOn {
			return
		}

		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

		if err != nil {
			http.Error(w, "Invalid Id", http.StatusInternalServerError)
			return
		}

		signon, signonContext, err := db.SignonRepo().GetById(context.Background(), id)

		if err != nil {
			http.Error(w, "No signon", http.StatusInternalServerError)
			return
		}

		Sidebar(signon, signonContext).Render(r.Context(), w)
	}
}

func UpdateSignonHandler(db database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

		if err != nil {
			http.Error(w, "Invalid Id", http.StatusInternalServerError)
			return
		}

		ctx := context.Background()

		signon, signonContext, err := db.SignonRepo().GetById(ctx, id)

		if err != nil {
			http.Error(w, "No signon", http.StatusInternalServerError)
			return
		}

		signonContext.Comment = r.FormValue("comment")

		if r.FormValue("state") != signonContext.State {
			signonContext.State = r.FormValue("state")

			RemoveTr(signon.ID).Render(r.Context(), w)
		}

		db.SignonRepo().UpdateSignonContext(ctx, signonContext)
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

		rows, err := csv.ParseCsv(file)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		signons := lo.Map(rows, func(row *csv.CsvRow, _ int) so.Signon {
			return row.ToSignOn()
		})

		db.SignonRepo().SaveAll(context.Background(), signons)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
