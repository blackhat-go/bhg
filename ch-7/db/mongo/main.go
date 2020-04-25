package main

import (
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/blackhat-go/bhg/ch-7/db/dbminer"
)

type MongoMiner struct {
	Host    string
	session *mgo.Session
}

func New(host string) (*MongoMiner, error) {
	m := MongoMiner{Host: host}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MongoMiner) connect() error {
	s, err := mgo.Dial(m.Host)
	if err != nil {
		return err
	}
	m.session = s
	return nil
}

func (m *MongoMiner) GetSchema() (*dbminer.Schema, error) {
	var s = new(dbminer.Schema)

	dbnames, err := m.session.DatabaseNames()
	if err != nil {
		return nil, err
	}

	for _, dbname := range dbnames {
		db := dbminer.Database{Name: dbname, Tables: []dbminer.Table{}}
		collections, err := m.session.DB(dbname).CollectionNames()
		if err != nil {
			return nil, err
		}

		for _, collection := range collections {
			table := dbminer.Table{Name: collection, Columns: []string{}}

			var docRaw bson.Raw
			err := m.session.DB(dbname).C(collection).Find(nil).One(&docRaw)
			if err != nil {
				return nil, err
			}

			var doc bson.RawD
			if err := docRaw.Unmarshal(&doc); err != nil {
				if err != nil {
					return nil, err
				}
			}

			for _, f := range doc {
				table.Columns = append(table.Columns, f.Name)
			}
			db.Tables = append(db.Tables, table)
		}
		s.Databases = append(s.Databases, db)
	}
	return s, nil
}

func main() {

	mm, err := New(os.Args[1])
	if err != nil {
		panic(err)
	}
	if err := dbminer.Search(mm); err != nil {
		panic(err)
	}
}
