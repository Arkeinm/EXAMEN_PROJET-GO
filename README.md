# EVALUATION - Projet FileOps & ProcessOps

## 📋 Description du projet

Ce projet implémente un système de gestion de fichiers et de processus avec l'utilisation du langage Go. Il contient plusieurs niveaux de fonctionnalités progressives, comme par exemple la manipulation des fichiers texte, analyser leur contenu, gérer des répertoires, récupérer des articles Wikipédia et gérer les processus système.

## 🎯 Niveau visé

Niveau 16/20 - ProcessOps : gestion des processus Windows + macOS

## 🚀 Procédure d'exécution

### Prérequis
- Go version 1.18 ou supérieur
- Accès à Internet pour le choix C correspondant à la fonctionnalité Wikipédia
- Droits administrateur/root pour la terminaison de processus

### Installation et lancement

1. Cloner ou télécharger le projet
2. Ouvrir un terminal dans le répertoire du projet
3. Exécuter le programme :
```bash
go run run-me.go
```

Ou compiler puis exécuter :
```bash
go build run-me.go
./run-me          # Linux/macOS
run-me.exe        # Windows
```

### Configuration

Le fichier `config.txt` contient les paramètres par défaut :
- Fichier d'entrée par défaut : `data/input.txt`
- Répertoire de sortie : `out/`
- Extension de fichier par défaut : `.txt`

Ils sont modifiables et utilisables selon vos attentes d'utilisation

## ✨ Fonctionnalités implémentées

### 📄 Choix A - FileOps : Analyse de fichier unique

**Fonctionnalités :**
- Saisie du nom de fichier (avec fichier par défaut si inexistant)
- Affichage des métadonnées (taille, date de modification)
- Comptage du nombre de lignes
- Comptage des mots (exclusion des valeurs numériques)
- Calcul de la longueur moyenne des mots
- Recherche de mot-clé dans les lignes
- Création de fichiers filtrés :
  - `filtered.txt` : lignes contenant le mot-clé
  - `filtered_not.txt` : lignes ne contenant pas le mot-clé
- Extraction des N premières et dernières lignes :
  - `head.txt` : N premières lignes
  - `tail.txt` : N dernières lignes

**Gestion d'erreurs :**
- Fichier inexistant → utilisation du fichier par défaut
- Nombre de lignes invalide → valeur par défaut de 4
- Nombre > 1000 → message d'avertissement

---

### 📁 Choix B - DirOps : Analyse de répertoire

**Fonctionnalités :**
- Saisie du nom de répertoire (avec répertoire par défaut)
- Saisie du mot-clé à rechercher
- Parcours récursif du répertoire
- Analyse de tous les fichiers `.txt`
- Génération de trois fichiers de sortie :

**1. `report.txt` :**
- Métadonnées de chaque fichier (nom, taille, date)
- Nombre de lignes
- Nombre de mots (sans numériques)
- Longueur moyenne des mots
- Nombre de lignes contenant le mot-clé
- Résumé global (total fichiers, mots, lignes avec mot-clé)

**2. `index.txt` :**
- Liste condensée : chemin | taille | date

**3. `merged.txt` :**
- Fusion du contenu de tous les fichiers analysés

**Gestion d'erreurs :**
- Répertoire inexistant → utilisation du répertoire par défaut
- Chemin n'est pas un répertoire → message d'erreur
- Erreur de lecture de fichier → message dans le rapport

---

### 🌐 Choix C - WebOps : Récupération et analyse d'articles Wikipédia

**Fonctionnalités :**
- Saisie du nom d'un article Wikipédia (exemple par défaut: `Go_(langage)`)
- Téléchargement des données de la page Wikipédia voulu
- Extraction du contenu des paragraphes grâce à la balise `<p>`
- Recherche de mot-clé dans les paragraphes
- Comptage des lignes contenant le mot-clé
- Création du fichier `wiki_<nom_article>.txt` avec les paragraphes filtrés

**Gestion d'erreurs :**
- Article vide → utilisation de `Go_(langage)` par défaut
- Page inexistante (erreur HTTP) → tentative avec article par défaut
- Erreur de téléchargement → message d'erreur détaillé
- Erreur de parsing HTML → message d'erreur

**Dépendances :**
- `github.com/PuerkitoBio/goquery` : parsing HTML

---

### ⚙️ Choix D - ProcessOps : Gestion des processus (Windows + macOS)

**Sous-menu avec 4 options :**

**1. Lister les processus (top N) :**
- Demande du nombre de processus à afficher
- Windows : utilise `tasklist /FO CSV`
- macOS/Linux : utilise `ps -Ao pid,comm` ou `ps -eo pid,comm`
- Affichage : PID + Nom du processus

**2. Rechercher/filtrer un processus :**
- Demande d'un mot-clé (ex: chrome, go, code)
- Filtrage et affichage des processus correspondants
- Recherche insensible à la casse

**3. Terminer un processus (kill sécurisé) :**
- Demande du PID
- Vérification de l'existence du processus
- Affichage du récapitulatif (PID + nom)
- Demande de confirmation explicite (yes/no)
- Option de forçage :
  - Windows : `taskkill /PID <pid> /T /F`
  - macOS/Linux : `kill -9 <pid>`
- Annulation possible à tout moment

**4. Retour au menu principal**

**Gestion d'erreurs :**
- ✅ PID invalide (non numérique)
- ✅ Processus inexistant ou déjà terminé
- ✅ Droits insuffisants (message pour exécuter en admin/root)
- ✅ Commandes non disponibles (`tasklist`, `ps`, `kill`, etc.)
- ✅ Détection automatique de l'OS (Windows/macOS/Linux)

---

## 📂 Structure du projet
```
EVALUATION/
├── data/
│   └── input.txt              # Fichier d'entrée par défaut
├── functions/
│   ├── choiceA.go            # Analyse de fichier unique
│   ├── choiceB.go            # Analyse de répertoire
│   ├── choiceC.go            # Récupération Wikipédia
│   └── choiceD.go            # Gestion des processus
├── out/                       # Répertoire de sortie (fichiers générés)
│   ├── filtered.txt
│   ├── filtered_not.txt
│   ├── head.txt
│   ├── tail.txt
│   ├── report.txt
│   ├── index.txt
│   ├── merged.txt
│   ├── wiki_Go_(langage).txt
|   └── wiki_Google.txt
├── config.txt                 # Configuration du programme
├── go.mod                     # Dépendances Go
├── go.sum                     # Checksums des dépendances
├── run-me.go                  # Point d'entrée du programme
└── README.md                  # Ce fichier
```

---

## 📝 Notes

- Les fichiers de sortie sont créés dans le répertoire par défaut `out/` - Si vous n'avez pas changé le fichier `config.txt`
- Le programme détecte automatiquement le système d'exploitation
- Les droits administrateur sont nécessaires pour terminer certains processus
- Le mot-clé par défaut est "Lorem" à chacune des demandes si aucun n'est fourni

---

## 🔧 Dépendances

Le projet utilise `github.com/PuerkitoBio/goquery` pour le parsing HTML.

Les dépendances sont gérées automatiquement via `go.mod` et `go.sum`. Aucune installation manuelle n'est nécessaire - Go téléchargera automatiquement les dépendances au premier lancement.
