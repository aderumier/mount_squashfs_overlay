package main

// checkWinFsp is a no-op. On Windows, go-winfsp handles DLL loading
// and reports errors directly. On Linux, cgofuse uses libfuse.
func checkWinFsp() error {
	return nil
}
