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
    decodeFile := flag.String("decode", "", "Image convert to file")
    outFile := flag.String("out", "out.txt", "File out from decode")
    h := flag.Bool("h", false, "Help by command")
    help := flag.Bool("help", false, "Help by command")
    flag.Parse()

    fmt.Println(*codeFile)

    if len(*codeFile) > 4 {
        imgCoder.ReadFileToImage(*codeFile, 100)
    }

    if len(*decodeFile) > 4 {
        imgDecoder.ReadImageToFile(*decodeFile, *outFile)
    }

    if *h || *help {
        flag.PrintDefaults()
    }
}
