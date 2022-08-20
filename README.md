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
