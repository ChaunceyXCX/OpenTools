<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { PluginService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'

interface PluginEntry {
  name: string
  title: string
  description: string
  version: string
  author: string
  logo?: string
  installed: boolean
}

const plugins = ref<PluginEntry[]>([])
const loading = ref(false)
const msg = ref('')

onMounted(async () => {
  loading.value = true
  try {
    const list = await PluginService.List()
    plugins.value = list.map((p: any) => ({ ...p, installed: true }))
  } catch {
    plugins.value = [
      { name: 'system', title: 'System', description: 'System commands (shutdown, reboot, etc)', version: '1.0.0', author: 'ZTools', installed: true },
      { name: 'translator', title: 'Translator', description: 'Offline multi-language translation', version: '1.0.0', author: 'Community', installed: false },
      { name: 'json-formatter', title: 'JSON Formatter', description: 'Format and validate JSON data', version: '0.9.0', author: 'Community', installed: false },
      { name: 'color-picker', title: 'Color Picker', description: 'Pick colors from screen', version: '1.2.0', author: 'Community', installed: false },
    ]
  }
  loading.value = false
})

async function install(name: string) {
  msg.value = `Installing ${name}...`
  try { await PluginService.Install(name); msg.value = `${name} installed` }
  catch { msg.value = `Failed to install ${name}` }
}

async function uninstall(name: string) {
  msg.value = `Uninstalling ${name}...`
  try { await PluginService.Uninstall(name); msg.value = `${name} uninstalled` }
  catch { msg.value = `Failed to uninstall ${name}` }
}
</script>

<template>
  <div class="plugin-market">
    <div class="market-header">
      <h2>Plugin Market</h2>
    </div>
    <div v-if="loading" class="loading">Loading plugins...</div>
    <div v-else class="plugin-list">
      <div v-for="p in plugins" :key="p.name" class="plugin-card">
        <div class="plugin-info">
          <div class="plugin-title">{{ p.title || p.name }}</div>
          <div class="plugin-desc">{{ p.description }}</div>
          <div class="plugin-meta">v{{ p.version }} · {{ p.author }}</div>
        </div>
        <div class="plugin-actions">
          <span v-if="p.installed" class="badge installed">Installed</span>
          <button v-else class="btn install-btn" @click="install(p.name)">Install</button>
          <button v-if="p.installed" class="btn danger" @click="uninstall(p.name)">Uninstall</button>
        </div>
      </div>
    </div>
    <div v-if="msg" class="toast">{{ msg }}</div>
  </div>
</template>

<style scoped>
.plugin-market { padding: 16px; overflow-y: auto; height: 100%; }
.market-header { margin-bottom: 12px; }
.market-header h2 { font-size: 16px; font-weight: 600; color: var(--text-color); }
.loading { text-align: center; padding: 24px; color: var(--text-secondary); font-size: 13px; }
.plugin-list { display: flex; flex-direction: column; gap: 8px; }
.plugin-card {
  display: flex; justify-content: space-between; align-items: center;
  padding: 12px; background: var(--card-bg); border-radius: var(--radius);
  border: 1px solid var(--border-color);
}
.plugin-info { flex: 1; min-width: 0; }
.plugin-title { font-size: 14px; font-weight: 600; color: var(--text-color); }
.plugin-desc { font-size: 12px; color: var(--text-secondary); margin: 2px 0; }
.plugin-meta { font-size: 11px; color: var(--text-secondary); }
.plugin-actions { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
.badge { font-size: 11px; padding: 3px 8px; border-radius: 4px; }
.badge.installed { background: color-mix(in srgb, var(--primary-color) 20%, transparent); color: var(--primary-color); }
.btn { padding: 6px 14px; border: none; border-radius: var(--radius-sm); color: var(--text-on-primary); cursor: pointer; font-size: 12px; }
.install-btn { background: var(--primary-color); }
.btn.danger { background: #c0392b; }
.toast {
  position: fixed; bottom: 20px; left: 50%; transform: translateX(-50%);
  padding: 8px 16px; background: var(--card-bg); border: 1px solid var(--border-color);
  border-radius: var(--radius); font-size: 13px; color: var(--text-color); z-index: 10;
}
</style>
