//go:build zstd && xz

package bench

import (
	"io"
	"io/fs"
	"os"
	"testing"

	caleb "github.com/CalebQ42/squashfs"
	karp "github.com/KarpelesLab/squashfs"
)

const testSquashfs = "../test.wsquashfs"

var benchFiles = []struct {
	name string
	path string
}{
	{"tiny_12B", ".update-timestamp"},
	{"small_17B", "game/autorun.cmd"},
	{"medium_27MB", "game/fnfmega_Data/level2"},
	{"large_97MB", "game/fnfmega_Data/level4.resS"},
	{"huge_504MB", "game/fnfmega_Data/sharedassets0.resource"},
}

type readers struct {
	kb *karp.Superblock
	cq caleb.Reader
}

func openSquash(tb testing.TB) readers {
	tb.Helper()
	f, err := os.Open(testSquashfs)
	if err != nil {
		tb.Skipf("test archive not found (%s): %v", testSquashfs, err)
	}
	kb, err := karp.New(f)
	if err != nil {
		tb.Fatalf("karp.New: %v", err)
	}
	cq, err := caleb.NewReader(f)
	if err != nil {
		tb.Fatalf("caleb.NewReader: %v", err)
	}
	return readers{kb: kb, cq: cq}
}

// ── CoW copy: CalebQ42 parallel WriteTo (current implementation) ─────────────
// WriteTo is used with a *os.File destination (implements io.WriterAt) which
// takes CalebQ42's parallel fast-path. io.Discard does NOT implement io.WriterAt
// and triggers the sequential path which deadlocks in v1.4.1 — use a temp file.

func BenchmarkCoW_Caleb_WriteTo(b *testing.B) {
	r := openSquash(b)
	tmp, err := os.CreateTemp("", "bench-caleb-*.tmp")
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { tmp.Close(); os.Remove(tmp.Name()) })

	for _, bf := range benchFiles {
		bf := bf
		b.Run(bf.name, func(b *testing.B) {
			fi, err := fs.Stat(r.kb, bf.path)
			if err != nil {
				b.Skipf("stat %s: %v", bf.path, err)
			}
			b.SetBytes(fi.Size())
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := tmp.Truncate(0); err != nil {
					b.Fatal(err)
				}
				if _, err := tmp.Seek(0, io.SeekStart); err != nil {
					b.Fatal(err)
				}
				f, err := r.cq.Open(bf.path)
				if err != nil {
					b.Fatal(err)
				}
				wt := f.(io.WriterTo)
				if _, err := wt.WriteTo(tmp); err != nil {
					b.Fatal(err)
				}
				f.Close()
			}
		})
	}
}

// ── CoW copy: KarpelesLab ReadAt 128KB aligned (previous implementation) ─────

func BenchmarkCoW_Karp_ReadAt(b *testing.B) {
	r := openSquash(b)
	for _, bf := range benchFiles {
		bf := bf
		b.Run(bf.name, func(b *testing.B) {
			fi, err := fs.Stat(r.kb, bf.path)
			if err != nil {
				b.Skipf("stat %s: %v", bf.path, err)
			}
			size := fi.Size()
			b.SetBytes(size)
			b.ResetTimer()
			buf := make([]byte, 128*1024)
			for i := 0; i < b.N; i++ {
				f, err := r.kb.Open(bf.path)
				if err != nil {
					b.Fatal(err)
				}
				ra := f.(io.ReaderAt)
				if _, err := io.CopyBuffer(io.Discard, io.NewSectionReader(ra, 0, size), buf); err != nil {
					b.Fatal(err)
				}
				f.Close()
			}
		})
	}
}

// ── FUSE Read: KarpelesLab single random-block ReadAt ────────────────────────

func BenchmarkFUSERead_Karp_RandomBlock(b *testing.B) {
	r := openSquash(b)
	for _, bf := range benchFiles {
		bf := bf
		b.Run(bf.name, func(b *testing.B) {
			fi, err := fs.Stat(r.kb, bf.path)
			if err != nil {
				b.Skipf("stat %s: %v", bf.path, err)
			}
			offset := fi.Size() / 2
			buf := make([]byte, 64*1024)
			b.SetBytes(int64(len(buf)))
			b.ResetTimer()

			f, err := r.kb.Open(bf.path)
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			ra := f.(io.ReaderAt)

			for i := 0; i < b.N; i++ {
				if _, err := ra.ReadAt(buf, offset); err != nil && err != io.EOF {
					b.Fatal(err)
				}
			}
		})
	}
}

// ── Stat ──────────────────────────────────────────────────────────────────────

func BenchmarkStat(b *testing.B) {
	r := openSquash(b)
	for _, bf := range benchFiles {
		bf := bf
		b.Run(bf.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := fs.Stat(r.kb, bf.path); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ── ReadDir ───────────────────────────────────────────────────────────────────

func BenchmarkReadDir(b *testing.B) {
	r := openSquash(b)
	dirs := []struct{ name, path string }{
		{"root", "."},
		{"game_data", "game/fnfmega_Data"},
		{"system32", "drive_c/windows/system32"},
	}
	for _, d := range dirs {
		d := d
		b.Run(d.name, func(b *testing.B) {
			entries, _ := fs.ReadDir(r.kb, d.path)
			b.SetBytes(int64(len(entries)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := fs.ReadDir(r.kb, d.path); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
