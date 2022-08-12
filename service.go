package mysql

// http://code.flickr.com/blog/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/
// https://github.com/go-sql-driver/mysql
// https://golang.org/pkg/database/sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/service"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

func init() {
	ctx := context.Background()
	service.RegisterService(ctx, "mysql", NewMySQLService)
}

type MySQLService struct {
	service.Service
	dsn       string
	key       string
	offset    int64
	increment int64
}

func NewMySQLService(ctx context.Context, uri string) (service.Service, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	dsn := q.Get("dsn")

	s := &MySQLService{
		dsn:       dsn,
		key:       "integers",
		offset:    1,
		increment: 2,
	}

	// maybe keep a permanent connection open?
	// (20180623/thisisaaronland)

	db, err := s.connect()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *MySQLService) SetLastInt(ctx context.Context, i int64) error {

	last, err := s.LastInt(ctx)

	if err != nil {
		return err
	}

	if i < last {
		return fmt.Errorf("integer value too small")
	}

	db, err := s.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	sql := fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT=%d", s.key, i)
	st, err := db.Prepare(sql)

	if err != nil {
		return err
	}

	_, err = st.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLService) SetOffset(ctx context.Context, i int64) error {
	s.offset = i
	return nil
}

func (s *MySQLService) SetIncrement(ctx context.Context, i int64) error {
	s.increment = i
	return nil
}

func (s *MySQLService) LastInt(ctx context.Context) (int64, error) {

	db, err := s.connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	sql := fmt.Sprintf("SELECT MAX(id) FROM %s", s.key)
	row := db.QueryRow(sql)

	var max int64

	err = row.Scan(&max)

	if err != nil {
		return -1, err
	}

	return max, nil
}

// https://dev.mysql.com/doc/refman/5.7/en/getting-unique-id.html

func (s *MySQLService) NextInt(ctx context.Context) (int64, error) {

	db, err := s.connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	err = s.set_autoincrement(db)

	if err != nil {
		return -1, err
	}

	sql := fmt.Sprintf("REPLACE INTO %s (stub) VALUES(?)", s.key)

	result, err := db.Exec(sql, "a")

	if err != nil {
		return -1, err
	}

	next, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return next, nil
}

func (s *MySQLService) set_autoincrement(db *sql.DB) error {

	sql_incr := fmt.Sprintf("SET @@auto_increment_increment=%d", s.increment)
	st_incr, err := db.Prepare(sql_incr)

	if err != nil {
		return err
	}

	defer st_incr.Close()

	_, err = st_incr.Exec()

	if err != nil {
		return err
	}

	sql_off := fmt.Sprintf("SET @@auto_increment_offset=%d", s.offset)
	st_off, err := db.Prepare(sql_off)

	if err != nil {
		return err
	}

	defer st_off.Close()

	_, err = st_off.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLService) Close(ctx context.Context) error {
	return nil
}

func (s *MySQLService) connect() (*sql.DB, error) {

	db, err := sql.Open("mysql", s.dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
