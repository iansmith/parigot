package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image/color"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/protobuf/proto"
)

var envVar string
var sockAddr string

var logCh = make(chan string, 32)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cleanupAndExit()
		}
	}()
	go processNetwork()
	launchUI()

}

func processNetwork() {
	sockAddr = ":4004"

	l, e := net.Listen("tcp", sockAddr)
	if e != nil {
		log.Printf("%v", e)
		cleanupAndExit()
	}
	internalMessage(fmt.Sprintf("waiting for connection on %s", sockAddr))
	for {
		conn, e := l.Accept()
		if e != nil {
			log.Printf("%v", e)
			cleanupAndExit()
		}
		internalMessage(fmt.Sprintf("handling new connection on %s", sockAddr))
		go handler(conn)
	}

}

func handler(conn net.Conn) {
	dataBuffer := make([]byte, netconst.FrontMatterSize+netconst.TrailerSize+netconst.ReadBufferSize)
	magicBuffer := dataBuffer[0:8]
	lenBuffer := dataBuffer[8:12]

	for {
		read, err := conn.Read(magicBuffer)
		if err != nil {
			internalMessage(fmt.Sprintf("[disconnected from %s", sockAddr))
			return
		}
		if read != 8 {
			internalMessage(fmt.Sprintf("[bad size of read for magic number: %d expected 8, closing %s]", read, sockAddr))
			return
		}
		magicNum := binary.LittleEndian.Uint64(magicBuffer)
		if magicNum != netconst.MagicStringOfBytes {
			internalMessage(fmt.Sprintf("[bad magic number, got %x expected %x, closing %s]", magicNum, netconst.MagicStringOfBytes, sockAddr))
		}
		read, err = conn.Read(lenBuffer)
		if err != nil {
			internalMessage(fmt.Sprintf("[disconnected from %s", sockAddr))
			return
		}
		if read != 4 {
			internalMessage(fmt.Sprintf("[bad size of read for length: %d expected 4, closing %s]", read, sockAddr))
			return
		}
		l := binary.LittleEndian.Uint32(lenBuffer)
		length := int(l)
		if length >= len(dataBuffer) {
			internalMessage(fmt.Sprintf("[data is too large: %d expected no more than %d, closing %s]", length, netconst.ReadBufferSize, sockAddr))
			return
		}
		objBuffer := dataBuffer[netconst.FrontMatterSize : length+netconst.FrontMatterSize]
		read, err = conn.Read(objBuffer)
		if err != nil {
			internalMessage(fmt.Sprintf("[disconnected from %s", sockAddr))
			return
		}
		var req logmsg.LogRequest
		err = proto.Unmarshal(objBuffer, &req)
		if err != nil {
			internalMessage(fmt.Sprintf("[unable to unmarshal data from socket: %v, closing %s]", err, sockAddr))
			return
		}
		crcBuffer := dataBuffer[length+netconst.FrontMatterSize : length+netconst.FrontMatterSize+4]

		read, err = conn.Read(crcBuffer)
		if err != nil {
			internalMessage(fmt.Sprintf("[disconnected from %s", sockAddr))
			return
		}
		result := crc32.Checksum(objBuffer, netconst.KoopmanTable)
		expected := binary.LittleEndian.Uint32(crcBuffer)
		if expected != result {
			internalMessage(fmt.Sprintf("[bad crc: expected %x but got %x, with CRC over %d bytes, closing %s]", expected, result, length, sockAddr))
			return
		}
		logMessage(&req)
	}
}

type myTheme struct{}

var myThemeInst myTheme

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	//log.Printf("color called with '%s', '%v'", name, variant)
	return theme.DefaultTheme().Color(name, variant)
}
func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	log.Printf("font request: %#v", style)
	return theme.DefaultTheme().Font(style)
}
func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	//log.Printf("size called with '%s'", name)
	if theme.SizeNamePadding == name {
		return 0.0
	}
	return theme.DefaultTheme().Size(name)
}
func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

const topBottomMargin = 10.0
const fontSize = 8.0 /* points?*/
const windowWidth = 600.0
const windowHeight = 400.0

func launchUI() {
	a := app.New()
	a.Settings().SetTheme(myThemeInst)
	w := a.NewWindow("Parigot Logviewer")

	boxCont := container.NewVBox()
	top := widget.NewLabel(" ")
	bottom := widget.NewLabel(" ")
	top.Resize(fyne.NewSize(100.0, topBottomMargin))
	bottom.Resize(fyne.NewSize(100.0, topBottomMargin))

	//boxCont.Resize(fyne.Size{Width: 600.0, Height: 400})

	go func() {
		for {
			s := <-logCh
			for strings.HasSuffix(s, "\n") {
				s = s[:len(s)-1]
			}
			label := widget.NewLabel(s)
			label.TextStyle = fyne.TextStyle{Monospace: true}
			boxCont.Add(label)
			//currentSize := fyne.MeasureText(s, fontSize, fyne.TextStyle{Monospace: true})
			boxCont.Refresh()

		}
	}()
	scrollCont := container.NewVScroll(boxCont)
	cont := container.NewBorder(top, bottom, nil, nil, scrollCont)

	w.SetContent(cont)
	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.ShowAndRun()
}

func internalMessage(s string) {
	log.Printf("xxx internal message %s", s)
}

func logMessage(req *logmsg.LogRequest) {
	s := fmt.Sprintf("%s:%d:%s", req.Stamp.AsTime().Format(time.RFC3339), req.Level, req.Message)
	logCh <- s
}

func cleanupAndExit() {
	log.Printf("closing socket: %s", sockAddr)
	os.Exit(5)
}
