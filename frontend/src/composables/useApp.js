import { ref, computed } from 'vue'

// Global state
const connections = ref([])
const settings = ref({})
const currentView = ref('connections')
const currentClientFilter = ref('all')
const isLoading = ref(false)
const isFirstRun = ref(false)

function hasUserCredentials(value) {
  return !!value?.user?.username
}

function normalizeConnectionsList(value) {
  if (Array.isArray(value)) return value
  if (value && typeof value === 'object') {
    const keys = Object.keys(value)
    if (keys.length && keys.every(key => /^\d+$/.test(key))) {
      return keys
        .sort((a, b) => Number(a) - Number(b))
        .map(key => value[key])
    }
  }
  return []
}

function unwrapIpc(res) {
  if (!res || typeof res !== 'object') return res
  if (res.success === false) {
    const err = new Error(res.error || 'IPC request failed')
    err.ipc = res
    throw err
  }
  if (res.success === true) {
    if (Object.prototype.hasOwnProperty.call(res, 'data')) {
      return res.data
    }
    return undefined
  }
  return res
}

/* --------------------------- LOAD DATA --------------------------- */

async function loadData() {
  isLoading.value = true

  try {
    if (!window.api) {
      connections.value = []
      settings.value = { user: { domain: '', username: '' } }
      isFirstRun.value = true
      return
    }

    const [conns, s] = await Promise.all([
      window.api.getConnections(),
      window.api.getSettings()
    ])

    let connsData = conns
    if (conns && typeof conns === 'object' && conns.success === true) {
      connsData = conns.data
    }

    let settingsData = s
    if (s && typeof s === 'object' && s.success === true) {
      settingsData = s.data
    }

    if (Array.isArray(connsData)) {
      connections.value = normalizeConnectionsList(connsData)
    } else {
      connections.value = []
    }
    settings.value = settingsData || {}

    isFirstRun.value = !hasUserCredentials(settingsData)

  } catch (error) {
    const errorMsg = error?.message || String(error)
    console.error('Error loading data:', errorMsg)
    if (window.api?.log) {
      window.api.log('error', `loadData failed: ${errorMsg}`)
    }
    // Keep the previous first-run state on load failures.
    // A transient backend error should not reopen the onboarding modal.

  } finally {

    isLoading.value = false

  }
}

/* ------------------------ CONNECTION OPS ------------------------ */

async function saveConnection(connection) {
  try {
    const isFactory = !!connection?.factoryId
    const isNew = !isFactory && !connection.id
    const result = await window.api.saveConnection(connection)

    const res = await window.api.getConnections()
    let data = res
    if (res && typeof res === 'object' && res.success === true) {
      data = res.data
    }
    if (Array.isArray(data)) {
      connections.value = normalizeConnectionsList(data)
    } else {
      connections.value = normalizeConnectionsList(unwrapIpc(res))
    }

    if (result && typeof result === 'object' && result.success && window.api?.trackEvent) {
      if (isFactory) {
        window.api.trackEvent('default_connection_rename', { factoryId: connection.factoryId })
      } else if (isNew) {
        window.api.trackEvent('connection_create', { type: connection.type })
      } else {
        window.api.trackEvent('connection_edit', { type: connection.type })
      }
    }
    
    return { success: true, result }
  } catch (error) {
    const errorMsg = error?.message || String(error)
    console.error('Failed to save connection:', errorMsg)
    if (window.api?.log) {
      window.api.log('error', `saveConnection failed: ${errorMsg}`)
    }
    return { success: false, error: errorMsg }
  }
}

async function deleteConnection(id) {
  try {
    const conn = connections.value.find(c => c.id === id)
    const connectionType = conn?.type || 'unknown'

    await window.api.deleteConnection(id)

    const res = await window.api.getConnections()
    let data = res
    if (res && typeof res === 'object' && res.success === true) {
      data = res.data
    }
    if (Array.isArray(data)) {
      connections.value = normalizeConnectionsList(data)
    } else {
      connections.value = normalizeConnectionsList(unwrapIpc(res))
    }
    
    if (window.api?.trackEvent) {
      window.api.trackEvent('connection_delete', { type: connectionType })
    }
    
    return { success: true }
  } catch (error) {
    const errorMsg = error?.message || String(error)
    console.error('Failed to delete connection:', errorMsg)
    if (window.api?.log) {
      window.api.log('error', `deleteConnection failed: ${errorMsg}`)
    }
    return { success: false, error: errorMsg }
  }
}

