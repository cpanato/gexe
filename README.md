# `Echo$` 
Script-like functionalities wrapped in the rigid safety of Go!

## Start
Create a session:
```
e := echo.New()
```
Then configure your session:
```
e.Conf.SetPanicOnError(true)
```

## Variables

### Var, SetVar
Echo supports storing values that can be accessed using method `e.Var()` for the duration of a session:
```
e.Var("Foo=Bar")
e.Var("Fuzz=${Foo} Buzz=Bazz") 
```
Method `e.SetVar(name, value string)` saves a named value one at a time.

### Env, SetEnv
Values can be made visible as environment variables for externally launced commands using method `e.Env()`. Both methods support value expansion as shown below:
```
e.Env("Foo=Bar")
e.Env("Fuzz=${Foo} Buzz=Bazz")
e.Env("BAZZ=$HOME")
```
Method `e.SetEnv(name, value string)` sets environment variables one value at a time.

### Expansion
All `echo` methods support variable value expansion using `$name` or `${name}` which are automatically replaced with the value of the named variable.

## Slices

### Split
A space-separated list can be turned into a native Go `[]string` using the `e.Split()` method as shown:
```
e.SetVar("list", "item3 item4")
for _, val := range e.Split("item0 item1 item2 $list"){  
    ...
}
```

An additional separator value may be provided to `e.Split()`:
```
e.SetVar("list", "item3;item4")
for _, val := range e.Arr("Hello;World;!;$list", ";"){  
    fmt.Println(val)
}
```

### Glob
Method `e.Glob()` expands the provided shell path pattern into a slice of matching file/directory names:
```
for _, f := range e.Glob("$HOME/go/src/*.com") {
  fmt.Println(f)
}
```
## Program
`Echo` exposes the `Prog` type which provides access to methods related to the running program as listed
below:
```
e.Prog.Name() - Then name of the running binary
e.Prog.Path() - the full path of the running binary
e.Prog.Args() - Slice of CLI parameters
e.Prog.Exit() - Exits the running program
e.Prog.Pid() - The process ID for the program
e.Prog.Ppid() - The program's parent process ID
e.Prog.Workdir() - the current working directory of the program
e.Prog.Avail(path) - Returns the full path of the specified program if available
```

## Processes
`Echo` can be used to start external processes by wraping the `os/exec` types.  Methods that start processes
return a `*Proc` value with information about the started process.

### StartProc
Method `e.StartProc(<command-string>)` starts a new process by running the provided command string.  The method 
returns a Proc which can be used to track the process as shown below:

```
proc := e.StartProc('echo "Hello World!"')
fmt.Println("Proces id", proc.ID)
```
Method `e.StartProc` returns immediately and does not wait for the process to complete.  This means Proc information
may be incomplete since the process may still be running (see Managing Processes for detail).

Each proc started is also saved in the Echo session and can be reached via slice `e.Procs[]`.


### Managing Processes 
When a process is launched with `StartProc`, it immediately returns a value of type `Proc` which exposes several methods to inspect the running 
(or completed) process.

```
proc.ExitCode() - the process exit code
proc.Exited() - true if program finished
proc.ID() - the process ID
proc.IsSuccess() - true if program finished with no error
proc.Peek() - Updates process state information
proc.SysTime() - the system time for the executed process
proc.UserTime() - the system user time for executed process
proc.Out() - exposes an io.Reader to stream result (generate error)
proc.Result() - reads all command result from proc.Out() and returns it as string (generate error)
proc.Wait() - Blocks, waits for a process to complete (generate error)
```

Type `Proc` stores the last error generated by invoking any of its methods.  The error can be accessed as `proc.Err()` as shown below:

```
proc := e.StartProc('echo "Hello world!"')
proc.Wait()
// optionally check for error
if proc.Err() != nil {
    fmt.Printl(proc.Err())
}
```

### RunProc
As a convenience, the `Proc` namespace exposes three methods that will start a process and wait for it to complete.  First, method
`e.RunProc` starts a process, wait for it to complete, then returns a value of type `Proc`:

```
e.Var("CONTAINER_IMG", "hello-world")
proc := e.RunProc('docker run --rm --detach $CONTAINER_IMG')
r := proc.Result()
if e.IsEmpty(r) || proc.Err() != nil {
    fmt.Println("found error:", proc.Err())
}
fmt.Println("Ran container container:", r)

```

### Run
Method `e.Run` starts a process, wait for it to complete, and then returns a string value as the result from the executed command:
```
r := e.Var("CONTAINER_IMG", "hello-world").RunProc('docker run --rm --detach $CONTAINER_IMG')
if e.IsEmpty(r) || proc.Err() != nil {
    fmt.Println("found error:", proc.Err())
}
fmt.Println("Ran container container:", r)

```
Note that `e.Run` is equivalent to calling `e.RunProc` and then calling method `Result` on the returned `Proc`.

### Runout
Method `e.Runout` starts a process, wait for its completion, and then prints the result to standard output:

```
e.Var("CONTAINER_IMG", "hello-world").Runout('docker run --rm --detach $CONTAINER_IMG')
```

## Strings
```
e.Empty(string) bool
e.Lower(string) string
e.Upper(string) string
e.Streq(string, string) bool // string equal
e.Trim(string)string 
```

## Files
File operation methods:

```
e.Abs()
e.Rel()
e.Base()
e.Dir()
e.PathSym()
e.Ext()
e.PathJoin()
e.PathMatched()
e.IsExit()
e.IsReg()
e.IsDir()
e.Mkdirs()
e.Rmdirs()
e.Chown()
e.Chmod()
e.AreSame() // Are files equal
```

## User

```
e.Username()
e.Home()
e.Gid()
e.Uid()
```

## Flagset
...

## More to come

## License
MIT