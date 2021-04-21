package main

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

var (
	fileTemplate = "{name}_{os}_{arch}"
)

// compileOpts represents a compile job's options
type compileOpts struct {
	os string
	arch string
	args []string
	env []string
}

// finState is a finished binary state
type finState struct {
	os string
	arch string
	filename string
	err      error
}

// compiler is a goroutine which receives options off of i, and outputs to o
func compiler(i chan compileOpts, o chan finState) {
	for opt := range i {
		var fs finState


		// Replace -o FILENAME with -o FILENAME_os_arch
		for i, v := range opt.args {
			if v == "-o" {
				replacer := strings.NewReplacer("{name}", opt.args[i+1], "{os}", opt.os, "{arch}", opt.arch)
				fs.filename = replacer.Replace(fileTemplate)

				if opt.os == "windows" {
					fs.filename += ".exe"
				}

				opt.args[i+1] = fs.filename
				break
			}
		}

		if strings.Contains(fs.filename, "/") {
			err := os.MkdirAll(path.Dir(fs.filename), 0655)

			if err != nil {
				fs.err = err
				o <- fs
				continue
			}
		}

		args := append([]string{"build"}, opt.args...)

		cmd := exec.Command("go", args...)

		log.WithField("args", args).Debug("Running compiler")

		cmd.Env = append(os.Environ(),
			"GOOS="+opt.os,
			"GOARCH="+opt.arch,
		)

		output, err := cmd.CombinedOutput()

		if err != nil {
			fs.err = fmt.Errorf("Failed to compile %s: %s: %s\n", fs.filename, bytes.TrimSpace(output), err)
		}

		o <- fs
	}
}

// supportedSystems runs `go tool dist list` to get supported operating systems + architectures
func supportedSystems(newEnv []string) (systemList, error) {
	cmd := exec.Command("go", "tool", "dist", "list")

	cmd.Env = newEnv

	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error:", string(out))
		return nil, err
	}

	list := make(systemList, 0)

	scanner := bufio.NewScanner(bytes.NewReader(out))

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "/")

		s := list.Find(parts[0])

		log.WithFields(log.Fields{
			"os": parts[0],
			"arch": parts[1],
		}).Debug("Adding OS/Arch")

		if s == nil {
			s = &system{
				name: parts[0],
				archs: []string{parts[1]},
			}

			list = append(list, s)
		} else {
			s.archs = append(s.archs, parts[1])
		}
	}

	return list, nil
}

func main() {
	parallel := runtime.NumCPU()

	newEnv := make([]string, 0)

	env := os.Environ()

	for _, v := range env {
		if strings.HasPrefix(v, "GOOS") || strings.HasPrefix(v, "GOARCH") {
			continue
		}

		newEnv = append(newEnv, v)
	}

	if binaryTpl := os.Getenv("GOBINARY"); binaryTpl != "" {
		fileTemplate = binaryTpl
	}

	systems, err := supportedSystems(newEnv)

	if err != nil {
		log.WithError(err).Fatalln("Error running go tool list")
	}

	// Supported = `go tool dist list`

	// Remove GOARCH/GOOS

	i := make(chan compileOpts, 512)
	o := make(chan finState, 512)
	defer close(i)
	defer close(o)

	log.WithField("threads", parallel).Info("Starting compiler routines")

	for x := 0; x < parallel; x++ {
		go compiler(i, o)
	}

	args := os.Args[1:]

	goOsStr := strings.TrimSpace(os.Getenv("GOOS"))

	if goOsStr == "" {
		log.Fatalln("Expected one or more specified GOOS")
	}

	goOs := strings.Split(goOsStr, ",")

	goArchStr := strings.TrimSpace(os.Getenv("GOARCH"))

	goArch := strings.Split(goArchStr, ",")

	if goArchStr == "" {
		goArch = []string{}
	}

	var numSys int

	for _, operatingSystem := range goOs {
		s := systems.Find(operatingSystem)

		if s == nil {
			log.WithField("os", operatingSystem).Warn("OS does not exist")
			continue
		}

		compileArches := goArch

		if len(compileArches) == 0 {
			compileArches = s.archs
		}

		for _, arch := range compileArches {
			if !s.HasArch(arch) {
				continue
			}

			numSys++

			newArgs := make([]string, len(args))

			copy(newArgs, args)

			i <- compileOpts{os: operatingSystem, arch: arch, env: env, args: newArgs}
		}
	}

	for x := 0; x < numSys; x++ {
		fs := <-o

		if fs.err != nil {
			log.WithFields(log.Fields{
				"os": fs.os,
				"arch": fs.arch,
				"error": fs.err,
			}).Warn("Compilation failed")
		} else {
			log.Printf("Finished compiling %s", fs.filename)
		}
	}
}