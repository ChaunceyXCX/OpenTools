<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { ClipboardService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'

const items = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')

const filtered = computed(() => {
  if (!searchQuery.value) return items.value
  const q = searchQuery.value.toLowerCase()
  return items.value.filter((i: any) => i.content?.toLowerCase().includes(q))
})

onMounted(() => load())

async function load() {
  loading.value = true
  try {
    items.value = await ClipboardService.GetHistory(1, 50)
  } catch { /* ignore */ }
  loading.value = false
}

async function remove(id: string) {
  await ClipboardService.Delete(id)
  await load()
}

async function clearAll() {
  await ClipboardService.Clear()
  items.value = []
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
}

function timeAgo(ts: number): string {
  const diff = Date.now() - ts
  if (diff < 60000) return 'just now'
  if (diff < 3600000) return Math.floor(diff / 60000) + 'm ago'
  if (diff < 86400000) return Math.floor(diff / 3600000) + 'h ago'
  return Math.floor(diff / 86400000) + 'd ago'
}
</script>

<template>
  <div class="clipboard-panel">
    <div class="panel-header">
      <h2>Clipboard</h2>
      <div class="header-actions">
        <button v-if="items.length" class="btn-text" @click="clearAll">Clear All</button>
      </div>
    </div>

    <div class="search-row">
      <input v-model="searchQuery" placeholder="Search clipboard..." class="search-clip" />
      <button class="btn-icon" @click="load" title="Refresh">↻</button>
    </div>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="filtered.length === 0" class="empty">
      {{ searchQuery ? 'No matches' : 'Clipboard is empty' }}
    </div>
    <div v-else class="list">
      <div v-for="item in filtered" :key="item.id" class="clip-item" @click="copy(item.content)">
        <div class="clip-preview">{{ item.content?.slice(0, 120) }}</div>
        <div class="clip-meta">
          <span class="clip-time">{{ timeAgo(item.timestamp) }}</span>
          <button class="btn-delete" @click.stop="remove(item.id)" title="Delete">✕</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.clipboard-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
}
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.panel-header h2 {
  font-size: 15px;
  font-weight: 600;
  color: #e0e0e0;
}
.search-row {
  display: flex;
  gap: 6px;
  margin-bottom: 8px;
}
.search-clip {
  flex: 1;
  padding: 7px 10px;
  border: 1px solid #3a4a5e;
  border-radius: 6px;
  background: #2a3748;
  color: #e0e0e0;
  font-size: 13px;
  outline: none;
}
.search-clip:focus {
  border-color: #4a90d9;
}
.btn-text {
  background: none;
  border: none;
  color: #6a7a8e;
  cursor: pointer;
  font-size: 12px;
}
.btn-text:hover { color: #e0e0e0; }
.btn-icon {
  background: #2a3748;
  border: 1px solid #3a4a5e;
  border-radius: 6px;
  color: #a0b0c0;
  cursor: pointer;
  width: 32px;
  font-size: 16px;
}
.list {
  flex: 1;
  overflow-y: auto;
}
.clip-item {
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  margin-bottom: 4px;
  transition: background 0.12s;
}
.clip-item:hover { background: #2a3748; }
.clip-preview {
  font-size: 13px;
  color: #c0d0e0;
  line-height: 1.4;
  word-break: break-all;
  margin-bottom: 4px;
}
.clip-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.clip-time {
  font-size: 11px;
  color: #5a6a7e;
}
.btn-delete {
  background: none;
  border: none;
  color: #5a6a7e;
  cursor: pointer;
  font-size: 12px;
  padding: 2px 4px;
  border-radius: 3px;
}
.btn-delete:hover {
  color: #e74c3c;
  background: #3a1a1a;
}
.loading, .empty {
  text-align: center;
  padding: 24px;
  color: #5a6a7e;
  font-size: 13px;
}
</style>
