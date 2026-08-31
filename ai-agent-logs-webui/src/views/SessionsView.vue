<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import GridPagination from '@/components/GridPagination.vue'
import {
  fetchSession,
  fetchSessions,
  formatDateTime,
  formatDuration,
  formatRate,
  type SessionDetail,
  type SessionPromptRow,
  type SessionRow,
} from '@/lib/auth'
import { usePagination } from '@/lib/usePagination'

const rows = ref<SessionRow[]>([])
const loading = ref(true)
const errorMessage = ref<string | null>(null)

const selectedId = ref<string | null>(null)
const detail = ref<SessionDetail | null>(null)
const detailLoading = ref(false)
const detailError = ref<string | null>(null)

const sortedRows = computed(() =>
  [...rows.value].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  ),
)

const { currentPage, totalPages, paginatedItems, pageSize, goToPage } = usePagination(sortedRows, 100)

const sortedDetailPrompts = computed((): SessionPromptRow[] => {
  if (!detail.value) return []
  return [...detail.value.prompts].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  )
})

const detailPromptsRef = computed(() => sortedDetailPrompts.value)
const {
  currentPage: detailPage,
  totalPages: detailTotalPages,
  paginatedItems: paginatedDetailPrompts,
  pageSize: detailPageSize,
  goToPage: goToDetailPage,
} = usePagination(detailPromptsRef, 100)

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

watch(selectedId, async (id) => {
  detail.value = null
  detailError.value = null
  detailPage.value = 1
  if (!id) return

  detailLoading.value = true
  try {
    detail.value = await fetchSession(id)
  } catch (err) {
    detailError.value = err instanceof Error ? err.message : 'Failed to load session detail'
  } finally {
    detailLoading.value = false
  }
})

function selectSession(id: string) {
  selectedId.value = selectedId.value === id ? null : id
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
            :class="{ 'bg-accent/50 even:bg-accent/50': selectedId === row.id }"
            :title="row.id"
            @click="selectSession(row.id)"
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

    <section v-if="selectedId" class="space-y-3" aria-label="Session detail">
      <p v-if="detailLoading" class="text-sm text-muted-foreground">Loading detail…</p>
      <p v-else-if="detailError" class="text-sm text-red-600 dark:text-red-400">{{ detailError }}</p>

      <div v-else-if="detail" class="overflow-x-auto rounded-lg border border-border">
        <table class="w-full min-w-[1200px] border-collapse text-left text-sm">
          <thead class="bg-muted/50">
            <tr>
              <th class="border-b border-border px-3 py-2.5 font-medium">Mode</th>
              <th class="border-b border-border px-3 py-2.5 font-medium">Prompt</th>
              <th class="border-b border-border px-3 py-2.5 font-medium">Response</th>
              <th class="w-16 border-b border-border px-3 py-2.5 font-medium">Rate</th>
              <th class="border-b border-border px-3 py-2.5 font-medium">Duration</th>
              <th class="border-b border-border px-3 py-2.5 font-medium">Tokens</th>
              <th class="border-b border-border px-3 py-2.5 font-medium">IP</th>
              <th class="border-b border-border px-3 py-2.5 font-medium">Created at</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="sortedDetailPrompts.length === 0">
              <td colspan="8" class="px-3 py-6 text-center text-muted-foreground">
                No prompts in this session.
              </td>
            </tr>
            <tr
              v-for="turn in paginatedDetailPrompts"
              :key="turn.id"
              class="align-top odd:bg-background even:bg-muted/20"
            >
              <td class="border-b border-border px-3 py-3 whitespace-nowrap">{{ turn.mode }}</td>
              <td class="border-b border-border px-3 py-3 max-w-sm whitespace-pre-wrap">
                {{ turn.prompt }}
              </td>
              <td class="border-b border-border px-3 py-3 max-w-md whitespace-pre-wrap">
                {{ turn.response }}
              </td>
              <td class="border-b border-border px-3 py-3 font-medium tabular-nums">
                {{ formatRate(turn.rate) }}
              </td>
              <td class="border-b border-border px-3 py-3 whitespace-nowrap tabular-nums">
                {{ formatDuration(turn.duration_ms) }}
              </td>
              <td class="border-b border-border px-3 py-3 tabular-nums">
                {{ turn.token_usage || '—' }}
              </td>
              <td class="border-b border-border px-3 py-3 whitespace-nowrap font-mono text-xs">
                {{ turn.user_ip || '—' }}
              </td>
              <td class="border-b border-border px-3 py-3 whitespace-nowrap tabular-nums">
                {{ formatDateTime(turn.created_at) }}
              </td>
            </tr>
          </tbody>
        </table>
        <GridPagination
          :current-page="detailPage"
          :total-pages="detailTotalPages"
          :total-items="sortedDetailPrompts.length"
          :page-size="detailPageSize"
          @go-to-page="goToDetailPage"
        />
      </div>
    </section>
  </div>
</template>
