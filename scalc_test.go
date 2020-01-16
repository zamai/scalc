package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
	"testing"
)

const binaryName = "scalc"

func TestMain(m *testing.M) {
	cmd := exec.Command("make")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("could not make binary: %+v \n %s", err, string(out))
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestE2E(t *testing.T) {
	tests := []struct {
		name     string
		args     string
		expected string
	}{
		{"SUM 2 files", "[ SUM a.txt b.txt ]", "1\n2\n3\n4\n"},
		{"SUM 2 expr", "[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]", "1\n3\n4\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command(path.Join(dir, binaryName), strings.Split(tt.args, " ")...)
			output, err := cmd.CombinedOutput()
			actual := string(output)
			if err != nil {
				fmt.Println(actual)
				t.Fatal(err)
			}

			if actual != tt.expected {
				t.Fatalf("actual = %s, expected = %s", actual, tt.expected)
			}
		})
	}
}

func TestCalc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{"good calc: SUM 2 files", "[ SUM a.txt b.txt ]", []int{1, 2, 3, 4}},
		{"good calc: DIF files", "[ DIF a.txt b.txt ]", []int{1}},
		{"good calc: INT files", "[ INT a.txt b.txt ]", []int{2, 3}},
		{"good calc: SUM 1 exp + 1 file", "[ SUM [ SUM a.txt b.txt ] a.txt ]", []int{1, 2, 3, 4}},
		{"good calc: SUM 2 exp ", "[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]", []int{1, 3, 4}},
		{"good calc: DIF 1 file 2 exp ", "[ DIF a.txt [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]", []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := calc(strings.Split(tt.input, " "))

			if !reflect.DeepEqual(r, tt.expected) {
				t.FailNow()
			}
		})
	}

}
