import { computed, ref } from 'vue'
import { clearSession, getStoredUser, getToken, setSession, type User } from './auth'

const user = ref<User | null>(getStoredUser())

export function useAuth() {
  const isAuthenticated = computed(() => Boolean(getToken()) && Boolean(user.value))

  function applySession(token: string, next: User) {
    setSession(token, next)
    user.value = next
  }

  function logout() {
    clearSession()
    user.value = null
  }

  async function hydrate() {
    user.value = getStoredUser()
  }

  return {
    user,
    isAuthenticated,
    applySession,
    logout,
    hydrate,
  }
}
