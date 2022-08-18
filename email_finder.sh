#!/usr/bin/env bash

domain=$1
startpage=1
skymem=$(curl -s http://www.skymem.info/srch?q=$1 | grep '<a href="/domain/' | awk '{print $2}' | cut -c15-41)
no_of_emails=$(curl -s http://www.skymem.info/domain/$skymem$startpage | grep "This is the preview" | awk '{print $6}')
let lastpage=$no_of_emails/5

company_name=${domain%%.*}
for (( i=$startpage; i<$lastpage; i++))
do  
    emails=$(curl -s http://www.skymem.info/domain/$skymem$i | grep "$company_name" | grep '<a href="/srch?q=' | sed '1d' | sed -e 's/<[^>]*>//g' | sed -e 's/^[ \t]*//')
    echo "$emails"
done
