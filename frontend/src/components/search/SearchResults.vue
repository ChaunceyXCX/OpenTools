<script setup lang="ts">
import { computed } from 'vue'
import { useCommandDataStore } from '../../stores/commandDataStore'
import AggregateView from './AggregateView.vue'

const store = useCommandDataStore()
const props = defineProps<{ selectedIndex: number }>()
const emit = defineEmits<{ launch: [cmd: any] }>()

const showAggregate = computed(() => !store.searchQuery.trim() && store.allCommands.length > 0)

const typeLabels: Record<string, string> = { app: 'Applications', plugin: 'Plugins', builtin: 'Commands' }

let globalIdx = 0
function onClick(cmd: any) { store.launch(cmd) }
</script>

<template>
  <div class="scrollable-content">
    <AggregateView v-if="showAggregate" @launch="onClick" />

    <div v-else-if="store.loading" class="state-msg"><div class="spinner" />{{ $t('search.loading') }}</div>

    <div v-else-if="store.results.length > 0" class="results-list">
      <template v-for="(cmds, type) in store.groupedResults" :key="type">
        <div class="group-header">{{ typeLabels[type] || type }}</div>
        <div v-for="cmd in cmds" :key="cmd.name"
          :class="['result-item', { active: globalIdx === selectedIndex }]"
          @click="onClick(cmd)">
          <div class="item-icon-col">
            <img v-if="cmd.icon" :src="cmd.icon" class="item-icon" />
            <span v-else class="item-icon-fallback">⚡</span>
          </div>
          <div class="item-body">
            <div class="item-name">{{ cmd.name }}</div>
            <div class="item-path">{{ cmd.path }}</div>
          </div>
          <div v-if="globalIdx === selectedIndex" class="item-hint">↵</div>
        </div>
        <span style="display:none">{{ globalIdx++ }}</span>
      </template>
    </div>

    <div v-else-if="store.searchQuery" class="state-msg">{{ $t('search.noResults', { query: store.searchQuery }) }}</div>
    <div v-else class="state-msg"><div class="welcome-icon">⌕</div>{{ $t('search.welcome') }}</div>
  </div>
</template>


<style scoped>
.scrollable-content { flex: 1; overflow-y: auto; padding: 4px 0; }
.results-list { padding: 0 8px; }
.group-header { font-size: 11px; font-weight: 600; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.5px; padding: 8px 8px 4px; }
.result-item { display: flex; align-items: center; gap: 10px; padding: 8px 10px; border-radius: var(--radius-sm); cursor: pointer; }
.result-item.active, .result-item:hover { background: var(--hover-bg); }
.item-icon-col { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.item-icon { width: 32px; height: 32px; object-fit: contain; }
.item-icon-fallback { font-size: 18px; }
.item-body { flex: 1; min-width: 0; }
.item-name { font-size: 14px; font-weight: 500; color: var(--text-color); }
.item-path { font-size: 11px; color: var(--text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.item-hint { color: var(--primary-color); font-size: 13px; }
.state-msg { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; padding: 40px 20px; color: var(--text-secondary); font-size: 13px; }
.welcome-icon { font-size: 32px; opacity: 0.3; }
.spinner { width: 18px; height: 18px; border: 2px solid var(--border-color); border-top-color: var(--primary-color); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
