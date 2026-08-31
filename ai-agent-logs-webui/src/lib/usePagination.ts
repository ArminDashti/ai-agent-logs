import { computed, ref, watch, type Ref } from 'vue'

export function usePagination<T>(items: Ref<T[]>, pageSize = 100) {
  const currentPage = ref(1)

  const totalPages = computed(() => Math.max(1, Math.ceil(items.value.length / pageSize)))

  const paginatedItems = computed(() => {
    const start = (currentPage.value - 1) * pageSize
    return items.value.slice(start, start + pageSize)
  })

  watch(
    () => items.value.length,
    () => {
      if (currentPage.value > totalPages.value) {
        currentPage.value = totalPages.value
      }
    },
  )

  function goToPage(page: number) {
    currentPage.value = Math.min(Math.max(1, page), totalPages.value)
  }

  return { currentPage, totalPages, paginatedItems, pageSize, goToPage }
}
