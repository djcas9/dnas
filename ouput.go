package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	lvl "github.com/syndtr/goleveldb/leveldb"
	"os"
	"sort"
	"text/tabwriter"
)

var (
	w     = new(tabwriter.Writer)
	count = 0
)

func MakeDB(path string) (db *lvl.DB, err error) {
	db, err = lvl.OpenFile(path, nil)

	return db, err
}

func (message *Message) ToLevelDB(db *lvl.DB, options *Options) (err error) {
	data, err := db.Get([]byte(message.Dns.Question), nil)
	var buf bytes.Buffer

	if err != nil {
		enc := gob.NewEncoder(&buf)
		eerr := enc.Encode(message.Dns.Answers)

		if eerr != nil {
			return eerr
		}

		err = db.Put([]byte(message.Dns.Question), buf.Bytes(), nil)

		return err
	}

	var a []Answer
	enc := gob.NewDecoder(bytes.NewReader(data))
	eerr := enc.Decode(&a)

	if eerr != nil {
		return eerr
	}

	// for i, aa := range message.Dns.Answers {
	// fmt.Println(i, aa)
	// }

	// fmt.Println("GOT DATA!@#!@#!@#!@#!@#!@#:::::::::::::::", a[0].Data, len(a))

	return nil
}

func FindKeyBy(database_path string, question string) {
	db, dberr := MakeDB(database_path)

	if dberr != nil {
		fmt.Println(dberr)
		os.Exit(-1)
	}

	iter := db.NewIterator(nil, nil)

	for iter.Next() {
		key := iter.Key()

		if string(key) != question {
			continue
		}

		data := iter.Value()

		var a []Answer
		enc := gob.NewDecoder(bytes.NewReader(data))
		eerr := enc.Decode(&a)

		if eerr != nil {
			panic(eerr)
		}

		w.Init(os.Stdout, 1, 2, 2, ' ', 0)

		fmt.Fprintf(w, "\n\tQuestion:\t\033[0;31;49m%s\033[0m\n\n", question)

		fmt.Fprintf(w, "\t\033[0;32;49mAnswers (%d):\033[0m\t\n\n", len(a))

		fmt.Fprintf(w, "\tRR\tName\tData\n")
		fmt.Fprintf(w, "\t----\t----\t----\n")

		for i := range a {
			fmt.Fprintf(w, "\t%s", a[i].Record)
			fmt.Fprintf(w, "\t%s", a[i].Name)
			fmt.Fprintf(w, "\t\033[0;32;49m%s\033[0m\n", a[i].Data)
		}

		fmt.Fprintf(w, "\n")
	}

	iter.Release()
	err := iter.Error()

	if err != nil {
		panic(err)
	}
}

func ListAllQuestions(database_path string) {
	list := make(map[string]int)
	var count int = 0
	var keys []string

	db, dberr := MakeDB(database_path)

	if dberr != nil {
		fmt.Println(dberr)
		os.Exit(-1)
	}

	iter := db.NewIterator(nil, nil)

	for iter.Next() {
		key := iter.Key()
		data := iter.Value()

		var a []Answer
		enc := gob.NewDecoder(bytes.NewReader(data))
		eerr := enc.Decode(&a)

		if eerr != nil {
			panic(eerr)
		}

		keyValue := string(key)
		length := len(a)

		keys = append(keys, keyValue)
		list[keyValue] = length
	}

	iter.Release()
	err := iter.Error()

	if err != nil {
		panic(err)
	}

	sort.Strings(keys)

	fmt.Printf("\n Questions:\n\n")

	for _, value := range keys {
		count = count + list[value]
		fmt.Printf(" * %s (\033[0;32;49m%d\033[0m)\n", value, list[value])
	}

	fmt.Printf("\n Questions: %d", len(keys))
	fmt.Printf("\n   Answers: %d\n\n", count)
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
