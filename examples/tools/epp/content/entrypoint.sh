#!/bin/bash
if [ -z "$1" ]
then
  files=$(find . -type f -name '*.epp')
elif [ "$1" = "help" ]
then
  echo "Please specify a file glob pattern to search for EPP files, such as '**/*.epp'"
  exit
else
  cmd="find ."
  for var in "$@"
  do
    cmd="$cmd -o -type f -name '$var'"
  done
  cmd=$(echo "$cmd" | sed '0,/-o/s///')
  files=$(eval ${cmd})
fi
echo "Validating EPP for ${files}"
puppet epp validate --continue_on_error $files
