<script setup>
import Sidebar from './components/Sidebar.vue'
import { RouterView, useRoute } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()
const showSidebar = computed(() => route.path !== '/')
</script>
<script>
export default {}
</script>

<template>
  <div class="app-shell" :class="{ 'no-sidebar': !showSidebar }">
    <Sidebar v-if="showSidebar" class="app-sidebar" />
    <main class="app-main">
      <RouterView />
    </main>
  </div>
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
}

.app-shell {
  display: flex;
  min-height: 100vh;
}

.app-shell.no-sidebar {
  display: block;
}

.app-sidebar {
  width: 240px;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
}

.app-main {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
}

/* Mobile: sidebar becomes a bottom tab bar (handled by Sidebar.vue itself) */
@media (max-width: 767px) {
  .app-shell {
    flex-direction: column;
  }
  .app-sidebar {
    display: none; /* desktop sidebar hidden on mobile; Sidebar.vue's mobile-tabs uses its own fixed positioning */
  }
  .app-main {
    padding-bottom: 68px; /* space for Sidebar.vue's fixed mobile tab bar */
  }
}
</style>
