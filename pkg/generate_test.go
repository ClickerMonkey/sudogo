package sudogo

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	g := Classic.Generator()

	start := time.Now()
	p := g.Generate(128)

	if p == nil {
		t.Fatalf("Failed to generate a puzzle in 128 tries")
	} else {
		p.Print()
	}

	duration := time.Since(start)
	fmt.Printf("TestGenerate in %s\n", duration)
}
