package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/1321822099/pdf_create/app/utils/config"
)

const (
	defaultTimeoutMillisecond = int64(60)
)

var (
	timeoutError = errors.New("timeout reached, process killed")
)

type CommandReq struct {
	WorkDir            string   `json:"work_dir"`
	Command            string   `json:"command"`
	Args               []string `json:"args"`
	TimeoutMillisecond int64    `json:"timeout_millisecond"`
	IngoreStdout       bool     `json:"ingore_stdout"`
	IngoreStderr       bool     `json:"ingore_stderr"`
}

type CommandResp struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func (req *CommandReq) validate() error {
	if req.Command == "" {
		return errors.New("miss field command")
	}
	if req.TimeoutMillisecond < 0 {
		return errors.New("invalid field timeout_millisecond")
	}
	return nil
}

func (req *CommandReq) Run(begin time.Time) (stdout string, stderr string, err error) {
	var cmd *exec.Cmd
	millisecond := getTimeOutMillisecond(req) - time.Since(begin).Milliseconds()
	if millisecond <= 0 {
		err = timeoutError
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(millisecond)*time.Millisecond)
	defer cancel()
	cmd = exec.CommandContext(ctx, req.Command, req.Args...)
	var stderrBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	cmd.Stdout = &stdoutBuf
	cmd.Dir = req.WorkDir
	err = cmd.Run()
	if err != nil {
		if err.Error() == "signal: killed" {
			err = timeoutError
			return
		}
		err = fmt.Errorf("failed to run cmd: %v, stderr: %s", err, stderrBuf.String())
		return
	}
	return stdoutBuf.String(), stderrBuf.String(), nil
}

func RunCommand(req *CommandReq) (interface{}, error) {
	err := req.validate()
	if err != nil {
		return nil, err
	}
	begin := time.Now()
	select {
	case <-pop():
		defer push()
		stdout, stderr, err := req.Run(begin)
		if err != nil {
			return nil, err
		}
		if req.IngoreStdout {
			stdout = ""
		}
		if req.IngoreStderr {
			stderr = ""
		}
		return &CommandResp{
			Stdout: stdout,
			Stderr: stderr,
		}, err
	case <-time.After(time.Duration(getTimeOutMillisecond(req)) * time.Millisecond):
		return nil, timeoutError
	}
}

func getTimeOutMillisecond(req *CommandReq) int64 {
	if req.TimeoutMillisecond > 0 {
		return req.TimeoutMillisecond
	}
	ms, _ := config.Int64("runcommand_max_request_timeout_millisecond", defaultTimeoutMillisecond)
	return ms
}
