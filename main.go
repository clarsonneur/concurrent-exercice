package main

import (
	"fmt"
	"time"
)

const (
	maxRoutines = 3
)

func main() {

	jobsStarted := make(chan int)
	jobsDone := make(chan int, maxRoutines)

	go func() {
		// On alimente le canal des messages
		for i := 1; i <= 10; i++ {
			jobsStarted <- i
			fmt.Printf("%d ", i)
		}
		close(jobsStarted)
	}()

	jobs := 0
	for {
		if jobs < maxRoutines {
			// On peut lancer une routine tant qu'on est inférieur à maxRoutine
			i, ok := <-jobsStarted
			if !ok {
				// La source d'information est tarie. On arrête de démarrer de nouvelles taches.
				break
			}

			// On démarre la tache en parallèle avec le message
			go func(value int) {
				fmt.Printf("(%d)", i)
				time.Sleep(time.Second)
				jobsDone <- value
			}(i)

			// On indique qu'on a démarré 1 tache de plus.
			jobs++
		} else {
			// On a atteint la limite du nombre tâches

			// On attend qu'une tache se libère
			id, ok := <-jobsDone
			if !ok {
				// Ce cas ne devrait pas apparaitre, car l'alimentation tarira avant pour sortir de la boucle.
				break
			}
			// On informe que la tache I est terminé.
			fmt.Printf("(*%d*)\n", id)

			// Le job s'est libéré. On pourra en démarrer 1 nouveau
			jobs--
		}
	}

	// Attendre la fin des derniers jobs
	for jobs != 0 {
		id := <-jobsDone

		// On informe que la tache I est terminé.
		fmt.Printf("(*%d*)\n", id)

		// Le job s'est libéré.
		jobs--
	}

	// On ferme le channel de fin d'execution par propreté
	close(jobsDone)

	fmt.Println("Fin du programme")
}
