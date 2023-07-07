package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// SceneID is a unique identifier for a scene.
// Negative values are reserved for internal use.
type SceneID int

// Scene is the minimal interface needed by a scene.
type Scene interface {
	Update(mgr *SceneMgr) error
	Draw(mgr *SceneMgr, screen *ebiten.Image)
}

// SceneWithLoader is a Scene that also requires loading that may
// take a non-trivial amount of time, e.g. loading assets from disk.
type SceneWithLoader interface {
	Scene

	// LoadSync is called in a Go routine (and thus must be thread safe).
	// Errors returned are handled in the main thread by SceneMgr's Update.
	LoadSync() error

	// PostLoad is called by the main thread (in SceneMgr's Update)
	// This func can't fail (anything that can fail should be in LoadSync)
	PostLoad()
}

type SceneMgr struct {
	scenes map[SceneID]Scene
	curID  SceneID

	loading int // 0 when no scene is loading; negative values should panic
	chLoad  chan loadResult
}

func NewSceneMgr() *SceneMgr {
	return &SceneMgr{
		scenes:  make(map[SceneID]Scene),
		curID:   -1,
		loading: 0,
		chLoad:  make(chan loadResult),
	}
}

type loadResult struct {
	ID    SceneID
	Scene SceneWithLoader
	Err   error
}

func (m *SceneMgr) AddScene(id SceneID, scene Scene) error {
	if id < 0 {
		return fmt.Errorf("invalid scene id: %d", id)
	}

	_, exists := m.scenes[id]
	if exists {
		return fmt.Errorf("duplicate scene: scene with id %d has already been added to SceneMgr", id)
	}

	loader, ok := scene.(SceneWithLoader)
	if !ok {
		m.scenes[id] = scene
		return nil
	}

	m.loading++
	go func() {
		err := loader.LoadSync()
		m.chLoad <- loadResult{ID: id, Scene: loader, Err: err}
	}()

	return nil
}

func (m *SceneMgr) UpdateLoading() error {
	select {
	case r := <-m.chLoad:
		m.loading--
		if m.loading < 0 {
			panic(fmt.Errorf("unexpected LoadResponse: scene_id=%d, err=%v", r.ID, r.Err))
		}
		if r.Err != nil {
			return fmt.Errorf("loading error: %w", r.Err)
		}
		m.scenes[r.ID] = r.Scene
		r.Scene.PostLoad()
	default:
	}
	return nil
}

func (m *SceneMgr) HasScene(id SceneID) bool {
	_, ok := m.scenes[id]
	return ok
}

func (m *SceneMgr) SwitchScene(id SceneID) error {
	if id < 0 {
		return fmt.Errorf("invalid scene id: %d", id)
	}

	if _, ok := m.scenes[id]; !ok {
		return fmt.Errorf("scene %d not found", id)
	}

	m.curID = id
	return nil
}

func (m *SceneMgr) MustSwitchScene(id SceneID) {
	err := m.SwitchScene(id)
	if err != nil {
		panic(err)
	}
}

//
// ebiten.Game interface

func (m *SceneMgr) Update() error {
	if err := m.UpdateLoading(); err != nil {
		return fmt.Errorf("loading error: %w", err)
	}

	if m.curID < 0 {
		return fmt.Errorf("no scene loaded (%d scene(s) waiting to load)", m.loading)
	}

	return m.scenes[m.curID].Update(m)
}

func (m *SceneMgr) Draw(screen *ebiten.Image) {
	m.scenes[m.curID].Draw(m, screen)
}

func (m *SceneMgr) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
