<script setup lang="ts">
import { useCommandStore } from '../stores/commandStore'

const store = useCommandStore()

const props = defineProps<{ selectedIndex: number }>()
const emit = defineEmits<{ launch: [cmd: any] }>()

function onClick(cmd: any) {
  store.launch(cmd)
}
</script>

<template>
  <div class="results">
    <div v-if="store.loading" class="state-msg">
      <div class="spinner" />
      <span>Scanning installed applications...</span>
    </div>

    <template v-else-if="store.results.length > 0">
      <div
        v-for="(cmd, i) in store.results"
        :key="cmd.name"
        :class="['result-item', { active: i === selectedIndex }]"
        @click="onClick(cmd)"
        @dblclick="onClick(cmd)"
      >
        <div class="item-icon">
          <img v-if="cmd.icon" :src="cmd.icon" class="icon-img" />
          <span v-else class="icon-fallback">⚡</span>
        </div>
        <div class="item-info">
          <div class="item-name">{{ cmd.name }}</div>
          <div class="item-path">{{ cmd.path }}</div>
        </div>
        <div v-if="i === selectedIndex" class="item-hint">↵ launch</div>
      </div>
      <div class="result-count">{{ store.results.length }} result{{ store.results.length !== 1 ? 's' : '' }}</div>
    </template>

    <div v-else-if="store.query" class="state-msg">
      No results for "{{ store.query }}"
    </div>

    <div v-else class="state-msg welcome">
      <div class="welcome-icon">⌨</div>
      <div>Type to search applications</div>
      <div class="welcome-hint">Use ↑↓ to navigate · Enter to launch · Esc to close</div>
    </div>
  </div>
</template>

<style scoped>
.results {
  flex: 1;
  overflow-y: auto;
  padding: 4px 8px;
}
.result-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.12s;
}
.result-item:hover,
.result-item.active {
  background: #2a3748;
}
.item-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.icon-img {
  width: 32px;
  height: 32px;
  object-fit: contain;
}
.icon-fallback {
  font-size: 18px;
}
.item-info {
  flex: 1;
  min-width: 0;
}
.item-name {
  font-size: 14px;
  font-weight: 500;
  color: #e0e0e0;
}
.item-path {
  font-size: 11px;
  color: #6a7a8e;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.item-hint {
  font-size: 11px;
  color: #4a90d9;
  padding: 2px 6px;
  border: 1px solid #4a90d933;
  border-radius: 4px;
}
.result-count {
  text-align: center;
  padding: 6px;
  color: #4a5a6e;
  font-size: 11px;
}
.state-msg {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 20px;
  color: #6a7a8e;
  font-size: 13px;
}
.welcome-icon {
  font-size: 32px;
  opacity: 0.4;
}
.welcome-hint {
  font-size: 11px;
  color: #4a5a6e;
}
.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #3a4a5e;
  border-top-color: #4a90d9;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
