/*
Copyright © 2025 David Ellefsen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/titan098/go-vcd2svg/waveform"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a VCD (Value Change Dump) file to an SVG",
	Long: `Converts a VCD (Value Change Dump) file to an SVG diagram.
	
Example:
go-vcd2svg convert -i input.vcd -o output.svg`,
	Run: func(cmd *cobra.Command, args []string) {
		input := cmd.Flags().Lookup("input").Value.String()
		output := cmd.Flags().Lookup("output").Value.String()

		// check if the input exists
		if !fileExists(input) {
			fmt.Println("File does not exist:", args[0])
			os.Exit(1)
		}

		// check if the output exists
		if output != "" && fileExists(output) {
			fmt.Println("File already exists:", output)
			os.Exit(1)
		}

		// generate the SVG
		outBytes, err := waveform.SvgFromFile(input)
		if err != nil {
			fmt.Printf("Error generating SVG: %s\n", err.Error())
		}

		// write the file to the specified file
		if output != "" {
			err := os.WriteFile(output, outBytes, 0644)
			if err != nil {
				fmt.Printf("Error writing to output file: %s\n", err.Error())
				os.Exit(1)
			}
		} else {
			// write the svg output to the console if no output is specified
			fmt.Println(string(outBytes))
		}
	},
}

func fileExists(filename string) bool {
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringP("input", "i", "", "Input VCD file path")
	convertCmd.Flags().StringP("output", "o", "", "Output SVG file path")
	convertCmd.MarkFlagRequired("input")
}
