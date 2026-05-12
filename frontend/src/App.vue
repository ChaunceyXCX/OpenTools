<script setup lang="ts">
import { ref, computed } from 'vue'
import { useCommandStore } from './stores/commandStore'
import SearchBox from './components/SearchBox.vue'
import SearchResults from './components/SearchResults.vue'
import ClipboardPanel from './components/ClipboardPanel.vue'
import SettingsPage from './components/SettingsPage.vue'

type Tab = 'search' | 'clipboard' | 'settings'
const activeTab = ref<Tab>('search')
const selectedIndex = ref(0)

const store = useCommandStore()
const isSearch = computed(() => activeTab.value === 'search')

function onKeynav(dir: 'up' | 'down') {
  const len = store.results.length
  if (len === 0) return
  if (dir === 'down') {
    selectedIndex.value = (selectedIndex.value + 1) % len
  } else {
    selectedIndex.value = (selectedIndex.value - 1 + len) % len
  }
}

function onConfirm() {
  const cmd = store.results[selectedIndex.value]
  if (cmd) store.launch(cmd)
}

function onGlobalKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && activeTab.value !== 'search') {
    activeTab.value = 'search'
    e.preventDefault()
  }
}
</script>

<template>
  <div class="window" @keydown="onGlobalKey">
    <div class="tab-bar">
      <button :class="['tab', { active: activeTab === 'search' }]" @click="activeTab = 'search'; selectedIndex = 0">
        <span class="tab-icon">⌕</span> Search
      </button>
      <button :class="['tab', { active: activeTab === 'clipboard' }]" @click="activeTab = 'clipboard'">
        <span class="tab-icon">📋</span> Clipboard
      </button>
      <button :class="['tab', { active: activeTab === 'settings' }]" @click="activeTab = 'settings'">
        <span class="tab-icon">⚙</span> Settings
      </button>
    </div>

    <div v-if="isSearch" class="search-view" tabindex="-1">
      <SearchBox @keynav="onKeynav" @confirm="onConfirm" />
      <SearchResults :selected-index="selectedIndex" />
    </div>
    <ClipboardPanel v-else-if="activeTab === 'clipboard'" />
    <SettingsPage v-else />
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body { width: 100%; height: 100%; overflow: hidden; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
  -webkit-font-smoothing: antialiased;
}
#app { width: 100%; height: 100%; }
.window {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #1b2636;
  color: #e0e0e0;
  user-select: none;
}
.tab-bar {
  display: flex;
  padding: 0 8px;
  background: #111820;
  border-bottom: 1px solid #1e2a38;
}
.tab {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 14px;
  border: none;
  background: transparent;
  color: #6a7a8e;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  border-bottom: 2px solid transparent;
  transition: all 0.15s;
}
.tab:hover { color: #a0b0c0; background: #1a2535; }
.tab.active { color: #4a90d9; border-bottom-color: #4a90d9; }
.tab-icon { font-size: 14px; }
.search-view { display: flex; flex-direction: column; flex: 1; min-height: 0; }
::-webkit-scrollbar { width: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #2a3748; border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: #3a4a5e; }
</style>
