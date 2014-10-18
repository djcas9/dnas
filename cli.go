package main

import (
	"fmt"
	"log"
	"os"

	"code.google.com/p/gopass"

	"github.com/sevlyar/go-daemon"
	"github.com/visionmedia/go-flags"

	"os/user"
	"strconv"
	"syscall"
)

func chuser(username string) (uid, gid int) {
	usr, err := user.Lookup(username)

	if err != nil {
		fmt.Printf("failed to find user %q: %s", username, err)
	}

	uid, err = strconv.Atoi(usr.Uid)

	if err != nil {
		fmt.Printf("bad user ID %q: %s", usr.Uid, err)
	}

	gid, err = strconv.Atoi(usr.Gid)

	if err != nil {
		fmt.Printf("bad group ID %q: %s", usr.Gid, err)
	}

	if err := syscall.Setgid(gid); err != nil {
		fmt.Printf("setgid(%d): %s", gid, err)
	}

	if err := syscall.Setuid(uid); err != nil {
		fmt.Printf("setuid(%d): %s", uid, err)
	}

	return uid, gid
}

// Options cli command options
type Options struct {
	Interface string `short:"i" long:"interface" description:"Interface to monitor" value-name:"eth0"`
	Port      int    `short:"p" long:"port" description:"DNS port" default:"53" value-name:"53"`
	Daemon    bool   `short:"D" long:"daemon" description:"Run DNAS in daemon mode"`
	Write     string `short:"w" long:"write" description:"Write JSON output to log file" value-name:"FILE"`
	User      string `short:"u" long:"user" description:"Drop privileges to this user" value-name:"USER"`
	Hexdump   bool   `short:"H" long:"hexdump" description:"Show hexdump of DNS packet"`

	Mysql    bool `long:"mysql" description:"Enable Mysql Output Support"`
	Postgres bool `long:"postgres" description:"Enable Postgres Output Support"`
	Sqlite3  bool `long:"sqlite3" description:"Enable Sqlite3 Output Support"`

	DbUser         string `long:"db-user" description:"Database User" value-name:"root"`
	DbPassword     string `long:"db-password" description:"Database Password" value-name:"PASSWORD"`
	DbDatabase     string `long:"db-database" description:"Database Database" value-name:"dnas"`
	DbHost         string `long:"db-host" description:"Database Host" value-name:"127.0.0.1"`
	DbPort         string `long:"db-port" description:"Database Port" value-name:"3306"`
	DbPath         string `long:"db-path" description:"Path to Database on disk. (sqlite3 only)" value-name:"~/.dnas.db"`
	DbTls          bool   `long:"db-tls" description:"Enable TLS / SSL encrypted connection to the database. (mysql/postgres only)" value-name:"false"`
	DbSkipVerify   bool   `long:"db-skip-verify" description:"Allow Self-signed or invalid certificate (mysql/postgres only)" value-name:"false"`
	DatabaseOutput bool   `long:"db-verbose" description:"Show database logs in STDOUT"`

	Quiet   bool `short:"q" long:"quiet" description:"Suppress DNAS output"`
	Version bool `short:"v" long:"version" description:"Show version information"`
}

func printUsage(p *flags.Parser) {
	fmt.Printf("\n  %s (%s) - %s\n",
		Name,
		Version,
		Description,
	)

	p.WriteHelp(os.Stdout)
	fmt.Printf("\n")
	os.Exit(1)
}

func printVersion() {
	fmt.Printf("%s - %s - Version: %s\n",
		Name,
		Description,
		Version,
	)

	os.Exit(1)
}

// CLIRun start DNAS and process all command-line options
func CLIRun(f func(options *Options)) {

	options := &Options{}

	var parser = flags.NewParser(options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		printUsage(parser)
	}

	if options.Version {
		printVersion()
	}

	if options.Mysql {
		if options.DbPassword == "" {
			password, err := gopass.GetPass("Database Password: ")

			if err != nil {
				panic(err)
			}

			options.DbPassword = password
		}
	}

	if options.DbHost != "" {

		options.DbHost = "tcp(" + options.DbHost + ":"

		if options.DbPort == "" {
			options.DbPort = "3306"
		}

		options.DbHost = options.DbHost + options.DbPort + ")"
	}

	if options.DbUser == "" {
		options.DbUser = "root"
	}

	if options.DbDatabase == "" {
		options.DbDatabase = "dnas"
	}

	if options.Interface == "" {
		printUsage(parser)
	}

	if options.Daemon {

		cntxt := &daemon.Context{
			PidFileName: "dnas.pid",
			PidFilePerm: 0644,
			LogFileName: "dnas.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
		}

		d, err := cntxt.Reborn()

		if err != nil {
			log.Fatalln(err)
		}

		if d != nil {
			return
		}

		defer cntxt.Release()

		go f(options)

		err = daemon.ServeSignals()

		if err != nil {
			log.Println("Error:", err)
		}
	} else {
		f(options)
	}
}
