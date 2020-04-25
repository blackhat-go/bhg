package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/blackhat-go/bhg/ch-7/db/dbminer"
)

type MySQLMiner struct {
	Host string
	Db   sql.DB
}

func New(host string) (*MySQLMiner, error) {
	m := MySQLMiner{Host: host}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MySQLMiner) connect() error {

	db, err := sql.Open("mysql", fmt.Sprintf("root:password@tcp(%s:3306)/information_schema", m.Host))
	if err != nil {
		log.Panicln(err)
	}
	m.Db = *db
	return nil
}

func (m *MySQLMiner) GetSchema() (*dbminer.Schema, error) {
	var s = new(dbminer.Schema)

	sql := `SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME FROM columns
	WHERE TABLE_SCHEMA NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
	ORDER BY TABLE_SCHEMA, TABLE_NAME`
	schemarows, err := m.Db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer schemarows.Close()

	var prevschema, prevtable string
	var db dbminer.Database
	var table dbminer.Table
	for schemarows.Next() {
		var currschema, currtable, currcol string
		if err := schemarows.Scan(&currschema, &currtable, &currcol); err != nil {
			return nil, err
		}

		if currschema != prevschema {
			if prevschema != "" {
				db.Tables = append(db.Tables, table)
				s.Databases = append(s.Databases, db)
			}
			db = dbminer.Database{Name: currschema, Tables: []dbminer.Table{}}
			prevschema = currschema
			prevtable = ""
		}

		if currtable != prevtable {
			if prevtable != "" {
				db.Tables = append(db.Tables, table)
			}
			table = dbminer.Table{Name: currtable, Columns: []string{}}
			prevtable = currtable
		}
		table.Columns = append(table.Columns, currcol)
	}
	db.Tables = append(db.Tables, table)
	s.Databases = append(s.Databases, db)
	if err := schemarows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func main() {
	mm, err := New(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer mm.Db.Close()

	if err := dbminer.Search(mm); err != nil {
		panic(err)
	}
}
