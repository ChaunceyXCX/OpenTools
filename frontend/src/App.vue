<script setup lang="ts">
import { ref, watch } from 'vue'
import { useCommandDataStore } from './stores/commandDataStore'
import { useWindowStore } from './stores/windowStore'
import SearchBox from './components/search/SearchBox.vue'
import SearchResults from './components/search/SearchResults.vue'
import SettingsPage from './components/SettingsPage.vue'
import ClipboardPanel from './components/ClipboardPanel.vue'
import PluginMarket from './components/PluginMarket.vue'

type Tab = 'search' | 'clipboard' | 'plugins' | 'settings'
const tabs: Tab[] = ['search', 'clipboard', 'plugins', 'settings']
const activeTab = ref<Tab>('search')
const selectedIndex = ref(0)
const windowVisible = ref(false)

const cmdStore = useCommandDataStore()
const winStore = useWindowStore()

setTimeout(() => windowVisible.value = true, 50)

watch(() => cmdStore.results.length, (n) => {
  if (n > 0 && activeTab.value === 'search') {
    try { (window as any).__resizeHeight?.(Math.min(n * 48 + 100, 600)) } catch { }
  }
})

function onKeynav(dir: 'up' | 'down') {
  const len = cmdStore.results.length
  if (len === 0) return
  selectedIndex.value = dir === 'down'
    ? (selectedIndex.value + 1) % len
    : (selectedIndex.value - 1 + len) % len
}

function onConfirm() {
  const cmd = cmdStore.results[selectedIndex.value]
  if (cmd) cmdStore.launch(cmd)
}

function onClosePlugin() { winStore.exitPlugin() }

function onGlobalKey(e: KeyboardEvent) {
  if (e.key === 'Tab') {
    e.preventDefault()
    const idx = tabs.indexOf(activeTab.value)
    activeTab.value = tabs[(idx + 1) % tabs.length]
    selectedIndex.value = 0
  }
}
</script>

<template>
  <div class="app-container" :class="{
    'app-container__plugin': winStore.viewMode === 'plugin',
    'visible': windowVisible
  }" @keydown="onGlobalKey">
    <div class="search-window">
      <div class="tab-bar">
        <button :class="['tab', { active: activeTab === 'search' }]"
          @click="activeTab = 'search'; selectedIndex = 0">
          ⌕ {{ $t('tab.search') }}
        </button>
        <button :class="['tab', { active: activeTab === 'clipboard' }]"
          @click="activeTab = 'clipboard'">
          📋 {{ $t('tab.clipboard') }}
        </button>
        <button :class="['tab', { active: activeTab === 'plugins' }]"
          @click="activeTab = 'plugins'">
          🧩 Plugins
        </button>
        <button :class="['tab', { active: activeTab === 'settings' }]"
          @click="activeTab = 'settings'">
          ⚙ {{ $t('tab.settings') }}
        </button>
      </div>

      <template v-if="activeTab === 'search'">
        <div class="search-box-wrapper">
          <SearchBox
            @keynav="onKeynav"
            @confirm="onConfirm"
            @close-plugin="onClosePlugin"
          />
        </div>

        <SearchResults
          v-if="winStore.viewMode === 'search'"
          :selected-index="selectedIndex"
        />

        <div v-else-if="winStore.viewMode === 'plugin'" class="plugin-placeholder">
          <div class="plugin-container">
            <div class="plugin-header">
              <span>{{ winStore.currentPlugin?.title || 'Plugin' }}</span>
            </div>
            <div class="plugin-body">
              <p style="color: var(--text-secondary);">Plugin content would render here</p>
            </div>
          </div>
        </div>
      </template>

      <ClipboardPanel v-else-if="activeTab === 'clipboard'" />
      <PluginMarket v-else-if="activeTab === 'plugins'" />
      <SettingsPage v-else-if="activeTab === 'settings'" />
    </div>
  </div>
</template>

<style>
@import './style.css';
#app { width: 100%; height: 100%; }
.app-container {
  width: 100%; height: 100vh;
  background: var(--bg-color);
  color: var(--text-color);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  opacity: 0;
  transform: translateY(4px);
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.app-container.visible {
  opacity: 1;
  transform: translateY(0);
}
.search-window {
  display: flex; flex-direction: column;
  flex: 1; min-height: 0;
}
.tab-bar {
  display: flex; padding: 0 8px;
  background: color-mix(in srgb, var(--bg-color) 97%, #000);
  border-bottom: 1px solid var(--border-color);
}
.tab {
  display: flex; align-items: center; gap: 4px;
  padding: 7px 14px;
  border: none; background: transparent;
  color: var(--text-secondary); cursor: pointer;
  font-size: 12px; font-weight: 500;
  border-bottom: 2px solid transparent;
  transition: all 0.15s;
}
.tab:hover { color: var(--text-color); background: var(--hover-bg); }
.tab.active { color: var(--primary-color); border-bottom-color: var(--primary-color); }
.search-box-wrapper { flex-shrink: 0; }
.plugin-placeholder { flex: 1; display: flex; }
.plugin-container {
  flex: 1; display: flex; flex-direction: column;
  padding: 16px;
}
.plugin-header { font-size: 14px; font-weight: 600; margin-bottom: 12px; }
.plugin-body { flex: 1; }
</style>
