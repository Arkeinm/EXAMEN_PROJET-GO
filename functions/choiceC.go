package functions

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HandleChoiceC(outDir string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Fournissez le nom d'un article Wikipédia (ex: Go_(langage)) : ")
	articleInput, _ := reader.ReadString('\n')
	article := strings.TrimSpace(articleInput)

	if article == "" {
		fmt.Println("Aucun article fourni, utilisation de 'Go_(langage)' comme article par défaut.")
		article = "Go_(langage)"
	}

	resp, err := fetchWikipediaPage(article)
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("Tentative avec l'article par défaut 'Go_(langage)'...")

		resp, err = fetchWikipediaPage("Go_(langage)")
		if err != nil {
			fmt.Printf("Impossible de récupérer la page : %v\n", err)
			return
		}
		article = "Go_(langage)" 
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Erreur lors de l'analyse du document HTML : %v\n", err)
		return
	}

	var lines []string
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text != "" {
			lines = append(lines, text)
		}
	})

	fmt.Print("Fournissez un mot-clé à rechercher : ")
	keywordInput, _ := reader.ReadString('\n')
	keyword := strings.TrimSpace(keywordInput)

	if keyword == "" {
		fmt.Println("Aucun mot-clé fourni, utilisation de 'Lorem' comme mot-clé.")
		keyword = "Lorem"
	}

	lineCountWithKeyword := 0
	var filteredLines []string

	for _, line := range lines {
		if strings.Contains(line, keyword) {
			lineCountWithKeyword++
			filteredLines = append(filteredLines, line)
		}
	}

	fmt.Printf("Nombre de lignes contenant le mot-clé '%s' : %d\n", keyword, lineCountWithKeyword)

	if lineCountWithKeyword > 0 {
		outputFileName := outDir + "/wiki_" + article + ".txt"
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Println("Erreur création fichier wiki :", err)
			return
		}
		defer outputFile.Close()

		for _, line := range filteredLines {
			outputFile.WriteString(line + "\n")
		}

		fmt.Printf("Fichier %s créé avec succès avec %d ligne(s).\n", outputFileName, lineCountWithKeyword)
	} else {
		fmt.Println("Aucune ligne ne contient le mot-clé, aucun fichier créé.")
	}
}

func fetchWikipediaPage(article string) (*http.Response, error) {
	url := "https://fr.wikipedia.org/wiki/" + article

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}

	req.Header.Set("User-Agent", "FileOps-StudentProject/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du téléchargement de la page : %v", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("erreur HTTP : statut %d lors de l'accès à %s", resp.StatusCode, url)
	}

	return resp, nil
}