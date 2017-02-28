package main

import (
    "github.com/IIISpikerIII/imgooder/imgCoder"
    "github.com/IIISpikerIII/imgooder/imgDecoder"
    "flag"
    "fmt"
)

const NCPU = 4
var countChanel int = 2


func main() {

    codeFile := flag.String("code", "", "File convert to image")
    decodeFile := flag.String("decode", "", "File convert to image")
    flag.Parse()

    fmt.Println(*codeFile)

    if len(*codeFile) > 4 {
        imgCoder.ReadFileToImage(*codeFile, 100)
    }

    if len(*decodeFile) > 4 {
        imgDecoder.ReadImageToFile(*decodeFile, "out.txt")
    }
}
