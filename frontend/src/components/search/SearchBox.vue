<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useCommandDataStore } from '../../stores/commandDataStore'
import { useWindowStore } from '../../stores/windowStore'

const commandStore = useCommandDataStore()
const windowStore = useWindowStore()
const inputRef = ref<HTMLInputElement>()

const emit = defineEmits<{
  keynav: [dir: 'up' | 'down']
  confirm: []
  closePlugin: []
}>()

onMounted(() => {
  commandStore.loadCommands()
  inputRef.value?.focus()
})

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') { e.preventDefault(); emit('keynav', 'down') }
  else if (e.key === 'ArrowUp') { e.preventDefault(); emit('keynav', 'up') }
  else if (e.key === 'Enter') { e.preventDefault(); emit('confirm') }
  else if (e.key === 'Escape') {
    if (windowStore.viewMode === 'plugin') emit('closePlugin')
  }
}
</script>

<template>
  <div class="search-box">
    <div class="search-input-container">
      <div v-if="windowStore.viewMode === 'plugin' && windowStore.currentPlugin" class="plugin-tag">
        <img
          v-if="windowStore.currentPlugin.logo"
          :src="windowStore.currentPlugin.logo"
          class="plugin-tag-icon"
        />
        <span class="plugin-tag-title">
          {{ windowStore.currentPlugin.title || windowStore.currentPlugin.name }}
        </span>
        <span v-if="windowStore.currentPlugin.cmdName" class="plugin-tag-cmd">
          {{ windowStore.currentPlugin.cmdName }}
        </span>
        <button class="plugin-tag-close" @click.stop="emit('closePlugin')">✕</button>
      </div>
      <input
        ref="inputRef"
        v-model="commandStore.searchQuery"
        type="text"
        class="search-input"
        :placeholder="windowStore.viewMode === 'plugin' ? $t('search.pluginPlaceholder') : $t('search.placeholder')"
        autofocus
        spellcheck="false"
        autocomplete="off"
        @keydown="onKeydown"
      />
    </div>
  </div>
</template>

<style scoped>
.search-box { padding: 0; }
.search-input-container {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-bottom: 1px solid var(--border-color);
}
.plugin-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 3px 6px 3px 8px;
  background: var(--primary-color);
  border-radius: 6px;
  font-size: 12px;
  color: white;
  white-space: nowrap;
}
.plugin-tag-icon { width: 16px; height: 16px; border-radius: 3px; }
.plugin-tag-cmd {
  padding: 1px 5px;
  background: rgba(255,255,255,0.2);
  border-radius: 4px;
  font-size: 11px;
}
.plugin-tag-close {
  background: none; border: none; color: white;
  cursor: pointer; font-size: 12px; padding: 2px;
  opacity: 0.7;
}
.plugin-tag-close:hover { opacity: 1; }
.search-input {
  flex: 1;
  border: none; outline: none;
  background: transparent;
  color: var(--text-color);
  font-size: 14px;
  line-height: 1.5;
  min-width: 0;
}
.search-input::placeholder { color: var(--text-secondary); }
</style>
