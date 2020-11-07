package md2hugo

import (
    "fmt"
    "io"
    "os"
)

func Debug(target interface{}) {
    Debugf("%+v", target)
}

func Debugf(format string, a ...interface{}) {
   writeLog(os.Stderr, "debug: ", format, a...)
}

func Warn(target interface{}) {
    Debugf("%+v", target)
}

func Warnf(format string, a ...interface{}) {
    writeLog(os.Stderr, "warn: ", format, a...)
}

func Log(target interface{}) {
    Logf("%+v", target)
}

func Logf(format string, a ...interface{}) {
    writeLog(os.Stdout, "", format, a...)
}

func writeLog(w io.Writer, prefix string, format string, a ...interface{}) {
    // print prefix
    _, err := fmt.Fprint(w, prefix)
    if err != nil {
        panic(err)
    }
    // print the main message with format
    _, err = fmt.Fprintf(w, format, a...)
    if err != nil {
        panic(err)
    }
    // append extra new line
    _, err = fmt.Fprint(w, "\n")
    if err != nil {
        panic(err)
    }
}

// IsDirectory checks whether path is a directory
func IsDirectory(path string) (bool, error) {
    stat, err := os.Stat(path)
    if err != nil {
        return false, err
    }
    if !stat.IsDir() {
        return false, nil
    }
    return true, nil
}

// IsFile checks whether path is a file
func IsFile(path string) (bool, error) {
    stat, err := os.Stat(path)
    if err != nil {
        return false, err
    }
    if stat.IsDir() {
        return false, nil
    }
    return true, nil
}
