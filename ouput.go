package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/AndreasBriese/bloom"
	"github.com/boltdb/bolt"
)

var (
	w      = tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	count  = 0
	dbname = []byte("dnas")
	layout = "Jan 2, 2006 at 03:04pm (MST)"
)

type Value struct {
	Answers []Answer
	Bloom   []byte
}

func EncodeDNS(data Value) (buff bytes.Buffer, err error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err = enc.Encode(data)

	if err != nil {
		return buf, err
	}

	return buf, nil
}

func DecodeDNS(data []byte) (buff Value, err error) {
	var a Value
	enc := gob.NewDecoder(bytes.NewReader(data))
	eerr := enc.Decode(&a)

	if eerr != nil {
		return a, eerr
	}

	return a, nil
}

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

func prettyPrint(question string, a Value) {

	fmt.Fprintf(w, "\n\tQuestion:\t\033[0;31;49m%s\033[0m\n\n", question)

	fmt.Fprintf(w, "\t\033[0;32;49mAnswers (%d):\033[0m\n\n", len(a.Answers))

	fmt.Fprintf(w, "\tRR\tName\tData\n")

	fmt.Fprintf(w, "\t--\t----\t----\n")

	for i := range a.Answers {
		fmt.Fprintf(w, "\t%s", a.Answers[i].Record)
		fmt.Fprintf(w, "\t%s", a.Answers[i].Name)
		fmt.Fprintf(w, "\t\033[0;32;49m%s\033[0m ", a.Answers[i].Data)
		fmt.Fprintf(w, "\t%s", a.Answers[i].UpdatedAt.Format(layout))

		if a.Answers[i].Active {
			fmt.Fprintf(w, "\t(Active: \033[0;32;49m%s\033[0m)\n", "Yes")
		} else {
			fmt.Fprintf(w, "\t(Active: \033[0;31;49m%s\033[0m)\n", "No")
		}
	}

	fmt.Fprintf(w, "\n")

	w.Flush()

}

func (message *Message) ToKVDB(options *Options) (err error) {

	db, dberr := MakeDB(options.Database)

	defer db.Close()

	if dberr != nil {
		fmt.Println(dberr)
		os.Exit(-1)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(dbname)

		if err != nil {
			return err
		}

		data := bucket.Get([]byte(message.Dns.Question))

		if data != nil {

			// Update Records

			a, _ := DecodeDNS(data)

			for _, aa := range message.Dns.Answers {
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
				index := contains(dr, message.Dns.Answers)

				a.Answers[i].UpdatedAt = time.Now()

				if index == -1 {
					a.Answers[i].Active = false
				} else {
					a.Answers[i].Active = true
				}

			}

			buf, _ := EncodeDNS(a)
			err = bucket.Put([]byte(message.Dns.Question), buf.Bytes())

			if err != nil {
				return err
			}

		} else {

			val := Value{Answers: message.Dns.Answers, Bloom: message.Bloom}

			buf, _ := EncodeDNS(val)

			err = bucket.Put([]byte(message.Dns.Question), buf.Bytes())

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func findByQuestion(database_path string, r *regexp.Regexp) {

	db, dberr := MakeDB(database_path)

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

			prettyPrint(string(key), a)

			return nil
		})

		return nil
	})
}

func findByAnswer(database_path string, find []byte) {

	db, dberr := MakeDB(database_path)

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

			isIn := bf.Has(find)

			if isIn {
				prettyPrint(string(key), a)
			}

			return nil
		})

		return nil
	})
}

func listAllQuestions(database_path string) {
	list := make(map[string]int)
	var count int = 0
	var keys []string

	db, dberr := MakeDB(database_path)

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

func (message *Message) ToStdout(options *Options) {
	w.Init(os.Stdout, 1, 2, 2, ' ', 0)

	count++

	fmt.Fprintf(w, "\t---\t%d\t---\n", count)
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\tQuestion:\t\033[0;31;49m%s\033[0m\n", message.Dns.Question)
	fmt.Fprintf(w, "\tTimestamp:\t%s\n", message.Timestamp)
	fmt.Fprintf(w, "\tSource:\t%s\n", message.SrcIp)
	fmt.Fprintf(w, "\tDestination:\t%s\n", message.DstIp)
	fmt.Fprintf(w, "\tLength:\t%d\n\n", message.Dns.Length)

	fmt.Fprintf(w, "\t\033[0;32;49mAnswers (%d):\033[0m\t\n\n", len(message.Dns.Answers))

	fmt.Fprintf(w, "\tRR\tName\tData\n")
	fmt.Fprintf(w, "\t----\t----\t----\n")

	for i := range message.Dns.Answers {
		fmt.Fprintf(w, "\t%s", message.Dns.Answers[i].Record)
		fmt.Fprintf(w, "\t%s", message.Dns.Answers[i].Name)
		fmt.Fprintf(w, "\t\033[0;32;49m%s\033[0m\n", message.Dns.Answers[i].Data)
	}

	if options.Hexdump {
		fmt.Fprintf(w, "\n\t\033[0;32;49mHexdump:\033[0m\n\n%s\n", hex.Dump(message.Packet))
	}

	fmt.Fprintf(w, "\n")
}

func (message *Message) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}
