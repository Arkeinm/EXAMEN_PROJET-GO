package functions

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

// Choix B - Analyse multi-fichiers
// L’utilisateur fourni le nom d’un repertoire
// 8.	Batch : analyser tous les .txt situé dans un emplacement demandé à l’utilisateur
// 9.	Rapport global : générer out/report.txt (format libre mais lisible)
// 10.	Indexation : générer out/index.txt listant (chemin, taille, date)
// 11.	Fusion : fusionner tous les .txt de base_dir → out/merged.txt
func HandleChoiceB(extensionsConfig string, outDir string, defaultExt string) {
	var reportLines []string
	var indexLines []string

	reader := bufio.NewReader(os.Stdin)
	extensions := strings.Split(extensionsConfig, ",")

	fmt.Print("Fournissez un nom de répertoire : ")
	input, _ := reader.ReadString('\n')
	dirName := strings.TrimSpace(input)

	dirInfo, err := os.Stat(dirName)
	if os.IsNotExist(err) || !dirInfo.IsDir() {
		fmt.Println("Le répertoire n'existe pas.")
		return
	}

	files, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du répertoire :", err)
		return
	}

	mergedFile, err := os.Create(outDir + "/merged" + defaultExt)
	if err != nil {
		fmt.Println("Erreur création fichier merged"+defaultExt+" :", err)
		return
	}
	defer mergedFile.Close()

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matched := false
		for _, ext := range extensions {
			if strings.HasSuffix(file.Name(), strings.TrimSpace(ext)) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		filePath := filepath.Join(dirName, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Erreur lors de l'obtention des informations du fichier :", err)
			continue
		}

		indexLines = append(indexLines, fmt.Sprintf("%s, %d octets, %s", filePath, fileInfo.Size(), fileInfo.ModTime().String()))

		f, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Erreur lors de l'ouverture du fichier :", err)
			continue
		}

		scanner := bufio.NewScanner(f)
		lineCount := 0
		for scanner.Scan() {
			line := scanner.Text()
			mergedFile.WriteString(line + "\n")
			lineCount++
		}
		f.Close()

		reportLines = append(reportLines, fmt.Sprintf("Fichier: %s, Lignes: %d", file.Name(), lineCount))
	}

	reportFile, err := os.Create(outDir + "/report" + defaultExt)
	if err != nil {
		fmt.Println("Erreur création fichier report"+defaultExt+" :", err)
		return
	}
	defer reportFile.Close()

	for _, line := range reportLines {
		reportFile.WriteString(line + "\n")
	}

	indexFile, err := os.Create(outDir + "/index" + defaultExt)
	if err != nil {
		fmt.Println("Erreur création fichier index"+defaultExt+" :", err)
		return
	}
	defer indexFile.Close()

	for _, line := range indexLines {
		indexFile.WriteString(line + "\n")
	}

	fmt.Println("Analyse terminée. Fichiers générés dans le répertoire ./" + outDir)
}
