# EmailFinder

## Install
```
git clone https://github.com/rix4uni/EmailFinder.git
cd EmailFinder
chmod +x email_finder.sh
```

## Usage

**Get Your Targets Email** change a `dell.com` with your target
```
curl -s https://raw.githubusercontent.com/rix4uni/EmailFinder/main/Emails/dell.com.txt
```

**Showing Emails**
```
bash email_finder.sh example.com
```

**Showing Keys**
```
bash email_finder.sh example.com | cut -f1 -d"@"
```

## Another website to collect Emails
```
for target in $(cat wildcards.txt | tr -d '\r');do curl -s "https://api.webscout.io/lookup/$target" | jq -r '.emails[].email' 2>/dev/null | anew $target.txt;done
```
