<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  fetchDashboardStats,
  formatDuration,
  formatRate,
  type DashboardStats,
} from '@/lib/auth'

const stats = ref<DashboardStats | null>(null)
const loading = ref(true)
const errorMessage = ref<string | null>(null)

onMounted(async () => {
  loading.value = true
  errorMessage.value = null
  try {
    stats.value = await fetchDashboardStats()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to load dashboard'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="mx-auto max-w-[1400px] space-y-6 px-4 py-8">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">Dashboard</h1>
      <p class="text-sm text-muted-foreground">
        Overview of agent log sessions and turns.
      </p>
    </div>

    <p v-if="loading" class="text-sm text-muted-foreground">Loading…</p>
    <p v-else-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>

    <div v-else-if="stats" class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <Card>
        <CardHeader class="pb-2">
          <CardDescription>Session count</CardDescription>
          <CardTitle class="text-3xl tabular-nums">{{ stats.session_count }}</CardTitle>
        </CardHeader>
        <CardContent>
          <p class="text-xs text-muted-foreground">Total conversation sessions</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="pb-2">
          <CardDescription>Avg rate</CardDescription>
          <CardTitle class="text-3xl tabular-nums">{{ formatRate(stats.avg_rate) }}</CardTitle>
        </CardHeader>
        <CardContent>
          <p class="text-xs text-muted-foreground">Average prompt rate (0–10)</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="pb-2">
          <CardDescription>Total duration</CardDescription>
          <CardTitle class="text-3xl tabular-nums">
            {{ formatDuration(stats.total_duration_ms) }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p class="text-xs text-muted-foreground">Sum of all turn durations</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="pb-2">
          <CardDescription>Token total</CardDescription>
          <CardTitle class="text-3xl tabular-nums">{{ stats.token_total }}</CardTitle>
        </CardHeader>
        <CardContent>
          <p class="text-xs text-muted-foreground">Parsed tokens across all turns</p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
