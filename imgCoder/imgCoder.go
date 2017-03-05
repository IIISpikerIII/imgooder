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

type Configuration struct {
    CountChanel int
    CountThread int
    WidthImg int
    OutImgFile string
    OutTxtFile string
}
var Conf Configuration

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

    bSize:= (bufSize* Conf.CountChanel)
    size:= fi.Size()
    lines:= math.Ceil(float64(size)/float64(bSize))

    imgCreateFromFile(int(lines), bufSize, file)
    fmt.Println("OK!");
    fmt.Println("File: ", Conf.OutImgFile);
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

func renderLine(c chan int, img *image.NRGBA, x int, y int, buf []byte)  {
    pix:=0
    countCh:= Conf.CountChanel
    //fmt.Println(buf);

    var alfa uint8 = 255
    for x1 := 0; x1 < x; x1++ {

        alfa = 255
        div:= (len(buf) - pix)
        if div <= 0 {
            continue
        }
        if div < countCh {
            countCh = div
            alfa = 0
        }

        switch countCh {
        case 1:
            img.Set(x1, y, color.NRGBA{buf[pix], 0, 0, alfa})
        case 2:
            img.Set(x1, y, color.NRGBA{buf[pix], buf[pix+1], 0, alfa})
        case 3:
            img.Set(x1, y, color.NRGBA{buf[pix], buf[pix+1], buf[pix+2], alfa})
        case 4:
            img.Set(x1, y, color.NRGBA{buf[pix], buf[pix+1], buf[pix+2], buf[pix+3]})
        }
        pix=pix+ Conf.CountChanel
    }
    fmt.Print(".");
    c <- y;
}

func imgCreateFromFile(lines int, bufSize int, file *os.File) {
    img := image.NewNRGBA(image.Rect(0, 0, bufSize, lines))

    bufSizeRead := bufSize* Conf.CountChanel
    var numCpu int = Conf.CountThread
    if(lines < Conf.CountThread){
        numCpu = lines
    }

    c := make(chan int, Conf.CountThread)
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
        //fmt.Printf("answer: %d\n", <-c);

        n, error := readFileToBuf(file, buf)
        if (len(n) > 0) && (error == "ok") {
            go renderLine(c, img, bufSize, yCount, n)
            yCount++
        } else {
            i++
        }
    }

    f, _ := os.OpenFile(Conf.OutImgFile, os.O_WRONLY|os.O_CREATE, 0600)
    defer f.Close()
    png.Encode(f, img)
}