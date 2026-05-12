<script setup lang="ts">
import { ref } from 'vue'

interface Action {
  label: string
  code: string
}

interface Detection {
  type: string
  content: string
  actions: Action[]
}

const props = defineProps<{ content: string }>()
const emit = defineEmits<{ close: [] }>()

const detection = ref<Detection | null>(null)

function detect(text: string) {
  const trimmed = text.trim()
  if (/^https?:\/\/\S+$/.test(trimmed)) {
    detection.value = {
      type: 'url',
      content: trimmed,
      actions: [
        { label: 'Open in Browser', code: 'open_url' },
        { label: 'Copy URL', code: 'copy' },
      ],
    }
  } else {
    detection.value = {
      type: 'text',
      content: trimmed.slice(0, 200),
      actions: [
        { label: 'Search', code: 'search' },
        { label: 'Copy', code: 'copy' },
      ],
    }
  }
}

function doAction(code: string) {
  switch (code) {
    case 'copy':
      navigator.clipboard.writeText(detection.value?.content || '')
      break
    case 'open_url':
      window.open(detection.value?.content, '_blank')
      break
    case 'search':
      // would trigger search tab
      break
  }
  emit('close')
}

import { onMounted } from 'vue'
onMounted(() => detect(props.content))
</script>

<template>
  <div class="overlay" @click.self="emit('close')">
    <div class="panel">
      <div class="panel-header">
        <span class="badge">{{ detection?.type }}</span>
        <button class="close-btn" @click="emit('close')">✕</button>
      </div>
      <div class="panel-body">{{ detection?.content }}</div>
      <div v-if="detection?.actions.length" class="panel-actions">
        <button
          v-for="a in detection.actions"
          :key="a.code"
          class="action-btn"
          @click="doAction(a.code)"
        >
          {{ a.label }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.panel {
  background: #1b2636;
  border: 1px solid #3a4a5e;
  border-radius: 12px;
  width: 400px;
  max-width: 90vw;
  box-shadow: 0 8px 32px rgba(0,0,0,0.5);
  overflow: hidden;
}
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  border-bottom: 1px solid #2a3748;
}
.badge {
  font-size: 11px;
  font-weight: 600;
  color: #4a90d9;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.close-btn {
  background: none;
  border: none;
  color: #6a7a8e;
  cursor: pointer;
  font-size: 14px;
  padding: 2px 6px;
  border-radius: 4px;
}
.close-btn:hover { color: #e0e0e0; background: #2a3748; }
.panel-body {
  padding: 14px;
  font-size: 13px;
  color: #c0d0e0;
  line-height: 1.5;
  word-break: break-all;
  max-height: 200px;
  overflow-y: auto;
}
.panel-actions {
  display: flex;
  gap: 6px;
  padding: 10px 14px;
  border-top: 1px solid #2a3748;
}
.action-btn {
  flex: 1;
  padding: 8px;
  border: 1px solid #3a4a5e;
  border-radius: 6px;
  background: #2a3748;
  color: #e0e0e0;
  cursor: pointer;
  font-size: 13px;
  text-align: center;
  transition: background 0.12s;
}
.action-btn:hover { background: #4a90d9; border-color: #4a90d9; }
</style>
