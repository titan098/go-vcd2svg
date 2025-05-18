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
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrawSVG_WireSignals(t *testing.T) {
	vcdData := &VcdData{
		Sim: map[uint64]map[string]string{
			0: {"clk": "0", "rst": "1"},
			1: {"clk": "1", "rst": "1"},
			2: {"clk": "0", "rst": "0"},
			3: {"clk": "1", "rst": "0"},
		},
		Decl: map[string]string{
			"!": "clk",
			"#": "rst",
		},
		Signals: []string{"clk", "rst"},
	}

	svgBytes := DrawSVG(vcdData)
	svgStr := string(svgBytes)

	assert.Contains(t, svgStr, "<svg")
	assert.Contains(t, svgStr, "clk")
	assert.Contains(t, svgStr, "rst")
}

func TestDrawSVG_BusSignal(t *testing.T) {
	vcdData := &VcdData{
		Sim: map[uint64]map[string]string{
			0: {"bus": "b1010"},
			1: {"bus": "b1010"},
			2: {"bus": "b1111"},
			3: {"bus": "b1111"},
		},
		Decl: map[string]string{
			"!": "bus",
		},
		Signals: []string{"bus"},
	}
	svgBytes := DrawSVG(vcdData)
	svgStr := string(svgBytes)

	assert.Contains(t, svgStr, "<svg")
	assert.Contains(t, svgStr, "b1010")
	assert.NotContains(t, svgStr, "0xAA")
}

func TestDrawSVG_ValidSVG(t *testing.T) {
	vcdData := &VcdData{
		Sim: map[uint64]map[string]string{
			0: {"sig": "0"},
			1: {"sig": "1"},
		},
		Decl: map[string]string{
			"!": "sig",
		},
		Signals: []string{"sig"},
	}

	svgBytes := DrawSVG(vcdData)

	// Parse SVG output as XML
	decoder := xml.NewDecoder(bytes.NewReader(svgBytes))
	foundSVG := false
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch el := token.(type) {
		case xml.StartElement:
			if el.Name.Local == "svg" {
				foundSVG = true
			}
		}
	}
	if !foundSVG {
		t.Errorf("SVG output does not appear to be valid XML or missing <svg>")
	}
}
