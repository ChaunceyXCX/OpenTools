<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ClipboardService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'

const { t } = useI18n()
const items = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')
let refreshTimer: any = null

const filtered = computed(() => {
  if (!searchQuery.value) return items.value
  const q = searchQuery.value.toLowerCase()
  return items.value.filter((i: any) => i.content?.toLowerCase().includes(q))
})

onMounted(() => { load(); refreshTimer = setInterval(load, 3000) })
onUnmounted(() => clearInterval(refreshTimer))

async function load() {
  loading.value = true
  try { items.value = await ClipboardService.GetHistory(1, 50) } catch { }
  loading.value = false
}

async function remove(id: string) { await ClipboardService.Delete(id); await load() }
async function clearAll() { await ClipboardService.Clear(); items.value = [] }
function copy(text: string) { navigator.clipboard.writeText(text) }

function timeAgo(ts: number): string {
  const diff = Date.now() - ts
  if (diff < 60000) return t('clipboard.justNow')
  if (diff < 3600000) return t('clipboard.minutesAgo', { m: Math.floor(diff / 60000) })
  if (diff < 86400000) return t('clipboard.hoursAgo', { h: Math.floor(diff / 3600000) })
  return t('clipboard.daysAgo', { d: Math.floor(diff / 86400000) })
}
</script>

<template>
  <div class="clipboard-panel">
    <div class="panel-header">
      <h2>{{ $t('clipboard.title') }}</h2>
      <div class="header-actions">
        <button v-if="items.length" class="btn-text" @click="clearAll">{{ $t('clipboard.clearAll') }}</button>
      </div>
    </div>

    <div class="search-row">
      <input v-model="searchQuery" :placeholder="$t('clipboard.searchPlaceholder')" class="search-clip" />
      <button class="btn-icon" @click="load" :title="$t('clipboard.refresh')">↻</button>
    </div>

    <div v-if="loading" class="loading">{{ $t('clipboard.loading') }}</div>
    <div v-else-if="filtered.length === 0" class="empty">
      {{ searchQuery ? $t('clipboard.noMatches') : $t('clipboard.empty') }}
    </div>
    <div v-else class="list">
      <div v-for="item in filtered" :key="item.id" class="clip-item" @click="copy(item.content)">
        <div class="clip-preview">{{ item.content?.slice(0, 120) }}</div>
        <div class="clip-meta">
          <span class="clip-time">{{ timeAgo(item.timestamp) }}</span>
          <button class="btn-delete" @click.stop="remove(item.id)">✕</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.clipboard-panel { display: flex; flex-direction: column; height: 100%; padding: 12px; }
.panel-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.panel-header h2 { font-size: 15px; font-weight: 600; color: var(--text-color); }
.search-row { display: flex; gap: 6px; margin-bottom: 8px; }
.search-clip { flex: 1; padding: 7px 10px; border: 1px solid var(--border-color); border-radius: var(--radius-sm); background: var(--control-bg); color: var(--text-color); font-size: 13px; outline: none; }
.search-clip:focus { border-color: var(--primary-color); }
.btn-text { background: none; border: none; color: var(--text-secondary); cursor: pointer; font-size: 12px; }
.btn-text:hover { color: var(--text-color); }
.btn-icon { background: var(--control-bg); border: 1px solid var(--border-color); border-radius: var(--radius-sm); color: var(--text-secondary); cursor: pointer; width: 32px; font-size: 16px; }
.list { flex: 1; overflow-y: auto; }
.clip-item { padding: 8px 10px; border-radius: var(--radius-sm); cursor: pointer; margin-bottom: 4px; }
.clip-item:hover { background: var(--hover-bg); }
.clip-preview { font-size: 13px; color: var(--text-color); line-height: 1.4; word-break: break-all; margin-bottom: 4px; }
.clip-meta { display: flex; justify-content: space-between; align-items: center; }
.clip-time { font-size: 11px; color: var(--text-secondary); }
.btn-delete { background: none; border: none; color: var(--text-secondary); cursor: pointer; font-size: 12px; padding: 2px 4px; border-radius: 3px; }
.btn-delete:hover { color: #e74c3c; background: #3a1a1a; }
.loading, .empty { text-align: center; padding: 24px; color: var(--text-secondary); font-size: 13px; }
</style>
