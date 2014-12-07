package engine

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/server"
	log "github.com/funkygao/log4go"
	"github.com/gorilla/mux"
	"net/http"
	"runtime"
	"syscall"
	"time"
)

func (this *Engine) stopHttpServ() {
	server.StopHttpServ()
}

func (this *Engine) launchHttpServ() {
	if this.conf.httpListenAddr == "" {
		return
	}

	server.LaunchHttpServ(this.conf.httpListenAddr, this.conf.pprofListenAddr)
	server.RegisterHttpApi("/admin/{cmd}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
	server.RegisterHttpApi("/h", func(w http.ResponseWriter, req *http.Request,
		params map[string]interface{}) (interface{}, error) {
		return this.handleHttpHelpQuery(w, req, params)
	}).Methods("GET")
	server.RegisterHttpApi("/help", func(w http.ResponseWriter, req *http.Request,
		params map[string]interface{}) (interface{}, error) {
		return this.handleHttpHelpQuery(w, req, params)
	}).Methods("GET")
}

func (this *Engine) handleHttpHelpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	output := make(map[string]interface{})
	if this.conf.pprofListenAddr != "" {
		output["pprof"] = "http://" + this.conf.pprofListenAddr + "/debug/pprof/"
	}

	output["uris"] = []string{"/admin/h", "/admin/help", "/admin/guide"}
	return output, nil
}

func (this *Engine) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	var (
		vars   = mux.Vars(req)
		cmd    = vars["cmd"]
		output = make(map[string]interface{})
	)

	switch cmd {
	case "ping":
		output["status"] = "ok"

	case "debug":
		stack := make([]byte, 1<<20)
		stackSize := runtime.Stack(stack, true)
		output["result"] = "go to global logger to see result"
		log.Info(string(stack[:stackSize]))

	case "stop":
		this.rpcServer.Stop()
		output["status"] = "stopped"

	case "stat":
		output["started"] = this.StartedAt
		output["elapsed"] = time.Since(this.StartedAt).String()
		output["pid"] = this.pid
		output["hostname"] = this.hostname
		output["stats"] = this.stats.String()
		rusage := syscall.Rusage{}
		syscall.Getrusage(0, &rusage)
		output["rusage"] = rusage

	case "runtime":
		output["runtime"] = this.stats.Runtime()

	case "mem":
		output["mem"] = *this.stats.memStats

	case "conf":
		output["engine"] = *this.conf
		output["servants"] = *config.Servants

	case "guide", "help", "h":
		output["uris"] = []string{
			"/admin/ping",
			"/admin/debug",
			"/admin/stop",
			"/admin/stat",
			"/admin/runtime",
			"/admin/mem",
			"/admin/conf",
		}
		if this.conf.pprofListenAddr != "" {
			output["pprof"] = "http://" + this.conf.pprofListenAddr + "/debug/pprof/"
		}

	default:
		return nil, server.ErrHttp404
	}

	return output, nil
}
