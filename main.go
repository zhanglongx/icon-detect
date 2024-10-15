package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/zhanglongx/icon-detect/pkg/detect"
	"github.com/zhanglongx/icon-detect/pkg/logfile"
	"github.com/zhanglongx/icon-detect/pkg/uri"
)

const (
	APPNAME = "icon-detect"
	TITLE   = "icon changes detected"
	VERSION = "1.1.0"

	PROGRAMTOKILL = "TOTALCMD64.EXE"
)

func main() {
	optVersion := flag.Bool("v", false, "show version")
	optBackUp := flag.Bool("b", false, "write to backup file")
	optUnReg := flag.Bool("u", false, "unregister URI scheme")
	optReg := flag.Bool("r", false, "register URI scheme")

	flag.Parse()

	err := logfile.InitLog(APPNAME + ".log")
	if err != nil {
		fmt.Printf("error initializing log file: %s\n", err)
		os.Exit(1)
	}

	defer logfile.DeInitLog()

	if *optVersion {
		fmt.Printf("%s %s", APPNAME, VERSION)
		return
	}

	if *optReg {
		exe, err := getCurrentExecutablePath()
		if err != nil {
			logfile.Fatal(err)
		}

		if uri.IsURISchemeRegistered(APPNAME) {
			logfile.Println("URI scheme already registered, re-registering")
			err = uri.UnRegisterRIScheme(APPNAME)
			if err != nil {
				logfile.Fatal(err)
			}
		}

		err = uri.RegisterURIScheme(APPNAME, APPNAME, exe)
		if err != nil {
			logfile.Fatal(err)
		}

		logfile.Println("URI scheme registered")

		return
	}

	if *optUnReg {
		err := uri.UnRegisterRIScheme(APPNAME)
		if err != nil {
			logfile.Fatal(err)
		}

		logfile.Println("URI scheme unregistered")

		return
	}

	if args := flag.Args(); len(args) > 0 {
		if len(args) != 1 {
			logfile.Fatal("invalid arguments")
		}

		uriParts := strings.Split(args[0], "://")
		if len(uriParts) < 2 {
			logfile.Fatal("invalid URI")
		}
		exe := strings.Split(uriParts[1], "/")[0]

		if exe != PROGRAMTOKILL {
			logfile.Fatal("should be " + PROGRAMTOKILL +
				" in " + args[0])
		}

		err := reStartProcess(PROGRAMTOKILL)
		if err != nil {
			logfile.Fatal(err)
		}

		logfile.Println("restarted", PROGRAMTOKILL)

		return
	}

	i := detect.NewIconDetect()

	isChanged, err := i.Detect()
	if err != nil {
		logfile.Fatal(err)
	}

	if isChanged {
		if *optBackUp {
			timeStr := time.Now().Format("20060102150405")
			fileName := "backup_" + timeStr + ".reg"

			err = i.WriteBackup(fileName)
			if err != nil {
				logfile.Fatal(err)
			}
		}

		err = i.Fix()
		if err != nil {
			logfile.Fatal(err)
		}

		scheme := ""
		if uri.IsURISchemeRegistered(APPNAME) {
			scheme = APPNAME + "://" + PROGRAMTOKILL
		}

		if err := detect.PushNotify(APPNAME, TITLE,
			scheme); err != nil {
			logfile.Fatal(err)
		}
	}
}

func getCurrentExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}

	return exePath, nil
}

func reStartProcess(name string) error {
	exe, err := detect.GetProcessPathByName(name)
	if err != nil {
		return err
	}

	err = detect.CloseProcessWindow(name)
	if err != nil {
		return err
	}

	cmd := exec.Command("cmd.exe", "/C", "start", "/MAX", exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
