package functions

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func HandleChoiceD() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== ProcessOps ===")
		fmt.Println("1) Lister les processus")
		fmt.Println("2) Rechercher/filtrer un processus")
		fmt.Println("3) Terminer un processus (kill)")
		fmt.Println("4) Retour au menu principal")
		fmt.Print("Votre choix : ")

		choiceInput, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(choiceInput)

		switch choice {
		case "1":
			listProcesses(reader)
		case "2":
			searchProcesses(reader)
		case "3":
			killProcess(reader)
		case "4":
			return
		default:
			fmt.Println("Choix invalide. Veuillez choisir entre 1 et 4.")
		}
	}
}

func checkCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func listProcesses(reader *bufio.Reader) {
	var cmdName string
	switch runtime.GOOS {
	case "windows":
		cmdName = "tasklist"
	default:
		cmdName = "ps"
	}

	if !checkCommandAvailable(cmdName) {
		fmt.Printf("Erreur : la commande '%s' n'est pas disponible sur ce système.\n", cmdName)
		return
	}

	fmt.Print("Combien de processus voulez-vous afficher ? (ex: 10, 20) : ")
	limitInput, _ := reader.ReadString('\n')
	limit, err := strconv.Atoi(strings.TrimSpace(limitInput))

	if err != nil || limit <= 0 {
		fmt.Println("Nombre invalide. Utilisation de 10 par défaut.")
		limit = 10
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/FO", "CSV", "/NH")
	case "darwin":
		cmd = exec.Command("ps", "-Ao", "pid,comm")
	default:
		cmd = exec.Command("ps", "-eo", "pid,comm")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Erreur lors de l'exécution de la commande : %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	count := 0

	fmt.Println("\nPID\t\tNom du processus")
	fmt.Println("---\t\t-----------------")

	for _, line := range lines {
		if line == "" {
			continue
		}

		switch runtime.GOOS {
		case "windows":
			fields := parseCSVLine(line)
			if len(fields) >= 2 {
				fmt.Printf("%s\t\t%s\n", fields[1], fields[0])
				count++
			}
		default:
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[0] != "PID" {
				fmt.Printf("%s\t\t%s\n", fields[0], strings.Join(fields[1:], " "))
				count++
			}
		}

		if count >= limit {
			break
		}
	}
}

func searchProcesses(reader *bufio.Reader) {
	var cmdName string
	switch runtime.GOOS {
	case "windows":
		cmdName = "tasklist"
	default:
		cmdName = "ps"
	}

	if !checkCommandAvailable(cmdName) {
		fmt.Printf("Erreur : la commande '%s' n'est pas disponible sur ce système.\n", cmdName)
		return
	}

	fmt.Print("Entrez un mot-clé à rechercher (ex: chrome, go, code) : ")
	keywordInput, _ := reader.ReadString('\n')
	keyword := strings.ToLower(strings.TrimSpace(keywordInput))

	if keyword == "" {
		fmt.Println("Aucun mot-clé fourni.")
		return
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/FO", "CSV", "/NH")
	case "darwin":
		cmd = exec.Command("ps", "-Ao", "pid,comm")
	default:
		cmd = exec.Command("ps", "-eo", "pid,comm")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Erreur lors de l'exécution de la commande : %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	found := false

	fmt.Println("\nPID\t\tNom du processus")
	fmt.Println("---\t\t-----------------")

	for _, line := range lines {
		if line == "" {
			continue
		}

		switch runtime.GOOS {
		case "windows":
			fields := parseCSVLine(line)
			if len(fields) >= 2 {
				processName := strings.ToLower(fields[0])
				if strings.Contains(processName, keyword) {
					fmt.Printf("%s\t\t%s\n", fields[1], fields[0])
					found = true
				}
			}
		default:
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[0] != "PID" {
				processName := strings.ToLower(strings.Join(fields[1:], " "))
				if strings.Contains(processName, keyword) {
					fmt.Printf("%s\t\t%s\n", fields[0], strings.Join(fields[1:], " "))
					found = true
				}
			}
		}
	}

	if !found {
		fmt.Printf("Aucun processus trouvé contenant '%s'.\n", keyword)
	}
}

func killProcess(reader *bufio.Reader) {
	var cmdName string
	switch runtime.GOOS {
	case "windows":
		cmdName = "taskkill"
	default:
		cmdName = "kill"
	}

	if !checkCommandAvailable(cmdName) {
		fmt.Printf("Erreur : la commande '%s' n'est pas disponible sur ce système.\n", cmdName)
		return
	}

	fmt.Print("Fournissez un PID de processus à terminer : ")
	pidInput, _ := reader.ReadString('\n')
	pidStr := strings.TrimSpace(pidInput)

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Println("PID invalide. Veuillez fournir un nombre entier.")
		return
	}

	processName := getProcessNameByPID(pid)
	if processName == "" {
		fmt.Printf("Aucun processus trouvé avec le PID %d. Le processus n'existe pas ou est déjà terminé.\n", pid)
		return
	}

	fmt.Printf("\nProcessus à terminer :\n")
	fmt.Printf("  PID  : %d\n", pid)
	fmt.Printf("  Nom  : %s\n", processName)

	fmt.Print("\nConfirmez-vous la terminaison de ce processus ? (yes/no) : ")
	confirmInput, _ := reader.ReadString('\n')
	confirmation := strings.TrimSpace(strings.ToLower(confirmInput))

	if confirmation != "yes" {
		fmt.Println("Opération annulée par l'utilisateur.")
		return
	}

	fmt.Print("Voulez-vous forcer la terminaison ? (yes/no) : ")
	forceInput, _ := reader.ReadString('\n')
	force := strings.TrimSpace(strings.ToLower(forceInput)) == "yes"

	err = terminateProcessByPID(pid, force)

	if err != nil {
		errMsg := err.Error()
		fmt.Printf("Erreur lors de la terminaison du processus : %v\n", err)
		
		if strings.Contains(errMsg, "Access is denied") || strings.Contains(errMsg, "Operation not permitted") {
			fmt.Println("Droits insuffisants. Essayez d'exécuter le programme en tant qu'administrateur/root.")
		} else if strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "No such process") {
			fmt.Println("Le processus n'existe plus ou est déjà terminé.")
		}
		return
	}

	fmt.Printf("Processus PID %d terminé avec succès.\n", pid)
}

func getProcessNameByPID(pid int) string {
	var cmd *exec.Cmd
	var cmdName string

	switch runtime.GOOS {
	case "windows":
		cmdName = "tasklist"
	default:
		cmdName = "ps"
	}

	if !checkCommandAvailable(cmdName) {
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/FO", "CSV", "/NH", "/FI", fmt.Sprintf("PID eq %d", pid))
	case "darwin":
		cmd = exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	default:
		cmd = exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		fields := parseCSVLine(result)
		if len(fields) >= 1 {
			return fields[0]
		}
		return ""
	default:
		return result
	}
}

func terminateProcessByPID(pid int, force bool) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		if force {
			cmd = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F")
		} else {
			cmd = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T")
		}
	default:
		if force {
			cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
		} else {
			cmd = exec.Command("kill", strconv.Itoa(pid))
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}

	return nil
}

func parseCSVLine(line string) []string {
	var fields []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		if char == '"' {
			inQuotes = !inQuotes
		} else if char == ',' && !inQuotes {
			fields = append(fields, current.String())
			current.Reset()
		} else {
			current.WriteByte(char)
		}
	}

	fields = append(fields, current.String())
	return fields
}