#!/bin/bash

SHELL_IDS=$(bin/test_shell_list)
for ID in ${SHELL_IDS}; do
	bin/test_shell_delete $ID
done