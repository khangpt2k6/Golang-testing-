package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Supported language configurations
type LanguageConfig struct {
	Extension  string
	Executable string
	RunArgs    []string
}

var languageConfigs = map[string]LanguageConfig{
	"python": {
		Extension:  ".py",
		Executable: "python",
		RunArgs:    []string{},
	},
	"javascript": {
		Extension:  ".js",
		Executable: "node",
		RunArgs:    []string{},
	},
	"ruby": {
		Extension:  ".rb",
		Executable: "ruby",
		RunArgs:    []string{},
	},
	"shell": {
		Extension:  ".sh",
		Executable: "bash",
		RunArgs:    []string{},
	},
	"php": {
		Extension:  ".php",
		Executable: "php",
		RunArgs:    []string{},
	},
}

func main() {
	// Set up command-line flags
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runLang := runCmd.String("lang", "", "Language to run (python, javascript, ruby, shell, php)")
	runFile := runCmd.String("file", "", "File to execute")

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createLang := createCmd.String("lang", "", "Language to create script for (python, javascript, ruby, shell, php)")
	createFile := createCmd.String("file", "", "Filename to create (without extension)")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Check if any arguments were provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse the subcommand
	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		if *runLang == "" || *runFile == "" {
			fmt.Println("Error: both -lang and -file are required for run command")
			runCmd.PrintDefaults()
			os.Exit(1)
		}
		runScript(*runLang, *runFile)
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createLang == "" || *createFile == "" {
			fmt.Println("Error: both -lang and -file are required for create command")
			createCmd.PrintDefaults()
			os.Exit(1)
		}
		createScript(*createLang, *createFile)
	case "list":
		listCmd.Parse(os.Args[2:])
		listLanguages()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("MultiLang CLI - Run scripts in multiple languages")
	fmt.Println("\nUsage:")
	fmt.Println("  multilang run -lang <language> -file <filename>")
	fmt.Println("  multilang create -lang <language> -file <filename>")
	fmt.Println("  multilang list")
	fmt.Println("\nExample:")
	fmt.Println("  multilang run -lang python -file hello")
	fmt.Println("  multilang create -lang javascript -file new_script")
}

func runScript(lang, file string) {
	config, ok := languageConfigs[strings.ToLower(lang)]
	if !ok {
		fmt.Printf("Unsupported language: %s\n", lang)
		listLanguages()
		os.Exit(1)
	}

	// Add extension if not already included
	if !strings.HasSuffix(file, config.Extension) {
		file = file + config.Extension
	}

	// Check if file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", file)
		os.Exit(1)
	}

	// Prepare command
	args := append(config.RunArgs, file)
	cmd := exec.Command(config.Executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the script
	fmt.Printf("Running %s script: %s\n", lang, file)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing script: %v\n", err)
		os.Exit(1)
	}
}

func createScript(lang, file string) {
	config, ok := languageConfigs[strings.ToLower(lang)]
	if !ok {
		fmt.Printf("Unsupported language: %s\n", lang)
		listLanguages()
		os.Exit(1)
	}

	// Add extension if not already included
	if !strings.HasSuffix(file, config.Extension) {
		file = file + config.Extension
	}

	// Check if file already exists
	if _, err := os.Stat(file); err == nil {
		fmt.Printf("File '%s' already exists. Overwrite? (y/n): ", file)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled")
			os.Exit(0)
		}
	}

	// Create basic template content based on language
	var content string
	switch strings.ToLower(lang) {
	case "python":
		content = `#!/usr/bin/env python
# -*- coding: utf-8 -*-

def main():
    print("Hello from Python!")

if __name__ == "__main__":
    main()
`
	case "javascript":
		content = `#!/usr/bin/env node

function main() {
    console.log("Hello from JavaScript!");
}

main();
`
	case "ruby":
		content = `#!/usr/bin/env ruby

def main
  puts "Hello from Ruby!"
end

main
`
	case "shell":
		content = `#!/bin/bash

echo "Hello from Bash!"
`
	case "php":
		content = `<?php

function main() {
    echo "Hello from PHP!\n";
}

main();
`
	}

	// Write content to file
	err := ioutil.WriteFile(file, []byte(content), 0755)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}

	absPath, _ := filepath.Abs(file)
	fmt.Printf("Created %s script: %s\n", lang, absPath)
}

func listLanguages() {
	fmt.Println("Supported languages:")
	for lang, config := range languageConfigs {
		fmt.Printf("  - %s (extension: %s, executable: %s)\n", 
			lang, config.Extension, config.Executable)
	}
}