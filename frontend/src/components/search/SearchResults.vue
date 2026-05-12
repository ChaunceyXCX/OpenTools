<script setup lang="ts">
import { computed } from 'vue'
import { useCommandDataStore } from '../../stores/commandDataStore'
import AggregateView from './AggregateView.vue'

const store = useCommandDataStore()

const props = defineProps<{ selectedIndex: number }>()
const emit = defineEmits<{ launch: [cmd: any] }>()

const showAggregate = computed(() => !store.searchQuery.trim() && store.allCommands.length > 0)

function onClick(cmd: any) {
  store.launch(cmd)
}
</script>

<template>
  <div class="scrollable-content">
    <AggregateView v-if="showAggregate" @launch="onClick" />

    <div v-else-if="store.loading" class="state-msg">
      <div class="spinner" />
      <span>Scanning...</span>
    </div>

    <div v-else-if="store.results.length > 0" class="results-list">
      <div
        v-for="(cmd, i) in store.results"
        :key="cmd.name"
        :class="['result-item', { active: i === selectedIndex }]"
        @click="onClick(cmd)"
      >
        <div class="item-icon-col">
          <img v-if="cmd.icon" :src="cmd.icon" class="item-icon" />
          <span v-else class="item-icon-fallback">⚡</span>
        </div>
        <div class="item-body">
          <div class="item-name">{{ cmd.name }}</div>
          <div class="item-path">{{ cmd.path }}</div>
        </div>
        <div v-if="i === selectedIndex" class="item-hint">↵</div>
      </div>
    </div>

    <div v-else-if="store.searchQuery" class="state-msg">
      No results for "{{ store.searchQuery }}"
    </div>

    <div v-else class="state-msg">
      <div class="welcome-icon">⌕</div>
      <span>Type to start searching</span>
    </div>
  </div>
</template>

<style scoped>
.scrollable-content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}
.results-list { padding: 0 8px; }
.result-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.1s;
}
.result-item.active, .result-item:hover { background: var(--hover-bg); }
.item-icon-col {
  width: 32px; height: 32px; display: flex;
  align-items: center; justify-content: center; flex-shrink: 0;
}
.item-icon { width: 32px; height: 32px; object-fit: contain; }
.item-icon-fallback { font-size: 18px; }
.item-body { flex: 1; min-width: 0; }
.item-name { font-size: 14px; font-weight: 500; color: var(--text-color); }
.item-path {
  font-size: 11px; color: var(--text-secondary);
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.item-hint { color: var(--primary-color); font-size: 13px; }
.state-msg {
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  gap: 8px; padding: 40px 20px;
  color: var(--text-secondary); font-size: 13px;
}
.welcome-icon { font-size: 32px; opacity: 0.3; }
.spinner {
  width: 18px; height: 18px;
  border: 2px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
