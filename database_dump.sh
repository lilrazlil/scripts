#!/bin/bash

PIDFILE=/tmp/backup.pid
TIMESTAMP=$(date '+%F-%H:%M')
LIMIT=1
DIRECTORY=./dump_db
DATABASES=(aboba biba)
DUMP=$1

if [ -f $PIDFILE ]; then
                if pgrep -F $PIDFILE &>/dev/null; then
                        exit 1
                        else echo $$ > $PIDFILE
                fi
        else echo $$ > $PIDFILE
fi

get_count_backups(){
COUNT=$(ls $DIRECTORY/$DUMP | grep "^20*" | wc -l)
}

echo "$TIMESTAMP:start dump"

mkdir -p $DIRECTORY/$DUMP/$TIMESTAMP > /dev/null 2>&1

echo "$TIMESTAMP:pg_dump start"

for i in "${!DATABASES[@]}"; do
DB=${DATABASES[$i]}
echo "$TIMESTAMP:${DB} DB start pg_dump"
pg_dump ${DB} | gzip > $DIRECTORY/$DUMP/$TIMESTAMP/${DB}_$(date '+%F-%H:%M').sql.gz
echo "$TIMESTAMP:${DB} DB finish pg_dump"
done

echo "$TIMESTAMP:dump end"

get_count_backups

echo "$TIMESTAMP:check count backups"

while [ "$COUNT" -gt "$LIMIT" ]
	do
		last_backup=$(ls $DIRECTORY/$DUMP | grep "^20*" | head -n1 )
		rm -rf $DIRECTORY/$DUMP/$last_backup;
		echo "remove $last_backup"
		get_count_backups
	done

echo "$TIMESTAMP:end backup"



