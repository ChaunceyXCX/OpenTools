import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import Fuse from 'fuse.js'

export interface Command {
  id: string
  name: string
  type: string
  subType?: string
  icon?: string
}

export const useCommandStore = defineStore('commands', () => {
  const commands = ref<Command[]>([])
  const query = ref('')

  const fuse = new Fuse<Command>([], {
    keys: ['name'],
    threshold: 0.4,
    includeScore: true,
  })

  const results = computed(() => {
    if (!query.value.trim()) return commands.value.slice(0, 8)
    return fuse.search(query.value).map(r => r.item).slice(0, 12)
  })

  function setCommands(cmdList: Command[]) {
    commands.value = cmdList
    fuse.setCollection(cmdList)
  }

  function setQuery(q: string) {
    query.value = q
  }

  return { commands, query, results, setCommands, setQuery }
})
