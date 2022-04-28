package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	name                 = filepath.Base(os.Args[0])
	logpath              = "/var/log/" + name
	logfile              = name + ".log"
	pf                   = fmt.Printf
	lf                   = log.Printf
	lln                  = log.Println
	pwd, _               = os.Getwd()
	I                    = "[INFO] "
	D                    = "[DEBUG] "
	E                    = "[ERROR] "
	debug                = false
	initLoggerStdoutOnly = false
	GetFunc              = func() string {
		pc, _, _, _ := runtime.Caller(1)
		return fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	}
	GetFuncLine = func() string {
		_, _, line, _ := runtime.Caller(1)
		return fmt.Sprintf("%v", line)
	}
	GetFuncDetails = func() string {
		pc, _, line, _ := runtime.Caller(1)
		return fmt.Sprintf("%v:[%+v]", runtime.FuncForPC(pc).Name(), line)
	}
	w io.Writer
)

func ldf(format string, args ...interface{}) {
	if debug {
		lf(D+format+"\n", args...)
	}
}

func lif(format string, args ...interface{}) {
	lf(I+format+"\n", args...)
}

func lef(format string, args ...interface{}) {
	lf(E+format+"\n", args...)
}

func liln(args ...interface{}) {
	var a []interface{}
	a = append(a, I)
	a = append(a, args...)
	lln(a...)
}

func leln(args ...interface{}) {
	var a []interface{}
	a = append(a, E)
	a = append(a, args...)
	lln(a...)
}

func ldln(args ...interface{}) {
	if debug {
		var a []interface{}
		a = append(a, D)
		a = append(a, args...)
		lln(a...)
	}
}

/*
func init() {
	log.SetPrefix(name + ": ")
	log.SetFlags(log.LstdFlags| log.Lmsgprefix)
	pathExists := createDirectory(logpath)
	if pathExists {
		f := CreateLogFile(logfile)
		w = io.MultiWriter(os.Stdout, f)
	} else {
		w = io.MultiWriter(os.Stdout)
	}
	log.SetOutput(w)

}
*/

func InitLogger() {
	if !initLoggerStdoutOnly {
		log.SetPrefix(name + ": ")
		log.SetFlags(log.LstdFlags | log.Lmsgprefix)
		pathExists := createDirectory(logpath)
		if pathExists {
			f := createLogFile(logfile)
			w = io.MultiWriter(os.Stdout, f)
		} else {
			w = io.MultiWriter(os.Stdout)
		}
		log.SetOutput(w)
	}
}

func SetDebug(b bool) {
	debug = b
}

func ChangeOutput(w io.Writer) {
	log.SetOutput(w)
}

func SetStdout() {
	log.SetOutput(os.Stdout)
}

func InitLoggerStdoutOnly(b bool) {
	if b {
		initLoggerStdoutOnly = true
		log.SetOutput(ioutil.Discard)
	}
}

func ldffunc(funcdetails, format string, args ...interface{}) {
	if debug {
		f := D + "[MSG]: func: < %v > "
		var a []interface{}
		a = append(a, funcdetails)
		a = append(a, args...)
		lf(f+format+"\n", a...)
	}
}

func createDirectory(path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		ldffunc(GetFuncDetails(), "Cannot create path [%v]", path)
		leln(err)
		return false
	}
	return true
}

func createLogFile(fname string) *os.File {
	f, err := os.OpenFile(logpath+"/"+fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		ldffunc(GetFuncDetails(), "Cannot create log file [%v]", fname)
		leln(err)
		return nil
	}
	return f
}
