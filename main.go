package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: gitCaseModifier <repo_path> <csv_file>")
		os.Exit(1)
	}

	repoPath := os.Args[1]
	mdPath := os.Args[2]

	file, err := os.Open(mdPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var records [][]string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		records = append(records, []string{parts[0], parts[1]})
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

       cmdLs := exec.Command("git", "ls-files")
       cmdLs.Dir = repoPath
       outLs, err := cmdLs.Output()
       if err != nil {
	       log.Fatalf("git ls-files failed: %v", err)
       }
       fileMap := make(map[string]string)
       scannerLs := bufio.NewScanner(strings.NewReader(string(outLs)))
       for scannerLs.Scan() {
	       rel := scannerLs.Text()
	       key := strings.ToLower(filepath.ToSlash(rel))
	       fileMap[key] = rel
       }
       if err := scannerLs.Err(); err != nil {
	       log.Fatalf("Failed to read ls-files: %v", err)
       }

       fmt.Println("--- fileMap keys (recognized files) ---")
       for k := range fileMap {
	       fmt.Println(k)
       }
       fmt.Println("--- end ---")

	var didRename bool
	for _, rec := range records {
	       if len(rec) < 2 {
		       continue
	       }
	       oldName := rec[0]
	       newName := rec[1]
	       oldKey := strings.ToLower(filepath.ToSlash(oldName))
	       actualOld, ok := fileMap[oldKey]
	       if !ok {
		       fmt.Printf("File not found: %s\n", oldName)
		       continue
	       }
	       if strings.EqualFold(actualOld, newName) {
		       tmpName := actualOld + ".__tmp__"
		       cmd1 := exec.Command("git", "mv", actualOld, tmpName)
		       cmd1.Dir = repoPath
		       out1, err1 := cmd1.CombinedOutput()
		       if err1 != nil {
			       fmt.Printf("git mv failed (temp): %s -> %s: %v\n%s\n", actualOld, tmpName, err1, string(out1))
			       continue
		       }
		       cmd2 := exec.Command("git", "mv", tmpName, newName)
		       cmd2.Dir = repoPath
		       out2, err2 := cmd2.CombinedOutput()
		       if err2 != nil {
			       fmt.Printf("git mv failed (target): %s -> %s: %v\n%s\n", tmpName, newName, err2, string(out2))
			       continue
		       }
		       fmt.Printf("rename success: %s -> %s -> %s\n", actualOld, tmpName, newName)
		       didRename = true
		       continue
	       }
	       cmd := exec.Command("git", "mv", actualOld, newName)
	       cmd.Dir = repoPath
	       out, err := cmd.CombinedOutput()
	       if err != nil {
		       fmt.Printf("git mv failed: %s -> %s: %v\n%s\n", actualOld, newName, err, string(out))
	       } else {
		       fmt.Printf("git mv success: %s -> %s\n", actualOld, newName)
		       didRename = true
	       }
       }

       if didRename {
	       cmdAdd := exec.Command("git", "add", "-A")
	       cmdAdd.Dir = repoPath
	       outAdd, errAdd := cmdAdd.CombinedOutput()
	       if errAdd != nil {
		       fmt.Printf("git add failed: %v\n%s\n", errAdd, string(outAdd))
		       return
	       }
	       cmdCommit := exec.Command("git", "commit", "-m", "case rename")
	       cmdCommit.Dir = repoPath
	       outCommit, errCommit := cmdCommit.CombinedOutput()
	       if errCommit != nil {
		       fmt.Printf("git commit failed: %v\n%s\n", errCommit, string(outCommit))
		       return
	       }
	       fmt.Println("Committed all renames: case rename")
       }
}
