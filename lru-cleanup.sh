#!/usr/bin/env bash

STOREFOLDER=/var/storedata/
OLDERTHEN=30 # in days

function cleanup() {
	echo "Cleaning up $1 , for files older then $2"
	if [ -d "$1" ];then
		find $1 -type f -atime -$2 -print |xargs -l1 rm -v
	else
		echo "$1 MISSING! , Aborted cleanup"
	fi
}

if [ -d "$STOREFOLDER" ];then
	echo "$STOREFOLDER exists"
else
	echo "$STOREFOLDER  MISSING!"
	exit
fi

case $1 in 
	requests)
		echo "cleaning up requets"
		cleanup $STOREFOLDER/request/ $OLDERTHEN
	;;
	responses)
	echo "responses"
		cleanup $STOREFOLDER/header/ $OLDERTHEN
		cleanup $STOREFOLDER/body/ $OLDERTHEN
	;;
	*)
	echo "The options are (requests|responses)"
	;;
esac
