package fs

import (
	"io"
	"os"
)

type fileWriter struct {
	path string
	err error
	finfo os.FileInfo
	mode os.FileMode
	flags int
}

// Write creates a new file or truncates an existing one
// and sets it up for writing.
func Write(path string) FileWriter {
	fw := &fileWriter{path:path, flags:os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode: 0644}
	info, err := os.Stat(fw.path)
	if err == nil {
		fw.finfo = info
	}
	return fw
}

func Append(path string) FileWriter {
	fw := &fileWriter{path:path, flags:os.O_CREATE|os.O_APPEND|os.O_WRONLY, mode:0644}
	info, err := os.Stat(fw.path)
	if err == nil {
		fw.finfo = info
	}

	return fw
}

func (fw *fileWriter) Err() error {
	return fw.err
}

func (fw *fileWriter) Info() os.FileInfo {
	return fw.finfo
}

func (fw *fileWriter) String(str string) FileWriter{
	file, err := os.OpenFile(fw.path, fw.flags, fw.mode)
	if err != nil {
		fw.err = err
		return fw
	}
	defer file.Close()
	if fw.finfo, fw.err = file.Stat(); fw.err != nil{
		return fw
	}

	if _, err := file.WriteString(str); err != nil {
		fw.err = err
	}
	return fw
}

func (fw *fileWriter) Lines(lines []string) FileWriter{
	file, err := os.OpenFile(fw.path, fw.flags, fw.mode)
	if err != nil {
		fw.err = err
		return fw
	}
	defer file.Close()
	if fw.finfo, fw.err = file.Stat(); fw.err != nil{
		return fw
	}

	len := len(lines)
	for i, line := range lines {
		if _, err := file.WriteString(line); err != nil {
			fw.err = err
			return fw
		}
		if len > (i+1){
			if _, err := file.Write([]byte{'\n'}); err != nil {
				fw.err = err
				return fw
			}
		}
	}
	return fw
}

func (fw *fileWriter) Bytes(data []byte) FileWriter{
	file, err := os.OpenFile(fw.path, fw.flags, fw.mode)
	if err != nil {
		fw.err = err
		return fw
	}
	defer file.Close()
	if fw.finfo, fw.err = file.Stat(); fw.err != nil{
		return fw
	}

	if _, err := file.Write(data); err != nil {
		fw.err = err
	}
	return fw
}

func (fw *fileWriter) ReadFrom(r io.Reader) FileWriter {
	file, err := os.OpenFile(fw.path, fw.flags, fw.mode)
	if err != nil {
		fw.err = err
		return fw
	}
	defer file.Close()
	if fw.finfo, fw.err = file.Stat(); fw.err != nil{
		return fw
	}

	if _, err := io.Copy(file, r); err != nil {
		fw.err = err
	}
	return fw
}
