package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"syscall"
)

const (
	maxRoutines = 3
	queuePath   = "queue"
	jobScript   = "task.sh"
)

// ListenDir vérifie le contenu d'un répertoire et transmet l'ID extrait du fichier.
func ListenDir(jobsStarted chan string) {
	reg, _ := regexp.Compile(`^job:([0-9]+):process$`)

	exit := false
	for !exit {
		err := filepath.Walk(queuePath, func(pathGiven string, info os.FileInfo, _ error) (_ error) {
			if info == nil {
				return
			}
			if info.IsDir() {
				return
			}
			file := path.Base(pathGiven)
			if file == "exit" {
				// Arrêt du processus d'alimentation
				exit = true
			}
			if v := reg.FindAllStringSubmatch(file, 1); v != nil {
				// On alimente le canal des messages
				jobsStarted <- v[0][1]
				err := os.Remove(pathGiven)
				if err != nil {
					return fmt.Errorf("Unable to remome message %s. %s", pathGiven, err)
				}
			}
			return
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
	close(jobsStarted)

}

// Run the job
func runJob(value string, jobsDone chan bool) {

	retryNum := 3
	for retryNum > 0 {
		task := exec.Command("bash", jobScript, value)
		outReader, _ := task.StdoutPipe()

		go func() {
			outScanner := bufio.NewScanner(outReader)
			for outScanner.Scan() {
				fmt.Println(outScanner.Text())
			}
		}()

		fmt.Printf("Trying job %s (%d/3)\n", value, 4-retryNum)
		task.Run()

		if status := task.ProcessState.Sys().(syscall.WaitStatus); status.ExitStatus() != 0 {
			retryNum--
			if retryNum == 0 {
				fmt.Printf("Unable to run task %s\n", value)
			} else {
				fmt.Printf("Retrying job %s (%d/3)\n", value, 4-retryNum)
			}
			continue
		}
		fmt.Printf("Job %s executed successfully.\n", value)
		break
	}
	jobsDone <- true
}

func main() {

	jobsStarted := make(chan string)
	jobsDone := make(chan bool, maxRoutines)

	go ListenDir(jobsStarted)

	jobs := 0

	fmt.Println("Create a 'queue' directory where you run the program and put one or more files in it. Only formatted files with 'job:xx:process' and 'exit' will be treated.")
	for {
		if jobs < maxRoutines {
			// On peut lancer une routine tant qu'on est inférieur à maxRoutine
			i, ok := <-jobsStarted
			if !ok {
				// La source d'information est tarie. On arrête de démarrer de nouvelles taches.
				break
			}

			// On démarre la tache en parallèle avec le message
			go runJob(i, jobsDone)

			// On indique qu'on a démarré 1 tache de plus.
			jobs++
		} else {
			// On a atteint la limite du nombre tâches

			// On attend qu'une tache se libère
			_, ok := <-jobsDone
			if !ok {
				// Ce cas ne devrait pas apparaitre, car l'alimentation tarira avant pour sortir de la boucle.
				break
			}

			// Le job s'est libéré. On pourra en démarrer 1 nouveau
			jobs--
		}
	}

	// Attendre la fin des derniers jobs
	for jobs != 0 {
		<-jobsDone

		// Le job s'est libéré.
		jobs--
	}

	// On ferme le channel de fin d'execution par propreté
	close(jobsDone)

	fmt.Println("Done")
}
