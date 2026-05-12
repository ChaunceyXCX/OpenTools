<script setup lang="ts">
import { ref } from 'vue'
import { useCommandDataStore } from '../../stores/commandDataStore'

const store = useCommandDataStore()
const recentExpanded = ref(true)
const emit = defineEmits<{ launch: [cmd: any] }>()

const recentCommands = store.allCommands.slice(0, 8)
</script>

<template>
  <div class="aggregate">
    <div v-if="recentCommands.length" class="section">
      <div class="section-header" @click="recentExpanded = !recentExpanded">
        <span class="section-title">Recent</span>
        <span class="section-toggle">{{ recentExpanded ? '−' : '+' }}</span>
      </div>
      <div v-show="recentExpanded" class="section-content">
        <div
          v-for="cmd in recentCommands"
          :key="cmd.name"
          class="result-item"
          @click="emit('launch', cmd)"
        >
          <div class="item-icon-col">
            <img v-if="cmd.icon" :src="cmd.icon" class="item-icon" />
            <span v-else class="item-icon-fallback">⚡</span>
          </div>
          <div class="item-body">
            <div class="item-name">{{ cmd.name }}</div>
            <div class="item-path">{{ cmd.path }}</div>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="section">
      <div class="section-header">
        <span class="section-title">Welcome</span>
      </div>
      <div class="section-content welcome-content">
        {{ $t('search.shortcut') }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.aggregate { padding: 4px 8px; }
.section { margin-bottom: 4px; }
.section-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 6px 8px; cursor: pointer; border-radius: var(--radius-sm);
}
.section-header:hover { background: var(--hover-bg); }
.section-title {
  font-size: 11px; font-weight: 600; color: var(--text-secondary);
  text-transform: uppercase; letter-spacing: 0.5px;
}
.section-toggle { font-size: 14px; color: var(--text-secondary); }
.section-content { padding: 2px 0; }
.result-item {
  display: flex; align-items: center; gap: 10px;
  padding: 6px 8px; border-radius: var(--radius-sm);
  cursor: pointer; transition: background 0.1s;
}
.result-item:hover { background: var(--hover-bg); }
.item-icon-col {
  width: 28px; height: 28px; display: flex;
  align-items: center; justify-content: center; flex-shrink: 0;
}
.item-icon { width: 28px; height: 28px; object-fit: contain; }
.item-icon-fallback { font-size: 16px; }
.item-body { flex: 1; min-width: 0; }
.item-name { font-size: 13px; color: var(--text-color); }
.item-path { font-size: 11px; color: var(--text-secondary); }
.welcome-content {
  padding: 8px; font-size: 13px; color: var(--text-secondary); text-align: center;
}
kbd {
  padding: 1px 6px; background: var(--control-bg);
  border-radius: 4px; font-size: 11px;
  color: var(--text-color); border: 1px solid var(--border-color);
}
</style>
