package util

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// log color tags
const (
	lct_reset = "\033[0m"
	lct_error = "\033[31m"
	lct_warn  = "\033[33m"
	lct_log   = "\033[32m"
	lct_debug = "\033[36m"
)

func Debug(args ...interface{}) {
	r := append([]interface{}{lct_debug, "[DEBUG] "}, args...)
	r = append(r, lct_reset)
	log.Print(r...)
}

func Log(args ...interface{}) {
	r := append([]interface{}{lct_log, "[LOG] "}, args...)
	r = append(r, lct_reset)
	log.Print(r...)
}

func Warn(args ...interface{}) {
	r := append([]interface{}{lct_warn, "[WARN] "}, args...)
	r = append(r, lct_reset)
	log.Print(r...)
}

func Fatal(args ...interface{}) {
	r := append([]interface{}{lct_error, "[FATAL] "}, args...)
	r = append(r, lct_reset)
	log.Print(r...)
}

func CheckErrPanic(err error, args ...interface{}) {
	if err != nil {
		r := append([]interface{}{"[PANIC] "}, args...)
		r = append(r, ": ")
		r = append(r, err)
		log.Print(r...)
		os.Exit(1)
	}
}

func CheckErr(err error, args ...interface{}) bool {
	if err != nil {
		r := append([]interface{}{"[ERROR] "}, args...)
		r = append(r, ": ")
		r = append(r, err)
		log.Print(r...)
		return true
	}
	return false
}

func Error(args ...interface{}) error {
	r := append([]interface{}{lct_error, "[ERROR] "}, args...)
	r = append(r, lct_reset)
	log.Print(r...)
	return errors.New(fmt.Sprint(r...))
}
