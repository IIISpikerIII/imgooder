package imgDecoder

import (
    "image"
    "os"
    "log"
    "math"
    "image/png"
)

const NCPU = 4
var countChanel int = 2

type ByteLine struct {
    Bt []byte
    NumLine int
}

func ReadImageToFile(fileName string, outFileName string)  {

    img, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer img.Close()

    pngImage, _ := png.Decode(img)
    if err != nil {
        log.Fatal(err)
    }

    fileCreateFromImg(pngImage, outFileName)
}

func fileCreateFromImg(image image.Image, outFileName string)  {
    bounds := image.Bounds()
    w, lines := bounds.Max.X, bounds.Max.Y

    var numCpu int = NCPU
    if(lines < NCPU){
        numCpu = lines
    }

    c := make(chan *ByteLine, NCPU)
    var yCount int = 0

    for i := 0; i < numCpu; i++ {
        go readLine(c, image, yCount, w)
        yCount++
    }

    f, _ := os.OpenFile(outFileName, os.O_RDWR|os.O_CREATE, 0600)
    defer f.Close()

    buf := make([][]byte , lines)
    //после ответа канала, проверяем очередь, отправляем еще задачу либо канал на выход
    for i := 0; i < numCpu; {
        //fmt.Printf("answer: %d\n", <-c);
        res:= <-c
        buf[res.NumLine] = res.Bt
        if (yCount < lines) {
            go readLine(c, image, yCount, w)
            yCount++
        } else {
            i++
        }
    }

    for _, n := range buf {
        f.Write(n)
    }
}

func readLine(c chan *ByteLine, image image.Image, i, w int) {
    buf := make([]byte, w)
    len:=0
    res := new(ByteLine)
    res.NumLine = i;
    for x := 0; x < w; x++ {
        pix:= image.At(x, i)
        r, _, _, _ := pix.RGBA()
        if r != 0 {
            len++
            bt:= math.Floor(float64(r)/float64(256))
            buf[x] = byte(bt)
        }
    }
    res.Bt = buf[:len]
    c<- res
}