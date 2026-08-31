<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import GridPagination from '@/components/GridPagination.vue'
import {
  fetchSession,
  formatDateTime,
  formatDuration,
  formatRate,
  type SessionDetail,
  type SessionPromptRow,
} from '@/lib/auth'
import { usePagination } from '@/lib/usePagination'

const route = useRoute()
const sessionId = computed(() => String(route.params.id ?? ''))

const detail = ref<SessionDetail | null>(null)
const loading = ref(true)
const errorMessage = ref<string | null>(null)

const sortedPrompts = computed((): SessionPromptRow[] => {
  if (!detail.value) return []
  return [...detail.value.prompts].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  )
})

const promptsRef = computed(() => sortedPrompts.value)
const { currentPage, totalPages, paginatedItems, pageSize, goToPage } = usePagination(promptsRef, 100)

onMounted(async () => {
  loading.value = true
  errorMessage.value = null
  try {
    detail.value = await fetchSession(sessionId.value)
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to load session'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="w-full space-y-4 p-4">
    <div class="flex flex-wrap items-center gap-3">
      <RouterLink
        to="/sessions"
        class="text-sm text-muted-foreground transition-colors hover:text-foreground"
      >
        ← Back to sessions
      </RouterLink>
    </div>

    <p v-if="loading" class="text-sm text-muted-foreground">Loading…</p>
    <p v-else-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>

    <template v-else-if="detail">
      <div class="space-y-1">
        <h1 class="text-xl font-semibold tracking-tight">{{ detail.title || 'Session' }}</h1>
        <p class="text-sm text-muted-foreground">
          {{ detail.app || '—' }} · {{ detail.project || '—' }} ·
          {{ detail.prompt_count }} prompt{{ detail.prompt_count === 1 ? '' : 's' }} ·
          {{ formatDateTime(detail.created_at) }}
        </p>
      </div>

      <div class="overflow-x-auto rounded-lg border border-border">
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
            <tr v-if="sortedPrompts.length === 0">
              <td colspan="8" class="px-3 py-6 text-center text-muted-foreground">
                No prompts in this session.
              </td>
            </tr>
            <tr
              v-for="turn in paginatedItems"
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
          :current-page="currentPage"
          :total-pages="totalPages"
          :total-items="sortedPrompts.length"
          :page-size="pageSize"
          @go-to-page="goToPage"
        />
      </div>
    </template>
  </div>
</template>
