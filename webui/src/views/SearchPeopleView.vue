<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">Find People</h1>
      <p class="page-sub">Search for users and start a conversation</p>
    </div>

    <div class="search-wrap">
      <form @submit.prevent="searchUsers" class="search-row">
        <input
          v-model="query"
          class="search-input"
          type="text"
          placeholder="Search by username…"
        />
        <button class="search-btn" type="submit" :disabled="loading">
          {{ loading ? 'Searching…' : 'Search' }}
        </button>
      </form>

      <div v-if="error" class="err-box">{{ error }}</div>

      <LoadingSpinner v-if="loading" />

      <div v-if="!loading && showResults">
        <p class="results-label">{{ users.length }} result{{ users.length !== 1 ? 's' : '' }} for "{{ lastQuery }}"</p>
        <div v-if="users.length === 0" class="empty-state">No users found.</div>
        <div v-for="user in users" :key="user.id" class="user-card">
          <div class="user-avatar">{{ getInitials(user.name) }}</div>
          <div class="user-info">
            <span class="user-name">{{ user.name }}</span>
          </div>
          <button class="chat-btn" @click="navigateToConversation(user.id, user.name)">
            Message
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
  name: "SearchPeopleView",
  components: { LoadingSpinner },
  data() {
    return {
      query: "",
      lastQuery: "",
      users: [],
      loading: false,
      showResults: false,
      error: "",
    };
  },
  methods: {
    async searchUsers() {
      if (!this.query.trim()) {
        this.error = "Please enter a search term.";
        return;
      }
      this.loading = true;
      this.error = "";
      this.users = [];
      this.showResults = false;
      try {
        const token = localStorage.getItem("token");
        const response = await axios.get(`/search`, {
          params: { username: this.query },
          headers: { Authorization: `Bearer ${token}` },
        });
        this.users = response.data;
        this.lastQuery = this.query;
        this.showResults = true;
      } catch (err) {
        this.error = `Error: ${err.response?.data || "Failed to search users."}`;
      } finally {
        this.loading = false;
      }
    },
    async navigateToConversation(recipientId, recipientName) {
      localStorage.setItem("conversationName", recipientName);
      const token = localStorage.getItem("token");
      try {
        const response = await axios.post(
          `/conversations`,
          { senderId: token, recipientId },
          { headers: { Authorization: `Bearer ${token}` } }
        );
        this.$router.push({ path: `/conversations/${response.data.conversationId}` });
      } catch (error) {
        console.error("Error starting conversation:", error);
        this.error = "Could not start conversation. Please try again.";
      }
    },
    getInitials(name) {
      if (!name) return "?";
      return name.slice(0, 2).toUpperCase();
    },
  },
  mounted() {
    if (!localStorage.getItem("token")) this.$router.push("/");
  },
};
</script>

<style scoped>
.page {
  max-width: 680px;
  margin: 0 auto;
  padding: 40px 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
.page-header { margin-bottom: 32px; }
.page-title { font-size: 28px; font-weight: 700; color: #1a1f36; margin: 0 0 6px; }
.page-sub { font-size: 15px; color: #718096; margin: 0; }

.search-wrap { background: #fff; border-radius: 16px; padding: 28px; box-shadow: 0 2px 16px rgba(0,0,0,0.06); }

.search-row { display: flex; gap: 10px; margin-bottom: 20px; }
.search-input {
  flex: 1;
  padding: 12px 16px;
  border: 1.5px solid #e2e8f0;
  border-radius: 10px;
  font-size: 15px;
  outline: none;
  transition: border-color 0.2s;
}
.search-input:focus { border-color: #4c6ef5; }
.search-btn {
  padding: 12px 24px;
  background: #4c6ef5;
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
  white-space: nowrap;
}
.search-btn:hover:not(:disabled) { background: #3b5bdb; }
.search-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.err-box {
  background: #fff5f5; border: 1px solid #fed7d7; color: #c53030;
  border-radius: 8px; padding: 10px 14px; font-size: 13px; margin-bottom: 16px;
}

.results-label { font-size: 13px; color: #718096; margin-bottom: 12px; }
.empty-state { text-align: center; color: #a0aec0; padding: 32px 0; font-size: 15px; }

.user-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 0;
  border-bottom: 1px solid #f0f2f5;
}
.user-card:last-child { border-bottom: none; }
.user-avatar {
  width: 44px; height: 44px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4c6ef5, #7c3aed);
  color: #fff;
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 16px;
  flex-shrink: 0;
}
.user-info { flex: 1; }
.user-name { font-size: 15px; font-weight: 600; color: #1a1f36; }

.chat-btn {
  padding: 8px 18px;
  background: #e8f0fe;
  color: #4c6ef5;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.chat-btn:hover { background: #c5d5fc; }
</style>
