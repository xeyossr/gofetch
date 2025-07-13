package main

import (
	"os"
	"os/user"
	"os/exec"
	"bufio"
	"path/filepath"
	"strings"
	"fmt"
	"strconv"
	"bytes"
	"sync"
)

type Colors struct {
	Red string
	Yellow string
	Green string
	Cyan string
	Blue string
	Magenta string
	Reset string
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func countLines(name string, args ...string) (int) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	outStr := string(out.String())

	lines := strings.Split(strings.TrimSpace(outStr), "\n")
	return len(lines)
}

func hname() string {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return hostName
}

func username() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	return currentUser.Username
}

func getOsName() string {
	f, err := os.Open(filepath.Join("/etc", "os-release"))
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
			if strings.HasPrefix(line, "NAME=") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					os := strings.Trim(parts[1], `"'`)
					return os
				}
				return "unknown"
			}
	}

	return "unknown"
}

func getKernel() string {
	cmd := exec.Command("uname", "-rs")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(out.String()))
}

func getPkgs() (string) {
	pkgCount := 0
	if fileExists("/usr/bin/pacman") {
		pkgCount = countLines("pacman", "-Qq")
	} else if fileExists("/usr/bin/dpkg") {
			pkgCount = countLines("dpkg-query", "-f", ".", "-W")
	} else if fileExists("/usr/bin/rpm") {
			pkgCount = countLines("rpm", "-qa")
	} else if fileExists("/sbin/apk") {
			pkgCount = countLines("apk", "info")
	}

	return fmt.Sprintf("%d", pkgCount)
}

func getShell() string {
	shellPath := os.Getenv("SHELL")
	return filepath.Base(shellPath)
}

func getTerm() string {
	return os.Getenv("TERM")
}

func getCurrDesktop() string {
	return os.Getenv("XDG_CURRENT_DESKTOP")
}

func getUptime() string {
	f, err := os.ReadFile(filepath.Join("/proc", "uptime"))
	if err != nil {
		panic(err)
	}
	
	uptime := strings.Split(string(f), " ")
	uptimeFloat, err := strconv.ParseFloat(uptime[0], 64)
	if err != nil {
		panic(err)
	}

	hours := int(uptimeFloat) / 3600
	minutes := (int(uptimeFloat) % 3600) / 60

	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func main() {
	colors := Colors{
		Red: "\033[31m",
		Yellow: "\033[33m",
		Green: "\033[32m",
		Cyan: "\033[36m",
		Blue: "\033[34m",
		Magenta: "\033[35m",
		Reset: "\033[0m",
	}

	var (
		osName string
		hName string
		userName string
		kernel string
		pkgsCount string
		shell string
		term string
		wm string
		uptime string
		wg sync.WaitGroup
	)
	
	wg.Add(9)
	
	go func() {
		defer wg.Done()
		osName = getOsName()
	}()

	go func() {
		defer wg.Done()
		hName = hname()
	}()

	go func() {
		defer wg.Done()
		userName = username()
	}()

	go func() {
		defer wg.Done()
		kernel = getKernel()
	}()
	
	go func() {
		defer wg.Done()
		pkgsCount = getPkgs()
	}()
	
	go func() {
		defer wg.Done()
		shell = getShell()
	}()
	
	go func() {
		defer wg.Done()
		term = getTerm()
	}()
	
	go func() {
		defer wg.Done()
		wm = getCurrDesktop()
	}()

	go func() {
		defer wg.Done()
		uptime = getUptime()
	}()

	wg.Wait()
	
	fmt.Printf("%s%s%s@%s%s%s\n", colors.Red, userName, colors.Reset, colors.Yellow, hName, colors.Reset)
	
	fmt.Printf("%sOS     >%s%s\n", colors.Green, colors.Reset, osName)
	fmt.Printf("%sKernel >%s%s\n", colors.Cyan, colors.Reset, kernel)
	fmt.Printf("%sPkgs   >%s%s\n", colors.Red, colors.Reset, pkgsCount)
	fmt.Printf("%sShell  >%s%s\n", colors.Magenta, colors.Reset, shell)
	fmt.Printf("%sTerm   >%s%s\n", colors.Blue, colors.Reset, term)
	fmt.Printf("%sWm     >%s%s\n", colors.Yellow, colors.Reset, wm)
	fmt.Printf("%sUptime >%s%s\n", colors.Green, colors.Reset, uptime)

}
