<script setup>
import Sidebar from './components/Sidebar.vue'
import { RouterView, RouterLink, useRoute, useRouter } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()
const router = useRouter()
const showNav = computed(() => route.path !== '/')

function logOut() {
  localStorage.clear()
  router.push('/')
}
</script>

<template>
  <div class="app-shell" :class="{ 'no-sidebar': !showNav }">
    <Sidebar v-if="showNav" />
    <main class="app-main">
      <RouterView />
    </main>
  </div>

  <!-- Mobile bottom tab bar: OUTSIDE app-shell, direct child of #app -->
  <nav v-if="showNav" class="mobile-tab-bar">
    <RouterLink to="/home" class="mtab" active-class="mtab-active">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
      <span>Chats</span>
    </RouterLink>
    <RouterLink to="/search" class="mtab" active-class="mtab-active">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
      <span>Search</span>
    </RouterLink>
    <RouterLink to="/me" class="mtab" active-class="mtab-active">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
      <span>Profile</span>
    </RouterLink>
    <button class="mtab mtab-logout" @click="logOut">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
      <span>Logout</span>
    </button>
  </nav>
</template>

<style>
*, *::before, *::after { box-sizing: border-box; }

html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f7f8fc;
  color: #1a1f36;
  -webkit-text-size-adjust: 100%;
}

#app {
  height: 100%;
}

/* ─── Desktop layout ─────────────────────────────── */
.app-shell {
  display: flex;
  min-height: 100vh;
}

.app-shell.no-sidebar {
  display: block;
}

.app-main {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
}

/* ─── Mobile tab bar — hidden on desktop ─────────── */
.mobile-tab-bar {
  display: none;
}

@media (max-width: 767px) {
  .app-shell {
    display: block;
  }

  .app-main {
    padding-bottom: 60px;
  }

  .mobile-tab-bar {
    display: flex;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 60px;
    background: #1a1f36;
    border-top: 1px solid rgba(255,255,255,0.1);
    flex-direction: row;
    align-items: stretch;
    z-index: 9999;
  }

  .mtab {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 3px;
    color: #718096;
    text-decoration: none;
    font-size: 10px;
    font-weight: 500;
    border: none;
    background: none;
    cursor: pointer;
    transition: color 0.15s;
    padding: 0;
  }

  .mtab svg {
    width: 20px;
    height: 20px;
  }

  .mtab:hover { color: #a0aec0; }
  .mtab-active { color: #7b9ff9 !important; }
  .mtab-logout { color: #fc8181; }
  .mtab-logout:hover { color: #feb2b2; }
}
</style>
