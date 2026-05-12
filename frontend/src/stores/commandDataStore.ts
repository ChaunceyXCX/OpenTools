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
}

export const useCommandDataStore = defineStore('commandData', () => {
  const allCommands = ref<Command[]>([])
  const searchQuery = ref('')
  const loading = ref(false)
  const selectedIndex = ref(0)

  let fuse: Fuse<Command> | null = null

  const results = computed(() => {
    if (!searchQuery.value.trim()) {
      return allCommands.value.slice(0, 10)
    }
    if (!fuse) return []
    return fuse.search(searchQuery.value).map(r => r.item).slice(0, 12)
  })

  async function loadCommands() {
    loading.value = true
    try {
      const cmds = await CommandsService.ScanApplications()
      allCommands.value = cmds.map((c: any) => ({
        ...c,
        id: c.name,
        type: 'direct',
        subType: 'app',
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
    } catch (err) {
      console.error('launch error:', err)
    }
  }

  function setQuery(q: string) {
    searchQuery.value = q
  }

  return {
    allCommands, searchQuery, loading, selectedIndex,
    results, loadCommands, launch, setQuery,
  }
})
