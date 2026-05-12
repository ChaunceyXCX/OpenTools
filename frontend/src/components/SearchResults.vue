<script setup lang="ts">
import { useCommandStore } from '../stores/commandStore'

const store = useCommandStore()
</script>

<template>
  <div class="results">
    <div
      v-for="cmd in store.results"
      :key="cmd.id"
      class="result-item"
    >
      <div class="item-icon">
        {{ cmd.icon ? '🔍' : '⚡' }}
      </div>
      <div class="item-info">
        <div class="item-name">{{ cmd.name }}</div>
        <div class="item-type">{{ cmd.type }}{{ cmd.subType ? ' · ' + cmd.subType : '' }}</div>
      </div>
    </div>
    <div v-if="store.results.length === 0 && store.query" class="empty">
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
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}
.item-name {
  font-size: 14px;
  font-weight: 500;
  color: #e0e0e0;
}
.item-type {
  font-size: 11px;
  color: #6a7a8e;
}
.empty {
  text-align: center;
  padding: 20px;
  color: #6a7a8e;
  font-size: 13px;
}
</style>
