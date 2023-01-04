package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"math"
	"math/rand"
	"testing"

	"golang.org/x/exp/slices"
	"pgregory.net/rapid"
)

func TestILog2(t *testing.T) {
	checkWithValues(t, ILog2)
}

func TestILog2B(t *testing.T) {
	checkWithValues(t, ILog2B)
}

func TestLog2ByAvernar(t *testing.T) {
	checkWithValues(t, Log2ByAvernar)
}

func checkWithValues(t *testing.T, f func(uint64) int64) {
	for x := uint64(0); x < 100000; x++ {
		check(t, f, x)
	}
	for i := 1; i < 64; i++ {
		x := uint64(1) << i
		check(t, f, x-1)
		check(t, f, x)
		check(t, f, x+1)
	}
	check(t, f, math.MaxUint64)
}

func check(t *testing.T, f func(uint64) int64, x uint64) {
	got := f(x)
	want := Log2ByStdlib(x)
	if got != want {
		t.Errorf("result mismatch, x=%x, got=%d, want=%d", x, got, want)
	}
}

//go:embed ilog2_c.txt
var ilog2CImplOutput []byte

func TestILog2_compareToCImpl(t *testing.T) {
	r := bytes.NewReader(ilog2CImplOutput)
	var x uint64
	var want int64
	for {
		n, err := fmt.Fscanf(r, "%d %d\n", &x, &want)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if n == 0 && err == io.EOF {
			break
		}
		if n != 2 {
			t.Fatalf("unexpected Fscanf count: got=%d, want=2", n)
		}

		got := ILog2(x)
		if got != want {
			t.Errorf("reuslt mismatch, x=%d, got=%d, want=%d", x, got, want)
		}
	}
}

func TestDebugPow2Minus1Good(t *testing.T) {
	fDebug := func(n uint64) uint64 {
		fmt.Printf("fDebug start n=%#0b\n", n)
		n |= n >> 1
		fmt.Printf("fDebug #1    n=%#0b\n", n)
		n |= n >> 2
		fmt.Printf("fDebug #2    n=%#0b\n", n)
		n |= n >> 4
		fmt.Printf("fDebug #3    n=%#0b\n", n)
		n |= n >> 8
		fmt.Printf("fDebug #4    n=%#0b\n", n)
		n |= n >> 16
		fmt.Printf("fDebug #5    n=%#0b\n", n)
		n |= n >> 32
		fmt.Printf("fDebug final n=%#0b\n", n)
		return n
	}
	fDebug(0x80000000000)
}

func TestDebugPow2Minus1Bad(t *testing.T) {
	gDebug := func(n uint64) uint64 {
		fmt.Printf("gDebug start n=%#0b\n", n)
		n1 := n >> 1
		fmt.Printf("gDebug      n1=%#0b, n|n1=%#0b\n", n1, n|n1)
		n2 := n >> 2
		fmt.Printf("gDebug      n2=%#0b, n|n1|n2=%#0b\n", n2, n|n1|n2)
		n4 := n >> 4
		fmt.Printf("gDebug      n4=%#0b, n|n1|n2|n4=%#0b\n", n4, n|n1|n2|n4)
		n8 := n >> 8
		fmt.Printf("gDebug      n8=%#0b, n|n1|n2|n4|n8=%#0b\n", n8, n|n1|n2|n4|n8)
		n16 := n >> 16
		fmt.Printf("gDebug     n16=%#0b, n|n1|n2|n4|n8|n16=%#0b\n", n16, n|n1|n2|n4|n8|n16)
		n32 := n >> 32
		fmt.Printf("gDebug     n32=%#0b, n|n1|n2|n4|n8|n16|n32=%#0b\n", n32, n|n1|n2|n4|n8|n16|n32)
		return n | n1 | n2 | n4 | n8 | n16 | n32
	}
	gDebug(0x80000000000)
}

func TestBinAvernarDeBruijn(t *testing.T) {
	got := fmt.Sprintf("%#064b", 0x07EDD5E59A4E28C2)
	want := "0b0000011111101101110101011110010110011010010011100010100011000010"
	if got != want {
		t.Errorf("result mismatch, got=%s, want=%s", got, want)
	}
}

func TestLog2ByAvernarPropertyEqualToStdlib(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.Uint64Range(1, 0xffffffffffff4bff).Draw(t, "x")
		got := Log2ByAvernar(x)
		want := Log2ByStdlib(x)
		if got != want {
			t.Fatalf("log2 mismatch, input=%d, got=%d, want=%d", x, got, want)
		}
	})
}

func TestLog2ByAvernarNotEqual(t *testing.T) {
	for x := uint64(0xffffffffffff4bff + 1); ; x++ {
		got := Log2ByAvernar(x)
		want := Log2ByStdlib(x)
		if got == want {
			t.Fatalf("log2 unexpected match, input=%x, got=%d, want=%d", x, got, want)
			// } else {
			// t.Logf("log2 expected unmatch, input=%x, got=%d, want=%d", x, got, want)
		}
		if x == math.MaxUint64 {
			break
		}
	}
}

func TestBuildTableAvernar(t *testing.T) {
	got := buildTable(0x03f6eaf2cd271461)
	want := u8Table
	if !slices.Equal(got, want) {
		t.Errorf("table mismatch, got=%v, want=%v", got, want)
	}
}

func buildInputValues(seed int64, n int) []uint64 {
	rnd := rand.New(rand.NewSource(seed))
	v := make([]uint64, n)
	for i := 0; i < n; i++ {
		for {
			x := rnd.Uint64()
			if x != 0 {
				v[i] = x
				break
			}
		}
	}
	return v
}

const seed = 12345
const dataCount = 1000

var inputValues = buildInputValues(seed, dataCount)

func nop(sum int64) {}

func BenchmarkILog2(b *testing.B) {
	sum := int64(0)
	for i := 0; i < b.N; i++ {
		for _, x := range inputValues {
			sum += ILog2(x)
		}
	}
	nop(sum)
}

func BenchmarkILog2B(b *testing.B) {
	sum := int64(0)
	for i := 0; i < b.N; i++ {
		for _, x := range inputValues {
			sum += ILog2B(x)
		}
	}
	nop(sum)
}

func BenchmarkLogByAvernar(b *testing.B) {
	sum := int64(0)
	for i := 0; i < b.N; i++ {
		for _, x := range inputValues {
			sum += Log2ByAvernar(x)
		}
	}
	nop(sum)
}

func BenchmarkLogByAvernarU8(b *testing.B) {
	sum := int64(0)
	for i := 0; i < b.N; i++ {
		for _, x := range inputValues {
			sum += Log2ByAvernarU8(x)
		}
	}
	nop(sum)
}

func BenchmarkLogByStdlib(b *testing.B) {
	sum := int64(0)
	for i := 0; i < b.N; i++ {
		for _, x := range inputValues {
			sum += Log2ByStdlib(x)
		}
	}
	nop(sum)
}
