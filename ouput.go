package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/AndreasBriese/bloom"
	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/karlbunch/tablewriter"
	"github.com/mgutz/ansi"
)

var (
	count  = 0
	dbname = []byte("dnas")
	layout = "01/02/06 03:04pm (MST)"
	lime   = ansi.ColorCode("green+h:black")
	red    = ansi.ColorCode("red")
	green  = ansi.ColorCode("green")
	reset  = ansi.ColorCode("reset")
)

// Value holds the dns answers and bloom filter fior database storage
type Value struct {
	Answers []Answer
	Bloom   []byte
}

func MysqlConnect(options string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", options)

	if err != nil {
		return db, err
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return db, err
	}

	return db, nil
}

// EncodeDNS encode the Value struct for storage in the database.
func EncodeDNS(data Value) (buff bytes.Buffer, err error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err = enc.Encode(data)

	if err != nil {
		return buf, err
	}

	return buf, nil
}

// DecodeDNS decode the Value struct from the database
func DecodeDNS(data []byte) (buff Value, err error) {
	var a Value
	enc := gob.NewDecoder(bytes.NewReader(data))
	eerr := enc.Decode(&a)

	if eerr != nil {
		return a, eerr
	}

	return a, nil
}

// MakeDB create or open an existing database file
func MakeDB(path string) (db *bolt.DB, err error) {
	db, err = bolt.Open(path, 0644, nil)

	return db, err
}

func contains(m Answer, list []Answer) int {
	for i, b := range list {
		if b.Record == m.Record && b.Name == m.Name && b.Data == m.Data {
			return i
		}
	}

	return -1
}

func prettyPrint(question string, a Value, count int) {

	fmt.Printf("---   %d   ---\n", count)

	fmt.Printf("\nQuestion: \033[0;31;49m%s\033[0m\n\n", question)

	fmt.Printf("\033[0;32;49mAnswers (%d):\033[0m\n\n", len(a.Answers))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"RR", "Name", "Data", "Last Seen", "Active"})
	// table.SetBorder(false)

	for i, aa := range a.Answers {

		active := green + "Yes" + reset

		if !a.Answers[i].Active {
			active = red + "No" + reset
		}

		table.Append([]string{
			aa.Record,
			aa.Name,
			lime + aa.Data + reset,
			aa.UpdatedAt.Format(layout),
			active,
		})
	}

	table.Render()
	fmt.Printf("\n")
}

func (message *Message) ToMysql(db *sql.DB, options *Options) (err error) {
	fmt.Println(db)
	return nil
}

// ToKVDB write data to key/value database
func (message *Message) ToKVDB(options *Options) (err error) {

	db, dberr := MakeDB(options.Database)

	defer db.Close()

	if dberr != nil {
		fmt.Println(dberr)
		os.Exit(-1)
	}

	bf := bloom.New(float64(1<<16), float64(0.01))

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(dbname)

		if err != nil {
			return err
		}

		data := bucket.Get([]byte(message.DNS.Question))

		if data != nil {

			// Update Records

			a, _ := DecodeDNS(data)

			for _, aa := range message.DNS.Answers {
				index := contains(aa, a.Answers)

				if index != -1 {
					// Update
					// fmt.Println("UPDATE RECORD!!!!", index)
					a.Answers[index].Active = true
					a.Answers[index].UpdatedAt = time.Now()
				} else {
					// Add New
					// fmt.Println("ADD NEW!!!!!", index)
					a.Answers = append(a.Answers, aa)
				}
			}

			for i, dr := range a.Answers {
				index := contains(dr, message.DNS.Answers)

				a.Answers[i].UpdatedAt = time.Now()

				if index == -1 {
					a.Answers[i].Active = false
				} else {
					a.Answers[i].Active = true
				}

				bf.Add([]byte(a.Answers[i].Data))
			}

			a.Bloom = bf.JSONMarshal()

			buf, _ := EncodeDNS(a)
			err = bucket.Put([]byte(message.DNS.Question), buf.Bytes())

			if err != nil {
				return err
			}

		} else {

			val := Value{Answers: message.DNS.Answers, Bloom: message.Bloom}

			buf, _ := EncodeDNS(val)

			err = bucket.Put([]byte(message.DNS.Question), buf.Bytes())

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func findByQuestion(database string, r *regexp.Regexp) {

	fmt.Printf("\n")
	var qcount int
	db, dberr := MakeDB(database)

	defer db.Close()

	if dberr != nil {
		fmt.Println(dberr)
		os.Exit(-1)
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbname)

		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", dbname)
		}

		bucket.ForEach(func(key, data []byte) error {

			if !r.Match(key) {
				return nil
			}

			a, _ := DecodeDNS(data)

			qcount++
			prettyPrint(string(key), a, qcount)

			return nil
		})

		return nil
	})
}

func findByAnswer(database string, find []byte) {

	fmt.Printf("\n")
	var acount int
	db, dberr := MakeDB(database)

	defer db.Close()

	if dberr != nil {
		fmt.Println(dberr)
		os.Exit(-1)
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbname)

		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", dbname)
		}

		bucket.ForEach(func(key, data []byte) error {

			a, _ := DecodeDNS(data)
			bf := bloom.JSONUnmarshal(a.Bloom)

			if len(a.Bloom) > 0 {
				if bf.Has(find) {
					acount++
					prettyPrint(string(key), a, acount)
				}
			}

			return nil
		})

		return nil
	})
}

func listAllQuestions(database string) {
	list := make(map[string]int)
	var count int
	var keys []string

	db, dberr := MakeDB(database)

	defer db.Close()

	if dberr != nil {
		fmt.Println(dberr)
		os.Exit(-1)
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(dbname)

		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", dbname)
		}

		bucket.ForEach(func(key, data []byte) error {
			a, _ := DecodeDNS(data)

			keyValue := string(key)
			length := len(a.Answers)

			keys = append(keys, keyValue)
			list[keyValue] = length

			return nil
		})

		sort.Strings(keys)

		fmt.Printf("\n Questions:\n\n")

		for _, value := range keys {
			count = count + list[value]
			fmt.Printf(" * %s (\033[0;32;49m%d\033[0m)\n", value, list[value])
		}

		fmt.Printf("\n Questions: %d", len(keys))
		fmt.Printf("\n   Answers: %d\n\n", count)

		return nil
	})
}

// ToStdout write data to standard out
func (message *Message) ToStdout(options *Options) {

	count++

	val := Value{Answers: message.DNS.Answers, Bloom: message.Bloom}

	prettyPrint(message.DNS.Question, val, count)

	if options.Hexdump {
		fmt.Printf("\033[0;32;49mHexdump:\033[0m\n\n%s\n", hex.Dump(message.Packet))
	}

}

// ToJSON convert Message struct to json
func (message *Message) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}
