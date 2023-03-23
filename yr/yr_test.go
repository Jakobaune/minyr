package yr_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/Jakobaune/minyr/yr"
)

func TestFileLineCount(t *testing.T) {
	file, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
		t.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if lineCount != 16756 {
		t.Errorf("unexpected number of lines in file: got %d, want %d", lineCount, 16756)
	}
}

func TestConvertLastField(t *testing.T) {
	input := "6"
	want := "42,8"

	got, err := yr.ConvertLastField(input)
	if err != nil {
		t.Errorf("convertLastField(%q) returned error: %v", input, err)
	}

	if got != want {
		t.Errorf("convertLastField(%q) = %q; want %q", input, got, want)
	}
}

func TestKonverter(t *testing.T) {
	// Set up input file
	inputFileName := "test_input.csv"
	inputFile, err := os.Create(inputFileName)
	if err != nil {
		t.Fatalf("Error creating input file: %v", err)
	}
	defer os.Remove(inputFileName)
	fmt.Fprintln(inputFile, "Kjevik;SN39040;18.03.2022 01:50;6")

	// Run Konverter function
	yr.Konverter()

	// Check output file
	outputFileName := "Resultat.txt"
	defer os.Remove(outputFileName)
	outputFile, err := os.Open(outputFileName)
	if err != nil {
		t.Fatalf("Error opening output file: %v", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(outputFile)
	if scanner.Scan() {
		got := scanner.Text()
		want := "Kjevik;SN39040;18.03.2022 01:50;42,8"
		if got != want {
			t.Errorf("Konverter() = %q; want %q", got, want)
		}
	}
}
