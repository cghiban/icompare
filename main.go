package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type entry struct {
	Name  string
	IsDir bool
}

func ilsExists() error {
	cmd := exec.Command("which", "ils")
	//cmd.Stdout = nil
	cmd.Stderr = os.Stdout
	//cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// terminated by Control-C so ignoring
			if exiterr.ExitCode() == 130 {
				return nil
			}
		}

		return err
	}

	return nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ils(remote string) ([]entry, error) {

	entries := []entry{}

	cmd := exec.Command("ils", remote)

	buff, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	//cmd.Stdout = buff
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	//cmd.Stdin = os.Stdin
	//fmt.Println(buff.Bytes())

	if err := cmd.Start(); err != nil {

		if exiterr, ok := err.(*exec.ExitError); ok {
			// terminated by Control-C so ignoring
			if exiterr.ExitCode() == 130 {
				return entries, err
			}
		}
		return entries, err
	}

	scanner := bufio.NewScanner(buff)
	var e entry
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			continue
		}
		//fmt.Println("line:", line, len(parts))
		if len(parts) == 3 {
			e.Name = parts[2]
		} else if parts[2] == "C-" {
			e.Name = path.Base(parts[3])
			e.IsDir = true
		}
		entries = append(entries, e)
	}
	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return entries, nil
	}

	// sort the entries
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	return entries, err
}

func ls(local string) ([]entry, error) {
	entries := []entry{}

	var e entry
	dirEntries, err := os.ReadDir(local)
	if err != nil {
		return entries, err
	}
	for _, de := range dirEntries {
		e = entry{de.Name(), de.IsDir()}
		entries = append(entries, e)
	}
	// sort the entries
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	return entries, nil
}

func main() {
	var localPath, remotePath string

	// set a and b as flag int vars
	flag.StringVar(&localPath, "l", "", "Local path")
	flag.StringVar(&remotePath, "r", "", "Remote path")

	// parse flags from command line
	flag.Parse()

	// check if local path exists
	if ok, err := exists(localPath); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	} else if !ok {
		log.Printf("Given local path %s does not exist!\n", localPath)
		os.Exit(1)
	}

	if err := ilsExists(); err != nil {
		log.Printf("ils icommand not found on this system: %s\n", err.Error())
		os.Exit(1)
	}

	// output
	//fmt.Println("+ local: ", localPath)
	localFiles, nil := ls(localPath)
	//fmt.Println(localFiles)
	//fmt.Println("+ remote: ", remotePath)

	remoteFiles, err := ils(remotePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	/*fmt.Println("----------")
	for _, fr := range remoteFiles {
		fType := " "
		if fr.IsDir {
			fType = "C"
		}
		fmt.Printf("\t[%s] %s\n", fType, fr.Name)
	}*/

	if diff := cmp.Diff(localFiles, remoteFiles); diff != "" {
		fmt.Printf("Oh no!!! Here's the diff:\n%s", diff)
	} else {
		//fmt.Printf("!! all's good ;)\n")
	}

}
