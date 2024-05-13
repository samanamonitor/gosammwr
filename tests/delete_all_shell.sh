#!/bin/bash

SHELL_IDS=$(./test_shell_list)
for ID in ${SHELL_IDS}; do
	./test_shell_delete $ID
done