async function resetDefaultConnections() {
  try {
    if (!window.api?.resetDefaultConnections) {
      return { success: false, error: 'Недоступно в браузере' }
    }

    await window.api.resetDefaultConnections()

    const res = await window.api.getConnections()
    let data = res
    if (res && typeof res === 'object' && res.success === true) {
      data = res.data
    }
    if (Array.isArray(data)) {
      connections.value = normalizeConnectionsList(data)
    } else {
      connections.value = []
    }

    if (window.api?.trackEvent) {
      window.api.trackEvent('default_connections_reset', {})
    }

    return { success: true }
  } catch (error) {
    const errorMsg = error?.message || String(error)
    console.error('Failed to reset default connections:', errorMsg)
    if (window.api?.log) {
      window.api.log('error', `resetDefaultConnections failed: ${errorMsg}`)
    }
    return { success: false, error: errorMsg }
  }
}

/* -------------------------- SETTINGS OPS ------------------------ */

async function saveSettings(newSettings) {
  try {
    const plainSettings = JSON.parse(JSON.stringify(newSettings))

    await window.api.saveSettings(plainSettings)

    const settingsRes = await window.api.getSettings()
    let settingsData = settingsRes
    if (settingsRes && typeof settingsRes === 'object' && settingsRes.success === true) {
      settingsData = settingsRes.data
    }
    settings.value = settingsData || {}

    const connsRes = await window.api.getConnections()
    let connsData = connsRes
    if (connsRes && typeof connsRes === 'object' && connsRes.success === true) {
      connsData = connsRes.data
    }
    if (Array.isArray(connsData)) {
      connections.value = normalizeConnectionsList(connsData)
    } else {
      connections.value = []
    }
    isFirstRun.value = !hasUserCredentials(settings.value)

    if (window.api?.trackEvent) {
      window.api.trackEvent('settings_save', {})
    }
    
    return { success: true }
  } catch (error) {
    const errorMsg = error?.message || String(error)
    console.error('Failed to save settings:', errorMsg)
    if (window.api?.log) {
      window.api.log('error', `saveSettings failed: ${errorMsg}`)
    }
    return { success: false, error: errorMsg }
  }
}

/* ---------------- USER CREDENTIALS ---------------- */

function getUserCredentials() {

  const s = settings.value

  return {
    domain: s.user?.domain || '',
    username: s.user?.username || ''
  }
}

function applyCredentialsToConnection(connection) {

  const creds = getUserCredentials()

  if (!connection.username?.trim() && (creds.domain || creds.username)) {

    const username = creds.domain
      ? `${creds.domain}\\${creds.username}`
      : creds.username

    return { ...connection, username }
  }

  return { ...connection }
}

function isPlainObject(v) {
  return v && typeof v === 'object' && !Array.isArray(v)
}

function deepMerge(a, b) {
  if (!isPlainObject(a)) return JSON.parse(JSON.stringify(b ?? {}))
  if (!isPlainObject(b)) return JSON.parse(JSON.stringify(a ?? {}))

  const out = { ...a }
  for (const k of Object.keys(b)) {
    const av = a[k]
    const bv = b[k]
    if (isPlainObject(av) && isPlainObject(bv)) out[k] = deepMerge(av, bv)
    else out[k] = bv
  }
  return out
}

/* --------------------- LAUNCH CONNECTION --------------------- */

