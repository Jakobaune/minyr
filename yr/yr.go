package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Jakobaune/funtemps/conv"
)

// Main-funksjonen åpner en input-fil for lesing, oppretter en output-fil for skriving,
// og skriver konvertert data til output-filen.
func main() {
	// Åpne input-filen for lesing
	file, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv")
	// Håndter eventuelle feil
	if err != nil {
		log.Fatal(err)
	}

	// Sørger for at input-filen blir lukket på slutten av funksjonen
	defer file.Close()

	// Opprett en ny fil for skriving
	output, err := os.Create("Resultat.txt")
	// Håndter eventuelle feil
	if err != nil {
		log.Fatal(err)
	}

	// Sørger for at output-filen blir lukket på slutten av funksjonen
	defer output.Close()

	// Opprett en writer fra output-filen
	writer := bufio.NewWriter(output)

	// Opprett en scanner fra input-filen
	scanner := bufio.NewScanner(file)

	// Skriv den første linjen fra input-filen til output-filen
	if scanner.Scan() {
		_, err := writer.WriteString(scanner.Text() + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	var outputLine string
	// Les hver linje fra input-filen
	for scanner.Scan() {
		// Hent den gjeldende linjen som en streng
		data := scanner.Text()

		// Del opp strengen ved semikolon til felt
		fields := strings.Split(data, ";")

		// Hent siste feltet
		var lastField string
		if len(fields) > 0 {
			lastField = fields[len(fields)-1]
		}

		var convertedField string
		if lastField != "" {
			// Konverter siste felt fra Celsius til Fahrenheit
			var err error
			convertedField, err = convertLastField(lastField)
			if err != nil {

				// Hvis konverteringen feiler, skriv ut en feilmelding og fortsett til neste linje
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				continue
			}
		}
		// Hvis konverteringen var vellykket, erstatt siste felt med den konverterte verdien
		if convertedField != "" {
			fields[len(fields)-1] = convertedField
		}

		// Hvis linjen starter med "Data er", skriv ut metadata og kildeinformasjon
		if data[0:7] == "Data er" {
			outputLine = "Endring gjort av Jakob Aune"
		} else {
			// Hvis linjen ikke starter med "Data er", skriv ut de konverterte feltene

			outputLine = strings.Join(fields, ";")
		}

		_, err = writer.WriteString(outputLine + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Sjekk om det var noen feil under scanningen av input-filen
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Tøm bufferen og skriv dataene til output-filen
	writer.Flush()
}

func convertLastField(lastField string) (string, error) {
	celsius, err := strconv.ParseFloat(lastField, 64)
	if err != nil {
		return "", err
	}

	fahrenheit := conv.CelsiusToFahrenheit(celsius)

	return fmt.Sprintf("%.1f", fahrenheit), nil
}

func average(scanner *bufio.Scanner, writer *bufio.Writer, unit string) {
	var sum, count float64

	// Loop over all lines in the input file
	for scanner.Scan() {
		// Parse the temperature value from the last field in the line
		data := scanner.Text()
		fields := strings.Split(data, ";")
		lastField := fields[len(fields)-1]
		temperature, err := strconv.ParseFloat(lastField, 64)
		if err != nil {
			// Ignore lines with invalid temperature values
			continue
		}

		// Accumulate the sum and count for the average calculation
		sum += temperature
		count++
	}

	// Calculate the average temperature
	average := sum / count

	// Print the result to the output file
	var outputLine string
	if unit == "c" {
		outputLine = fmt.Sprintf("Gjennomsnittstemperatur for perioden: %.1f °C", average)
	} else if unit == "f" {
		averageF := conv.CelsiusToFahrenheit(average)
		outputLine = fmt.Sprintf("Gjennomsnittstemperatur for perioden: %.1f °F", averageF)
	} else {
		outputLine = "Ugyldig enhet valgt. Velg 'c' eller 'f'."
	}
	_, err := writer.WriteString(outputLine + "\n")
	if err != nil {
		log.Fatal(err)
	}
}
