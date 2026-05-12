<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useCommandDataStore, type PastedFile } from '../../stores/commandDataStore'
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

function onPaste(e: ClipboardEvent) {
  const items = e.clipboardData?.items
  if (!items) return
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item.type.startsWith('image/')) {
      const blob = item.getAsFile()
      if (blob) {
        const reader = new FileReader()
        reader.onload = () => commandStore.pastedImage = reader.result as string
        reader.readAsDataURL(blob)
        e.preventDefault()
        return
      }
    }
  }
}

function onDrop(e: DragEvent) {
  e.preventDefault()
  const files = e.dataTransfer?.files
  if (!files || files.length === 0) return
  const fileList: PastedFile[] = []
  for (let i = 0; i < files.length; i++) {
    fileList.push({
      name: files[i].name,
      path: (files[i] as any).path || files[i].name,
      isDirectory: false,
    })
  }
  commandStore.pastedFiles = fileList
}
</script>

<template>
  <div class="search-box" @drop="onDrop" @dragover.prevent>
    <div class="search-input-container">
      <div v-if="commandStore.pastedImage" class="pasted-preview">
        <img :src="commandStore.pastedImage" class="paste-img" />
        <button class="paste-close" @click="commandStore.pastedImage = null">✕</button>
      </div>
      <div v-if="commandStore.pastedFiles" class="pasted-files">
        <span class="file-badge">📎 {{ commandStore.pastedFiles.length }} files</span>
        <button class="paste-close" @click="commandStore.pastedFiles = null">✕</button>
      </div>
      <div v-if="windowStore.viewMode === 'plugin' && windowStore.currentPlugin" class="plugin-tag">
        <img v-if="windowStore.currentPlugin.logo" :src="windowStore.currentPlugin.logo" class="plugin-tag-icon" />
        <span class="plugin-tag-title">{{ windowStore.currentPlugin.title || windowStore.currentPlugin.name }}</span>
        <span v-if="windowStore.currentPlugin.cmdName" class="plugin-tag-cmd">{{ windowStore.currentPlugin.cmdName }}</span>
        <button class="plugin-tag-close" @click.stop="emit('closePlugin')">✕</button>
      </div>
      <input ref="inputRef" v-model="commandStore.searchQuery" type="text" class="search-input"
        :placeholder="windowStore.viewMode === 'plugin' ? $t('search.pluginPlaceholder') : $t('search.placeholder')"
        autofocus spellcheck="false" autocomplete="off"
        @keydown="onKeydown" @paste="onPaste"
      />
    </div>
  </div>
</template>

<style scoped>
.search-box { padding: 0; }
.search-input-container { display: flex; align-items: center; gap: 6px; padding: 6px 12px; border-bottom: 1px solid var(--border-color); }
.search-input { flex: 1; border: none; outline: none; background: transparent; color: var(--text-color); font-size: 14px; min-width: 0; }
.search-input::placeholder { color: var(--text-secondary); }
.plugin-tag { display: flex; align-items: center; gap: 6px; padding: 3px 6px 3px 8px; background: var(--primary-color); border-radius: 6px; font-size: 12px; color: white; white-space: nowrap; }
.plugin-tag-icon { width: 16px; height: 16px; border-radius: 3px; }
.plugin-tag-cmd { padding: 1px 5px; background: rgba(255,255,255,0.2); border-radius: 4px; font-size: 11px; }
.plugin-tag-close { background: none; border: none; color: white; cursor: pointer; font-size: 12px; padding: 2px; opacity: 0.7; }
.plugin-tag-close:hover { opacity: 1; }
.pasted-preview { display: flex; align-items: center; gap: 4px; padding: 2px 4px; background: var(--control-bg); border-radius: var(--radius-sm); }
.paste-img { width: 24px; height: 24px; object-fit: cover; border-radius: 4px; }
.paste-close { background: none; border: none; color: var(--text-secondary); cursor: pointer; font-size: 10px; }
.pasted-files { display: flex; align-items: center; gap: 4px; }
.file-badge { font-size: 12px; color: var(--text-secondary); padding: 2px 6px; background: var(--control-bg); border-radius: var(--radius-sm); }
</style>
