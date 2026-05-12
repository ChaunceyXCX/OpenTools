import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import Fuse from 'fuse.js'
import { CommandsService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'
import type { Command } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/core/scanner/models'

function buildFuseIndex(items: Command[]) {
  return new Fuse<Command>(items, {
    keys: [
      { name: 'name', weight: 2 },
      { name: 'aliases', weight: 1.5 },
      { name: 'acronym', weight: 1 },
    ],
    threshold: 0.4,
    includeScore: true,
  })
}

export const useCommandStore = defineStore('commands', () => {
  const commands = ref<Command[]>([])
  const query = ref('')
  const loading = ref(false)

  let fuse = buildFuseIndex([])

  const results = computed(() => {
    if (!query.value.trim()) return commands.value.slice(0, 8)
    return fuse.search(query.value).map(r => r.item).slice(0, 12)
  })

  async function loadCommands() {
    loading.value = true
    try {
      const cmds = await CommandsService.ScanApplications()
      commands.value = cmds
      fuse = buildFuseIndex(cmds)
    } catch (err) {
      console.error('scan error:', err)
    } finally {
      loading.value = false
    }
  }

  async function launch(cmd: Command) {
    try {
      const result = await CommandsService.Launch(cmd.path)
      if (!result.success) {
        console.error('launch error:', result.error)
      }
    } catch (err) {
      console.error('launch error:', err)
    }
  }

  function setQuery(q: string) {
    query.value = q
  }

  return { commands, query, loading, results, loadCommands, launch, setQuery }
})
