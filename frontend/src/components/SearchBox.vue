<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useCommandStore } from '../stores/commandStore'

const store = useCommandStore()
const input = ref<HTMLInputElement>()

const emit = defineEmits<{
  keynav: [dir: 'up' | 'down']
  confirm: []
}>()

onMounted(() => {
  store.loadCommands()
  input.value?.focus()
})

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    emit('keynav', 'down')
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    emit('keynav', 'up')
  } else if (e.key === 'Enter') {
    e.preventDefault()
    emit('confirm')
  }
}
</script>

<template>
  <div class="search-box">
    <input
      ref="input"
      v-model="store.query"
      type="text"
      class="search-input"
      :placeholder="store.loading ? 'Scanning applications...' : 'Search apps, type a command...'"
      autofocus
      spellcheck="false"
      autocomplete="off"
      @keydown="onKeydown"
    />
  </div>
</template>

<style scoped>
.search-box {
  padding: 8px 12px;
}
.search-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #3a4a5e;
  border-radius: 8px;
  background: #2a3748;
  color: #e0e0e0;
  font-size: 15px;
  outline: none;
  transition: border-color 0.2s;
}
.search-input:focus {
  border-color: #4a90d9;
}
.search-input::placeholder {
  color: #6a7a8e;
}
</style>
