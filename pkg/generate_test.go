package sudogo

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	g := Classic.Generator()

	start := time.Now()
	p, attempts := g.Generate()

	if p == nil {
		t.Fatalf("Failed to generate a puzzle")
	} else {
		p.Print()
	}

	duration := time.Since(start)
	fmt.Printf("TestGenerate in %s after %d attempts.\n", duration, attempts)
}
