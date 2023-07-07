package main

import (
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SplashScene struct {
	start        time.Time
	chFromScript chan string
	chToScript   chan string
	msgs         []string
	runScript    bool
}

func NewSplashScene() *SplashScene {
	return &SplashScene{
		start:        time.Time{},
		chFromScript: make(chan string),
		chToScript:   make(chan string),
		msgs:         make([]string, 0, 10),
		runScript:    false,
	}
}

func (s *SplashScene) AddMsg(msg string) {
	s.msgs = append(s.msgs, msg)
	if len(s.msgs) > 10 {
		s.msgs = s.msgs[1:]
	}
}

func (s *SplashScene) Update(mgr *SceneMgr) error {
	if !s.runScript {
		mgr.AddScene(SceneGame, NewGameScene(nil))
		s.runScript = true
		go s.Script(s.chFromScript)
	}

loop:
	for {
		select {
		case msg, ok := <-s.chFromScript:
			if !ok {
				// BUG: this will panic if AddScene isn't finished when Script() closes the channel
				if err := mgr.SwitchScene(SceneGame); err != nil {
					return err
				}
				break loop
			}
			s.AddMsg(msg)
		default:
			break loop
		}
	}

	return nil
}

func (s *SplashScene) Draw(mgr *SceneMgr, screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, strings.Join(s.msgs, "\n"))
}

func (s *SplashScene) Script(ch chan string) {
	start := time.Now()
	end := start.Add(time.Second * 2)

	ch <- "Loading..."

	time.Sleep(time.Second)
	ch <- "Debug 1..."
	ch <- "Debug 2..."
	ch <- "Debug 3..."
	ch <- "Debug 4..."
	ch <- "Debug 5..."
	ch <- "Debug 6..."
	ch <- "Debug 7..."
	ch <- "Debug 8..."
	ch <- "Debug 9..."
	ch <- "Still Loading..."

	if time.Now().Before(end) {
		time.Sleep(end.Sub(time.Now()))
	}

	close(ch)
}
