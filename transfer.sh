#!/bin/sh
#system:paperino@192.168.1.123/xe
ORA_URI=$1
INPUT_FOLDER=$2
OUTPUT_FOLDER=$3
PG_USER=$4
PG_PWD=$5
PG_HOST=$6
PG_SCHEMA=$7
LOG_FILE=$8

if [ -z "$LOG_FILE" ] 
then
   LOG_FILE=transfer.log
fi

Help()
{
   # Display Help
   echo "Export utility"
   echo
   echo "Syntax: $0 <ora-user:ora-password@ora-host/ora-service> <input-folder> <output-folder> <pg-user> <pg-pwd> <pg-host> <pg-schema> <log-file>"
   echo
}

# Get the options
while getopts ":h" option; do
   case $option in
      h) # display Help
         Help
         exit;;
   esac
done

now=`date '+%d/%m/%Y_%H:%M:%S'`
echo "[$(date '+%d/%m/%Y %H:%M:%S')] Start transfer job" > $LOG_FILE
echo "[$(date '+%d/%m/%Y %H:%M:%S')] Export oracle data" >> $LOG_FILE
echo "   Input folder: $INPUT_FOLDER" >> $LOG_FILE
echo "   Output folder: $OUTPUT_FOLDER" >> $LOG_FILE
./simple-db-exporter export $ORA_URI $INPUT_FOLDER $OUTPUT_FOLDER >> $LOG_FILE

if [ $? -ne 0 ]
then
   echo "Error exporting from oracle" >> $LOG_FILE
   exit 1
fi
echo
echo "[$(date '+%d/%m/%Y %H:%M:%S')] Import exported sql file from $INPUT_FOLDER" >> $LOG_FILE
for filename in output/*.sql; do
   REAL_FILE_PATH=$(cd $(dirname $filename); pwd)/$(basename $filename);
   echo "[$(date '+%d/%m/%Y %H:%M:%S')] Import $REAL_FILE_PATH" >> $LOG_FILE
   # PGPASSWORD=$PG_PWD psql -U $PG_USER -h $PG_HOST -d $PG_SCHEMA < $REAL_FILE_PATH
   # if [ $? -ne 0 ]
   # then
   #    echo "[$(date '+%d/%m/%Y %H:%M:%S')] error during data import" >> $LOG_FILE
   #    exit 1
   # fi
done

echo "[$(date '+%d/%m/%Y %H:%M:%S')] work done" >> $LOG_FILE