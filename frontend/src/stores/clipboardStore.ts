import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ClipboardService } from '../../bindings/github.com/ChaunceyXCX/OpenTools/internal/api/index'

export const useClipboardStore = defineStore('clipboard', () => {
  const items = ref<any[]>([])
  const loading = ref(false)

  async function loadHistory(page = 1, pageSize = 50) {
    loading.value = true
    try {
      items.value = await ClipboardService.GetHistory(page, pageSize)
    } catch (err) {
      console.error('clipboard history error:', err)
    } finally {
      loading.value = false
    }
  }

  async function deleteItem(id: string) {
    await ClipboardService.Delete(id)
    await loadHistory()
  }

  async function clearAll() {
    await ClipboardService.Clear()
    items.value = []
  }

  return { items, loading, loadHistory, deleteItem, clearAll }
})
