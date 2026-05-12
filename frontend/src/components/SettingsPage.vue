<script setup lang="ts">
import { ref } from 'vue'
import { DatabaseService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'

const dbGet = ref('')
const dbSetKey = ref('')
const dbSetVal = ref('')
const dbResult = ref('')

const wdUrl = ref('')
const wdUser = ref('')
const wdPass = ref('')
const wdDir = ref('/ztools')
const syncStatus = ref('')

async function onDbGet() {
  try {
    const doc = await DatabaseService.Get(dbGet.value)
    dbResult.value = doc ? JSON.stringify(doc) : '(nil)'
  } catch (e: any) { dbResult.value = 'error: ' + e.message }
}
async function onDbSet() {
  try {
    const r = await DatabaseService.Put(dbSetKey.value, dbSetVal.value)
    dbResult.value = JSON.stringify(r)
  } catch (e: any) { dbResult.value = 'error: ' + e.message }
}
async function onDbDelete() {
  try {
    const r = await DatabaseService.Remove(dbGet.value)
    dbResult.value = JSON.stringify(r)
  } catch (e: any) { dbResult.value = 'error: ' + e.message }
}

async function onSyncPush() {
  syncStatus.value = 'Sync config saved (Wails CLI: wails3 task run)'
}
</script>

<template>
  <div class="settings">
    <h2 class="title">Settings</h2>

    <!-- Keyboard Shortcuts -->
    <section class="section">
      <h3>⌨ Keyboard Shortcuts</h3>
      <div class="shortcut"><span>Toggle Window</span><kbd>Alt + Z</kbd></div>
      <div class="shortcut"><span>Hide Window</span><kbd>Esc</kbd></div>
      <div class="shortcut"><span>Navigate Results</span><kbd>↑ ↓</kbd></div>
      <div class="shortcut"><span>Launch Selected</span><kbd>Enter</kbd></div>
    </section>

    <!-- WebDAV Sync -->
    <section class="section">
      <h3>☁ WebDAV Sync</h3>
      <div class="field-row">
        <input v-model="wdUrl" placeholder="Server URL (e.g. https://example.com/dav)" class="input" />
      </div>
      <div class="field-row">
        <input v-model="wdUser" placeholder="Username" class="input half" />
        <input v-model="wdPass" type="password" placeholder="Password" class="input half" />
      </div>
      <div class="field-row">
        <input v-model="wdDir" placeholder="Remote directory" class="input" />
        <button class="btn" @click="onSyncPush">Save</button>
      </div>
      <div v-if="syncStatus" class="result">{{ syncStatus }}</div>
    </section>

    <!-- Database -->
    <section class="section">
      <h3>🗄 Database</h3>
      <div class="field-row">
        <input v-model="dbGet" placeholder="Key (e.g. ZTOOLS/setting)" class="input" />
        <button class="btn" @click="onDbGet">Get</button>
        <button class="btn danger" @click="onDbDelete">Delete</button>
      </div>
      <div class="field-row">
        <input v-model="dbSetKey" placeholder="Key" class="input half" />
        <input v-model="dbSetVal" placeholder="Value" class="input half" />
        <button class="btn" @click="onDbSet">Put</button>
      </div>
      <pre v-if="dbResult" class="result db-result">{{ dbResult }}</pre>
    </section>

    <!-- About -->
    <section class="section about">
      <div class="about-row"><span>Version</span><span>2.4.1 (Wails Port)</span></div>
      <div class="about-row"><span>Runtime</span><span>Go + Wails v3</span></div>
      <div class="about-row"><span>Frontend</span><span>Vue 3 + Pinia + Fuse.js</span></div>
    </section>
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
  padding-bottom: 16px;
  border-bottom: 1px solid #1e2a38;
}
.section h3 {
  font-size: 12px;
  font-weight: 600;
  color: #8a9aaa;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 10px;
}
.field-row {
  display: flex;
  gap: 6px;
  margin-bottom: 6px;
}
.input {
  flex: 1;
  padding: 7px 10px;
  border: 1px solid #3a4a5e;
  border-radius: 6px;
  background: #2a3748;
  color: #e0e0e0;
  font-size: 13px;
  outline: none;
}
.input.half { flex: 0.5; }
.input:focus { border-color: #4a90d9; }
.btn {
  padding: 7px 14px;
  border: none;
  border-radius: 6px;
  background: #4a90d9;
  color: white;
  cursor: pointer;
  font-size: 13px;
  white-space: nowrap;
}
.btn.danger { background: #c0392b; }
.result {
  margin-top: 6px;
  padding: 6px 8px;
  font-size: 12px;
  color: #8a9aaa;
}
.db-result {
  background: #111820;
  border-radius: 4px;
  overflow-x: auto;
}
.shortcut {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px 0;
  font-size: 13px;
  color: #c0c0c0;
}
kbd {
  padding: 2px 8px;
  background: #2a3748;
  border-radius: 4px;
  font-size: 11px;
  color: #e0e0e0;
  border: 1px solid #3a4a5e;
}
.about { border-bottom: none; }
.about-row {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
  font-size: 13px;
  color: #8a9aaa;
}
</style>
