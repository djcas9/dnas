package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogits/gogs/modules/log"
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

func MysqlConnect(options *Options) (db *sql.DB, err error) {

	var connect string

	connect = options.MysqlUser + ":" +
		options.MysqlPassword + "@" + options.MysqlHost + "/" +
		options.MysqlDatabase

	db, err = sql.Open("mysql", connect)

	if err != nil {
		return db, err
	}

	// defer db.Close()

	err = db.Ping()

	if err != nil {
		return db, err
	}

	checkQuestionTable := `SELECT COUNT(*)
	FROM information_schema.tables 
	WHERE table_schema = '` + options.MysqlDatabase + `' 
	AND table_name = 'questions';`

	checkAnswerTable := `SELECT COUNT(*)
	FROM information_schema.tables 
	WHERE table_schema = '` + options.MysqlDatabase + `' 
	AND table_name = 'answers';`

	var countQ int
	var countA int

	err = db.QueryRow(checkQuestionTable).Scan(&countQ)

	if err != nil {
		panic(err)
	}

	err = db.QueryRow(checkAnswerTable).Scan(&countA)

	if err != nil {
		panic(err)
	}

	fmt.Println(countQ, countA)

	// var questionTable string
	// var answerTable string
	// var r *sql.Result

	if countQ == 0 {
		questionTable := `
		CREATE TABLE ` + "`" + `questions` + "`" + ` (
			` + "`" + `id` + "`" + ` int(11) unsigned NOT NULL AUTO_INCREMENT,
			` + "`" + `packet` + "`" + ` text,
			` + "`" + `dst_ip` + "`" + ` text,
			` + "`" + `src_ip` + "`" + ` text,
			` + "`" + `timestamp` + "`" + ` datetime DEFAULT NULL,
			` + "`" + `protocol` + "`" + ` text,
			` + "`" + `name` + "`" + ` text,
			PRIMARY KEY (` + "`" + `id` + "`" + `)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1;
		`

		_, err := db.Exec(questionTable)

		if err != nil {
			panic(err)
		}
	}

	if countA == 0 {
		answerTable := `
		CREATE TABLE ` + "`" + `answers` + "`" + ` (
			` + "`" + `id` + "`" + ` int(11) unsigned NOT NULL AUTO_INCREMENT,
			` + "`" + `question_id` + "`" + ` int(11) NOT NULL,
			` + "`" + `name` + "`" + ` tinytext,
			` + "`" + `record` + "`" + ` tinytext,
			` + "`" + `data` + "`" + ` text,
			` + "`" + `created_at` + "`" + ` datetime DEFAULT NULL,
			` + "`" + `updated_at` + "`" + ` datetime DEFAULT NULL,
			` + "`" + `active` + "`" + ` tinyint(1) NOT NULL,
			PRIMARY KEY (` + "`" + `id` + "`" + `)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1;
		`

		_, err := db.Exec(answerTable)
		if err != nil {
			panic(err)
		}
	}

	// check for tables or create
	// CREATE TABLE questions (id INT);

	return db, nil
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

	_, err = db.Exec("INSERT INTO questions (name, packet, src_ip, dst_ip, timestamp, protocol) VALUES (?, ?, ?, ?, ?, ?);",
		message.DNS.Question, message.Packet, message.SrcIP, message.DstIP, message.Timestamp, message.Protocol)

	if err != nil {
		log.Error(err.Error())
	}

	return nil
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
