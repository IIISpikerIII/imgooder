package imgDecoder

import (
    "image"
    "os"
    "log"
    "image/png"
    "image/color"
    "fmt"
)

type Configuration struct {
    CountChanel int
    CountThread int
    WidthImg int
    OutImgFile string
    OutTxtFile string
}
var Conf Configuration

type ByteLine struct {
    Bt []byte
    NumLine int
}

/*
    Main function by convert imagefile to textfile
 */
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
    fmt.Println("OK!");
    fmt.Println("File: ", outFileName);
}

/*
    Create file from image
 */
func fileCreateFromImg(image image.Image, outFileName string)  {
    bounds := image.Bounds()
    w, lines := bounds.Max.X, bounds.Max.Y

    // Определяем количество каналов в зависимости от количества строк
    var numCpu int = Conf.CountThread
    if(lines < Conf.CountThread){
        numCpu = lines
    }

    c := make(chan *ByteLine, Conf.CountThread)
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

/*
    Reading line from imagefile
 */
func readLine(c chan *ByteLine, image image.Image, i, w int) {
    // Буфер для декодированных байт
    buf := make([]byte, w* Conf.CountChanel)
    len:=0
    res := new(ByteLine)
    res.NumLine = i;

    // Обход по пикселям
    for x := 0; x < w; x++ {
        pix:= image.At(x, i)
        //r, g, b, a := pix.RGBA()
        //fmt.Print(r, " ", g, " ", b, " ", a, " ");
        //rgbacol := color.NRGBAModel.Convert(pix.(color.NRGBA))

        r := pix.(color.NRGBA).R
        g := pix.(color.NRGBA).G
        b := pix.(color.NRGBA).B
        a := pix.(color.NRGBA).A

        switch Conf.CountChanel {
            case 1 :
                buf, len = addByteBuf(buf, []uint8{r}, len)
            case 2 :
                buf, len = addByteBuf(buf, []uint8{r, g}, len)
            case 3 :
                buf, len = addByteBuf(buf, []uint8{r, g, b}, len)
            case 4 :
                buf, len = addByteBuf(buf, []uint8{r, g, b, a}, len)
        }
    }
    fmt.Print(".")
    res.Bt = buf[:len]
    c<- res
}

func addByteBuf(buf []byte, bytes []uint8, start int) ([]byte, int){
    for  _, value := range bytes {
        if value != 0 {
            buf[start] = byte(value)
            start++
        }
    }
    return buf, start;
}