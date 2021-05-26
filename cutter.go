package main

import (
	"github.com/oliamb/cutter"
	"image"
	"image/png"
	"os"
	"log"
	"math"
	"strings"
	"bufio"
	"strconv"
	"io/ioutil"
	"container/list"
	"fmt"
)
const errorMAX = 3

func verificadorDeCasco(xmin float64, ymin float64, xmax float64, ymax float64, nameTexFile string, classFilter string) (estado bool){
	text, err := os.Open(nameTexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer text.Close()

	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		if splitLine := strings.Split(scanner.Text(), " "); splitLine[0] == classFilter{
			xmin_aux,_:= strconv.ParseFloat(splitLine[4], 64)
			ymin_aux,_:= strconv.ParseFloat(splitLine[5], 64)
			xmax_aux,_:= strconv.ParseFloat(splitLine[6], 64)		
			ymax_aux,_:= strconv.ParseFloat(splitLine[7], 64)
			errorExp_xmin := math.Abs(((xmin - xmin_aux)/xmin)*100)
			errorExp_ymin := math.Abs(((ymin - ymin_aux)/ymin)*100)
			errorExp_xmax := math.Abs(((xmax - xmax_aux)/xmax)*100)
			if (xmin <= xmin_aux || errorExp_xmin <= errorMAX) && (ymin <= ymin_aux || errorExp_ymin <= errorMAX) && 
				(xmax >= xmax_aux || errorExp_xmax <= errorMAX) && (ymax >= ymax_aux){
					return true
					break
			}
		}
	}
	return false
}

func magic(imagen string, classFilter string){
	imagenFile, err := os.Open(imagen)
	if err != nil {
		log.Fatal(err)
	}
	defer imagenFile.Close()

	imagenText := strings.Replace(strings.Replace(imagen,".png",".txt",1),"dataJoin","labels",1)
	text, err := os.Open(imagenText)
	if err != nil {
		log.Fatal(err)
	}
	defer text.Close()

	img, err := png.Decode(imagenFile)
	if err != nil {
		log.Fatal(err)
	}
	cnt := 0
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		if splitLine := strings.Split(scanner.Text(), " "); splitLine[0] == "persona"{
			xmin,_:= strconv.ParseFloat(splitLine[4], 64)
			ymin,_:= strconv.ParseFloat(splitLine[5], 64)
			xman,_:= strconv.ParseFloat(splitLine[6], 64)		
			ymax,_:= strconv.ParseFloat(splitLine[7], 64)
			width := int(math.Round(xman - xmin ))
			height := int(math.Round(ymax-ymin))
			w_init := int(math.Round(xmin))
			h_init := int(math.Round(ymin))

			cImg, err := cutter.Crop(img, cutter.Config{
				Width:  width,
				Height: height,
			  Anchor:  image.Point{w_init, h_init},
				Options: 0,
			  })
			  if err != nil {
				  log.Fatal("Cannot crop image:", err)
			  }


			nameImage := strings.Split(imagen, "/")
			newName := strings.Replace(nameImage[len(nameImage)-1],".png","",1)+"_cut"+strconv.Itoa(cnt)+".png"

			directory:= "sin_"+classFilter+"/"
			if verificadorDeCasco(xmin,ymin,xman,ymax,imagenText, classFilter){
				directory = classFilter+"/"
			}

			outfile, err := os.Create(directory+newName)

			if err != nil {
				log.Fatal(err)
			}
			defer outfile.Close()

			err = png.Encode(outfile, cImg)
			if err != nil {
				log.Fatal(err)
			}
			cnt+=1
		}
	}
}

func listDirectoryRecursive(src string) (l_images *list.List) {
	l_img := list.New()
	archivos, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}
	for _, archivo := range archivos {
		if archivo.IsDir() {
			l_img.PushBackList(listDirectoryRecursive(src + "/" + archivo.Name()))
		} else {
			split_r := strings.Split(archivo.Name(), ".")
			extension := strings.ToLower(split_r[len(split_r)-1])
			if extension == "png" || extension == "jpg" || extension == "jpeg" {
				l_img.PushBack(src + "/" + archivo.Name())
			}
		}
	}
	return l_img
}
//Antes de hecharlo andar hay que crear dos carpetas conCasco y sin_conCasco
func main(){
	l := listDirectoryRecursive("dataset")
	largoLista:= l.Len()
	cnt_fotos:=1
	for e := l.Front(); e != nil; e = e.Next(){
		magic(e.Value.(string),"conCasco")
		if cnt_fotos%1000 == 0 {
			fmt.Println(strconv.Itoa((cnt_fotos*100)/largoLista) + "%")
		}
		
		cnt_fotos+=1
	}
}