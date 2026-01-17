package functions

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func HandleChoiceB(extensionsConfig string, outDir string, defaultExt string) {
	count := 0
	totalWordCount := 0
	totalWordLength := 0
	totalLinesWithKeyword := 0
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Fournissez un nom de répertoire : ")
	repertoryInput, _ := reader.ReadString('\n')
	repertoryPath := strings.TrimSpace(repertoryInput)

	if repertoryPath == "" {
		fmt.Printf("Aucun répertoire fourni, utilisation de ./%s comme répertoire par défaut.\n", outDir)
		repertoryPath = outDir
	}

	repertoryInfo, err := os.Stat(repertoryPath)
	if os.IsNotExist(err) {
		fmt.Println("Le répertoire n'existe pas. Nous prendrons par défaut :", outDir)
		repertoryPath = outDir
		repertoryInfo, err = os.Stat(repertoryPath)
		if os.IsNotExist(err) {
			fmt.Println("Problème : le répertoire par défaut n'existe pas non plus.")
			return
		}
	}

	if err == nil && !repertoryInfo.IsDir() {
		fmt.Printf("Erreur : '%s' existe mais n'est pas un répertoire.\n", repertoryPath)
		return
	}

	fmt.Print("Fournissez un mot-clé à rechercher dans les lignes : ")
	keywordInput, _ := reader.ReadString('\n')
	keyword := strings.TrimSpace(keywordInput)

	if keyword == "" {
		fmt.Println("Aucun mot-clé fourni, utilisation de 'Lorem' comme mot-clé.")
		keyword = "Lorem"
	}

	reportFile, err := os.Create(outDir + "/report.txt")
	if err != nil {
		fmt.Println("Erreur création fichier report.txt :", err)
		return
	}
	defer reportFile.Close()

	fmt.Println("Fichier report.txt créé avec succès.")

	indexFile, err := os.Create(outDir + "/index.txt")
	if err != nil {
		fmt.Println("Erreur création fichier index.txt :", err)
		return
	}
	defer indexFile.Close()

	fmt.Println("Fichier index.txt créé avec succès.")

	mergedFile, err := os.Create(outDir + "/merged.txt")
	if err != nil {
		fmt.Println("Erreur création fichier merged.txt :", err)
		return
	}
	defer mergedFile.Close()

	fmt.Println("Fichier merged.txt créé avec succès.")

	error := filepath.WalkDir(repertoryPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(d.Name())

		if ext == defaultExt {
			count++

			fileInfo, err := d.Info()
			if err != nil {
				fmt.Fprintf(reportFile, "Erreur lors de la récupération des informations pour %s: %v\n", path, err)
				return nil
			}

			fmt.Fprintf(reportFile, "Fichier : %s\n", path)
			fmt.Fprintf(reportFile, "Taille : %d octets\n", fileInfo.Size())
			fmt.Fprintf(reportFile, "Date de modification : %s\n", fileInfo.ModTime().String())

			fmt.Fprintf(indexFile, "%s | %d octets | %s\n", path, fileInfo.Size(), fileInfo.ModTime().String())

			file, err := os.Open(path)
			if err != nil {
				fmt.Fprintf(reportFile, "Erreur ouverture dufichier : %v\n\n", err)
				return nil
			}
			defer file.Close()

			var lines []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintf(reportFile, "Erreur lecture du fichier : %v\n\n", err)
				return nil
			}

			fmt.Fprintf(reportFile, "Nombre de lignes : %d\n", len(lines))

			wordCount := 0
			totalLength := 0

			for _, line := range lines {
				for word := range strings.FieldsSeq(line) {
					if _, err := fmt.Sscanf(word, "%f", new(float64)); err != nil {
						wordCount++
						totalLength += len(word)
					}
				}
			}

			avgLength := 0
			if wordCount > 0 {
				avgLength = totalLength / wordCount
			}

			fmt.Fprintf(reportFile, "Nombre de mots (sans numériques) : %d\n", wordCount)
			fmt.Fprintf(reportFile, "Longueur moyenne des mots : %d\n", avgLength)

			totalWordCount += wordCount
			totalWordLength += totalLength

			lineCountWithKeyword := 0
			for _, line := range lines {
				if strings.Contains(line, keyword) {
					lineCountWithKeyword++
				}
			}

			fmt.Fprintf(reportFile, "Nombre de lignes contenant le mot-clé '%s' : %d\n\n", keyword, lineCountWithKeyword)
			totalLinesWithKeyword += lineCountWithKeyword

			fmt.Fprintf(mergedFile, "=== Contenu de %s ===\n", path)
			for _, line := range lines {
				fmt.Fprintln(mergedFile, line)
			}
			fmt.Fprintln(mergedFile, "")
		}

		return nil
	})

	if error != nil {
		fmt.Fprintln(os.Stderr, "error:", error)
		os.Exit(1)
	}

	fmt.Fprintf(reportFile, "=== RÉSUMÉ GLOBAL ===\n")
	fmt.Fprintf(reportFile, "Nombre de fichiers avec l'extension %s : %d\n", defaultExt, count)
	fmt.Fprintf(reportFile, "Total de mots (tous fichiers) : %d\n", totalWordCount)

	globalAvgLength := 0
	if totalWordCount > 0 {
		globalAvgLength = totalWordLength / totalWordCount
	}

	fmt.Fprintf(reportFile, "Longueur moyenne globale des mots : %d\n", globalAvgLength)
	fmt.Fprintf(reportFile, "Total de lignes contenant '%s' : %d\n", keyword, totalLinesWithKeyword)

}
