import { createApp } from 'vue'
import App from './App.vue'
import * as Backend from '../wailsjs/go/app/App'
import { BrowserOpenURL, EventsOn } from '../wailsjs/runtime/runtime'
import './styles.css'

function isWailsAvailable() {
  return typeof window !== 'undefined' && typeof window.go === 'object'
}

async function invoke(binding, ...args) {
  if (typeof binding !== 'function' || !isWailsAvailable()) {
    return { success: false, error: 'Недоступно в текущем окружении' }
  }

  try {
    return await binding(...args)
  } catch (error) {
    return { success: false, error: error?.message || String(error) }
  }
}

window.api = {
  getConnections: () => invoke(Backend.GetConnections),
  saveConnection: (connection) => invoke(Backend.SaveConnection, connection),
  deleteConnection: (id) => invoke(Backend.DeleteConnection, id),
  resetDefaultConnections: () => invoke(Backend.ResetDefaultConnections),

  getSettings: () => invoke(Backend.GetSettings),
  saveSettings: (settings) => invoke(Backend.SaveSettings, settings),

  launchRdp: (connection, settings) => invoke(Backend.LaunchRdp, connection, settings),
  launchHorizon: (connection, settings) => invoke(Backend.LaunchHorizon, connection, settings),
  launchCitrix: (connection, settings) => invoke(Backend.LaunchCitrix, connection, settings),
  launchVpn: () => invoke(Backend.LaunchVpn),

  networkGeo: () => invoke(Backend.NetworkGeo),
  networkPing: (host, packets) => invoke(Backend.NetworkPing, host, packets),

  checkForUpdates: () => invoke(Backend.CheckForUpdates),
  downloadUpdate: () => invoke(Backend.DownloadUpdate),
  installUpdate: () => invoke(Backend.InstallUpdate),
  getUpdateStatus: () => invoke(Backend.GetUpdateStatus),
  onAutoUpdateEvent: (callback) => {
    if (!isWailsAvailable()) {
      return () => {}
    }

    return EventsOn('auto-update-event', callback)
  },

  getVersion: () => invoke(Backend.GetVersion),
  getPlatform: () => invoke(Backend.GetPlatform),
  openExternal: async (url) => {
    if (!isWailsAvailable()) {
      window.open(url, '_blank', 'noopener')
      return { success: true, data: true }
    }

    try {
      await BrowserOpenURL(url)
      return { success: true, data: true }
    } catch (error) {
      return { success: false, error: error?.message || String(error) }
    }
  },
  log: (level, message) => invoke(Backend.Log, level, message),

  trackEvent: () => Promise.resolve({ success: true, data: true }),
  trackConnectionLaunch: () => Promise.resolve({ success: true, data: true }),
  trackTabView: () => Promise.resolve({ success: true, data: true }),
  trackHelpView: () => Promise.resolve({ success: true, data: true }),
  trackNetworkCheck: () => Promise.resolve({ success: true, data: true }),
  trackError: () => Promise.resolve({ success: true, data: true })
}

createApp(App).mount('#app')
