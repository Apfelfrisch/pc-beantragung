package signon

import (
	"context"
	"database/sql"
)

type SignOnRepo interface {
	ListAll(ctx context.Context) ([]Signon, error)
	ListForState(ctx context.Context, state string) ([]Signon, error)
	GetById(ctx context.Context, id int64) (Signon, SignonContext, error)
	UpdateSignonContext(ctx context.Context, context SignonContext) error
	SaveAll(ctx context.Context, signons []Signon) error
}

type signOnRepo struct {
	db      *sql.DB
	queries *Queries
}

func SignonRepo(queries *Queries, db *sql.DB) SignOnRepo {
	return signOnRepo{queries: queries, db: db}
}

func (self signOnRepo) ListAll(ctx context.Context) ([]Signon, error) {
	return self.queries.ListSignOns(ctx)
}

func (self signOnRepo) ListForState(ctx context.Context, state string) ([]Signon, error) {
	return self.queries.ListSignOnsByState(ctx, state)
}

func (self signOnRepo) GetById(ctx context.Context, id int64) (Signon, SignonContext, error) {
	row, err := self.queries.GetSignOn(ctx, id)

	if err != nil {
		return Signon{}, SignonContext{}, err
	}

	return row.Signon, row.SignonContext, nil
}

func (self signOnRepo) SaveAll(ctx context.Context, signons []Signon) error {
	tx, err := self.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("delete from signons")
	if err != nil {
		return err
	}

	queries := self.queries.WithTx(tx)

	for _, signon := range signons {
		queries.CreateSignOn(ctx, CreateSignOnParams{
			IDPc:                 signon.IDPc,
			Company:              signon.Company,
			Firstname:            signon.Firstname,
			Lastname:             signon.Lastname,
			Zip:                  signon.Zip,
			City:                 signon.City,
			Street:               signon.Street,
			HouseNo:              signon.HouseNo,
			PcState:              signon.PcState,
			DesiredDeliveryStart: signon.DesiredDeliveryStart,
			MeterNo:              signon.MeterNo,
			Malo:                 signon.Malo,
			Melo:                 signon.Melo,
			ConfigID:             signon.ConfigID,
			CreatedAt:            signon.CreatedAt,
		})
	}

	queries.FillupContext(ctx)

	return tx.Commit()
}

func (self signOnRepo) UpdateSignonContext(ctx context.Context, context SignonContext) error {
	return self.queries.UpdateContext(ctx, UpdateContextParams{
		State:      context.State,
		Comment:    context.Comment,
		SignonIDPc: context.SignonIDPc,
	})
}