async function launchConnection(conn) {

  if (!conn)
    return { success: false, error: 'Connection not found' }

  try {

    const plainConnection = JSON.parse(JSON.stringify(conn))

    const connectionWithCreds = applyCredentialsToConnection(plainConnection)

    const globalSettings = JSON.parse(JSON.stringify(settings.value))

    let clientSettings = {}

    if (plainConnection.clientSettings) {
      clientSettings = plainConnection.clientSettings
    }

    const mergedSettings = {
      ...globalSettings,
      [plainConnection.type]: {
        ...deepMerge((globalSettings[plainConnection.type] || {}), clientSettings)
      }
    }

    let result

    switch (plainConnection.type) {

      case 'rdp':
        result = await window.api.launchRdp(connectionWithCreds, mergedSettings)
        // Трекинг запуска подключения
        if (result?.success && window.api?.trackConnectionLaunch) {
          window.api.trackConnectionLaunch('rdp')
        }
        break

      case 'horizon':
        result = await window.api.launchHorizon(connectionWithCreds, mergedSettings)
        // Трекинг запуска подключения
        if (result?.success && window.api?.trackConnectionLaunch) {
          window.api.trackConnectionLaunch('horizon')
        }
        break

      case 'citrix':
        result = await window.api.launchCitrix(connectionWithCreds, mergedSettings)
        // Трекинг запуска подключения
        if (result?.success && window.api?.trackConnectionLaunch) {
          window.api.trackConnectionLaunch('citrix')
        }
        break

      default:
        return {
          success: false,
          error: `Unsupported connection type: ${plainConnection.type}`
        }
    }

    return result || { success: false, error: 'Unknown error' }

  } catch (error) {
    const errorMsg = error?.message || String(error)
    console.error('Launch error:', errorMsg)
    if (window.api?.log) {
      window.api.log('error', `launchConnection failed: ${errorMsg}`)
    }
    // Трекинг ошибки
    if (window.api?.trackError) {
      window.api.trackError(errorMsg)
    }
    return {
      success: false,
      error: errorMsg
    }
  }
}

/* ---------------- FILTERED CONNECTIONS ---------------- */

const filteredConnections = computed(() =>
  currentClientFilter.value === 'all'
    ? connections.value
    : connections.value.filter(c => c.type === currentClientFilter.value)
)

/* ---------------- AUTO-UPDATER ---------------- */

// Update state
const updateStatus = ref({
  updateAvailable: false,
  updateDownloaded: false,
  version: null,
  releaseUrl: null,
  downloadUrl: null,
  installOnQuit: false
})

const updateProgress = ref({
  percent: 0,
  bytesPerSecond: 0,
  transferred: 0,
  total: 0
})

const updateError = ref(null)

// Initialize auto-update event listener
function initAutoUpdater() {
  if (!window.api?.onAutoUpdateEvent) return

  window.api.onAutoUpdateEvent(({ event, data }) => {
    switch (event) {
      case 'download-progress':
        updateProgress.value.percent = data
        break
      case 'updater:downloaded':
        updateStatus.value.updateDownloaded = true
        break
      case 'updater:error':
        updateError.value = data
        break
    }
  })
}

// Check for updates manually
async function checkForUpdates() {
  updateError.value = null
  try {
    const result = await window.api.checkForUpdates()
    const data = unwrapIpc(result)
    if (data?.updateAvailable) {
      updateStatus.value.updateAvailable = true
      updateStatus.value.version = data.version
      updateStatus.value.downloadUrl = data.downloadUrl
      updateStatus.value.releaseUrl = data.releaseUrl
    } else {
      updateStatus.value.updateAvailable = false
    }
    return result
  } catch (error) {
    const errorMsg = error.message
    updateError.value = errorMsg
    return { success: false, error: errorMsg }
  }
}

// Download available update
async function downloadUpdate() {
  updateError.value = null
  updateProgress.value.percent = 0
  try {
    const result = await window.api.downloadUpdate()
    if (result?.success === false) updateError.value = result.error || 'Download failed'
    return result
  } catch (error) {
    const errorMsg = error.message
    updateError.value = errorMsg
    return { success: false, error: errorMsg }
  }
}

// Install downloaded update now
async function installNow() {
  await window.api.installUpdate()
}

// Install on quit (Windows/macOS)
async function installOnQuit() {
  try {
    await window.api.installOnQuit()
    updateStatus.value.installOnQuit = true
  } catch (error) {
    updateError.value = error.message
  }
}

// Cancel download
async function cancelDownload() {
  try {
    await window.api.cancelDownload()
  } catch (error) {
    // Ignore
  }
}

/* ---------------- EXPORT ---------------- */

export function useApp() {

  return {

    connections,
    settings,
    currentView,
    currentClientFilter,
    isLoading,
    isFirstRun,
    filteredConnections,

    // Auto-updater state
    updateStatus,
    updateProgress,
    updateError,

    loadData,
    saveConnection,
    deleteConnection,
    resetDefaultConnections,
    saveSettings,
    launchConnection,
    getUserCredentials,
    applyCredentialsToConnection,
    // Auto-updater methods
    initAutoUpdater,
    checkForUpdates,
    downloadUpdate,
    installNow,
    installOnQuit,
    cancelDownload
  }
}
