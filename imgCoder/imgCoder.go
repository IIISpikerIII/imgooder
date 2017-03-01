package imgCoder

import (
    "image"
    "image/color"
    "fmt"
    "os"
    "log"
    "io"
    "math"
    "image/png"
)

const NCPU = 4
var countChanel int = 3

func ReadFileToImage(fileName string, bufSize int)  {

    file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    fi, err := file.Stat()
    if err != nil {
        log.Fatal(err)
    }

    bSize:= (bufSize*countChanel)
    size:= fi.Size()
    lines:= math.Ceil(float64(size)/float64(bSize))

    imgCreateFromFile(int(lines), bufSize, file)
}

func readFileToBuf(file *os.File, buf []byte) ([]byte, string) {
    n, err := file.Read(buf)

    if n > 0 {
        //str := string(buf[:n])
        //fmt.Println(str)

        return buf[:n], "ok"
    }

    if err == io.EOF {
        return make([]byte, 0), "ok"
    }
    if err != nil {
        log.Printf("read %d bytes: %v", n, err)
        return make([]byte, 0), "err"
    }

    return make([]byte, 0), "err"
}

func renderLine(c chan int, img *image.RGBA, x int, y int, buf []byte)  {
    pix:=0
    countCh:= countChanel
    for x1 := 0; x1 < x; x1++ {

        div:= (len(buf) - pix)
        if div <= 0 {
            continue
        }
        if div < countCh {
            countCh = div
        }

        switch countCh {
        case 1:
            img.Set(x1, y, color.RGBA{buf[pix], 0, 0, 255})
        case 2:
            img.Set(x1, y, color.RGBA{buf[pix], buf[pix+1], 0, 255})
        case 3:
            img.Set(x1, y, color.RGBA{buf[pix], buf[pix+1], buf[pix+2], 255})
        case 4:
            img.Set(x1, y, color.RGBA{buf[pix], buf[pix+1], buf[pix+2], buf[x1+3]})
        }
        pix=pix+countChanel
    }
    c <- y;
}

func imgCreateFromFile(lines int, bufSize int, file *os.File) {
    img := image.NewRGBA(image.Rect(0, 0, bufSize, lines))

    bufSizeRead := bufSize*countChanel
    var numCpu int = NCPU
    if(lines < NCPU){
        numCpu = lines
    }

    c := make(chan int, NCPU)
    var yCount int = 0

    for i := 0; i < numCpu; i++ {
        buf := make([]byte, bufSizeRead)
        n, error := readFileToBuf(file, buf)
        if (len(n) > 0) && (error == "ok") {
            go renderLine(c, img, bufSize, yCount, n)
            yCount++
        }
    }

    //после ответа канала, проверяем очередь, отправляем еще задачу либо канал на выход
    for i := 0; i < numCpu; {
        buf := make([]byte, bufSizeRead)
        fmt.Printf("answer: %d\n", <-c);

        n, error := readFileToBuf(file, buf)
        if (len(n) > 0) && (error == "ok") {
            go renderLine(c, img, bufSize, yCount, n)
            yCount++
        } else {
            i++
        }
    }

    f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
    defer f.Close()
    png.Encode(f, img)
}