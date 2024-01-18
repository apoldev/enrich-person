package repo

import (
	"database/sql"
	"errors"
	"fio/internal/person"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type PostgresRepo struct {
	Logger          *logrus.Entry
	Db              *sqlx.DB
	PaginationLimit int
}

func (r *PostgresRepo) Update(p *person.Person) error {

	res, err := r.Db.NamedExec(`Update person set 
                  name=:name, 
                  surname=:surname,
                  patronymic=:patronymic,
                  age=:age,
                  gender=:gender,
                  nationality=:nationality
              where id=:id`, p)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()

	if affected == 0 {
		return person.ErrNotFound
	}

	return nil

}

func (r *PostgresRepo) Get(id int) (*person.Person, error) {

	p := person.Person{}

	err := r.Db.Get(
		&p,
		`select id, name, surname, patronymic, age, gender, nationality from person where id = $1`,
		id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, person.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PostgresRepo) Delete(id int) error {

	res, err := r.Db.Exec(`delete from person where id = $1`, id)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()

	if affected == 0 {
		return person.ErrNotFound
	}

	return err
}

func (r *PostgresRepo) Create(person *person.Person) error {

	err := r.Db.QueryRow(`insert into person 
    (name, surname, patronymic, age, gender, nationality) values 
    ($1, $2, $3, $4, $5, $6) 
    RETURNING id;`,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	).Scan(&person.ID)

	return err
}

func (r *PostgresRepo) createQueryWithFilters(filters map[string]string) (sql string, args []interface{}) {

	s := strings.Builder{}
	args = make([]interface{}, 0, len(filters))

	s.WriteString(`select id, name, surname, patronymic, age, gender, nationality from person`)

	j := 1

	if len(filters) > 0 {
		s.WriteString(" WHERE ")
	}

	for i := range filters {

		if j > 1 {
			s.WriteString(" and ")
		}

		s.WriteString(fmt.Sprintf("%s = $%d", i, j))
		args = append(args, filters[i])
		j++
	}

	argsCount := len(args)
	s.WriteString(fmt.Sprintf(" order by id desc limit $%d offset $%d", argsCount+1, argsCount+2))

	return s.String(), args

}

func (r *PostgresRepo) GetPersons(page int, filters map[string]string) ([]person.Person, error) {

	persons := make([]person.Person, 0)

	offset := r.PaginationLimit * (page - 1)

	sql, args := r.createQueryWithFilters(filters)
	args = append(args, r.PaginationLimit, offset)

	r.Logger.Debugf("sql: %s, args: %v", sql, args)

	err := r.Db.Select(&persons, sql, args...)

	if err != nil {
		return nil, err
	}

	return persons, nil
}
