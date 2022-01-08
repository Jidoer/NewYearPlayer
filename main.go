// @Title main.go
// @Description  我用Golang摆烂了你的AE作品
// @Author  FlyKO  2022/01/08
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)


func init() {
	logFile, err := os.OpenFile("Images.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetPrefix(">")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

}


func main() {
	var images []image.Image

	dir_list, e := ioutil.ReadDir("res")
	if e != nil {
		fmt.Println("read dir error")
	}
	for i, v := range dir_list {

		log.Println(i, v.Name()+strconv.Itoa(i))

		file, err := os.Open("res/" + v.Name())
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		f, err := os.Open("res/" + v.Name())
		if err != nil {
			panic(err)
		}
		img, formame, err := image.Decode(f)
		if err != nil {
			panic(err)
		}
		//image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(GOPHER_IMAGE)))
		images = append(images, img)
		fmt.Println(strconv.Itoa(i)+": 欢迎 《点赞.投币.收藏》 ADD:[ok]:", formame)
	}

	//Play MP3
	
	f, err := os.Open("audio.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))



	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	img := widgets.NewImage(nil)
	img.SetRect(0, 0, 150, 50)
	index := 0
	render := func() {
		img.Image = images[index]
		if !img.Monochrome {
			img.Title = fmt.Sprintf("schedule %d/%d", index+1, len(images))
		} else if !img.MonochromeInvert {
			img.Title = fmt.Sprintf("Monochrome(%d) %d/%d", img.MonochromeThreshold, index+1, len(images))
		} else {
			img.Title = fmt.Sprintf("InverseMonochrome(%d) %d/%d", img.MonochromeThreshold, index+1, len(images))
		}
		ui.Render(img)
	}
	render()

	go func() {
		for {
			index = (index + 1) % len(images)
			time.Sleep(time.Microsecond*60)
			render()
			log.Println(index)
			if(index >= 864){
				os.Exit(0)
			}
		}
	}()

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		//按Q||Ctrl+c退出播放
		case "q", "<C-c>":
			return
		}

	}

}
