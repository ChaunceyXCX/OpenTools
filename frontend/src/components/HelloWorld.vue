<script setup lang="ts">
import { ref } from 'vue'
import { Window } from '@wailsio/runtime'
import { GreetService } from '../../bindings/github.com/ChaunceyXCX/OpenTools'

defineProps<{ msg: string }>()

const name = ref('ZTools')
const result = ref('')

async function greet() {
  try {
    result.value = await GreetService.Greet(name.value)
  } catch (err) {
    console.error(err)
  }
}

function closeWindow() {
  Window.Hide()
}
</script>

<template>
  <div class="container">
    <h1>{{ msg }}</h1>
    <p>A high-performance launcher and plugin platform</p>
    <p>Wails v3 + Vue 3 + TypeScript</p>
    <div class="actions">
      <input v-model="name" @keyup.enter="greet" placeholder="Enter name..." />
      <button @click="greet">Greet</button>
      <button @click="closeWindow">Close</button>
    </div>
    <p v-if="result" class="result">{{ result }}</p>
  </div>
</template>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  gap: 12px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
  color: #e0e0e0;
  background: #1b2636;
}
h1 {
  font-size: 2em;
  font-weight: 600;
}
.actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
}
input {
  padding: 8px 12px;
  border: 1px solid #444;
  border-radius: 6px;
  background: #2a3748;
  color: #e0e0e0;
  font-size: 14px;
}
button {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  background: #4a90d9;
  color: white;
  cursor: pointer;
  font-size: 14px;
}
button:hover {
  background: #357abd;
}
.result {
  margin-top: 12px;
  padding: 8px 16px;
  background: #2a3748;
  border-radius: 6px;
}
</style>
