<template>
  <div class="sidebar-root">
  <!-- Desktop: vertical sidebar -->
  <nav class="sidebar desktop-sidebar">
    <div class="brand">
      <span class="brand-icon">🔐</span>
      <span class="brand-name">CipherTalk</span>
    </div>

    <div class="user-block">
      <div class="user-avatar">{{ initials }}</div>
      <div class="user-info">
        <span class="user-name">{{ userName }}</span>
        <span class="user-status">Online</span>
      </div>
    </div>

    <ul class="nav-list">
      <li>
        <RouterLink to="/home" class="nav-item" active-class="nav-active">
          <span class="nav-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
          </span>
          <span class="nav-label">Conversations</span>
        </RouterLink>
      </li>
      <li>
        <RouterLink to="/search" class="nav-item" active-class="nav-active">
          <span class="nav-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          </span>
          <span class="nav-label">Search People</span>
        </RouterLink>
      </li>
      <li>
        <RouterLink to="/me" class="nav-item" active-class="nav-active">
          <span class="nav-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
          </span>
          <span class="nav-label">Profile</span>
        </RouterLink>
      </li>
    </ul>

    <div class="sidebar-footer">
      <button class="logout-btn" @click="logOut">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="width:16px;height:16px"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
        Log Out
      </button>
    </div>
  </nav>

  <!-- Mobile: bottom tab bar -->
  <nav class="sidebar mobile-tabs">
    <RouterLink to="/home" class="tab-item" active-class="tab-active">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
      <span>Chats</span>
    </RouterLink>
    <RouterLink to="/search" class="tab-item" active-class="tab-active">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
      <span>Search</span>
    </RouterLink>
    <RouterLink to="/me" class="tab-item" active-class="tab-active">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
      <span>Profile</span>
    </RouterLink>
    <button class="tab-item tab-logout" @click="logOut">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
      <span>Logout</span>
    </button>
  </nav>
  </div>
</template>

<script>
import { RouterLink } from 'vue-router';
export default {
  name: 'Sidebar',
  components: { RouterLink },
  data() {
    return {
      userName: localStorage.getItem('name') || 'User',
    };
  },
  computed: {
    initials() {
      return (this.userName || 'U').slice(0, 2).toUpperCase();
    },
  },
  methods: {
    logOut() {
      localStorage.clear();
      this.$router.push('/');
    },
  },
};
</script>

<style scoped>
/* ─── Root wrapper ───────────────────────────────────────────────── */
.sidebar-root {
  display: contents;
}

/* ─── Desktop Sidebar ────────────────────────────────────────────── */
.desktop-sidebar {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: linear-gradient(180deg, #1a1f36 0%, #2d3461 100%);
  padding: 0;
  overflow: hidden;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 22px 20px 18px;
  border-bottom: 1px solid rgba(255,255,255,0.08);
}
.brand-icon { font-size: 22px; }
.brand-name {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.3px;
}

.user-block {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}
.user-avatar {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4c6ef5, #7c3aed);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  flex-shrink: 0;
}
.user-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.user-name {
  font-size: 13px;
  font-weight: 600;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.user-status {
  font-size: 11px;
  color: #68d391;
}

.nav-list {
  list-style: none;
  padding: 12px 10px;
  margin: 0;
  flex: 1;
}
.nav-list li { margin-bottom: 2px; }

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 12px;
  border-radius: 10px;
  color: #a0aec0;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: background 0.15s, color 0.15s;
}
.nav-item:hover {
  background: rgba(255,255,255,0.08);
  color: #e2e8f0;
}
.nav-active {
  background: rgba(76,110,245,0.25) !important;
  color: #7b9ff9 !important;
}

.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.nav-icon svg {
  width: 18px;
  height: 18px;
}

.sidebar-footer {
  padding: 12px 10px 16px;
  border-top: 1px solid rgba(255,255,255,0.06);
}
.logout-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  background: rgba(229,62,62,0.12);
  color: #fc8181;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s;
}
.logout-btn:hover {
  background: rgba(229,62,62,0.22);
  color: #feb2b2;
}

/* ─── Mobile Bottom Tabs ─────────────────────────────────────────── */
.mobile-tabs {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: #1a1f36;
  border-top: 1px solid rgba(255,255,255,0.1);
  flex-direction: row;
  align-items: stretch;
  z-index: 100;
}

.tab-item {
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
}
.tab-item svg { width: 20px; height: 20px; }
.tab-item:hover { color: #a0aec0; }
.tab-active { color: #7b9ff9 !important; }
.tab-logout { color: #fc8181; }
.tab-logout:hover { color: #feb2b2; }

/* ─── Responsive ─────────────────────────────────────────────────── */
@media (max-width: 767px) {
  .desktop-sidebar { display: none; }
  .mobile-tabs { display: flex; }
}
@media (min-width: 768px) {
  .mobile-tabs { display: none; }
}
</style>
