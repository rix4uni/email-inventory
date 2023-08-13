# EmailFinder

## Install
```
git clone https://github.com/rix4uni/EmailFinder.git
cd EmailFinder
chmod +x email_finder.sh
```

## Usage

**Showing Emails**
```
./email_finder.sh example.com
```
**Showing Keys**
```
./email_finder.sh example.com | cut -f1 -d"@"
```

## Another website to collect Emails
```
for target in $(cat wildcards.txt | tr -d '\r');do curl -s "https://api.webscout.io/lookup/$target" | jq -r '.emails[].email' 2>/dev/null | anew $target.txt;done
```
