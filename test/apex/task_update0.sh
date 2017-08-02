#/bin/bash

# Assumed that this will be invoked from repo root directory

ITEM=`test/apex/task_add0.sh`
echo $ITEM | sed 's/cool/sweet/g' | sed 's/bar/baz/g'| sed 's/foo/bar/g' | apex invoke task_update