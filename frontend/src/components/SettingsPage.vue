<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLanguage, setTheme, setColor, getStoredTheme, getStoredColor } from '../i18n'
import { DatabaseService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'

const { t, locale } = useI18n()
const lang = locale as any
const currentTheme = ref(getStoredTheme())
const currentColor = ref(getStoredColor())
const colors = ['blue', 'purple', 'green', 'orange', 'red', 'pink']

function switchTheme(theme: string) { currentTheme.value = theme; setTheme(theme) }
function switchColor(color: string) { currentColor.value = color; setColor(color) }

const dbGet = ref('')
const dbSetKey = ref('')
const dbSetVal = ref('')
const dbResult = ref('')
const wdUrl = ref('')
const wdUser = ref('')
const wdPass = ref('')
const wdDir = ref('/opentools')
const syncStatus = ref('')

async function onDbGet() {
  try { const doc = await DatabaseService.Get(dbGet.value); dbResult.value = doc ? JSON.stringify(doc) : '(nil)' } catch (e: any) { dbResult.value = 'error: ' + e.message }
}
async function onDbSet() {
  try { const r = await DatabaseService.Put(dbSetKey.value, dbSetVal.value); dbResult.value = JSON.stringify(r) } catch (e: any) { dbResult.value = 'error: ' + e.message }
}
async function onDbDelete() {
  try { const r = await DatabaseService.Remove(dbGet.value); dbResult.value = JSON.stringify(r) } catch (e: any) { dbResult.value = 'error: ' + e.message }
}
async function onSyncPush() { syncStatus.value = t('settings.syncConfigured') }

function switchLang(lang: string) { setLanguage(lang) }
</script>

<template>
  <div class="settings">
    <h2 class="title">{{ $t('settings.title') }}</h2>

    <!-- Language -->
    <section class="section">
      <h3>{{ $t('settings.language') }}</h3>
      <div class="btn-row">
        <button :class="['opt-btn', { active: lang === 'zh-CN' }]" @click="switchLang('zh-CN')">中文</button>
        <button :class="['opt-btn', { active: lang === 'en-US' }]" @click="switchLang('en-US')">English</button>
      </div>
    </section>

    <!-- Theme -->
    <section class="section">
      <h3>Theme</h3>
      <div class="btn-row">
        <button :class="['opt-btn', { active: currentTheme === 'dark' }]" @click="switchTheme('dark')">Dark</button>
        <button :class="['opt-btn', { active: currentTheme === 'light' }]" @click="switchTheme('light')">Light</button>
      </div>
      <div class="color-row">
        <button v-for="c in colors" :key="c"
          :class="['color-dot', `theme-${c}`, { active: currentColor === c }]"
          @click="switchColor(c)"></button>
      </div>
    </section>

    <!-- Shortcuts -->
    <section class="section">
      <h3>⌨ {{ $t('settings.shortcuts') }}</h3>
      <div class="shortcut"><span>{{ $t('settings.toggleWindow') }}</span><kbd>Alt + Z</kbd></div>
      <div class="shortcut"><span>{{ $t('settings.hideWindow') }}</span><kbd>Esc</kbd></div>
      <div class="shortcut"><span>{{ $t('settings.navigateResults') }}</span><kbd>↑ ↓</kbd></div>
      <div class="shortcut"><span>{{ $t('settings.launchSelected') }}</span><kbd>Enter</kbd></div>
    </section>

    <!-- WebDAV -->
    <section class="section">
      <h3>☁ {{ $t('settings.webdav') }}</h3>
      <div class="field-row">
        <input v-model="wdUrl" :placeholder="t('settings.serverUrl')" class="input" />
      </div>
      <div class="field-row">
        <input v-model="wdUser" :placeholder="t('settings.username')" class="input half" />
        <input v-model="wdPass" type="password" :placeholder="t('settings.password')" class="input half" />
      </div>
      <div class="field-row">
        <input v-model="wdDir" :placeholder="t('settings.remoteDir')" class="input" />
        <button class="btn" @click="onSyncPush">{{ $t('settings.save') }}</button>
      </div>
      <div v-if="syncStatus" class="result">{{ syncStatus }}</div>
    </section>

    <!-- Database -->
    <section class="section">
      <h3>🗄 {{ $t('settings.database') }}</h3>
      <div class="field-row">
        <input v-model="dbGet" :placeholder="t('settings.key')" class="input" />
        <button class="btn" @click="onDbGet">{{ $t('settings.get') }}</button>
        <button class="btn danger" @click="onDbDelete">{{ $t('settings.delete') }}</button>
      </div>
      <div class="field-row">
        <input v-model="dbSetKey" :placeholder="t('settings.key')" class="input half" />
        <input v-model="dbSetVal" :placeholder="t('settings.value')" class="input half" />
        <button class="btn" @click="onDbSet">{{ $t('settings.put') }}</button>
      </div>
      <pre v-if="dbResult" class="result db-result">{{ dbResult }}</pre>
    </section>

    <!-- About -->
    <section class="section about">
      <div class="about-row"><span>{{ $t('settings.version') }}</span><span>2.4.1 (OpenTools)</span></div>
      <div class="about-row"><span>{{ $t('settings.runtime') }}</span><span>Go + Wails v3</span></div>
      <div class="about-row"><span>{{ $t('settings.frontend') }}</span><span>Vue 3 + Pinia + i18n</span></div>
    </section>
  </div>
</template>

<style scoped>
.settings { padding: 16px; height: 100%; overflow-y: auto; }
.title { font-size: 18px; font-weight: 600; margin-bottom: 16px; color: var(--text-color); }
.section { margin-bottom: 20px; padding-bottom: 16px; border-bottom: 1px solid var(--border-color); }
.section h3 { font-size: 12px; font-weight: 600; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 10px; }
.field-row { display: flex; gap: 6px; margin-bottom: 6px; }
.input { flex: 1; padding: 7px 10px; border: 1px solid var(--border-color); border-radius: var(--radius-sm); background: var(--control-bg); color: var(--text-color); font-size: 13px; outline: none; }
.input.half { flex: 0.5; }
.input:focus { border-color: var(--primary-color); }
.btn { padding: 7px 14px; border: none; border-radius: var(--radius-sm); background: var(--primary-color); color: var(--text-on-primary); cursor: pointer; font-size: 13px; white-space: nowrap; }
.btn.danger { background: #c0392b; }
.result { margin-top: 6px; padding: 6px 8px; font-size: 12px; color: var(--text-secondary); }
.db-result { background: color-mix(in srgb, var(--bg-color) 90%, #000); border-radius: 4px; overflow-x: auto; }
.shortcut { display: flex; justify-content: space-between; align-items: center; padding: 5px 0; font-size: 13px; color: var(--text-color); }
kbd { padding: 2px 8px; background: var(--control-bg); border-radius: 4px; font-size: 11px; color: var(--text-color); border: 1px solid var(--border-color); }
.about { border-bottom: none; }
.about-row { display: flex; justify-content: space-between; padding: 4px 0; font-size: 13px; color: var(--text-secondary); }
.btn-row { display: flex; gap: 8px; }
.opt-btn { padding: 6px 16px; border: 1px solid var(--border-color); border-radius: var(--radius-sm); background: var(--control-bg); color: var(--text-color); cursor: pointer; font-size: 13px; }
.opt-btn.active { border-color: var(--primary-color); color: var(--primary-color); background: color-mix(in srgb, var(--primary-color) 15%, transparent); }
.color-row { display: flex; gap: 8px; margin-top: 8px; }
.color-dot { width: 28px; height: 28px; border-radius: 50%; border: 2px solid transparent; cursor: pointer; transition: transform 0.15s; }
.color-dot.active { border-color: var(--text-color); transform: scale(1.15); }
.color-dot.theme-blue { background: #4a90d9; }
.color-dot.theme-purple { background: #7c3aed; }
.color-dot.theme-green { background: #059669; }
.color-dot.theme-orange { background: #ea580c; }
.color-dot.theme-red { background: #dc2626; }
.color-dot.theme-pink { background: #db2777; }
</style>
