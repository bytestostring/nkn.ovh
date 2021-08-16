package nknovh_engine

import (
		"fmt"
		"os"
		"time"
		"sync"
		)

type logger struct {
	dir string
	list map[string]*os.File
	mux sync.Mutex
}

func (o *logger) Syslog(text string, filename string) error {

	var f string = o.dir + "/" + filename + ".log"
	var err error
	t := time.Now()
	ws := t.Format("[2006-01-02 / 15:04:05.000]") + " " + text + "\n"
	o.mux.Lock()
	defer o.mux.Unlock()
	if _, ok := o.list[filename]; ok {
		if _, err = o.list[filename].WriteString(ws); err != nil {
			fmt.Printf("WriteString has returned error: (%s)\n", err.Error())
			return err
		}
		return nil
	}
	o.list[filename], err = os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		fmt.Printf("Can't create/open the file \"%s.log\": (%s)\n", filename, err.Error())
		return err
	}
	if _, err = o.list[filename].WriteString(ws); err != nil {
		fmt.Printf("WriteString has returned error: (%s)\n", err.Error())
		return err
	}
	return nil
}

func (o *logger) Init() {
	o.list = map[string]*os.File{}
}
