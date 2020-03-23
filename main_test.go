package main

import (
	"strings"
	"testing"
)

const (
	alpha = "A\nB\nC\nD\nE\nF\nG\nH\nI\nJ\nK\nL\nM\nN\nO\nP\nQ\nR\nS\nT\nU\nV\nW\nX\nY\nZ\n"
	num   = "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n"
)

func TestMain(t *testing.T) {
	tests := []struct {
		in   string
		out  string
		exit int
	}{
		{"-v", version + "\n", 0},
		{"-i testdata/alpha.txt .", alpha, 0},
		{"-i testdata/num.txt .", num, 0},
		{"-i testdata/alpha.txt -i testdata/num.txt .", alpha + num, 0},
		{"-i testdata/alpha.txt -i testdata/num.txt [A-Z]", alpha, 0},
		{"-i testdata/alpha.txt -i testdata/num.txt [0-9]", num, 0},
		{"-i testdata/alpha.txt -i testdata/num.txt N", "N\n", 0},
		{"-i testdata/alpha.txt -i testdata/num.txt [a-z]", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var w strings.Builder
			exit := realMain(strings.Split(tt.in, " "), &w)
			if exit != tt.exit {
				t.Fatalf("exit: got %d, want %d", exit, tt.exit)
			}
			got := w.String() // strings.TrimSpace(w.String())
			if tt.out != got {
				t.Errorf("got\n%s\nwant\n%s", got, tt.out)
			}
		})
	}
}
