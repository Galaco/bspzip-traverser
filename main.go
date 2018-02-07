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

//@TODO
//extension whitelists
//output directory+filename arguments

/**
 * @brief      Recursive read directory, build filepaths
 *
 * @param      path         The path
 * @param      directory    The directory
 * @param      filesToPack  The files to pack
 *
 * @return     { description_of_the_return_value }
 */
func ParseDirectory(path string, directory string, filesToPack []string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// Recursively traverse directory
		if f.IsDir() {
			filesToPack = ParseDirectory(path+f.Name(), directory+f.Name()+"/", filesToPack)
		} else {
			filesToPack = append(filesToPack, directory+f.Name())
		}
	}

	return filesToPack
}

func CreateOutputFile(filename string) *os.File {
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

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

func WriteEntry(writer *bufio.Writer, line string) {
	writer.WriteString(line + "\n")
	fmt.Println("Add: " + line)
}

func ShouldDiscardFile(filename string, useStrict bool) bool {
	// Always ignore
	if strings.HasPrefix(filename, ".") {
		return true
	}

	baseDirectory := strings.Split(filename, "/")
	// We're in the base directory
	if baseDirectory[0] == filename {
		baseDirectory[0] = ""
	}

	//lookup valid types

	return false
}

/**
 * @brief      Create a bspzip filelist.txt from a directory
 *
 * @return     { description_of_the_return_value }
 */
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
