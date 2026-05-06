package lua

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

// HotReload manages hot-reloading of Lua scripts
type HotReload struct {
	watcher *fsnotify.Watcher
	loader  *Loader
	logger  *zap.Logger
	running bool
	mu      sync.RWMutex
}

// NewHotReload creates a new hot reload manager
func NewHotReload(loader *Loader, logger *zap.Logger) (*HotReload, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &HotReload{
		watcher: watcher,
		loader:  loader,
		logger:  logger.Named("lua-hot-reload"),
		running: false,
	}, nil
}

// Start begins watching for file changes
func (h *HotReload) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.running {
		return nil
	}

	// Watch the plugin directory
	if h.loader.scriptDir != "" {
		if err := h.watcher.Add(h.loader.scriptDir); err != nil {
			return err
		}
	}

	h.running = true

	go h.watch()

	h.logger.Info("hot reload started", zap.String("dir", h.loader.scriptDir))
	return nil
}

// Stop stops watching for file changes
func (h *HotReload) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.running {
		return
	}

	h.running = false

	if h.watcher != nil {
		h.watcher.Close()
	}

	h.logger.Info("hot reload stopped")
}

// watch watches for file system events
func (h *HotReload) watch() {
	for {
		select {
		case event, ok := <-h.watcher.Events:
			if !ok {
				return
			}
			h.handleEvent(event)

		case err, ok := <-h.watcher.Errors:
			if !ok {
				return
			}
			h.logger.Error("watcher error", zap.Error(err))
		}
	}
}

// handleEvent handles a file system event
func (h *HotReload) handleEvent(event fsnotify.Event) {
	h.mu.RLock()
	if !h.running {
		h.mu.RUnlock()
		return
	}
	h.mu.RUnlock()

	// Only handle .lua files
	if !strings.HasSuffix(event.Name, ".lua") {
		return
	}

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		h.logger.Info("new script detected", zap.String("path", event.Name))
		h.loader.ReloadScript(event.Name)

	case event.Op&fsnotify.Write == fsnotify.Write:
		h.logger.Info("script modified", zap.String("path", event.Name))
		h.loader.ReloadScript(event.Name)

	case event.Op&fsnotify.Remove == fsnotify.Remove:
		baseName := filepath.Base(event.Name)
		name := strings.TrimSuffix(baseName, ".lua")
		h.logger.Info("script removed", zap.String("path", event.Name))
		h.loader.RemoveScript(name)

	case event.Op&fsnotify.Rename == fsnotify.Rename:
		// Check if the file still exists (might be a rename to new name)
		if _, ok := h.loader.Get(filepath.Base(event.Name)); ok {
			// File was renamed away
			baseName := filepath.Base(event.Name)
			name := strings.TrimSuffix(baseName, ".lua")
			h.loader.RemoveScript(name)
		}
	}
}