package main

import (
	"bytes"
	"io"
	"testing"

	m "github.com/xemoe/go-monitor/monitor"
)

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func TestValidate(t *testing.T) {
	monitor := m.Monitor{}

	err := monitor.Validate()
	if err == nil {
		t.Errorf("Looking for %v, got %v", "We need to monitor at least one process", nil)
	}

	monitor.Processes = []string{"test"}
	//
	// @TODO validate required config
	//
	/**
	err = monitor.Validate()
	if err == nil {
		t.Errorf("Looking for %v, got %v", "Not all config variables present", nil)
	}
	**/

	monitor.Config.DefaultTTLSeconds = 1

	err = monitor.Validate()
	if err != nil {
		t.Errorf("Looking for %v, got %v", nil, err)
	}
}

func TestGetServerInfo(t *testing.T) {
	// @TODO
}

func TestCheckProc(t *testing.T) {
	// @TODO Not sure how to test this without involving setting up a channel
}

func TestLineCount(t *testing.T) {
	line := bytes.NewBufferString("test one line\n")
	lines, err := lineCounter(line)
	if err != nil {
		t.Errorf("Looking for %v, got %v", nil, err)
	}
	if lines != 1 {
		t.Errorf("Looking for %v, got %v", 1, lines)
	}

	line = bytes.NewBufferString("test one line\ntest two lines\n")
	lines, err = lineCounter(line)
	if err != nil {
		t.Errorf("Looking for %v, got %v", nil, err)
	}
	if lines != 2 {
		t.Errorf("Looking for %v, got %v", 2, lines)
	}

	line = bytes.NewBufferString("test one line\ntest two lines\nthree\nfour\nfive\n")
	lines, err = lineCounter(line)
	if err != nil {
		t.Errorf("Looking for %v, got %v", nil, err)
	}
	if lines != 5 {
		t.Errorf("Looking for %v, got %v", 5, lines)
	}
}

func TestNotifyProcError(t *testing.T) {
	// @TODO
}
