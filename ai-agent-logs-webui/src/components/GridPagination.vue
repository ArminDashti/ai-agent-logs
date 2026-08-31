<script setup lang="ts">
import { Button } from '@/components/ui/button'

defineProps<{
  currentPage: number
  totalPages: number
  totalItems: number
  pageSize: number
}>()

const emit = defineEmits<{
  goToPage: [page: number]
}>()
</script>

<template>
  <div
    v-if="totalItems > pageSize"
    class="flex flex-wrap items-center justify-between gap-3 border-t border-border px-3 py-2.5 text-sm text-muted-foreground"
  >
    <span class="tabular-nums">
      {{ (currentPage - 1) * pageSize + 1 }}–{{ Math.min(currentPage * pageSize, totalItems) }}
      of {{ totalItems }}
    </span>
    <div class="flex items-center gap-2">
      <Button
        variant="outline"
        size="sm"
        :disabled="currentPage <= 1"
        @click="emit('goToPage', currentPage - 1)"
      >
        Previous
      </Button>
      <span class="tabular-nums">Page {{ currentPage }} / {{ totalPages }}</span>
      <Button
        variant="outline"
        size="sm"
        :disabled="currentPage >= totalPages"
        @click="emit('goToPage', currentPage + 1)"
      >
        Next
      </Button>
    </div>
  </div>
</template>
