<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

const localPart = ['arminonline', '71']
const domainParts = ['gmail', 'com']

const mailtoHref = computed(() => {
  const user = localPart.join('')
  const domain = `${domainParts[0]}.${domainParts[1]}`
  return `mailto:${user}@${domain}`
})

const visibleEmail = computed(
  () => `${localPart.join('')} [at] ${domainParts[0]} [dot] ${domainParts[1]}`,
)

const hobbies = [
  'cars',
  'movies',
  'vibe coding',
  'coding',
  'politics',
  'military',
  'jet fighters',
  'aircraft',
]

const photoSrc = ref('/about-me/armin.jpg')
const photoOk = ref(false)

onMounted(() => {
  const img = new Image()
  img.onload = () => {
    photoOk.value = true
  }
  img.onerror = () => {
    photoOk.value = false
  }
  img.src = photoSrc.value
})
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-6 px-4 py-10">
    <Card>
      <CardHeader>
        <p class="mb-1 text-sm font-medium uppercase tracking-[0.2em] text-primary">About Me</p>
        <CardTitle class="text-3xl tracking-tight">Armin Dashti</CardTitle>
        <CardDescription class="text-base">
          Software engineer and vibe coder — conductor of craft, not a solitary typist.
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-5 text-sm leading-relaxed text-muted-foreground">
        <div
          v-if="photoOk"
          class="overflow-hidden rounded-lg border border-border"
        >
          <img :src="photoSrc" alt="Armin Dashti" class="h-auto w-full object-cover" />
        </div>
        <div
          v-else
          class="flex h-40 items-center justify-center rounded-lg border border-dashed border-muted-foreground/40 bg-muted/30 text-xs text-muted-foreground"
        >
          Photo placeholder — add <code class="mx-1">public/about-me/armin.jpg</code>
        </div>

        <p>
          Armin Dashti is a software engineer and vibe coder. A vibe coder is not a spectator of
          the craft, but its conductor: he shapes intent with clarity, sets the boundaries of taste
          and constraint, then lets Cursor’s agents compose the implementation.
        </p>

        <p>
          Contact:
          <a class="text-primary underline-offset-4 hover:underline" :href="mailtoHref">
            {{ visibleEmail }}
          </a>
        </p>

        <ul class="flex flex-wrap gap-2 pt-1">
          <li
            v-for="hobby in hobbies"
            :key="hobby"
            class="rounded-md border border-border bg-muted/40 px-2.5 py-1 text-xs text-foreground"
          >
            {{ hobby }}
          </li>
        </ul>

        <div class="flex flex-wrap gap-3 pt-1">
          <a href="https://github.com/ArminDashti" target="_blank" rel="noopener noreferrer">
            <Button variant="outline">GitHub</Button>
          </a>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
