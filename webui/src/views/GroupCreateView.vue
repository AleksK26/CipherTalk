<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">Create Group</h1>
      <p class="page-sub">Add a name, find members and start chatting</p>
    </div>

    <div class="card">
      <!-- Group Name -->
      <div class="section">
        <label class="field-label">Group Name <span class="required">*</span></label>
        <input
          v-model="groupName"
          class="field-input"
          type="text"
          placeholder="e.g. Project Team (3–16 chars)"
          maxlength="16"
        />
        <p class="field-hint">{{ groupName.length }}/16 characters</p>
      </div>

      <!-- Group Image (optional) -->
      <div class="section">
        <label class="field-label">Group Photo <span class="optional">(optional)</span></label>
        <div class="photo-row">
          <div class="photo-preview">
            <img v-if="previewImage" :src="previewImage" alt="Preview" class="preview-img" />
            <div v-else class="preview-placeholder">{{ groupName ? getInitials(groupName) : '?' }}</div>
          </div>
          <label class="upload-btn">
            Choose Image
            <input type="file" ref="fileInput" @change="handleImageUpload" accept="image/*" style="display:none" />
          </label>
          <button v-if="file" class="clear-btn" @click="clearImage">Remove</button>
        </div>
      </div>

      <!-- Member Search -->
      <div class="section">
        <label class="field-label">Add Members <span class="required">*</span></label>
        <div class="search-row">
          <input
            v-model="query"
            class="field-input"
            type="text"
            placeholder="Search by username…"
            @keyup.enter="searchUsers"
          />
          <button class="search-btn" @click="searchUsers" :disabled="loading">Search</button>
        </div>

        <div v-if="error" class="err-box">{{ error }}</div>
        <LoadingSpinner v-if="loading" />

        <div v-if="!loading && showResults && users.length > 0" class="results-list">
          <div v-for="user in users" :key="user.id" class="result-item">
            <div class="user-avatar">{{ getInitials(user.name) }}</div>
            <span class="user-name">{{ user.name }}</span>
            <button
              class="add-btn"
              @click="addUserToGroup(user)"
              :disabled="isUserAdded(user)"
            >
              {{ isUserAdded(user) ? 'Added' : 'Add' }}
            </button>
          </div>
        </div>
        <p v-if="!loading && showResults && users.length === 0" class="empty-state">No users found for "{{ lastQuery }}"</p>
      </div>

      <!-- Selected Members -->
      <div v-if="selectedUsers.length > 0" class="section">
        <label class="field-label">Members ({{ selectedUsers.length }})</label>
        <div class="chips">
          <div v-for="user in selectedUsers" :key="user.id" class="chip">
            <span>{{ user.name }}</span>
            <button @click="removeUserFromGroup(user)" class="chip-remove">×</button>
          </div>
        </div>
      </div>

      <!-- Validation Summary -->
      <div class="validation-list">
        <span :class="groupName.trim().length >= 3 ? 'ok' : 'missing'">
          {{ groupName.trim().length >= 3 ? '✓' : '○' }} Group name (min 3 chars)
        </span>
        <span :class="selectedUsers.length > 0 ? 'ok' : 'missing'">
          {{ selectedUsers.length > 0 ? '✓' : '○' }} At least one member
        </span>
      </div>

      <div v-if="createError" class="err-box">{{ createError }}</div>

      <button
        class="create-btn"
        @click="createGroup"
        :disabled="!canCreateGroup || creating"
      >
        {{ creating ? 'Creating…' : 'Create Group' }}
      </button>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
  name: "GroupCreateView",
  components: { LoadingSpinner },
  data() {
    return {
      groupName: "",
      query: "",
      lastQuery: "",
      users: [],
      loading: false,
      showResults: false,
      error: "",
      createError: "",
      selectedUsers: [],
      previewImage: null,
      file: null,
      creating: false,
    };
  },
  computed: {
    canCreateGroup() {
      return this.groupName.trim().length >= 3 && this.selectedUsers.length > 0;
    },
  },
  methods: {
    async searchUsers() {
      if (!this.query.trim()) { this.error = "Enter a username to search."; return; }
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
        this.users = response.data.filter(u => u.id !== token);
        this.lastQuery = this.query;
        this.showResults = true;
      } catch (err) {
        this.error = "Failed to search users.";
      } finally {
        this.loading = false;
      }
    },
    addUserToGroup(user) {
      if (!this.isUserAdded(user)) this.selectedUsers.push(user);
    },
    isUserAdded(user) {
      return this.selectedUsers.some(u => u.id === user.id);
    },
    removeUserFromGroup(user) {
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id);
    },
    handleImageUpload(event) {
      const f = event.target.files[0];
      if (!f) return;
      this.file = f;
      const reader = new FileReader();
      reader.onload = e => { this.previewImage = e.target.result; };
      reader.readAsDataURL(f);
    },
    clearImage() {
      this.file = null;
      this.previewImage = null;
      this.$refs.fileInput.value = "";
    },
    getInitials(name) {
      if (!name) return "?";
      return name.slice(0, 2).toUpperCase();
    },
    async createGroup() {
      if (!this.canCreateGroup) return;
      this.creating = true;
      this.createError = "";
      const token = localStorage.getItem("token");
      const formData = new FormData();
      formData.append("name", this.groupName.trim());
      // Include current user in members
      const memberIds = [...this.selectedUsers.map(u => u.id), token];
      formData.append("members", JSON.stringify(memberIds));
      if (this.file) {
        formData.append("image", this.file);
      } else {
        // Generate a default 1x1 transparent PNG as placeholder
        const canvas = document.createElement("canvas");
        canvas.width = 1; canvas.height = 1;
        canvas.toBlob(blob => {}, "image/png");
        // Use a fetch of a default image instead
        try {
          const res = await fetch('/nopfp.jpg');
          const blob = await res.blob();
          formData.append("image", blob, "default.jpg");
        } catch {
          // If that fails, create a minimal placeholder
          const arr = new Uint8Array([137,80,78,71,13,10,26,10,0,0,0,13,73,72,68,82,0,0,0,1,0,0,0,1,8,2,0,0,0,144,119,83,222,0,0,0,12,73,68,65,84,8,215,99,248,207,192,0,0,0,2,0,1,226,33,188,51,0,0,0,0,73,69,78,68,174,66,96,130]);
          const blob = new Blob([arr], { type: "image/png" });
          formData.append("image", blob, "default.png");
        }
      }
      try {
        await axios.post(`/groups`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
            Authorization: `Bearer ${token}`,
          },
        });
        this.$router.push("/home");
      } catch (err) {
        const msg = err.response?.data || "Failed to create group.";
        this.createError = `Error: ${msg}`;
      } finally {
        this.creating = false;
      }
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
.page-header { margin-bottom: 28px; }
.page-title { font-size: 28px; font-weight: 700; color: #1a1f36; margin: 0 0 6px; }
.page-sub { font-size: 15px; color: #718096; margin: 0; }

.card {
  background: #fff;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.07);
}
.section { margin-bottom: 28px; }

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #4a5568;
  margin-bottom: 8px;
}
.required { color: #e53e3e; }
.optional { color: #a0aec0; font-weight: 400; }
.field-input {
  width: 100%;
  padding: 11px 14px;
  border: 1.5px solid #e2e8f0;
  border-radius: 10px;
  font-size: 15px;
  outline: none;
  box-sizing: border-box;
  transition: border-color 0.2s;
}
.field-input:focus { border-color: #4c6ef5; }
.field-hint { font-size: 12px; color: #a0aec0; margin: 4px 0 0; }

/* Photo upload */
.photo-row { display: flex; align-items: center; gap: 14px; }
.photo-preview {
  width: 60px; height: 60px; border-radius: 50%;
  overflow: hidden; background: #e8f0fe;
  display: flex; align-items: center; justify-content: center;
  font-size: 20px; font-weight: 700; color: #4c6ef5;
  flex-shrink: 0;
}
.preview-img { width: 100%; height: 100%; object-fit: cover; }
.preview-placeholder {}
.upload-btn {
  padding: 8px 18px;
  background: #f0f2f5;
  color: #4a5568;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}
.upload-btn:hover { background: #e2e8f0; }
.clear-btn {
  padding: 8px 14px;
  background: #fff5f5;
  color: #c53030;
  border: 1px solid #fed7d7;
  border-radius: 8px;
  font-size: 13px;
  cursor: pointer;
}

/* Member search */
.search-row { display: flex; gap: 10px; }
.search-btn {
  padding: 11px 20px;
  background: #4c6ef5;
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
}
.search-btn:hover:not(:disabled) { background: #3b5bdb; }
.search-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.err-box {
  background: #fff5f5; border: 1px solid #fed7d7; color: #c53030;
  border-radius: 8px; padding: 10px 14px; font-size: 13px; margin: 12px 0;
}

.results-list { margin-top: 12px; border: 1px solid #e2e8f0; border-radius: 10px; overflow: hidden; }
.result-item {
  display: flex; align-items: center; gap: 12px;
  padding: 12px 14px;
  border-bottom: 1px solid #f0f2f5;
}
.result-item:last-child { border-bottom: none; }
.user-avatar {
  width: 36px; height: 36px; border-radius: 50%;
  background: linear-gradient(135deg, #4c6ef5, #7c3aed);
  color: #fff; display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 700; flex-shrink: 0;
}
.user-name { flex: 1; font-size: 14px; font-weight: 500; color: #1a1f36; }
.add-btn {
  padding: 6px 14px; background: #e8f0fe; color: #4c6ef5;
  border: none; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer;
}
.add-btn:disabled { background: #f0faf0; color: #38a169; cursor: default; }

.empty-state { text-align: center; color: #a0aec0; font-size: 14px; padding: 16px 0; }

/* Selected chips */
.chips { display: flex; flex-wrap: wrap; gap: 8px; }
.chip {
  display: flex; align-items: center; gap: 6px;
  background: #e8f0fe; color: #3b5bdb;
  border-radius: 20px; padding: 6px 12px;
  font-size: 13px; font-weight: 500;
}
.chip-remove {
  background: none; border: none; color: #4c6ef5;
  font-size: 16px; cursor: pointer; line-height: 1; padding: 0;
}

/* Validation */
.validation-list {
  display: flex; gap: 20px; margin-bottom: 20px;
  font-size: 13px;
}
.ok { color: #38a169; font-weight: 500; }
.missing { color: #a0aec0; }

.create-btn {
  width: 100%;
  padding: 14px;
  background: #4c6ef5;
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.create-btn:hover:not(:disabled) { background: #3b5bdb; }
.create-btn:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
