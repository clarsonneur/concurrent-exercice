# Introduction

Ce programme est écrit pour répondre à un test de programmation 
de type messaging avec eventuellement la parallisation du traitement de 
ces messages


## Execution du test

1. lancez le programme avec un `go run main.go`

2. Créez autant de fichier que voulu dans le sous répertoire `queue`

## Output typique

```
Create a 'queue' directory where you run the program and put one or more files in it. Only formatted files with 'job:xx:process' and 'exit' will be treated.
Trying job 24 (1/3)
Trying job 34 (1/3)
Trying job 26 (1/3)
Launching job 34
Launching job 24
Launching job 26
Job 24 executed successfully.
Trying job 424 (1/3)
Launching job 424
Retrying job 26 (2/3)
Trying job 26 (2/3)
Launching job 26
Job 34 executed successfully.
Trying job 425 (1/3)
Launching job 425
Retrying job 26 (3/3)
Trying job 26 (3/3)
Launching job 26
Retrying job 424 (2/3)
Trying job 424 (2/3)
Launching job 424
Unable to run task 26
Job 425 executed successfully.
Retrying job 424 (3/3)
Trying job 424 (3/3)
Launching job 424
Unable to run task 424
Done
```