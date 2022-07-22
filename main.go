package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Note struct {
	Volume  int
	Samples []int
}

type Song struct {
	Name      string
	samplemap []beep.Buffer
	Samplemap []string
	Pattern   []Note
	Bpm       uint
	spb       time.Duration
	//currentnote int
}

// func input(start chan struct{}, stop chan struct{}) {
// 	var input string
// 	for {
// 		fmt.Scanln(&input)
// 		switch input {
// 		case "stop":
// 			stop <- struct{}{}
// 		case "start":
// 			start <- struct{}{}
// 		}
// 	}
// }

func (s *Song) load() Song {
	s.spb = time.Minute / time.Duration(s.Bpm)

	var samples []beep.Buffer

	for _, samplefile := range s.Samplemap {
		f, err := os.Open(path.Join("Samples", samplefile))
		if err != nil {
			log.Fatal(err)
		}
		streamer, format, err := wav.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		buf := beep.NewBuffer(format)
		//beep.Resample(3,format.SampleRate)
		buf.Append(streamer)
		samples = append(samples, *buf)
	}
	s.samplemap = samples
	return *s
}

func (s Song) play() {
	format := beep.SampleRate(44100)
	speaker.Init(format, format.N(time.Second/10))

	for i := 0; true; i = (i + 1) % len(s.Pattern) {
		fmt.Println(i)
		for _, value := range s.Pattern[i].Samples {
			//sound := s.samplemap[value].Streamer(0,s.samplemap[value].Len())
			speaker.Play(s.samplemap[value].Streamer(0, s.samplemap[value].Len()))

		}
		time.Sleep(s.spb)
	}
}

func main() {
	//file loading
	var songs []Song

	file, err := os.Open("pattern.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &songs)

	var currentsong = songs[0].load()
	//fmt.Println(currentsong)
	//fmt.Printf("%#v", currentsong)
	//sequencer
	// start := make(chan struct{})
	// stop := make(chan struct{})
	//go input(start, stop)

	currentsong.play()
}
