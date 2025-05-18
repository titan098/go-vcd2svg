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
package waveform

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"sort"

	"github.com/filmil/go-vcd-parser/vcd"
)

type VcdData struct {
	Sim    map[uint64]map[string]string
	Decl    map[string]string
	Signals    []string
}

// ParseVCD parses a VCD  file from the provided bytes.Reader.
// The 'name' parameter is used to identify the file (may be used in errors).
// It returns a pointer to a VcdData struct containing the parsed simulation data,
// or an error if parsing fails.
func ParseVCD(reader *bytes.Reader, name string) (*VcdData, error) {
	parser := vcd.NewParser[vcd.File]()
	ast, err := parser.Parse(name, reader)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}
	return ProcessVcd(ast), nil
}

// ParseVcdAndGenerateSvg parses a VCD file from the provided bytes.Reader with the given name,
// and generates an SVG waveform representation of the signal data.
// It returns the generated SVG as a []byte slice, or an error if parsing fails.
func ParseVcdAndGenerateSvg(reader *bytes.Reader, name string) ([]byte, error) {
	vcdData, err := ParseVCD(reader, name)
	if err != nil {
		return nil, err
	}
	return DrawSVG(vcdData), nil
}

// SvgFromFile reads a VCD (Value Change Dump) file from the given filename,
// parses its contents, and generates an SVG waveform representation.
// Returns the SVG as a []byte slice, or an error if the file cannot be read or parsed.
func SvgFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Read file into memory (for *bytes.Reader compatibility)
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return ParseVcdAndGenerateSvg(bytes.NewReader(content), filename)
}

// SvgFromBytes parses VCD data provided as a byte slice, and generates
// an SVG waveform representation. Returns the SVG as a []byte slice,
// or an error if parsing fails.
func SvgFromBytes(content []byte) ([]byte, error) {
	return ParseVcdAndGenerateSvg(bytes.NewReader(content), "noname.vcd")
}

// processVcd processes a parsed VCD AST (Abstract Syntax Tree) and returns a
// Structure to represent the signal changes over time.
func ProcessVcd(ast *vcd.File) *VcdData {
	vcdData := VcdData{
		Sim: map[uint64]map[string]string{
			0: {},
		},
		Decl: map[string]string{},
	}

	// Determine the signal names from the signal codes
	// keep track of the scope for the signals
	scope := []string{""}
	for _, v1 := range ast.DeclarationCommand {
		if v1.Scope != nil {
			scope = append(scope, fmt.Sprintf("%s ", v1.Scope.Id))
		}
		if v1.Upscope != nil {
			scope = scope[0 : len(scope)-1]
		}
		if v1.Var != nil {
			vcdData.Decl[v1.Var.Code] = fmt.Sprintf("%s%s", scope[len(scope)-1], v1.Var.Id.Name)
		}
	}

	// for each simulation time period keep track of which signals changes
	// we keep track of every signal at each time period so that it easier
	// render
	var s uint64
	for _, d := range ast.SimulationCommand {
		if d.SimulationTime != nil {
			s = d.SimulationTime.Value()
			_, ok := vcdData.Sim[s]
			if !ok {
				vcdData.Sim[s] = maps.Clone(vcdData.Sim[s-1])
			}
		}

		if d.ValueChange != nil {
			if d.ValueChange.ScalarValueChange != nil {
				vcdData.Sim[s][vcdData.Decl[d.ValueChange.ScalarValueChange.GetIdCode()]] = d.ValueChange.ScalarValueChange.GetValue()
			} else if d.ValueChange.VectorValueChange != nil {
				vcdData.Sim[s][vcdData.Decl[d.ValueChange.VectorValueChange.GetCode()]] = d.ValueChange.VectorValueChange.GetValue()
			}
		}
	}

	// Collect the signal names so they are consistent
	seen := map[string]bool{}
	for _, step := range vcdData.Sim {
		for sig := range step {
			if !seen[sig] {
				vcdData.Signals = append(vcdData.Signals, sig)
				seen[sig] = true
			}
		}
	}
	sort.Strings(vcdData.Signals)
	return &vcdData
}
