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
  -webkit-text-size-adjust: 100%;
}

#app {
  height: 100%;
}

.app-shell {
  display: flex;
  min-height: 100vh;
  min-height: 100dvh;
}

.app-shell.no-sidebar {
  display: block;
}

.app-sidebar {
  /* display:contents wrapper — sizing handled by nav elements inside Sidebar.vue */
}

.app-main {
  flex: 1;
  min-width: 0;
  min-height: 0;
  overflow-y: auto;
}

@media (max-width: 767px) {
  .app-shell {
    flex-direction: column;
    height: 100vh;
    height: 100dvh;
  }
  .app-sidebar {
    display: contents;
  }
  .app-main {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding-bottom: calc(60px + env(safe-area-inset-bottom, 0px) + 8px);
  }
}
</style>
