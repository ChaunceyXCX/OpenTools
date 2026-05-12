<script setup lang="ts">
import { useCommandStore } from '../stores/commandStore'

const store = useCommandStore()

function onClick(cmd: any) {
  store.launch(cmd)
}
</script>

<template>
  <div class="results">
    <div
      v-for="cmd in store.results"
      :key="cmd.name"
      class="result-item"
      @click="onClick(cmd)"
    >
      <div class="item-icon">
        <img v-if="cmd.icon" :src="cmd.icon" class="icon-img" />
        <span v-else class="icon-fallback">⚡</span>
      </div>
      <div class="item-info">
        <div class="item-name">{{ cmd.name }}</div>
        <div class="item-path">{{ cmd.path }}</div>
      </div>
    </div>
    <div v-if="store.loading" class="status">Scanning applications...</div>
    <div v-else-if="store.results.length === 0 && store.query" class="status">
      No results for "{{ store.query }}"
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
  transition: background 0.15s;
}
.result-item:hover {
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
  max-width: 300px;
}
.status {
  text-align: center;
  padding: 20px;
  color: #6a7a8e;
  font-size: 13px;
}
</style>
