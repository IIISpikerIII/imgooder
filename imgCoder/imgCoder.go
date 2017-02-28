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
var countChanel int = 2

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

    size:= fi.Size()
    lines:= math.Ceil(float64(size)/float64(bufSize))

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
    for x1 := 0; x1 < x; x1++ {
        if len(buf) > x1 {
            img.Set(x1, y, color.RGBA{buf[x1], 0, 0, 255})
        }
    }
    c <- y;
}

func imgCreateFromFile(lines int, bufSize int, file *os.File) {
    img := image.NewRGBA(image.Rect(0, 0, bufSize, lines))

    var numCpu int = NCPU
    if(lines < NCPU){
        numCpu = lines
    }

    c := make(chan int, NCPU)
    var yCount int = 0

    for i := 0; i < numCpu; i++ {
        buf := make([]byte, bufSize)
        n, error := readFileToBuf(file, buf)
        if (len(n) > 0) && (error == "ok") {
            go renderLine(c, img, bufSize, yCount, n)
            yCount++
        }
    }

    //после ответа канала, проверяем очередь, отправляем еще задачу либо канал на выход
    for i := 0; i < numCpu; {
        buf := make([]byte, bufSize)
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