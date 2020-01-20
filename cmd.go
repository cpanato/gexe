package echo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// RunProc executes command in cmdStr and waits for the result.
// It returns a *Proc with information about the executed process.
func (e *Echo) RunProc(cmdStr string) *Proc {
	cmdStr = lineRgx.ReplaceAllString(cmdStr, " ")
	e.shouldLog(cmdStr)

	proc := e.StartProc(cmdStr)
	proc.Wait()
	return proc
}

// StartProc executes the command in cmdStr and returns immediately
// without waiting. Information about the running process is stored in *Proc.
func (e *Echo) StartProc(cmdStr string) *Proc {
	return e.startProc(cmdStr)
}

// Run parses and executes cmdStr and returns the result as a string.
// This uses RunProc and extracts the result as a string.
func (e *Echo) Run(cmdStr string) string {
	p := e.RunProc(cmdStr)
	if p.Err() != nil {
		e.shouldLog(p.Err().Error())
		e.shouldPanic(p.Err().Error())
		return ""
	}

	data := p.Result()
	if p.Err() != nil {
		e.shouldLog(p.Err().Error())
		e.shouldPanic(p.Err().Error())
		return ""
	}

	return os.Expand(data, e.Val)
}

// Runout parses and executes command cmdStr and prints out the result
func (e *Echo) Runout(cmdStr string) {
	fmt.Print(e.Run(cmdStr))
}

func (e *Echo) startProc(cmdStr string) *Proc {
	words := e.splitWords(e.Eval(cmdStr))

	output := new(bytes.Buffer)
	command := exec.Command(words[0], words[1:]...)
	command.Stdout = output
	command.Stderr = output

	if err := command.Start(); err != nil {
		proc := Proc{id: command.Process.Pid, cmd: command, err: err, state: command.ProcessState}
		e.Procs = append(e.Procs, proc)
		return &proc
	}

	proc := Proc{id: command.Process.Pid, cmd: command, state: command.ProcessState, output: output}
	e.Procs = append(e.Procs, proc)
	return &proc
}
