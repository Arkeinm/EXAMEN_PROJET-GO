package functions

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func HandleChoiceA(defaultFile string, outDir string) {
	var lines []string

	wordCount := 0
	totalLength := 0
	avgLength := 0
	lineCountWithKeyword := 0

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Fournissez un nom de fichier : ")
	input, _ := reader.ReadString('\n')
	fileName := strings.TrimSpace(input)

	fileInfo, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		fmt.Println("Le fichier n'existe pas. Nous prendrons par défaut :", defaultFile)
		fileName = defaultFile
		fileInfo, err = os.Stat(fileName)
		if os.IsNotExist(err) {
			fmt.Println("Problème : le fichier par défaut n'existe pas non plus.")
			return
		}
	}

	fmt.Printf("Taille du fichier : %d octets\n", fileInfo.Size())
	fmt.Printf("Date de création ou de dernière modification : %s\n", fileInfo.ModTime().String())

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lecture fichier :", err)
		return
	}

	fmt.Printf("Nombre de lignes dans le fichier : %d\n", len(lines))

	for _, line := range lines {
		for word := range strings.FieldsSeq(line) {
			if _, err := fmt.Sscanf(word, "%f", new(float64)); err != nil {
				wordCount++
				totalLength += len(word)
			}
		}
	}

	if wordCount > 0 {
		avgLength = totalLength / wordCount
	}

	fmt.Printf("Nombre de mots (sans numériques) : %d\n", wordCount)
	fmt.Printf("Longueur moyenne des mots : %d\n", avgLength)

	fmt.Print("Fournissez un mot-clé à rechercher dans les lignes : ")
	keywordInput, _ := reader.ReadString('\n')
	keyword := strings.TrimSpace(keywordInput)

	if keyword == "" {
		fmt.Println("Aucun mot-clé fourni, utilisation de 'Lorem' comme mot-clé.")
		keyword = "Lorem"
	}

	for _, line := range lines {
		if strings.Contains(line, keyword) {
			lineCountWithKeyword++
		}
	}

	fmt.Printf("Nombre de lignes contenant le mot-clé '%s' : %d\n", keyword, lineCountWithKeyword)

	filteredFile, err := os.Create(outDir + "/filtered.txt")
	if err != nil {
		fmt.Println("Erreur création fichier filtered.txt :", err)
		return
	}
	defer filteredFile.Close()

	fmt.Println("Fichier filtered.txt créé avec succès.")

	filteredNotFile, err := os.Create(outDir + "/filtered_not.txt")
	if err != nil {
		fmt.Println("Erreur création fichier filtered_not.txt :", err)
		return
	}
	defer filteredNotFile.Close()

	fmt.Println("Fichier filtered_not.txt créé avec succès.")

	for _, line := range lines {
		if strings.Contains(line, keyword) {
			filteredFile.WriteString(line + "\n")
		} else {
			filteredNotFile.WriteString(line + "\n")
		}
	}

	fmt.Print("Fournissez un nombre de lignes que vous voudriez compris entre 0 et " + strconv.Itoa(len(lines)) + " : ")
	input, _ = reader.ReadString('\n')
	numberInput, err := strconv.Atoi(strings.TrimSpace(input))

	if err != nil {
		fmt.Println("Erreur de conversion en entier. Valeur par défaut de 4 sera utilisée.")
		numberInput = 4
	}

	if numberInput > 1000 {
		fmt.Println("Je ne suis pas écrivain mais informaticien, mettez une valeur plus petite s'il vous plaît.")
	} else if numberInput > 0 {
		if numberInput > len(lines) {
			numberInput = len(lines)
		}

		headFile, err := os.Create(outDir + "/head.txt")

		if err != nil {
			fmt.Println("Erreur création fichier head.txt :", err)
			return
		}
		defer headFile.Close()

		for i := 0; i < numberInput; i++ {
			headFile.WriteString(lines[i] + "\n")
		}

		tailFile, err := os.Create(outDir + "/tail.txt")

		if err != nil {
			fmt.Println("Erreur création fichier tail.txt :", err)
			return
		}
		defer tailFile.Close()

		for i := len(lines) - numberInput; i < len(lines); i++ {
			tailFile.WriteString(lines[i] + "\n")
		}

		fmt.Println("Fichiers head.txt et tail.txt créés avec succès dans le répertoire ./" + outDir)
		return
	}
	fmt.Println("Le nombre fourni est invalide. Aucun fichier head.txt ou tail.txt n'a été créé.")
}
