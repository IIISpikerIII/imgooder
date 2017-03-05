# imgooder
Code/decode img files

##Quick Start

###File to image
```
cd path/to/imgooder
./bin/imgooder_darwin_amd64 -code doc/test.txt
```
Console output
```
...........................OK!
File:  out.png
```
**out.png**

!["graph data buble viewer"](https://github.com/IIISpikerIII/imgooder/tree/master/doc/out.png?raw=true)

###Image to file
```
cd path/to/imgooder
./bin/imgooder_darwin_amd64 -decode doc/out.png
```
Console output
```
...........................OK!
File:  out.txt
```
**out.txt**

!["graph data buble viewer"](https://github.com/IIISpikerIII/imgooder/tree/master/doc/outtxt.png?raw=true)

##Params

**imgooder**

Convert file to image

**-code** File convert to image

**-w**  Width out file (default 100)

**-img** Image out file (default "out.png")

Convert image to file

**-decode** Image convert to file

**-out** File out from decode (default "out.txt")

Common params

**-h**  or **-help**  Help by command

**-ch** Count chanels (default 3)