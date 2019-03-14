#!/bin/sh
# 
# This script aims to add a key to all values of the JSON that was output by Safety
# 

file=$1
tmpFile=tmp.json
outputFile=output.json
i=0
k=0

while read -r line; do
    if [ $i -eq 0 ]; then
        printf "{\"issues\":[" > $tmpFile
        read -r line
        if [ "$line" == "" ]; then
            printf "]}" >> $tmpFile
            break
        fi
    fi

    if [ "$line" == "[" ]; then
        printf "{" >> $tmpFile
        for j in 0 1 2 3 4; do
            if [ $j -eq 0 ]; then
                read -r line
                printf "\"dependency\":$line" >> $tmpFile
            elif [ $j -eq 1 ]; then
                read -r line
                printf "\"vulnerable_below\":$line" >> $tmpFile
            elif [ $j -eq 2 ]; then
                read -r line
                printf "\"installed_version\":$line" >> $tmpFile
            elif [ $j -eq 3 ]; then
                read -r line
                printf "\"description\":$line" >> $tmpFile
            else
                read -r line
                printf "\"id\":$line" >> $tmpFile
            fi
        done
	    read -r line
        printf "}," >> $tmpFile
    fi

    if [ "$line" == "]" ] && [ $k == 0 ]; then
        printf "]}" >> $tmpFile
        k=$((k+1))
    fi

    if [ $i == 0 ]; then
        i=$((i+1))
    fi

done < "$file"

cat $tmpFile | sed 's/\(.*\),/\1/' > $outputFile

rm $tmpFile

