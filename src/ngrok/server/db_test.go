package server

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"testing"
	"errors"
)

func TestDBWrite(t *testing.T) {
	db, err := bolt.Open("ngrok_test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("ports"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("dbcuser1"), []byte("100"))
		err = b.Put([]byte("dbcuser2"), []byte("200"))

		return err

	})
}

func TestDBRead(t *testing.T) {
	db, err := bolt.Open("ngrok_test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ports"))
		if b == nil {
			return errors.New("bucket ports not exist")
		}

		k := "dbcuser1"
		v := b.Get([]byte(k))
		fmt.Printf("(%s %s)\n", k,v)
		return nil
	})
}

func TestDBView(t *testing.T) {
	db, err := bolt.Open("ngrok_test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ports"))

		if b == nil {
			return errors.New("bucket ports not exist")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
}