<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import GridPagination from '@/components/GridPagination.vue'
import {
  fetchSessions,
  formatDateTime,
  formatDuration,
  formatRate,
  type SessionRow,
} from '@/lib/auth'
import { usePagination } from '@/lib/usePagination'

const router = useRouter()
const rows = ref<SessionRow[]>([])
const loading = ref(true)
const errorMessage = ref<string | null>(null)

const sortedRows = computed(() =>
  [...rows.value].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  ),
)

const { currentPage, totalPages, paginatedItems, pageSize, goToPage } = usePagination(sortedRows, 100)

onMounted(async () => {
  loading.value = true
  errorMessage.value = null
  try {
    rows.value = await fetchSessions()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to load sessions'
  } finally {
    loading.value = false
  }
})

function openSession(id: string) {
  void router.push({ name: 'session-detail', params: { id } })
}
</script>

<template>
  <div class="w-full space-y-4 p-4">
    <p v-if="loading" class="text-sm text-muted-foreground">Loading…</p>
    <p v-else-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>

    <div v-else class="overflow-x-auto rounded-lg border border-border">
      <table class="w-full min-w-[1100px] border-collapse text-left text-sm">
        <thead class="bg-muted/50">
          <tr>
            <th class="border-b border-border px-3 py-2.5 font-medium">App</th>
            <th class="border-b border-border px-3 py-2.5 font-medium">Project</th>
            <th class="border-b border-border px-3 py-2.5 font-medium">Title</th>
            <th class="border-b border-border px-3 py-2.5 font-medium">Prompt</th>
            <th class="w-16 border-b border-border px-3 py-2.5 font-medium">Rate</th>
            <th class="border-b border-border px-3 py-2.5 font-medium">Duration</th>
            <th class="border-b border-border px-3 py-2.5 font-medium">Tokens</th>
            <th class="border-b border-border px-3 py-2.5 font-medium">Device</th>
            <th class="border-b border-border px-3 py-2.5 font-medium">Created at</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="sortedRows.length === 0">
            <td colspan="9" class="px-3 py-6 text-center text-muted-foreground">
              No sessions yet.
            </td>
          </tr>
          <tr
            v-for="row in paginatedItems"
            :key="row.id"
            class="cursor-pointer align-top odd:bg-background even:bg-muted/20 hover:bg-accent/40"
            :title="row.id"
            @click="openSession(row.id)"
          >
            <td class="border-b border-border px-3 py-3 whitespace-nowrap">{{ row.app || '—' }}</td>
            <td class="border-b border-border px-3 py-3">{{ row.project || '—' }}</td>
            <td class="border-b border-border px-3 py-3">{{ row.title }}</td>
            <td class="border-b border-border px-3 py-3 tabular-nums">{{ row.prompt_count }}</td>
            <td class="border-b border-border px-3 py-3 font-medium tabular-nums">
              {{ formatRate(row.rate) }}
            </td>
            <td class="border-b border-border px-3 py-3 whitespace-nowrap tabular-nums">
              {{ formatDuration(row.duration_ms) }}
            </td>
            <td class="border-b border-border px-3 py-3 tabular-nums">
              {{ row.token_total ?? 0 }}
            </td>
            <td class="border-b border-border px-3 py-3 whitespace-nowrap">{{ row.device || '—' }}</td>
            <td class="border-b border-border px-3 py-3 whitespace-nowrap tabular-nums">
              {{ formatDateTime(row.created_at) }}
            </td>
          </tr>
        </tbody>
      </table>
      <GridPagination
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-items="sortedRows.length"
        :page-size="pageSize"
        @go-to-page="goToPage"
      />
    </div>
  </div>
</template>
