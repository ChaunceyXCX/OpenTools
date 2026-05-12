<script setup lang="ts">
import { ref } from 'vue'
import SearchBox from './components/SearchBox.vue'
import SearchResults from './components/SearchResults.vue'
import SettingsPage from './components/SettingsPage.vue'

const activeTab = ref<'search' | 'settings'>('search')
</script>

<template>
  <div class="window">
    <div class="tab-bar">
      <button :class="['tab', { active: activeTab === 'search' }]" @click="activeTab = 'search'">
        Search
      </button>
      <button :class="['tab', { active: activeTab === 'settings' }]" @click="activeTab = 'settings'">
        Settings
      </button>
    </div>
    <template v-if="activeTab === 'search'">
      <SearchBox />
      <SearchResults />
    </template>
    <SettingsPage v-else />
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
html, body {
  width: 100%;
  height: 100%;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
}
#app {
  width: 100%;
  height: 100%;
}
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
  gap: 0;
  padding: 4px 8px 0;
  background: #111820;
}
.tab {
  padding: 6px 16px;
  border: none;
  background: transparent;
  color: #6a7a8e;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  border-bottom: 2px solid transparent;
}
.tab:hover { color: #a0b0c0; }
.tab.active {
  color: #4a90d9;
  border-bottom-color: #4a90d9;
}
</style>
