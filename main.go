package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Defines a list of valid filetypes for standard Source Engine directories
var strictFiletypeStructure = map[string][]string{
	"cfg" : {".cfg", ".txt"},
	"maps" : {".txt"},
	"materials" : {".vmt", ".vtf"},
	"models" : {".mdl", ".phy", ".vtx", ".vvd"},
	"particles" : {".pcf", ".txt"},
	"resources" : {".dds", ".txt"},
	"scripts" : {".lua", ".nut"},
	"sound" : {".mp3", ".txt", ".wav"},
}

/**
	Iterate through all files and child directories to build a list of files
 */
func ParseDirectory(path string, directory string, filesToPack []string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// Recursively traverse directory
	for _, f := range files {
		if f.IsDir() {
			filesToPack = ParseDirectory(path+f.Name() + "/", directory+f.Name()+"/", filesToPack)
		} else {
			filesToPack = append(filesToPack, directory+f.Name())
		}
	}

	return filesToPack
}

/**
	Create a new output file.
 */
func CreateOutputFile(filename string) *os.File {
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

/**
	Validates each file list entry before instructing to write
 */
func WriteFile(file *os.File, absoluteBasePath string, filesToPack []string, useStrict bool) {
	// Write all lines to file
	writer := bufio.NewWriter(file)

	for _, relativePath := range filesToPack {
		if ShouldDiscardFile(relativePath, useStrict) {
			continue
		}
		WriteEntry(writer, relativePath)
		WriteEntry(writer, absoluteBasePath+relativePath)
	}

	writer.Flush()
}

/**
	Write a single entry to file
 */
func WriteEntry(writer *bufio.Writer, line string) {
	writer.WriteString(line + "\n")
	fmt.Println("Add: " + line)
}

/**
	Determines whether an entry is suitable to be written
 */
func ShouldDiscardFile(filename string, useStrict bool) bool {
	// Always ignore
	if strings.HasPrefix(filename, ".") {
		return true
	}

	baseDirectory := strings.Split(filename, "/")

	// Treat all files as valid at top-level
	if baseDirectory[0] == filename {
		return false
	}

	// Lookup if current directory has extension restrictions
	if useStrict == true {
		if _, ok := strictFiletypeStructure[baseDirectory[0]]; ok {
			for _, validExtension := range strictFiletypeStructure[baseDirectory[0]] {
				if filepath.Ext(filename) == validExtension {
					return false
				}
			}
			fmt.Println("Strict Ignored: " + filename + ". Found in directory: " + baseDirectory[0])
			return true
		}
	}

	return false
}

func main() {
	// Intro
	fmt.Println("Galaco/DormantLemon's Bspzip Traversal tool")
	fmt.Println("Generate bspzip filelist's from any directory")

	// Parse flags
	targetFlag := flag.String("target", "", "Directory to generate filelist from.")
	outputFlag := flag.String("output", "", "Output filelist path/filename")
	useStrictFlag := flag.Bool("strict", true, "Ignore all unexpected filetypes within each directory")
	flag.Parse()

	// Assert valid usage
	if *targetFlag == "" || *outputFlag == "" {
		fmt.Println("Missing/invalid parameters specified.")
		log.Fatal("Correct Usage: bspzip-traverser.exe -target=<target_directory> -output=<path/filename> [-strict]")
	}

	fmt.Println(*targetFlag)

	targetDir, err := filepath.Abs(*targetFlag)
	if err != nil {
		log.Fatal("Specified target is not a valid directory")
	}
	outputFile := *outputFlag

	//Validate flags
	if targetDir[len(targetDir)-1:] != "/" {
		targetDir = targetDir + "/"
	}

	// Construct filelist
	fmt.Println("Parsing target directory...")
	files := ParseDirectory(targetDir, "", []string{})

	//Fetch top-level directory absolute path
	fmt.Println("Creating output file...")
	output := CreateOutputFile(outputFile)
	fmt.Println("Output file created!")

	// Write file
	fmt.Println("Writing filelist...")
	WriteFile(output, targetDir, files, *useStrictFlag)

	// Report success
	fmt.Println("Completed without issues")
}
