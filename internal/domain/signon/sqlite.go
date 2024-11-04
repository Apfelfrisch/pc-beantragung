package signon

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const SELECT = `SELECT id, 
	id_pc, 
	company, 
	firstname, 
	lastname, 
	zip, 
	city,
	street, 
	house_no, 
	pc_state, 
	desired_delivery_start, 
	meter_no, 
	malo, 
	melo, 
	config_id, 
	state, 
	comment 
	from signons 
	left join signon_context on signons.id_pc = signon_context.signon_id_pc`

const INSERT = `INSERT INTO signons 
	(
		id_pc, 
		company, 
		firstname, 
		lastname, 
		zip, 
		city, 
		street, 
		house_no,
		pc_state, 
		desired_delivery_start,
		meter_no, malo,
		melo, 
		config_id,
		created_at
	) 
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

func NewSignSqliteOnRepo(db *sql.DB) SignOnRepo {
	return &signOnSqLiteRepo{
		db: db,
	}
}

type signOnSqLiteRepo struct {
	db *sql.DB
}

func (self *signOnSqLiteRepo) GetById(id int) (*SignOn, error) {
	row := self.db.QueryRow(SELECT+" where id = ?", id)

	var signOn SignOn

	err := scanToSignOn(row, &signOn)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &signOn, nil
}

func (self *signOnSqLiteRepo) GetByState(state ProcessingState) ([]SignOn, error) {
	return self.signOnsFromQuery(SELECT+" where state = ?", state)
}

func (self *signOnSqLiteRepo) GetAll() ([]SignOn, error) {
	return self.signOnsFromQuery(SELECT)
}

func (self *signOnSqLiteRepo) UpdateContext(signon *SignOn) error {
	tx, err := self.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE signon_context SET state = ?, comment = ? WHERE signon_id_pc = ?", signon.MyState, signon.MyComment, signon.IdPc)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (self *signOnSqLiteRepo) SaveAll(signons []SignOn) error {
	tx, err := self.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from signons")
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare(INSERT)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, signon := range signons {
		err := self.insertSignOn(tx, signon)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return self.insertMissingState()
}

func (self *signOnSqLiteRepo) insertSignOn(tx *sql.Tx, signon SignOn) error {
	stmt, err := tx.Prepare(INSERT)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		signon.IdPc,
		signon.Company,
		signon.Firstname,
		signon.Lastname,
		signon.Zip,
		signon.City,
		signon.Street,
		signon.HouseNo,
		signon.PCState,
		signon.DesiredDeliveryStart,
		signon.MeterNo,
		signon.Malo,
		signon.Melo,
		signon.ConfigId,
		time.Now().Format(time.DateOnly),
	)

	if err != nil {
		log.Print(err)
		tx.Rollback()
		return err
	}

	return nil
}

func (self *signOnSqLiteRepo) insertMissingState() error {
	signons, err := self.signOnsFromQuery(SELECT + " where signon_context.state is null")

	if err != nil {
		fmt.Println(err)
		return err
	}

	tx, err := self.db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	query := "INSERT INTO signon_context (signon_id_pc, state) values (?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, signon := range signons {
		_, err := stmt.Exec(signon.IdPc, StateProcessing)

		if err != nil {
			log.Print(err)
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (self *signOnSqLiteRepo) signOnsFromQuery(query string, args ...any) ([]SignOn, error) {
	rows, err := self.db.Query(query, args...)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	var signOns []SignOn

	for rows.Next() {
		var signOn SignOn
		err = scanToSignOn(rows, &signOn)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		signOns = append(signOns, signOn)
	}

	return signOns, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanToSignOn(row scanner, signon *SignOn) error {
	return row.Scan(
		&signon.Id,
		&signon.IdPc,
		&signon.Company,
		&signon.Firstname,
		&signon.Lastname,
		&signon.Zip,
		&signon.City,
		&signon.Street,
		&signon.HouseNo,
		&signon.PCState,
		&signon.DesiredDeliveryStart,
		&signon.MeterNo,
		&signon.Malo,
		&signon.Melo,
		&signon.ConfigId,
		&signon.MyState,
		&signon.MyComment,
	)
}
