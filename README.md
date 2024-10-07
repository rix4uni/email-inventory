## EmailFinder

Email OSINT tool

## Installation
```
go install github.com/rix4uni/emailfinder@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/emailfinder/releases/download/v0.0.1/emailfinder-linux-amd64-0.0.1.tgz
tar -xvzf emailfinder-linux-amd64-0.0.1.tgz
rm -rf emailfinder-linux-amd64-0.0.1.tgz
mv emailfinder ~/go/bin/emailfinder
```
Or download [binary release](https://github.com/rix4uni/emailfinder/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/emailfinder.git
cd emailfinder; go install
```

#### Filter domains.txt
```
cat bbp.txt | sed 's/^\*\.//' | sed 's/^\*//' | sed 's/^\.//' | grep -v "*" | unew -q domains.txt
```

## Usage
```

                           _  __ ____ _             __
  ___   ____ ___   ____ _ (_)/ // __/(_)____   ____/ /___   _____
 / _ \ / __  __ \ / __  // // // /_ / // __ \ / __  // _ \ / ___/
/  __// / / / / // /_/ // // // __// // / / // /_/ //  __// /
\___//_/ /_/ /_/ \__,_//_//_//_/  /_//_/ /_/ \__,_/ \___//_/

                                         Current emailfinder version v0.0.1

This CLI tool allows users to fetch emails from Google, DuckDuckGo, Bing, Yahoo, Yandex, Github search results.

Examples:
echo "domain.com" | emailfinder
cat domains.txt | emailfinder
cat domains.txt | emailfinder --exact-match
cat domains.txt | emailfinder --search-engine google, yandex

Usage:
  emailfinder [flags]
  emailfinder [command]

Available Commands:
  bing        Fetch emails from Bing search results for a given domain.
  completion  Generate the autocompletion script for the specified shell
  duckduckgo  Fetch emails from DuckDuckGo search results for a given domain.
  github      A brief description of your command
  google      Fetch emails from Google search results for a given domain.
  help        Help about any command
  saved       Fetches emails for given domains using https://github.com/rix4uni/EmailFinder/tree/main/Emails
  skymem      Fetch emails from a domain using skymem
  yahoo       Fetch emails from Yahoo search results for a given domain.
  yandex      Fetch emails from Yandex search results for a given domain.

Flags:
  -h, --help      help for emailfinder
  -u, --update    update emailfinder to latest version
  -v, --version   Print the version of the tool and exit.

Use "emailfinder [command] --help" for more information about a command.
```

## Usage Examples
add more email files to https://github.com/rix4uni/EmailFinder/tree/main/Emails
```
cat domains.txt | emailfinder skymem -s
```

## TODO
- cmd/github.go
- add flags for cmd/root.go