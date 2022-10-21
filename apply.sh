#!/bin/bash

if [ "${BASH_SOURCE[0]}" -ef "$0" ]
then
    echo "This script should be sourced, not run."
    exit 1
fi

apply() {
	remote=$1
	config=$2
	tmpdir=$(mktemp -d --suffix=gen)
	git clone $remote $tmpdir
	./replace -config $config -dir $tmpdit
}
