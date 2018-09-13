package main

import "testing"

func TestRootCmd(t *testing.T) {
	rootCmd.SetArgs([]string{"server"})
}
