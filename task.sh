#!/bin/bash

echo "Launching job $1"
sleep $(((RANDOM % 10 + 3)))
exit $(((RANDOM % 2)))