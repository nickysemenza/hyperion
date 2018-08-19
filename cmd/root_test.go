package cmd

import "testing"

func TestRootCmd(t *testing.T) {
	rootCmd.SetArgs([]string{"server"})
}
