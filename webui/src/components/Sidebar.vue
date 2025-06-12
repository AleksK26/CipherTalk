<template>
  <nav class="sidebar d-md-block">
    <div class="sidebar-sticky pt-3">
      <div class="px-3 mb-3">
        <strong v-if="userName">Hello, {{ userName }}!</strong>
        <strong v-else>Hello, user!</strong>
      </div>
      <ul class="nav flex-column">
        <li class="nav-item">
          <RouterLink to="/home" class="nav-link" active-class="active">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
            Conversations
          </RouterLink>
        </li>
        <li class="nav-item">
          <RouterLink to="/me" class="nav-link" active-class="active">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
            Profile
          </RouterLink>
        </li>
        <li class="nav-item">
          <RouterLink to="/search" class="nav-link" active-class="active">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
            Search People
          </RouterLink>
        </li>
      </ul>
      <div class="px-3 mt-4">
        <button class="btn btn-outline-danger btn-sm w-100" @click="logOut">Log Out</button>
      </div>
    </div>
  </nav>
</template>

<script>
import { RouterLink } from 'vue-router';
export default {
  name: 'Sidebar',
  components: { RouterLink },
  data() {
    return {
      userName: localStorage.getItem('name') || '',
    };
  },
  methods: {
    logOut() {
      localStorage.clear();
      this.$router.push({ path: '/' });
    },
  },
};
</script>

<style>
.sidebar {
  min-height: 100vh;
  border-right: 1px solid #eee;
  background: #2c8b47; 
  padding-top: 1rem;
}
.nav-link.active {
  font-weight: bold;
  color: #fff !important;
  background: #3da75d;
}
.nav-link {
  transition: background 0.2s, color 0.2s;
  color: #e0ffe0;
}
.nav-link:hover {
  background: #4fcb74;
  color: #fff;
}
.btn-outline-danger {
  margin-top: 1rem;
}
.sidebar-sticky strong,
.sidebar-sticky {
  color: #e0ffe0;
}
@media (max-width: 767.98px) {
  .sidebar {
    min-height: auto;
    padding-top: 0.5rem;
  }
}
</style>
