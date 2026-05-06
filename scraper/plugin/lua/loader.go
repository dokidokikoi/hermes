package lua

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Loader manages Lua plugin scrapers
type Loader struct {
	vmPool        *Pool
	scripts       map[string]*LuaScraper
	scriptDir     string
	logger        *zap.Logger
	executionTimeoutMu sync.RWMutex
	executionTimeout   string // For parsing from config
	mu            sync.RWMutex
}

// NewLoader creates a new Lua plugin loader
func NewLoader(logger *zap.Logger, vmPoolSize int, executionTimeout string) *Loader {
	return &Loader{
		vmPool:        NewPool(logger, vmPoolSize),
		scripts:       make(map[string]*LuaScraper),
		logger:        logger.Named("lua-loader"),
		executionTimeout: executionTimeout,
	}
}

// LoadFromDir loads all Lua plugins from a directory
func (l *Loader) LoadFromDir(dir string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.scriptDir = dir

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		l.logger.Info("plugin directory does not exist, creating it", zap.String("dir", dir))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrap(err, "failed to create plugin directory")
		}
		return nil
	}

	// Read all .lua files
	files, err := os.ReadDir(dir)
	if err != nil {
		return errors.Wrap(err, "failed to read plugin directory")
	}

	// Sort for consistent loading order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".lua") {
			continue
		}

		scriptPath := filepath.Join(dir, file.Name())
		if err := l.loadScript(scriptPath); err != nil {
			l.logger.Error("failed to load script", zap.String("path", scriptPath), zap.Error(err))
			continue
		}
	}

	l.logger.Info("loaded Lua plugins", zap.Int("count", len(l.scripts)))
	return nil
}

// loadScript loads a single Lua script
func (l *Loader) loadScript(scriptPath string) error {
	// Read the script file
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return errors.Wrap(err, "failed to read script file")
	}

	script := string(content)

	// Create a temporary scraper to validate the script
	tempScraper := NewLuaScraper(
		filepath.Base(scriptPath),
		script,
		l.vmPool,
		l.logger,
		DefaultExecutionTimeout,
	)

	// Get plugin info to validate
	name, version, err := tempScraper.GetInfo()
	if err != nil {
		return errors.Wrap(err, "script validation failed")
	}

	if name == "" {
		name = strings.TrimSuffix(filepath.Base(scriptPath), ".lua")
	}

	// Store the scraper
	l.scripts[name] = tempScraper

	l.logger.Info("loaded plugin",
		zap.String("name", name),
		zap.String("version", version),
		zap.String("path", scriptPath))

	return nil
}

// ReloadScript reloads a specific script
func (l *Loader) ReloadScript(scriptPath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get the name from the path
	baseName := filepath.Base(scriptPath)
	name := strings.TrimSuffix(baseName, ".lua")

	// Read the script file
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return errors.Wrap(err, "failed to read script file")
	}

	script := string(content)

	// If scraper exists, update it
	if existing, ok := l.scripts[name]; ok {
		if err := existing.UpdateScript(script); err != nil {
			return errors.Wrap(err, "failed to update script")
		}
		_, version, _ := existing.GetInfo()
		l.logger.Info("reloaded plugin",
			zap.String("name", name),
			zap.String("version", version),
			zap.String("path", scriptPath))
		return nil
	}

	// Create new scraper
	scraper := NewLuaScraper(
		name,
		script,
		l.vmPool,
		l.logger,
		DefaultExecutionTimeout,
	)

	pluginName, version, err := scraper.GetInfo()
	if err != nil {
		return errors.Wrap(err, "script validation failed")
	}

	if pluginName != "" {
		name = pluginName
	}

	l.scripts[name] = scraper

	l.logger.Info("loaded new plugin",
		zap.String("name", name),
		zap.String("version", version),
		zap.String("path", scriptPath))

	return nil
}

// RemoveScript removes a script by name
func (l *Loader) RemoveScript(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.scripts, name)
	l.logger.Info("removed plugin", zap.String("name", name))
}

// GetAll returns all loaded scrapers
func (l *Loader) GetAll() map[string]*LuaScraper {
	l.mu.RLock()
	defer l.mu.RUnlock()

	result := make(map[string]*LuaScraper, len(l.scripts))
	for k, v := range l.scripts {
		result[k] = v
	}
	return result
}

// Get returns a scraper by name
func (l *Loader) Get(name string) (*LuaScraper, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	scraper, ok := l.scripts[name]
	return scraper, ok
}

// GetNames returns all scraper names
func (l *Loader) GetNames() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	names := make([]string, 0, len(l.scripts))
	for name := range l.scripts {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Close closes the loader and releases resources
func (l *Loader) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.vmPool.Close()
	l.scripts = make(map[string]*LuaScraper)
}