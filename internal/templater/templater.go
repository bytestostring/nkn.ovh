package templater

import (
		"fmt"
		"os"
		"errors"
		"sync"
		"io/ioutil"
		"regexp"
		"bytes"
	)


func FileExists(f string) bool {
	x, err := os.Stat(f);
	if os.IsNotExist(err) {
		return false
	}
	if x.IsDir() {
		return false
	}
	return true
}


type Templater struct {
	dir string
	locker sync.RWMutex
	Templates map[string]*Temp
}

type Temp struct {
	Name string
	Call *Templater
	page []byte
	rawPage []byte
	rawPages map[string][]byte
	complete bool
	mu_complete sync.RWMutex
	mu sync.RWMutex
}

func NewTemplater(directory string) *Templater {
	return &Templater{
						dir: directory,
						locker: sync.RWMutex{},
						Templates: map[string]*Temp{},
					}
}

func (t *Templater) New(name string) *Temp {
	if _, ok := t.Templates[name]; ok {
		return t.Templates[name]
	}
	x := &Temp{
				Name: name,
				Call: t,
				rawPages: map[string][]byte{},
				mu_complete: sync.RWMutex{},
				mu: sync.RWMutex{},
			}

	t.Templates[name] = x
	return x
}

func (t *Temp) Complete() {
	t.mu_complete.Lock()
	defer t.mu_complete.Unlock()
	t.complete = true
}

func (t *Temp) Uncomplete() {
	t.mu_complete.Lock()
	defer t.mu_complete.Unlock()
	t.complete = false
}

func (t *Temp) IsComplete() bool {
	t.mu_complete.RLock()
	defer t.mu_complete.RUnlock()
	return t.complete
}

func (t *Temp) Clear() {
	t.mu.Lock()
	t.page = []byte{}
	t.mu.Unlock()
	return
}

func (t *Temp) Flush() {
	t.mu.Lock()
	t.page = t.rawPage
	t.mu.Unlock()
	return
}

func (t *Temp) GetPage(page string, inside_path string, recursive ...bool) error {
	t.mu_complete.RLock()
	if t.complete {
		t.mu_complete.RUnlock()
		return nil
	}
	t.mu_complete.RUnlock()
	path := fmt.Sprintf("%s/%s/%s.tpl", t.Call.dir, inside_path, page)
	if !FileExists(path) {
		return errors.New("File not found: " + path)
	}
	t.mu.Lock()
	if _, ok := t.rawPages[page]; !ok {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			t.mu.Unlock()
			return errors.New("Cannot read the file: " + path)
		}
		t.rawPages[page] = read
	}
	t.mu.Unlock()
	err, buf := t.handlingIncludes(page)
	if err != nil {
		return err
	}
	if len(recursive) > 0 {
		return nil
	}
	t.mu.Lock()
	t.page = append(t.page, buf...)
	t.rawPage = t.page
	t.mu.Unlock()
	return nil
}


func (t *Temp) handlingIncludes(page string) (error, []byte) {

	t.mu.RLock()
	includes := bytes.Split(t.rawPages[page], []byte("\n"))
	t.mu.RUnlock()
	buf := make([]byte, 0)
	re := regexp.MustCompile(`{{{include "(.+)\/(.+)\.tpl"}}}\|m`)
	for i, _:= range includes {
		s := re.FindSubmatch(includes[i])
		if len(s) == 0 {
			buf = append(buf, includes[i]...)
			buf = append(buf, []byte("\n")...)
			continue
		}
		x := [2]string{string(s[2]), string(s[1])}
		if err := t.GetPage(x[0], x[1], true); err != nil {
			return errors.New(fmt.Sprintf("Cannot get the template file: %s/%s.tpl", x[1], x[0])), []byte{}
		}
		t.mu.Lock()
		buf = append(buf, t.rawPages[x[0]]...)
		if _, ok := t.rawPages[x[0]]; ok {
			delete(t.rawPages, x[0]);
		}
		t.mu.Unlock()
	}

	return nil, buf
}

func (t *Temp) Set(marker string, replace string) {
	t.page = bytes.ReplaceAll(t.page, []byte(fmt.Sprintf("{{%s}}|m", marker)), []byte(replace))
	return
}

func (t *Temp) View() []byte {
	return t.page
}