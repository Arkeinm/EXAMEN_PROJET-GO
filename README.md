EXAMEN — PROJET GO
FileOps / WebOps / ProcOps / SecureOps
M1 Devops — Projet individuel – 16/01/2026
Notation graduée : 10 / 13 / 16 / 18 (sur 20)
________________________________________
Objectif
Développer un outil Go en console qui manipule des fichiers texte, récupère du contenu web (Wikipédia), gère les processus système, et ajoute des mécanismes de sécurité (droits / verrouillage) aux niveaux avancés.
________________________________________
Contraintes générales
•	Go uniquement
•	Standard library autorisée
•	Bibliothèques externes autorisées uniquement lorsqu’elles sont imposées (Wikipédia)
•	Le programme doit gérer les erreurs proprement
•	Les fichiers générés doivent être écrits dans out/
________________________________________
Livrables du projet (idéalement lien vers github/gitlab, sinon zip)
A rendre sur Cesar
•	code source Go 
•	config.txt - config.json
•	data/ tous les fichiers d’input nécessaires
•	out/ fichiers créés par votre programme
•	README.md :
o	Procédure d’exécution
o	fonctionnalités implémentées
o	niveau visé (10/13/16/18)
o	Description du travail effectué
________________________________________
NIVEAU 10/20 — FileOps (fichiers & données) + config TXT
1) Menu interactif
•	menu en boucle + choix utilisateur + quitter
2) Configuration au format .txt
Au démarrage, le programme lit un fichier config.txt.
Format imposé
Fichier texte avec une clé par ligne :
default_file=data/input.txt
base_dir=data
out_dir=out
default_ext=.txt
•	les lignes vides ou commençant par # sont ignorées
•	si une clé manque, vous utilisez une valeur par défaut
3) Fichier texte courant (fournir un ou plusieurs fichier - lorem ipsum par exemple)
•	charger le fichier par défaut depuis la config
•	possibilité de choisir un autre fichier via le menu
•	vérifier existence et type (fichier)
4) Liste des fonctionnalités – choix menu
Sélectionnable depuis le menu Choix A et B
Chaque choix exécutera toutes les sous-fonctionnalités décrites
Le programme demandera à l’utilisateur un chemin. Si l’utilisateur ne le fourni pas => chemin par défaut dans le fichier de config.
Choix A - Analyse sur fichier courant
L’utilisateur fourni le nom d’un fichier
1.	Infos fichier : taille, date création/modif, nb lignes
2.	Stats mots : nb mots en ignorant les numériques + longueur moyenne
3.	Compter lignes contenant un mot-clé
4.	Filtrer lignes contenant un mot-clé demandé → out/filtered.txt
5.	Filtrer lignes ne contenant pas le mot-clé demandé → out/filtered_not.txt
6.	Head : N premières lignes → out/head.txt
7.	Tail : N dernières lignes → out/tail.txt
Choix B - Analyse multi-fichiers
L’utilisateur fourni le nom d’un repertoire
8.	Batch : analyser tous les .txt situé dans un emplacement demandé à l’utilisateur
9.	Rapport global : générer out/report.txt (format libre mais lisible)
10.	Indexation : générer out/index.txt listant (chemin, taille, date)
11.	Fusion : fusionner tous les .txt de base_dir → out/merged.txt
________________________________________
NIVEAU 13/20 — WebOps : Wikipédia avec goquery
Ajouter au menu :
Choix C - Analyser une page Wikipédia
Bibliothèque conseillée : goquery
Installation :
go get github.com/PuerkitoBio/goquery
Comportement attendu
1.	demander un article (ex: Go_(langage))
2.	télécharger :
https://fr.wikipedia.org/wiki/<ARTICLE>
3.	extraire le texte des paragraphes
4.	appliquer au moins 2 traitements déjà présents (mots / filtres / stats)
5.	écrire dans out/wiki_<article>.txt (ou afficher)
Bonus : possibilité de traiter plusieurs articles en même temps.
Extraction HTML via regex déconseillée.
________________________________________
NIVEAU 16/20 — ProcOps : gestion des processus Windows + macOS
Ce niveau doit fonctionner sur Windows et sur macOS.
Piste : OS-aware (détection runtime.GOOS).
Fonctionnalités ajoutée au menu principal :
Choix D . Créer un sous-menu ProcessOps :
1) Lister les processus (top N)
•	Windows : utiliser tasklist (par ex. tasklist /FO CSV)
•	macOS : utiliser ps -Ao pid,comm
Afficher au minimum : PID + nom.
2) Rechercher / filtrer
•	demander un mot (ex: “chrome”, “go”, “code”)
•	afficher les processus correspondants
3) Kill sécurisé
•	demander un PID
•	afficher un récapitulatif (PID + nom si possible)
•	demander une confirmation explicite (yes/no)
•	exécuter :
o	Windows : taskkill /PID <pid> /T (et éventuellement /F en option)
o	macOS : kill <pid> (option avancée : kill -9)
4) Gestion d’erreurs
•	PID invalide
•	droits insuffisants
•	process déjà terminé
•	commande non disponible
________________________________________
NIVEAU 18/20 — SecureOps + config JSON + droits + verrouillage
À ce niveau :
1.	la configuration doit être migrée vers JSON
2.	vous ajoutez une fonctionnalité de gestion des droits / verrouillage
________________________________________
A) Configuration en JSON
Remplacer config.txt par config.json.
Exemple :
{
  "default_file": "data/input.txt",
  "base_dir": "data",
  "out_dir": "out",
  "default_ext": ".txt",
  "wiki_lang": "fr",
  "process_top_n": 10
}
Le programme doit accepter un flag optionnel :
•	--config config.json
________________________________________
Choix E Gestion des droits
Ajouter au menu “SecureOps” une option :
Verrouiller / Déverrouiller un fichier
Objectif : empêcher l’écriture concurrente (ou simuler un verrou).
Deux options acceptées :
Option 1 (recommandée, portable, simple) : Lockfile
•	créer un fichier out/<nom>.lock
•	si le lock existe → fichier considéré “verrouillé”
•	déverrouiller = supprimer le lock
•	avantage : marche partout, simple à corriger
Option 2 (avancée, vraie lock OS - Bonus) : File lock OS
•	macOS/Linux : syscall.Flock
•	Windows : LockFileEx via syscall
•	plus complexe...

________________________________________
C) Bonus sécurité (au choix, mais “système”)
Ajouter au moins 1 :
•	rendre un fichier “read-only” (si OS le permet) via os.Chmod (macOS) / attribut Windows via commande
•	vérifier permissions (macOS) et signaler WARN
•	journaliser les actions sensibles (kill, lock) dans out/audit.log
•	Toute action destructive (kill, lock) doit être confirmée et loggée dans out/audit.log (fortement recommandé)


BONUS : Dockeriser ?

