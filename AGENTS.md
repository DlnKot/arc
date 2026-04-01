# Agent Coding Guidelines for ARC

ARC (Alfa Remote Client) is a Wails-based desktop application for remote connections (RDP, VMware Horizon, Citrix, VPN).

## Project Structure

```
ARC/
├── main.go                          # Entry point, Wails initialization
├── wails.json                       # Wails build configuration
├── internal/
│   ├── analytics/                   # Analytics/telemetry service
│   ├── app/                         # Main app logic, Wails bindings
│   ├── config/                      # App constants (name, version)
│   ├── domain/                      # Domain types, settings structs
│   ├── launchers/                   # RDP, Horizon, Citrix, VPN launchers
│   ├── logging/                     # File logger interface
│   ├── network/                     # Ping, geo network checks
│   ├── store/                       # JSON file persistence
│   └── updater/                    # Auto-update service (GitHub/internal)
├── frontend/
│   └── src/
│       ├── main.js                 # API bindings, Wails invoke wrapper
│       ├── App.vue                 # Root component
│       ├── components/              # Vue components
│       ├── composables/            # Vue composables (useApp.js)
│       └── styles.css              # Global styles, CSS variables
└── build/
    ├── installer.nsi                # NSIS installer script
    └── windows/icon.ico            # Application icon
```

## Build Commands

### Development
```bash
# Backend + Frontend dev with hot reload
wails dev

# Frontend only (requires running backend separately)
cd frontend && npm run dev
```

### Production Build
```bash
# Full build (frontend + backend)
wails build

# Platform-specific
wails build -platform=darwin/arm64
wails build -platform=windows/amd64

# Frontend only
cd frontend && npm run build
```

### Windows NSIS Installer
```bash
# Install NSIS locally, then:
cd build
makensis installer.nsi

# Or via Chocolatey in CI:
choco install nsis -y
```

### Go Commands
```bash
go mod tidy
go build ./...
go run .
```

### Running Tests
```bash
# All tests
go test ./...

# Single package
go test ./internal/launchers/...
go test ./internal/domain/...

# With verbose output
go test ./... -v

# Run specific test
go test ./internal/launchers/... -v -run TestBuildHorizonArgsWindows
```

## Code Style

### Go

#### Imports
- Standard library first, then third-party, then internal packages
- Grouped with blank lines between groups
- Use short aliases for packages with common names (e.g., `appsvc` for `github.com/DlnKot/arc/internal/app`)

```go
import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/wailsapp/wails/v2"

    "github.com/DlnKot/arc/internal/domain"
    "github.com/DlnKot/arc/internal/logging"
)
```

#### Naming Conventions
- Structs: `PascalCase` (e.g., `App`, `Settings`, `RdpSettings`)
- Interfaces: `PascalCase` with `er` suffix when idiomatic (e.g., `Logger`)
- Private fields: `camelCase` (e.g., `ctx`, `store`, `logger`)
- Methods on structs: `PascalCase`
- Acronyms: Keep original casing (e.g., `RDP`, `URL` not `Rdp`, `Url`)

#### Error Handling
- Return errors from functions; use `domain.Result[T]` for Wails bindings
- Never silently ignore errors without comment
- Log errors with context using the logger
- Use `fmt.Errorf("context: %w", err)` for error wrapping

```go
func (a *App) DeleteConnection(id string) domain.Result[bool] {
    if err := a.store.DeleteConnection(id); err != nil {
        a.logger.Errorf("delete connection failed: %v", err)
        return fail[bool](err)
    }
    return ok(true)
}
```

#### Result Type Pattern
All Wails-bound methods return `domain.Result[T]` using helper functions:

```go
func ok[T any](data T) domain.Result[T] {
    return domain.Result[T]{Success: true, Data: data}
}

func fail[T any](err error) domain.Result[T] {
    return domain.Result[T]{Success: false, Error: err.Error()}
}
```

#### Struct Tags
- JSON tags use `json:"camelCase"` for frontend compatibility
- Omit empty fields: `json:"fieldName,omitempty"` when appropriate

