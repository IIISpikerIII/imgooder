package main

import (
    "github.com/IIISpikerIII/imgooder/imgCoder"
    "github.com/IIISpikerIII/imgooder/imgDecoder"
    "flag"
)

type Configuration struct {
    countChanel int
    countThread int
    widthImg int
    outImgFile string
    outTxtFile string
}
var conf Configuration

func init() {
    conf.countChanel = 3
    conf.countThread = 4
    conf.widthImg = 100
    conf.outImgFile = "out.png"
    conf.outTxtFile = "out.txt"
}

func main() {

    // Coder config
    codeFile := flag.String("code", "", "File convert to image")
    width := flag.Int("w", conf.widthImg, "Width out file")
    outImgFile := flag.String("img", conf.outImgFile, "Image out file")

    // Decode config
    decodeFile := flag.String("decode", "", "Image convert to file")
    outFile := flag.String("out", conf.outTxtFile, "File out from decode")

    h := flag.Bool("h", false, "Help by command")
    help := flag.Bool("help", false, "Help by command")
    ch := flag.Int("ch", conf.countChanel, "Count chanels")
    flag.Parse()

    conf.countChanel = *ch
    conf.widthImg = *width
    conf.outTxtFile = *outFile
    conf.outImgFile = *outImgFile

    if len(*codeFile) > 4 {
        imgCoder.Conf.OutImgFile = conf.outImgFile;
        imgCoder.Conf.CountChanel = conf.countChanel;
        imgCoder.Conf.CountThread = conf.countThread;
        imgCoder.ReadFileToImage(*codeFile, conf.widthImg)
    }

    if len(*decodeFile) > 4 {
        imgDecoder.Conf.CountChanel = conf.countChanel;
        imgDecoder.Conf.CountThread = conf.countThread;
        imgDecoder.ReadImageToFile(*decodeFile, *outFile)
    }

    if *h || *help {
        flag.PrintDefaults()
    }
}
