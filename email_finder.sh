#!/usr/bin/env bash

domain=$1
startpage=1
skymem=$(curl -s "http://www.skymem.info/srch?q=$domain" | grep '<a href="/domain/' | awk '{print $2}' | cut -c15-41)

# Check if skymem is empty
if [ -z "$skymem" ]; then
    exit 1
fi

no_of_emails=$(curl -s "http://www.skymem.info/domain/$skymem$startpage" | grep "This is the preview" | awk '{print $6}')

# Check if no_of_emails is empty or not a number
if ! [[ "$no_of_emails" =~ ^[0-9]+$ ]]; then
    exit 1
fi

lastpage=$((no_of_emails / 5))

company_name=${domain%%.*}
for (( i = startpage; i < lastpage; i++ ))
do
    emails=$(curl -s "http://www.skymem.info/domain/$skymem$i" | grep "$company_name" | grep '<a href="/srch?q=' | sed '1d' | sed -e 's/<[^>]*>//g' | sed -e 's/^[ \t]*//')
    echo "$emails"
done