#### Logging
- Use the `logging.Logger` interface (`Infof`, `Warnf`, `Errorf`)
- Log at appropriate levels: info for operations, warn for recoverable issues, error for failures
- Include context in log messages

### Vue/Frontend

#### Script Setup Syntax
- Use `<script setup>` with Composition API
- Import Vue primitives explicitly: `import { ref, computed, watch } from 'vue'`

#### API Calls
- All backend calls go through `window.api` (defined in `main.js`)
- Handle both success and failure responses
- Use `unwrapIpc()` helper to extract data from `domain.Result[T]`

```javascript
async function loadData() {
  try {
    const result = await window.api.getSettings()
    settings.value = unwrapIpc(result)
  } catch (error) {
    console.error('Error:', error.message)
    window.api?.log?.('error', `loadData failed: ${error.message}`)
  }
}
```

#### Response Handling
The backend returns `domain.Result[T]`:
```javascript
{ success: true, data: <actual_data> }
{ success: false, error: "error message" }
```

Use the `unwrapIpc()` helper to extract data:
```javascript
function unwrapIpc(res) {
  if (!res || typeof res !== 'object') return res
  if (res.success === false) {
    const err = new Error(res.error || 'IPC request failed')
    err.ipc = res
    throw err
  }
  if (res.success === true) {
    return Object.prototype.hasOwnProperty.call(res, 'data') ? res.data : undefined
  }
  return res
}
```

#### State Management
- Use Vue `ref()` for primitives, `reactive()` for objects
- Global state via composables (e.g., `useApp.js`)
- Avoid directly mutating props

#### Template Guidelines
- Use `v-if`/`v-else` for conditional rendering
- Bind methods with arrow functions only when passing arguments
- Use `:key` with `v-for`

### CSS/Styles

#### CSS Variables
Define in `styles.css` and use throughout:

```css
:root {
  --accent-primary: #3b82f6;
  --accent-danger: #ef4444;
  --bg-primary: #1f2937;
  --bg-secondary: #374151;
  --bg-tertiary: #4b5563;
  --text-primary: #f9fafb;
  --text-inverse: #111827;
  --text-muted: #9ca3af;
  --border-color: #4b5563;
  --border-light: #6b7280;
  --radius: 8px;
  --radius-sm: 6px;
  --radius-xl: 16px;
  --transition: all 0.2s ease;
}
```

#### Component Styles
- Use `<style scoped>` for component-specific styles
- Keep styles organized by section

## Wails Specific

### Regenerate Bindings
After modifying backend methods, regenerate JS bindings:
```bash
wails generate bindings
```
This updates `frontend/wailsjs/go/app/App.js`

### Platform-Specific Code
- Use `platform_*.go` suffix for OS-specific implementations
- Examples: `network/platform_windows.go`, `launchers/platform_windows.go`

### Events (Frontend Communication)
Use Wails `EventsOn` for backend-to-frontend events:
```javascript
EventsOn('event-name', (data) => {
  // Handle event
})
```

## Testing

### Manual Testing Checklist
- [ ] Delete connection button works
- [ ] Add/edit connection works
- [ ] Launch RDP/Horizon/Citrix connections
- [ ] Settings save and persist
- [ ] Network ping works
- [ ] Auto-updater checks for updates
- [ ] Analytics tracking works
- [ ] Light/dark theme renders correctly
- [ ] First-run modal shows when no credentials

### Cross-Platform Testing
Test on:
- macOS (ARM64)
- Windows (x64)
- Verify file paths, encoding (UTF-8 for Windows)

## Common Issues & Solutions

### Windows: Ping shows question marks
- Use `HideWindow: true` attribute on `exec.Cmd`
- Set `cmd.Stdout = os.Stdout` with UTF-8 encoding

### Windows: NSIS installer paths
- Paths in NSIS script are relative to the directory where `makensis` is run
- Run from `build/` directory with paths like `bin\ARC.exe`

### Wails: native confirm() doesn't work
- Create custom `ConfirmDialog.vue` component
- Use `wails.Runtime.Events.Emit()` for frontend events

### Frontend: settings not saving
- Check `defaultSettings` in Vue component doesn't override backend values
- Ensure merge logic includes all fields
