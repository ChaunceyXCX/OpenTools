<script setup lang="ts">
import { ref } from 'vue'
import { DatabaseService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'

const dbGet = ref('')
const dbSetKey = ref('')
const dbSetVal = ref('')
const dbResult = ref('')

async function onDbGet() {
  try {
    const doc = await DatabaseService.Get(dbGet.value)
    dbResult.value = doc ? JSON.stringify(doc) : '(nil)'
  } catch (e: any) {
    dbResult.value = 'error: ' + e.message
  }
}

async function onDbSet() {
  try {
    const result = await DatabaseService.Put(dbSetKey.value, dbSetVal.value)
    dbResult.value = JSON.stringify(result)
  } catch (e: any) {
    dbResult.value = 'error: ' + e.message
  }
}

async function onDbDelete() {
  try {
    const result = await DatabaseService.Remove(dbGet.value)
    dbResult.value = JSON.stringify(result)
  } catch (e: any) {
    dbResult.value = 'error: ' + e.message
  }
}
</script>

<template>
  <div class="settings">
    <h2 class="title">Settings</h2>

    <section class="section">
      <h3>Database</h3>
      <div class="row">
        <input v-model="dbGet" placeholder="Key (e.g. ZTOOLS/test)" class="input" />
        <button @click="onDbGet" class="btn">Get</button>
        <button @click="onDbDelete" class="btn danger">Delete</button>
      </div>
      <div class="row">
        <input v-model="dbSetKey" placeholder="Key" class="input" />
        <input v-model="dbSetVal" placeholder="Value" class="input" />
        <button @click="onDbSet" class="btn">Put</button>
      </div>
      <pre v-if="dbResult" class="result">{{ dbResult }}</pre>
    </section>

    <section class="section">
      <h3>Keyboard Shortcuts</h3>
      <div class="shortcut">
        <span>Toggle Window</span>
        <kbd>Alt+Z</kbd>
      </div>
      <div class="shortcut">
        <span>Hide Window</span>
        <kbd>Esc</kbd>
      </div>
    </section>

    <div class="version">
      ZTools Wails Port · v2.4.1
    </div>
  </div>
</template>

<style scoped>
.settings {
  padding: 16px;
  height: 100%;
  overflow-y: auto;
}
.title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #e0e0e0;
}
.section {
  margin-bottom: 20px;
}
.section h3 {
  font-size: 13px;
  font-weight: 600;
  color: #8a9aaa;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}
.row {
  display: flex;
  gap: 6px;
  margin-bottom: 6px;
}
.input {
  flex: 1;
  padding: 6px 10px;
  border: 1px solid #3a4a5e;
  border-radius: 6px;
  background: #2a3748;
  color: #e0e0e0;
  font-size: 13px;
}
.btn {
  padding: 6px 14px;
  border: none;
  border-radius: 6px;
  background: #4a90d9;
  color: white;
  cursor: pointer;
  font-size: 13px;
}
.btn.danger {
  background: #c0392b;
}
.result {
  margin-top: 8px;
  padding: 8px;
  background: #111820;
  border-radius: 4px;
  font-size: 12px;
  color: #8a9aaa;
  overflow-x: auto;
}
.shortcut {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
  font-size: 13px;
  color: #c0c0c0;
}
kbd {
  padding: 2px 8px;
  background: #2a3748;
  border-radius: 4px;
  font-size: 12px;
  color: #e0e0e0;
  border: 1px solid #3a4a5e;
}
.version {
  text-align: center;
  padding: 16px;
  color: #5a6a7e;
  font-size: 12px;
}
</style>
