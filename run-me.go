package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"myapp/functions"
)

func main() {
	filePath := "config.txt"
	defaultFile := "data/input.txt"
	extensionsConfig := "data"
	outDir := "out"
	defaultExt := ".txt"

	reader := bufio.NewReader(os.Stdin)

	configFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer configFile.Close()

	scanner := bufio.NewScanner(configFile)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "default_file":
			defaultFile = value
		case "extensions":
			extensionsConfig = value
		case "out_dir":
			outDir = value
		case "default_ext":
			defaultExt = value
		}
	}

	for {
		fmt.Println("/!\\ Pour arrêter le programme, appuyez sur la touche 'Entrée' sans rien taper. /!\\")
		fmt.Print("Quel choix veux-tu faire ? (A ou B) : ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		if choice == "" {
			fmt.Println("Entrée vide, arrêt du programme.")
			break
		}

		if choice != "A" && choice != "B" {
			fmt.Println("Choix invalide. Veuillez choisir A ou B.")
			continue
		}

		switch choice {
		case "A":
			functions.HandleChoiceA(defaultFile, outDir)
		case "B":
			functions.HandleChoiceB(extensionsConfig, outDir, defaultExt)
		}

	}

}