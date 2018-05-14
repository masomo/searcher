package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/inconshreveable/log15"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/tidwall/redcon"
)

var (
	// metric registry keys
	searcherCommandMetricPrefix = "searcher.command"
	redisMonitorCh              = make(chan string)
)

func onRedisCommand(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToLower(string(cmd.Args[0]))

	start := time.Now()

	redisCommandNext(conn, cmd)

	commandMetric := metrics.GetOrRegisterTimer(fmt.Sprintf("%s.%s", searcherCommandMetricPrefix, command), nil)
	commandMetric.UpdateSince(start)

	redisMonitor(conn, cmd)
}

// monitor middleware
func redisMonitor(conn redcon.Conn, cmd redcon.Command) {
	select {
	case redisMonitorCh <- fmt.Sprintf("- %s [%s] |%s|",
		time.Now().Format("2006/01/02 15:04:05.00"),
		conn.RemoteAddr(), string(bytes.Join(cmd.Args, []byte(" ")))):
	default:
	}
}

func redisCommandNext(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToLower(string(cmd.Args[0]))

	switch command {
	case "ping":
		if len(cmd.Args) >= 2 {
			conn.WriteBulkString(string(cmd.Args[1]))
			return
		}

		conn.WriteString("PONG")
	case "monitor":
		dc := conn.Detach()
		go func() {
			defer dc.Close()

			dc.WriteString("OK")
			dc.Flush()

			for {
				dc.WriteString(<-redisMonitorCh)

				err := dc.Flush()
				if err != nil {
					break
				}
			}
		}()
	case "metrics":
		data, err := json.Marshal(metrics.DefaultRegistry.GetAll())
		if err != nil {
			conn.WriteNull()
			return
		}

		conn.WriteBulkString(string(data))
	case "set":
		if len(cmd.Args) != 4 {
			writeRedisError(conn, errors.New("ERR wrong number of arguments for 'SET' command"))
			return
		}

		key := string(cmd.Args[1])
		id := string(cmd.Args[2])
		value := string(cmd.Args[3])

		search.Set(key, id, value)

		conn.WriteString("OK")
	case "del":
		if len(cmd.Args) != 3 {
			writeRedisError(conn, errors.New("ERR wrong number of arguments for 'DEL' command"))
			return
		}

		key := string(cmd.Args[1])
		id := string(cmd.Args[2])

		search.Del(key, id)

		conn.WriteString("OK")
	case "search":
		if len(cmd.Args) < 3 {
			writeRedisError(conn, errors.New("ERR wrong number of arguments for 'SEARCH' command"))
			return
		}

		key := string(cmd.Args[1])
		query := string(cmd.Args[2])
		start := 0
		stop := 0

		if len(cmd.Args) == 5 {
			var err error
			start, err = strconv.Atoi(string(cmd.Args[3]))
			if err != nil {
				start = 0
			}

			stop, err = strconv.Atoi(string(cmd.Args[4]))
			if err != nil {
				stop = 0
			}
		}

		result := search.Search(key, query, start, stop)

		data, err := json.Marshal(result)
		if err != nil {
			conn.WriteNull()
			return
		}

		conn.WriteBulkString(string(data))
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	}
}

func writeRedisError(conn redcon.Conn, err error) {
	conn.WriteError(fmt.Sprintf("ERR on command: %s", err.Error()))
}

func onRedisConnect(conn redcon.Conn) bool {
	log.Info("Redis new connection", "remote", conn.RemoteAddr())
	return true
}

func onRedisClose(conn redcon.Conn, err error) {
	log.Info("Redis connection closed", "remote", conn.RemoteAddr())
}
