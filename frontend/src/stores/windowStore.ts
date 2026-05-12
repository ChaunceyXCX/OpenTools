import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface CurrentPlugin {
  name: string
  title?: string
  logo?: string
  cmdName?: string
}

export const useWindowStore = defineStore('window', () => {
  const currentPlugin = ref<CurrentPlugin | null>(null)
  const viewMode = ref<'search' | 'plugin'>('search')

  function enterPlugin(plugin: CurrentPlugin) {
    currentPlugin.value = plugin
    viewMode.value = 'plugin'
  }

  function exitPlugin() {
    currentPlugin.value = null
    viewMode.value = 'search'
  }

  return { currentPlugin, viewMode, enterPlugin, exitPlugin }
})
