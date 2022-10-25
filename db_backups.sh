#!/bin/bash

PIDFILE=/tmp/backup.pid
TIMESTAMP=$(date '+%F-%H-%M')
LIMIT=1
if [ -f $PIDFILE ]; then
                if pgrep -F $PIDFILE &>/dev/null; then
                        exit 1
                        else echo $$ > $PIDFILE
                fi
        else echo $$ > $PIDFILE
fi

get_count_backups(){
COUNT=$(ls /u00/basebackup/ | grep "^20*" | wc -l)
}

echo "$TIMESTAMP:start backup"

rm -rf /u00/archive_wal/*

echo "$TIMESTAMP:delete wal"

mkdir -p /u00/basebackup/$TIMESTAMP > /dev/null 2>&1

echo "$TIMESTAMP:pg_basebackup start"

pg_basebackup -Ft -z -D /u00/basebackup/$TIMESTAMP

echo "$TIMESTAMP:pg_basebackup end"

get_count_backups

echo "$TIMESTAMP:check count backups"

while [ "$COUNT" -gt "$LIMIT" ]
	do
		last_backup=$(ls /u00/basebackup/ | grep "^20*" | head -n1 )
		rm -rf /u00/basebackup/$last_backup;
		echo "remove $last_backup"
		get_count_backups
	done

echo "$TIMESTAMP:end backup"
