package ports_db

import (
	"errors"
	"github.com/boltdb/bolt"
	"log"
)

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open("./db/ngrok_client_ports.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}


// k: :dbcuser1",  v: "client-id-tcp:20001"
func WriteToDB(k string, v string) error {

	if db == nil {
		return errors.New("db is nil")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("ports"))
		if err != nil {
			return err
		}

		return b.Put([]byte(k), []byte(v))
	})
}

// k: token
func ReadFromDB(k string) (v string, err error) {

	if db == nil {
		return "", errors.New("db is nil")
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ports"))
		if b == nil {
			return errors.New("bucket not exist")
		}
		v_ := b.Get([]byte(k))
		v = string(v_[:])
		//fmt.Printf("(%s %s)\n", k,v)
		return nil
	})

	return
}

// k: token
func DeleteFromDB(k string) (err error) {
	if db == nil {
	return errors.New("db is nil")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("ports"))
		if err != nil {
			return err
		}

		return b.Delete([]byte(k))
	})
}

func ReadAllFromDB() (kvs map[string]string, err error) {

	if db == nil {
		return nil, errors.New("db is nil")
	}

	kvs = make(map[string]string)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ports"))
		if b == nil {
			return errors.New("bucket not exist")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			kvs[string(k[:])] = string(v[:])

		}
		return nil
	})

	return
}
