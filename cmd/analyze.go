// Copyright © 2019 Rodney Rodriguez
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/rodneyxr/ffatoolkit/ffa"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var fileTypeFlag string
var filepathFlag string
var resultsDir string

// analyzeCmd represents the list command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze a dockerfile or directory full of dockerfiles",
	Run: func(cmd *cobra.Command, args []string) {
		var files []string

		// Stat the file
		info, err := os.Stat(filepathFlag)
		if err != nil {
			panic(err)
		}

		if info.IsDir() {
			// If the file is a directory, add all files to the files list
			if err := filepath.Walk(filepathFlag, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, path)
				}
				return err
			}); err != nil {
				panic(err)
			}
		} else {
			// if it is not a directory, the file will be the only one in the list
			abs, _ := filepath.Abs(filepathFlag)
			files = append(files, abs)
		}

		// Create the results directory
		_ = os.Mkdir(resultsDir, os.ModeDir)

		for i, filename := range files {
			fmt.Printf("%d: %s\n", i, filename)

			// Read the file data
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Println(err)
				continue
			}

			var ffaScript []string

			switch fileTypeFlag {
			case "docker":
				// Parse the Dockerfile
				commandList, err := ffa.ExtractAllCommandsFromDockerfile(string(data))
				if err != nil {
					log.Print(err)
				}

				for _, cmd := range commandList {
					switch cmd.Cmd {
					case "run":
						results, err := ffa.AnalyzeShellCommand(strings.Join(cmd.Value, " "))
						if err != nil {
							log.Println(err)
							continue
						}
						ffaScript = append(ffaScript, results...)
						break
					case "workdir":
						ffaScript = append(ffaScript, "cd '"+cmd.Value[0]+"';")
						break
					case "copy":
						if len(cmd.Value) == 2 {
							ffaScript = append(ffaScript, "cp '"+cmd.Value[0]+"' '"+cmd.Value[1]+"';")
						}
						break
					}
				}
				break
			case "shell":
				results, err := ffa.AnalyzeShellCommand(string(data))
				if err != nil {
					log.Println(err)
					break
				}
				ffaScript = append(ffaScript, results...)
				break
			default:
				log.Fatal("unsupported file type")
			}

			// Save the ffa script to a file
			ffaFilename := filepath.Join(resultsDir, filepath.Base(filename)+".ffa")
			ffaScriptData := []byte(strings.Join(ffaScript, "\n"))
			if err = ioutil.WriteFile(ffaFilename, ffaScriptData, os.ModePerm); err != nil {
				log.Print(err)
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVar(&fileTypeFlag, "type", "", "type of file to analyze (shell or docker)")
	analyzeCmd.Flags().StringVar(&filepathFlag, "filepath", "", "path to file or directory to analyze")
	analyzeCmd.Flags().StringVar(&resultsDir, "results", "results", "directory to save results")
}