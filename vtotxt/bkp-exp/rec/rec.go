package rec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	wave "github.com/zenwerk/go-wave"
	"math/rand"
	"os"
	"strings"
	"time"
)

func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}

var waveWriter wave.Writer
var stream *portaudio.Stream

func Record(fn string) {

	audioFileName := fn

	if !strings.HasSuffix(audioFileName, ".wav") {
		audioFileName += ".wav"
	}
	waveFile, err := os.Create(audioFileName)
	errCheck(err)

	// www.people.csail.mit.edu/hubert/pyaudio/  - under the Record tab
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	framesPerBuffer := make([]byte, 1024)

	// init PortAudio

	portaudio.Initialize()
	//defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	errCheck(err)
	//defer stream.Close()

	// setup Wave file writer

	param := wave.WriterParam{
		Out:           waveFile,
		Channel:       inputChannels,
		SampleRate:    sampleRate,
		BitsPerSample: 16, // if 16, change to WriteSample16()
	}

	waveWriter, err := wave.NewWriter(param)
	errCheck(err)

	//defer waveWriter.Close()

	st := time.Now()
	c := make(chan bool, 1)
	d := false
	go func() {
		diffTime(st, c)
		<-c
		d = true
		fmt.Println("Stopping!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		// better to control
		// how we close then relying on defer
		waveWriter.Close()
		stream.Close()
		portaudio.Terminate()
		//fmt.Println("Play", audioFileName, "with a audio player to hear the result.")
		//                         os.Exit(0)

	}()

	/*
	   // recording in progress ticker. From good old DOS days.
	   ticker := []string{
	           "-",
	           "\\",
	           "/",
	           "|",
	   }
	*/
	rand.Seed(time.Now().UnixNano())

	// start reading from microphone
	errCheck(stream.Start())
	for {
		errCheck(stream.Read())

		//fmt.Printf("\rRecording is live now. Say something to your microphone! [%v]", ticker[rand.Intn(len(ticker)-1)])

		// write to wave file
		//_, err := waveWriter.Write([]byte(framesPerBuffer)) // WriteSample16 for 16 bits
		c16, _ := convert16(framesPerBuffer)
		_, err := waveWriter.Write(c16) // WriteSample16 for 16 bits
		//_, err := waveWriter.WriteSample16(convert(framesPerBuffer)) // WriteSample16 for 16 bits
		errCheck(err)

		if d == true {
			fmt.Println("Break recording")
			break
		}
	}
	errCheck(stream.Stop())
}

func convert16(samples []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	for i := 0; i < len(samples); i++ {
		err := binary.Write(buf, binary.LittleEndian, samples[i])
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func convert(raw []byte) []uint16 {
	sizeuint := 2 // bytes

	data := make([]uint16, len(raw)/sizeuint)
	for i := range data {
		// assuming little endian
		data[i] = uint16(binary.LittleEndian.Uint16(raw[i*sizeuint : (i+1)*sizeuint]))
	}
	return data
}

func diffTime(st time.Time, c chan bool) {

	for {
		d := time.Since(st)
		//fmt.Println(d.Seconds())
		if d.Seconds() > 2.0 {
			c <- true
			break
		}
	}
}

func Stop() {
	waveWriter.Close()
	stream.Close()
	portaudio.Terminate()
	//    os.Exit(0)
}

/*
func main() {
  Record("hello.wav")
}
*/
