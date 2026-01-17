package functions

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Choix B - Analyse multi-fichiers
// 8. Batch : analyser tous les fichiers avec extensions autorisées
// 9. Rapport global : générer out/report.txt
// 10. Indexation : générer out/index.txt (chemin, taille, date)
// 11. Fusion : fusionner tous les fichiers → out/merged.txt
func HandleChoiceB(extensionsConfig string, outDir string, defaultExt string) {
	count := 0
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

	reportFile, err := os.Create(outDir + "/report.txt")
	if err != nil {
		fmt.Println("Erreur création fichier filtered.txt :", err)
		return
	}
	defer reportFile.Close()

	fmt.Println("Fichier report.txt créé avec succès.")

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
			fmt.Fprintf(reportFile, "Date de modification : %s\n\n", fileInfo.ModTime().String())
		}

		return nil
	})

	if error != nil {
		fmt.Fprintln(os.Stderr, "error:", error)
		os.Exit(1)
	}

	fmt.Fprintf(reportFile, "Nombre de fichiers avec l'extension %s : %d\n\n", defaultExt, count)
}
