<template>
  <div class="home">
    <div class="home-header">
      <div class="header-left">
        <h1 class="header-title">Messages</h1>
        <p class="header-sub">{{ currentUsername }}</p>
      </div>
      <div class="header-actions">
        <button class="action-btn new-group-btn" @click="newGroup" title="New Group">
          + New Group
        </button>
        <button class="action-btn icon-btn" @click="refresh" title="Refresh">↻</button>
        <button class="action-btn logout-btn" @click="logOut" title="Log Out">Log Out</button>
      </div>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />

    <div v-if="conversations.length === 0 && !loading" class="empty-state">
      <div class="empty-icon">💬</div>
      <p class="empty-title">No conversations yet</p>
      <p class="empty-sub">Search for people to start chatting</p>
      <button class="find-btn" @click="$router.push('/search')">Find People</button>
    </div>

    <div v-else class="conv-list">
      <div
        v-for="conv in conversations"
        :key="conv.id"
        class="conv-item"
        @click="viewConversation(conv.id, conv.name)"
      >
        <div class="conv-photo">
          <img
            v-if="conv.conversationPhoto && conv.conversationPhoto.String"
            :src="'data:image/jpeg;base64,' + conv.conversationPhoto.String"
            alt="Photo"
            class="conv-img"
          />
          <div v-else class="conv-initials">{{ getInitials(conv.name) }}</div>
        </div>
        <div class="conv-body">
          <div class="conv-top">
            <span class="conv-name">{{ conv.name }}</span>
            <span v-if="conv.lastMessage" class="conv-time">
              {{ formatTime(conv.lastMessage.timestamp) }}
            </span>
          </div>
          <div class="conv-preview">
            <span v-if="conv.lastMessage">
              <strong v-if="conv.lastMessage.senderName">{{ conv.lastMessage.senderName }}: </strong>
              <span v-if="conv.lastMessage.attachment" class="attach-hint">📷 </span>
              <span>{{ truncate(conv.lastMessage.content) }}</span>
            </span>
            <span v-else class="no-msg">No messages yet</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "HomeView",
  components: { ErrorMsg },
  data() {
    localStorage.removeItem("recipientId");
    return {
      currentUsername: localStorage.getItem("name") || "",
      errormsg: null,
      loading: false,
      conversations: [],
      pollIntervalId: null,
    };
  },
  methods: {
    async loadConversations() {
      this.loading = true;
      this.errormsg = null;
      try {
        const token = localStorage.getItem("token");
        if (!token) { this.$router.push("/"); return; }
        const response = await this.$axios.get("/conversations", {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.conversations = response.data || [];
      } catch (error) {
        this.errormsg = "Failed to load conversations.";
      } finally {
        this.loading = false;
      }
    },
    viewConversation(id, name) {
      localStorage.setItem("conversationName", name);
      this.$router.push({ path: `/conversations/${id}` });
    },
    refresh() { this.loadConversations(); },
    logOut() { this.$router.push("/"); },
    newGroup() { this.$router.push("/new-group"); },
    truncate(text, max = 45) {
      if (!text) return "";
      return text.length > max ? text.slice(0, max) + "…" : text;
    },
    getInitials(name) {
      if (!name) return "?";
      return name.slice(0, 2).toUpperCase();
    },
    formatTime(ts) {
      if (!ts) return "";
      const d = new Date(ts);
      const now = new Date();
      if (d.toDateString() === now.toDateString()) {
        return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      }
      return d.toLocaleDateString([], { month: 'short', day: 'numeric' });
    },
  },
  mounted() {
    this.loadConversations();
    this.pollIntervalId = setInterval(this.loadConversations, 5000);
  },
  beforeUnmount() {
    clearInterval(this.pollIntervalId);
  },
};
</script>

<style scoped>
.home {
  max-width: 720px;
  margin: 0 auto;
  padding: 32px 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 28px;
}
.header-title { font-size: 28px; font-weight: 700; color: #1a1f36; margin: 0 0 2px; }
.header-sub { font-size: 14px; color: #718096; margin: 0; }

.header-actions { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.action-btn {
  padding: 8px 14px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: 1.5px solid #e2e8f0;
  background: #fff;
  color: #4a5568;
  transition: all 0.2s;
}
.action-btn:hover { background: #f7f8fc; }
.new-group-btn { background: #4c6ef5; color: #fff; border-color: #4c6ef5; }
.new-group-btn:hover { background: #3b5bdb; border-color: #3b5bdb; }
.logout-btn { color: #c53030; border-color: #fed7d7; }
.logout-btn:hover { background: #fff5f5; }
.icon-btn { padding: 8px 12px; font-size: 16px; }

/* Empty state */
.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #a0aec0;
}
.empty-icon { font-size: 56px; margin-bottom: 16px; }
.empty-title { font-size: 20px; font-weight: 600; color: #4a5568; margin: 0 0 8px; }
.empty-sub { font-size: 15px; margin: 0 0 24px; }
.find-btn {
  padding: 11px 28px;
  background: #4c6ef5; color: #fff;
  border: none; border-radius: 10px;
  font-size: 15px; font-weight: 600; cursor: pointer;
}
.find-btn:hover { background: #3b5bdb; }

/* Conversation list */
.conv-list {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.06);
  overflow: hidden;
}
.conv-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 20px;
  cursor: pointer;
  border-bottom: 1px solid #f0f2f5;
  transition: background 0.15s;
}
.conv-item:last-child { border-bottom: none; }
.conv-item:hover { background: #f7f8fc; }

.conv-photo { flex-shrink: 0; }
.conv-img {
  width: 52px; height: 52px;
  border-radius: 50%;
  object-fit: cover;
}
.conv-initials {
  width: 52px; height: 52px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4c6ef5, #7c3aed);
  color: #fff;
  display: flex; align-items: center; justify-content: center;
  font-size: 18px; font-weight: 700;
}

.conv-body { flex: 1; min-width: 0; }
.conv-top {
  display: flex; justify-content: space-between; align-items: baseline;
  margin-bottom: 3px;
}
.conv-name { font-size: 15px; font-weight: 600; color: #1a1f36; }
.conv-time { font-size: 12px; color: #a0aec0; flex-shrink: 0; margin-left: 8px; }
.conv-preview { font-size: 13px; color: #718096; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.no-msg { color: #cbd5e0; font-style: italic; }
.attach-hint { }

@media (max-width: 480px) {
  .home { padding: 20px 12px; }
  .home-header { flex-direction: column; align-items: flex-start; gap: 12px; margin-bottom: 16px; }
  .header-title { font-size: 22px; }
  .header-actions { width: 100%; justify-content: flex-start; }
  .action-btn { font-size: 13px; padding: 7px 12px; }
}
</style>
