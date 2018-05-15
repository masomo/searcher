package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/tidwall/redcon"
	"github.com/tinylib/msgp/msgp"
)

var (
	flagaddr   = flag.String("addr", "0.0.0.0:6480", "redis server listen addr")
	flagdbdir  = flag.String("dbdir", "./", "Database folder")
	flagcpus   = flag.Int("C", 1, "Set the maximum number of CPUs to use")
	flagLogLvl = flag.String("L", "info", "Log verbosity level [crit,error,warn,info,debug]")
	flagsave   = flag.Duration("save", time.Duration(5*time.Minute), "Dump memory to db file interval")
	flagpprof  = flag.Bool("pprof", false, "Debug information on http port :6060")
)

var search *Search

func randomInt(min, max int) int {
	if min == max {
		return min
	}

	return rand.Intn(max-min) + min
}

func randomString(n int) string {
	var letters = []rune("123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func syncDBProcess() {
	runtime.LockOSThread()

	ret, ret2, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)

	if errno != 0 || ret2 < 0 {
		log.Error("Search DB sync failed", "error", "fork process failed")
		return
	}

	if runtime.GOOS == "darwin" && ret2 == 1 {
		ret = 0
	}

	if ret > 0 {
		// Parent process

		doneCh := make(chan bool)
		timeout := time.NewTimer(5 * time.Minute)
		killed := false

		go func() {
			select {
			case <-doneCh:
				log.Debug("Sync process received done")
			case <-timeout.C:
				log.Error("Sync process waited too long, sending kill process..")

				killed = true

				proc, err := os.FindProcess(int(ret))
				if err != nil {
					return
				}

				proc.Kill()
			}
		}()

		status := syscall.WaitStatus(0)

		_, err := syscall.Wait4(int(ret), &status, 0, nil)
		if err != nil {
			log.Debug("Sync process wait4 failed", "error", err.Error())
		}

		if !killed {
			timeout.Stop()
			doneCh <- true
		}

		return
	}

	log.Info("Sync process started", "pid", syscall.Getpid(), "parent", syscall.Getppid())

	// Child process
	tempFileName := fmt.Sprintf("%s/search.db-temp-%d", *flagdbdir, randomInt(1000, 9999))

	searchDB, err := os.OpenFile(tempFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Error("Search DB sync failed", "error", err.Error())
		goto done
	}

	err = search.Sync(searchDB, true)
	if err != nil {
		log.Error("Search DB sync failed", "error", err.Error())
		goto done
	}

	searchDB.Close()

	err = os.Rename(tempFileName, fmt.Sprintf("%s/search.db", *flagdbdir))
	if err != nil {
		log.Error("Search DB sync failed", "error", err.Error())
	}

done:
	log.Info("Sync process exiting...")
	os.Exit(0)
}

func syncDBService(done chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			syncDBProcess()
		}
	}
}

func main() {
	flag.Parse()

	if *flagcpus == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(*flagcpus)
	}

	lvl, err := log.LvlFromString(*flagLogLvl)
	if err != nil {
		log.Crit("Log verbosity level unknown")
		os.Exit(1)
	}

	log.Root().SetHandler(log.LvlFilterHandler(lvl, log.StdoutHandler))

	startup := time.Now()

	searchDB, err := os.OpenFile(fmt.Sprintf("%s/search.db", *flagdbdir), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Crit("Search DB open failed", "error", err.Error())
		os.Exit(1)
	}

	search = NewSearch()

	err = msgp.ReadFile(search, searchDB)
	if err != nil && err.Error() == os.ErrInvalid.Error() {
		log.Warn("Search data not found on DB file, creating new one...")

		search = NewSearch()

		err = search.Sync(searchDB, false)
		if err != nil {
			log.Crit("Search encode to DB file failed", "error", err.Error())
			os.Exit(1)
		}

	} else if err != nil {
		log.Crit("Search decode from DB file failed", "error", err.Error())
		os.Exit(1)
	}

	searchDB.Close()

	syncTicker := time.NewTicker(*flagsave)
	syncDone := make(chan bool)

	go syncDBService(syncDone, syncTicker)

	go func() {
		err = redcon.ListenAndServe(*flagaddr, onRedisCommand, onRedisConnect, onRedisClose)
		if err != nil {
			log.Crit("Redis server startup failed", "error", err.Error())
			os.Exit(1)
		}
	}()

	if *flagpprof {
		go func() {
			err := http.ListenAndServe(":6060", nil)
			if err != nil {
				log.Error("http listener for pprof failed", "error", err.Error())
			}
		}()
	}

	log.Info("Searcher service started", "addr", *flagaddr, "startup", time.Since(startup).Round(time.Millisecond))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	<-c

	log.Info("Searcher service stopping")

	syncTicker.Stop()
	syncDone <- true

	syncDBProcess()
}
