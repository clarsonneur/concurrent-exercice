package main

import (
	"fmt"
	"strings"
	"time"

	tm "github.com/buger/goterm"
)

const (
	maxRoutines = 6
)

// printStatus affiche le status instantanné de l'ensemble des jobs et affiche ceux terminé en cumulatif.
func printStatus(id, running int, jobsCtrl []int) {
	if id > 0 {
		fmt.Printf("Message %d job done.\n", id)
	}
	tm.Printf("%s%s\n", strings.Repeat("#", running), strings.Repeat(" ", maxRoutines))
	for index, curId := range jobsCtrl {
		if curId == 0 {
			tm.Printf("%d:       \n", index)
			continue
		}
		tm.Printf("%d: %d     \n", index, curId)
	}
	tm.MoveCursorUp(maxRoutines + 2)
	tm.Flush()
}

func main() {

	tm.Clear()

	jobsStarted := make(chan int)
	jobsDone := make(chan int, maxRoutines)

	go func() {
		// On alimente le canal des messages
		for i := 1; i <= 40; i++ {
			jobsStarted <- i
		}
		close(jobsStarted)
	}()

	jobsCtrl := make([]int, maxRoutines)
	jobs := 0
	for {
		if jobs < maxRoutines {
			// On peut lancer une routine tant qu'on est inférieur à maxRoutine
			i, ok := <-jobsStarted
			if !ok {
				// La source d'information est tarie. On arrête de démarrer de nouvelles taches.
				break
			}

			jobIndex := 0
			for index, id := range jobsCtrl {
				if id == 0 {
					jobIndex = index
					break
				}
			}

			// On assigne un index au job pour affichage
			jobsCtrl[jobIndex] = i

			// On démarre la tache en parallèle avec le message
			go func(value, index int) {
				time.Sleep(time.Second)
				jobsDone <- index
			}(i, jobIndex)

			// On indique qu'on a démarré 1 tache de plus.
			jobs++
			printStatus(0, jobs, jobsCtrl)
		} else {
			// On a atteint la limite du nombre tâches

			// On attend qu'une tache se libère
			jobIndex, ok := <-jobsDone
			if !ok {
				// Ce cas ne devrait pas apparaitre, car l'alimentation tarira avant pour sortir de la boucle.
				break
			}
			id := jobsCtrl[jobIndex]

			// On libère l'index du job
			jobsCtrl[jobIndex] = 0

			// Le job s'est libéré. On pourra en démarrer 1 nouveau
			jobs--

			// On informe que la tache est terminé.
			printStatus(id, jobs, jobsCtrl)
		}
	}

	// Attendre la fin des derniers jobs
	for jobs != 0 {
		jobIndex := <-jobsDone
		id := jobsCtrl[jobIndex]

		// On libère l'index du job
		jobsCtrl[jobIndex] = 0

		// Le job s'est libéré.
		jobs--

		// On informe que la tache est terminé.
		printStatus(id, jobs, jobsCtrl)
	}

	// On ferme le channel de fin d'execution par propreté
	close(jobsDone)

	fmt.Printf("%sFin du programme\n", strings.Repeat("\n", maxRoutines+1))
}
