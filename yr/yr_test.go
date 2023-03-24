package yr_test

import (
	"bufio"
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
func TestConvertLine(t *testing.T) {

	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "Kjevik;SN39040;18.03.2022 01:50;6", want: "Kjevik;SN39040;18.03.2022 01:50;42.8"},
		{input: "Kjevik;SN39040;07.03.2023 18:20;0", want: "Kjevik;SN39040;07.03.2023 18:20;32.0"},
		{input: "Kjevik;SN39040;08.03.2023 02:20;-11", want: "Kjevik;SN39040;08.03.2023 02:20;12.2"},
		{input: "Data er gyldig per 18.03.2023 (CC BY 4.0), Meteorologisk institutt (MET);;;", want: "Data er basert på gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Jakob Aune"},
	}

	for _, tc := range tests {
		got, _ := yr.ConvertLine(tc.input)
		if !(tc.want == got) {
			t.Errorf("expected: %v, got: %v", tc.want, got)
		}
	}
}
func TestAverageCelcius(t *testing.T) {
	expected := "gjennomsnittstemperatur 8.56"
	result := yr.AverageCelcius()
	if result != expected {
		t.Errorf("expected %q but got %q", expected, result)
	}
}
