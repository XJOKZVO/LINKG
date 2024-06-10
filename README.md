# LINKG
Tool that allows you to extract URLs from robots.txt, sitemap.xml, and all links from a given webpage, based on command-line options.

# Installation
```
go install github.com/XJOKZVO/LINKG@latest
```

# Options:
```
         _       ___   _   _   _  __   ____  
        | |     |_ _| | \ | | | |/ /  / ___| 
        | |      | |  |  \| | | ' /  | |  _  
        | |___   | |  | |\  | | . \  | |_| | 
        |_____| |___| |_| \_| |_|\_\  \____| 
                                                                           

Usage of ./main:
  -links
        Extract all links from the webpage
  -output string
        Output file name prefix
  -robots
        Extract URLs from robots.txt
  -sitemap
        Extract URLs from sitemap.xml
  -url string
        Website URL
```
# Usage:
```
./LINKG -url http://example.com -output example -robots -sitemap -links
```
