import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import Fuse from 'fuse.js'
import { CommandsService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'

export interface Command {
  id?: string
  name: string
  path: string
  icon?: string
  aliases?: string[]
  acronym?: string
  type?: string
  subType?: string
  pluginName?: string
  featureCode?: string
  cmdType?: string
}

export interface PastedFile {
  name: string
  path: string
  isDirectory: boolean
}

export const useCommandDataStore = defineStore('commandData', () => {
  const allCommands = ref<Command[]>([])
  const searchQuery = ref('')
  const loading = ref(false)
  const selectedIndex = ref(0)
  const pastedImage = ref<string | null>(null)
  const pastedFiles = ref<PastedFile[] | null>(null)
  const pastedText = ref<string | null>(null)

  let fuse: Fuse<Command> | null = null

  const results = computed(() => {
    if (!searchQuery.value.trim()) {
      return allCommands.value.slice(0, 10)
    }
    if (!fuse) return []
    return fuse.search(searchQuery.value).map(r => r.item).slice(0, 12)
  })

  const groupedResults = computed(() => {
    const groups: Record<string, Command[]> = {}
    for (const r of results.value) {
      const key = r.type === 'direct' ? (r.subType || 'app') : (r.type || 'other')
      if (!groups[key]) groups[key] = []
      groups[key].push(r)
    }
    return groups
  })

  async function loadCommands() {
    loading.value = true
    try {
      const cmds = await CommandsService.ScanApplications()
      allCommands.value = cmds.map((c: any) => ({
        ...c, id: c.name,
        type: c.type || 'direct',
        subType: c.subType || 'app',
      }))
      fuse = new Fuse(allCommands.value, {
        keys: [
          { name: 'name', weight: 2 },
          { name: 'aliases', weight: 1.5 },
          { name: 'acronym', weight: 1 },
        ],
        threshold: 0.4,
        includeScore: true,
      })
    } catch (err) {
      console.error('load commands error:', err)
    } finally {
      loading.value = false
    }
  }

  async function launch(cmd: Command) {
    try {
      const r = await CommandsService.Launch(cmd.path)
      if (!r.success) console.error('launch error:', r.error)
      HistoryService.addEntry(cmd.name, cmd.path)
    } catch (err) {
      console.error('launch error:', err)
    }
  }

  function setQuery(q: string) { searchQuery.value = q }
  function clearPaste() {
    pastedImage.value = null
    pastedFiles.value = null
    pastedText.value = null
  }

  return {
    allCommands, searchQuery, loading, selectedIndex,
    pastedImage, pastedFiles, pastedText,
    results, groupedResults, loadCommands, launch, setQuery, clearPaste,
  }
})

class HistoryService {
  static addEntry(name: string, path: string) {
    try {
      const key = 'opentools-recent'
      const raw = localStorage.getItem(key) || '[]'
      const entries = JSON.parse(raw) as { name: string; path: string; ts: number }[]
      const filtered = entries.filter(e => e.name !== name)
      filtered.unshift({ name, path, ts: Date.now() })
      localStorage.setItem(key, JSON.stringify(filtered.slice(0, 20)))
    } catch { }
  }
  static getRecent(): { name: string; path: string; ts: number }[] {
    try {
      return JSON.parse(localStorage.getItem('opentools-recent') || '[]')
    } catch { return [] }
  }
}
export { HistoryService }
