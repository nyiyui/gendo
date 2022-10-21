#!/bin/bash

if [ "${BASH_SOURCE[0]}" -ef "$0" ]
then
    echo "This script should be sourced, not run."
    exit 1
fi

apply() {
	remote=$1
	config=$2
	commit_message=$3
	1>&2 echo "=== $remote"
	tmpdir=$(mktemp -d --suffix=gen)
	(
		git clone --quiet "$remote" $tmpdir
		./replace -config "$config" -dir $tmpdir
		cd $tmpdir
		git add .
		git commit -m "$commit_message"
		git push
	)
	rm -rf $tmpdir
}
