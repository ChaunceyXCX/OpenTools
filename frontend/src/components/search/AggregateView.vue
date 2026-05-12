<script setup lang="ts">
import { ref } from 'vue'
import { HistoryService } from '../../stores/commandDataStore'

const recentExpanded = ref(true)
const emit = defineEmits<{ launch: [cmd: any] }>()
const recentItems = HistoryService.getRecent()
</script>

<template>
  <div class="aggregate">
    <div v-if="recentItems.length" class="section">
      <div class="section-header" @click="recentExpanded = !recentExpanded">
        <span class="section-title">Recent</span>
        <span class="section-toggle">{{ recentExpanded ? '−' : '+' }}</span>
      </div>
      <div v-show="recentExpanded">
        <div v-for="item in recentItems" :key="item.name" class="result-item" @click="emit('launch', item)">
          <div class="item-icon-col">⚡</div>
          <div class="item-body"><div class="item-name">{{ item.name }}</div><div class="item-path">{{ item.path }}</div></div>
        </div>
      </div>
    </div>
    <div v-else class="section">
      <div class="section-header"><span class="section-title">Welcome</span></div>
      <div class="welcome-content">{{ $t('search.shortcut') }}</div>
    </div>
  </div>
</template>

<style scoped>
.aggregate { padding: 4px 8px; }
.section { margin-bottom: 4px; }
.section-header { display: flex; justify-content: space-between; align-items: center; padding: 6px 8px; cursor: pointer; border-radius: var(--radius-sm); }
.section-header:hover { background: var(--hover-bg); }
.section-title { font-size: 11px; font-weight: 600; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.5px; }
.section-toggle { font-size: 14px; color: var(--text-secondary); }
.result-item { display: flex; align-items: center; gap: 10px; padding: 6px 8px; border-radius: var(--radius-sm); cursor: pointer; }
.result-item:hover { background: var(--hover-bg); }
.item-icon-col { width: 28px; height: 28px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.item-body { flex: 1; min-width: 0; }
.item-name { font-size: 13px; color: var(--text-color); }
.item-path { font-size: 11px; color: var(--text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.welcome-content { padding: 8px; font-size: 13px; color: var(--text-secondary); text-align: center; }
</style>
