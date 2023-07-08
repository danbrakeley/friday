package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	// This is hacky, but WritePixels is better than Fill in term of automatic texture packing.
	whiteImage.WritePixels(pix)
}

func drawVerticesForUtil(dst *ebiten.Image, vs []ebiten.Vertex, is []uint16, clr color.Color) {
	r, g, b, a := clr.RGBA()
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	op.AntiAlias = true
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}

func drawShape(dst *ebiten.Image, shape []Vec2D, strokeWidth float32, clr color.Color) {
	if len(shape) == 0 {
		return
	}

	var path vector.Path

	end := len(shape) - 1
	path.MoveTo(shape[end].X, shape[end].Y)
	for _, v := range shape {
		path.LineTo(v.X, v.Y)
	}

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)

	drawVerticesForUtil(dst, vs, is, clr)
}

func main() {
	defer midi.CloseDriver()

	cfg, err := LoadConfig("config.json")
	if os.IsNotExist(err) || len(cfg.MidiDevice) == 0 {
		cfg, err = CreateConfig()
		if err != nil {
			log.Fatal(err)
		}
		err = SaveConfig("config.json", cfg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Loaded config.json")
	}

	mgr := NewSceneMgr()
	midiMgr, err := NewMidiMgr(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer midiMgr.Close()
	// mgr.AddScene(SceneSplash, NewSplashScene())
	// mgr.SwitchScene(SceneSplash)
	mgr.AddScene(SceneGame, NewGameScene(midiMgr))
	mgr.SwitchScene(SceneGame)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("spin click (friday 02)")
	if err := ebiten.RunGame(mgr); err != nil {
		log.Fatal(err)
	}
}

func CreateConfig() (Config, error) {
	fmt.Println("MIDI device not configured.")

	fmt.Printf("\nListing MIDI devices...\n")
	inPorts := midi.GetInPorts()

	for i, port := range inPorts {
		fmt.Printf("%2d: %s (#%d)\n", i, port.String(), port.Number())
	}

	n := MustReadNumber(0, len(inPorts)-1, "\nChoose your MIDI device")
	port := inPorts[n]

	stop, err := midi.ListenTo(port, func(msg midi.Message, timestampms int32) {
		var ch, controller, value uint8
		switch {
		case msg.GetControlChange(&ch, &controller, &value):
			fmt.Printf("control change: channel=%v, controller=%v, value=%v\n", ch, controller, value)
		default:
			// ignore
		}
	})
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		MidiDevice: port.String(),
		Knobs:      make([]KnobConfig, KNOB_COUNT),
	}

	fmt.Printf("\nMIDI device %s active. Turn knobs to print control change messages.\n", port.String())

	for i := 0; i < KNOB_COUNT; i++ {
		cfg.Knobs[i].Channel = MustReadNumber(0, 15, fmt.Sprintf("\nChoose the Channel for Knob %d", i))
		cfg.Knobs[i].Controller = MustReadNumber(0, 127, fmt.Sprintf("\nChoose the Controller for Knob %d", i))
	}

	stop()

	b, err := json.MarshalIndent(cfg, "  ", "  ")
	if err != nil {
		return Config{}, err
	}
	fmt.Printf("\nConfig created: \n%s\n", string(b))

	return cfg, nil
}

func MustReadNumber(min, max int, msg string) int {
	punctuation := ":"
	if strings.HasSuffix(msg, "?") {
		punctuation = "?"
		msg = msg[:len(msg)-1]
	}

	fmt.Printf("%s [%d-%d]%s ", msg, min, max, punctuation)

try_again:
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// remove the delimeter from the string
	line = strings.TrimSpace(line)

	n, err := strconv.Atoi(line)
	if err != nil {
		fmt.Printf("Invalid input: %s\nPlease enter a number in the range [%d-%d]: ", err.Error(), min, max)
		goto try_again
	}

	if n < min || n > max {
		fmt.Printf("Invalid input: %d\nPlease enter a number in the range [%d-%d]: ", n, min, max)
		goto try_again
	}

	return n
}
