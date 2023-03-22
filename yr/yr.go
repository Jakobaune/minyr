package yr

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Jakobaune/funtemps/conv"
)

// Main-funksjonen åpner en input-fil for lesing, oppretter en output-fil for skriving,
// og skriver konvertert data til output-filen.
func Konverter() {
	// Sjekk om filen allerede finnes
	if _, err := os.Stat("Resultat.txt"); err == nil {
		var input string
		fmt.Print("Filen finnes allerede. Vil du opprette en ny fil? (ja/nei): ")
		fmt.Scanln(&input)
		if input == "nei" {
			return // avslutt funksjonen uten å opprette en ny fil
		}
	}
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

func Average() {
	// Open the csv file
	file, err := os.OpenFile("kjevik-temp-celsius-20220318-20230318.csv", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the lines from the csv file
	scanner := bufio.NewScanner(file)

	var sum float64
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if count == 0 {
			count++
			continue // ignore header line
		}
		fields := strings.Split(line, ";")
		if len(fields) != 4 {
			log.Fatalf("unexpected number of fields in line %d: %d", count, len(fields))
		}
		if fields[3] == "" {
			continue // ignore line with empty temperature field
		}
		temperature, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			log.Fatalf("could not parse temperature in line %d: %s", count, err)
		}
		sum += temperature
		count++
	}
	average := sum / float64(count)
	average = math.Round(average*100) / 100 // round to two decimal places

	fmt.Println("Vil du skrive ut temperaturen i grader Celsius eller Fahrenheit? (c/f)")
	var choice string
	fmt.Scanln(&choice)
	if choice == "c" {
		fmt.Printf("Gjennomsnittlig temperatur: %.2f°C\n", average)
	} else if choice == "f" {
		fahrenheit := conv.CelsiusToFahrenheit(average)
		fmt.Printf("Gjennomsnittlig temperatur: %.2f°F\n", fahrenheit)
	} else {
		fmt.Println("Ugyldig valg")
	}
}
