package husky

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

var whileChar = regexp.MustCompile("^[0-9a-zA-Z_-.]+$")

type Database struct {
	conn      *sql.DB
	limit     int
	offset    int
	where     []string
	fields    []string
	bind      []interface{}
	enableLog bool
	queryLog  []string
}

func NewDb(dsn string) *Database {
	d := &Database{}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("DataBase Connection Error: %s\n", err)
		return d
	}

	d.conn = db
	return d
}

func (d *Database) EnableLog(enable bool) {
	d.enableLog = enable
}

func (d *Database) Select(columns ...string) *Database {
	for _, c := range columns {
		if !whileChar.MatchString(c) {
			panic("Invalid columns name specified.")
		}
		d.fields = append(d.fields, c)
	}

	return d
}

func (d *Database) Limit(limit int) *Database {
	d.limit = limit

	return d
}

func (d *Database) Offset(offset int) *Database {
	d.offset = offset

	return d
}

func (d *Database) Where(field, operator string, bind interface{}) *Database {
	d.where = append(d.where, field+" "+operator+" ?")
	d.bind = append(d.bind, bind)

	return d
}

func (d *Database) Get(table string) (rows *sql.Rows, err error) {
	query := d.buildSelectQuery(table)
	defer d.clear()

	d.log(query, d.bind)
	if len(d.bind) > 0 {
		return d.conn.Query(query, d.bind...)
	} else {
		return d.conn.Query(query)
	}
}

func (d *Database) GetRow(table string) (row *sql.Row) {
	query := d.buildSelectQuery(table)
	defer d.clear()

	d.log(query, d.bind)
	if len(d.bind) > 0 {
		return d.conn.QueryRow(query, d.bind...)
	} else {
		return d.conn.QueryRow(query)
	}
}

func (d *Database) Insert(table string, values map[string]interface{}) (result sql.Result, err error) {
	query, bind := d.buildInsertQuery(table, values)
	defer d.clear()

	d.log(query, d.bind)
	return d.conn.Exec(query, bind...)
}

func (d *Database) Update(table string, values map[string]interface{}) (result sql.Result, err error) {
	query, bind := d.buildUpdateQuery(table, values)
	defer d.clear()

	d.log(query, d.bind)
	return d.conn.Exec(query, bind...)
}

func (d *Database) buildSelectQuery(table string) (query string) {
	if !whileChar.MatchString(table) {
		panic("Invalid table name specified.")
	}

	query = "SELECT "
	if len(d.fields) == 0 {
		query += "*"
	} else {
		query += strings.Join(d.fields, ", ")
	}
	query += " FROM " + table
	if len(d.where) > 0 {
		query += " WHERE " + strings.Join(d.where, " AND ")
	}

	if d.limit > 0 {
		query += " LIMIT " + fmt.Sprint(d.limit)
	}
	if d.offset > 0 {
		query += " OFFSET " + fmt.Sprint(d.offset)
	}
	return
}

func (d *Database) buildInsertQuery(table string, values map[string]interface{}) (query string, bind []interface{}) {
	var (
		fields    []string
		statement []string
	)

	for f, val := range values {
		fields = append(fields, f)
		statement = append(statement, "?")
		bind = append(bind, val)
	}

	query = fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(fields, ", "),
		strings.Join(statement, ", "),
	)
	return
}

func (d *Database) buildUpdateQuery(table string, values map[string]interface{}) (query string, bind []interface{}) {
	var fields []string

	for f, val := range values {
		fields = append(fields, f+" = ?")
		bind = append(bind, val)
	}

	query = "UPDATE " + table + " SET  (" + strings.Join(fields, ", ") + ")"
	if len(d.where) > 0 {
		query += " WHERE " + strings.Join(d.where, " AND ")
	}

	if d.limit > 0 {
		query += " LIMIT " + fmt.Sprint(d.limit)
	}

	return
}

func (d *Database) log(query string, params []interface{}) {
	if d.enableLog {
		log := fmt.Sprintf("%s, %v", query, params)
		d.queryLog = append(d.queryLog, log)
	}
}

func (d *Database) LastQuery() string {
	index := len(d.queryLog) - 1
	return d.queryLog[index]
}

func (d *Database) AllQuery() string {
	return strings.Join(d.queryLog, "\n")
}

func (d *Database) clear() {
	d.limit = 0
	d.offset = 0
	d.where = []string{}
	d.fields = []string{}
	d.bind = []interface{}{}
}
