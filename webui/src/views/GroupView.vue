<template>
  <div class="page">
    <div class="page-header">
      <button class="back-btn" @click="$router.back()">← Back</button>
    </div>

    <!-- Group identity card -->
    <div class="group-card">
      <div class="group-photo-wrap">
        <img
          v-if="group.photo"
          :src="`data:image/jpeg;base64,${group.photo}`"
          class="group-photo"
          alt="Group"
        />
        <div v-else class="group-initials">{{ group.name ? group.name.slice(0,2).toUpperCase() : '?' }}</div>
        <label v-if="isAdmin" class="photo-edit-btn" title="Change photo">
          ✏️
          <input type="file" @change="handleGroupPhotoUpload" accept="image/*" style="display:none" />
        </label>
      </div>

      <div class="group-identity">
        <div v-if="editingName" class="name-edit-row">
          <input v-model="newGroupName" class="name-input" maxlength="16" />
          <button class="save-btn" @click="updateGroupName" :disabled="!newGroupName.trim()">Save</button>
          <button class="cancel-btn" @click="editingName = false">Cancel</button>
        </div>
        <div v-else class="name-row">
          <h2 class="group-name">{{ group.name }}</h2>
          <button v-if="isAdmin" class="edit-icon-btn" @click="startEditName" title="Rename group">✏️</button>
        </div>
        <p class="member-count">{{ group.members.length }} member{{ group.members.length !== 1 ? 's' : '' }}</p>
      </div>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />

    <!-- Members list -->
    <div class="section-card">
      <div class="section-head">
        <h3 class="section-title">Members</h3>
        <button v-if="isAdmin" class="add-member-trigger" @click="showAddMember = !showAddMember">
          + Add Member
        </button>
      </div>

      <!-- Add member search (admin only) -->
      <div v-if="isAdmin && showAddMember" class="add-member-form">
        <div class="search-row">
          <input v-model="newMemberQuery" class="field-input" placeholder="Search by username…" @keyup.enter="searchNewMember" />
          <button class="search-btn" @click="searchNewMember">Search</button>
        </div>
        <div v-if="addMemberResults.length > 0" class="search-results">
          <div v-for="u in addMemberResults" :key="u.id" class="result-row">
            <div class="member-avatar sm">{{ u.name.slice(0,2).toUpperCase() }}</div>
            <span class="result-name">{{ u.name }}</span>
            <button class="add-btn" @click="addMember(u)" :disabled="isMember(u.id)">
              {{ isMember(u.id) ? 'Already in group' : 'Add' }}
            </button>
          </div>
        </div>
      </div>

      <ul class="members-list">
        <li v-for="member in group.members" :key="member.id" class="member-row">
          <div class="member-avatar-wrap">
            <img
              v-if="member.photo"
              :src="`data:image/jpeg;base64,${member.photo}`"
              class="member-photo"
              alt=""
              @error="e => e.target.style.display='none'"
            />
            <div v-else class="member-avatar">{{ member.name.slice(0,2).toUpperCase() }}</div>
          </div>
          <div class="member-info">
            <span class="member-name">{{ member.name }}</span>
            <span v-if="member.id === currentUserId" class="you-tag">(you)</span>
          </div>
          <span :class="['role-tag', member.role === 'admin' ? 'role-admin' : 'role-member']">
            {{ member.role === 'admin' ? 'Admin' : 'Member' }}
          </span>
          <button
            v-if="isAdmin && member.id !== currentUserId && member.role !== 'admin'"
            class="remove-btn"
            @click="removeMember(member)"
            title="Remove from group"
          >Remove</button>
        </li>
      </ul>
    </div>

    <!-- Danger zone -->
    <div class="danger-zone">
      <button @click="leaveGroup" class="leave-btn">Leave Group</button>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "GroupView",
  components: { ErrorMsg },
  data() {
    return {
      group: { id: "", name: "", photo: null, members: [] },
      newGroupName: "",
      editingName: false,
      newGroupPhoto: null,
      newMemberQuery: "",
      addMemberResults: [],
      showAddMember: false,
      errormsg: null,
      currentUserId: localStorage.getItem("token"),
    };
  },
  computed: {
    isAdmin() {
      return this.group.members.some(m => m.id === this.currentUserId && m.role === "admin");
    },
  },
  methods: {
    async fetchGroupDetails() {
      try {
        const token = localStorage.getItem("token");
        const response = await axios.get(`/groups/${this.$route.params.groupId}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        const data = response.data;
        this.group = {
          id: data.id,
          name: data.name,
          photo: data.groupPhoto || null,
          members: data.members || [],
        };
      } catch {
        this.errormsg = "Failed to load group details.";
      }
    },
    startEditName() {
      this.newGroupName = this.group.name;
      this.editingName = true;
    },
    async updateGroupName() {
      if (!this.newGroupName.trim()) return;
      try {
        const token = localStorage.getItem("token");
        await axios.put(`/groups/${this.group.id}/name`, { groupName: this.newGroupName }, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.group.name = this.newGroupName;
        this.editingName = false;
      } catch (err) {
        this.errormsg = err.response?.data || "Failed to update group name.";
      }
    },
    async handleGroupPhotoUpload(event) {
      const file = event.target.files[0];
      if (!file) return;
      const token = localStorage.getItem("token");
      const formData = new FormData();
      formData.append("photo", file);
      try {
        await axios.put(`/groups/${this.group.id}/photo`, formData, {
          headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'multipart/form-data' },
        });
        const reader = new FileReader();
        reader.onload = e => { this.group.photo = e.target.result.split(',')[1]; };
        reader.readAsDataURL(file);
      } catch {
        this.errormsg = "Failed to update group photo.";
      }
    },
    async searchNewMember() {
      if (!this.newMemberQuery.trim()) return;
      try {
        const token = localStorage.getItem("token");
        const response = await axios.get(`/search`, {
          params: { username: this.newMemberQuery },
          headers: { Authorization: `Bearer ${token}` },
        });
        this.addMemberResults = response.data;
      } catch {
        this.errormsg = "Failed to search users.";
      }
    },
    isMember(userId) {
      return this.group.members.some(m => m.id === userId);
    },
    async addMember(user) {
      try {
        const token = localStorage.getItem("token");
        await axios.post(`/groups/${this.group.id}`, { userId: user.id }, {
          headers: { Authorization: `Bearer ${token}` },
        });
        await this.fetchGroupDetails();
        this.addMemberResults = [];
        this.newMemberQuery = "";
        this.showAddMember = false;
      } catch (err) {
        this.errormsg = err.response?.data || "Failed to add member.";
      }
    },
    async removeMember(member) {
      if (!confirm(`Remove ${member.name} from the group?`)) return;
      try {
        const token = localStorage.getItem("token");
        await axios.delete(`/groups/${this.group.id}/members/${member.id}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.group.members = this.group.members.filter(m => m.id !== member.id);
      } catch (err) {
        this.errormsg = err.response?.data || "Failed to remove member.";
      }
    },
    async leaveGroup() {
      if (!confirm("Are you sure you want to leave this group?")) return;
      try {
        const token = localStorage.getItem("token");
        await axios.delete(`/groups/${this.group.id}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.$router.push("/home");
      } catch (err) {
        this.errormsg = err.response?.data || "Failed to leave group.";
      }
    },
  },
  mounted() {
    this.fetchGroupDetails();
  },
};
</script>

<style scoped>
.page {
  max-width: 680px;
  margin: 0 auto;
  padding: 24px 20px 48px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
.page-header { margin-bottom: 20px; }
.back-btn {
  background: none; border: none; color: #4c6ef5;
  font-size: 15px; font-weight: 500; cursor: pointer; padding: 0;
}
.back-btn:hover { text-decoration: underline; }

/* Group card */
.group-card {
  background: #fff;
  border-radius: 16px;
  padding: 28px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.07);
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
}

.group-photo-wrap { position: relative; flex-shrink: 0; }
.group-photo {
  width: 80px; height: 80px; border-radius: 50%; object-fit: cover;
}
.group-initials {
  width: 80px; height: 80px; border-radius: 50%;
  background: linear-gradient(135deg, #4c6ef5, #7c3aed);
  color: #fff; display: flex; align-items: center; justify-content: center;
  font-size: 28px; font-weight: 700;
}
.photo-edit-btn {
  position: absolute; bottom: 0; right: 0;
  background: #fff; border-radius: 50%;
  width: 28px; height: 28px; display: flex; align-items: center; justify-content: center;
  box-shadow: 0 1px 6px rgba(0,0,0,0.15); cursor: pointer; font-size: 13px;
}

.group-identity { flex: 1; min-width: 0; }
.name-row { display: flex; align-items: center; gap: 10px; }
.group-name { font-size: 22px; font-weight: 700; color: #1a1f36; margin: 0 0 4px; }
.edit-icon-btn { background: none; border: none; font-size: 15px; cursor: pointer; color: #a0aec0; }
.edit-icon-btn:hover { color: #4c6ef5; }
.member-count { font-size: 14px; color: #718096; margin: 0; }

.name-edit-row { display: flex; gap: 8px; align-items: center; }
.name-input {
  flex: 1; padding: 8px 12px; border: 1.5px solid #4c6ef5;
  border-radius: 8px; font-size: 15px; outline: none;
}
.save-btn { padding: 8px 14px; background: #4c6ef5; color: #fff; border: none; border-radius: 8px; cursor: pointer; font-weight: 600; }
.cancel-btn { padding: 8px 14px; background: #f0f2f5; color: #4a5568; border: none; border-radius: 8px; cursor: pointer; }

/* Section card */
.section-card {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.07);
  margin-bottom: 16px;
}
.section-head {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 16px;
}
.section-title { font-size: 15px; font-weight: 700; color: #1a1f36; margin: 0; }
.add-member-trigger {
  padding: 7px 14px; background: #e8f0fe; color: #4c6ef5;
  border: none; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer;
}
.add-member-trigger:hover { background: #c5d5fc; }

/* Add member search */
.add-member-form { margin-bottom: 16px; padding: 14px; background: #f7f8fc; border-radius: 10px; }
.search-row { display: flex; gap: 8px; margin-bottom: 10px; }
.field-input {
  flex: 1; padding: 9px 12px; border: 1.5px solid #e2e8f0;
  border-radius: 8px; font-size: 14px; outline: none;
}
.field-input:focus { border-color: #4c6ef5; }
.search-btn { padding: 9px 16px; background: #4c6ef5; color: #fff; border: none; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; }
.search-results { border: 1px solid #e2e8f0; border-radius: 8px; overflow: hidden; }
.result-row { display: flex; align-items: center; gap: 10px; padding: 10px 12px; border-bottom: 1px solid #f0f2f5; }
.result-row:last-child { border-bottom: none; }
.result-name { flex: 1; font-size: 14px; }
.add-btn { padding: 6px 12px; background: #e8f0fe; color: #4c6ef5; border: none; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; }
.add-btn:disabled { opacity: 0.5; cursor: default; }

/* Members */
.members-list { list-style: none; padding: 0; margin: 0; }
.member-row {
  display: flex; align-items: center; gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #f0f2f5;
}
.member-row:last-child { border-bottom: none; }

.member-avatar-wrap { flex-shrink: 0; }
.member-photo { width: 40px; height: 40px; border-radius: 50%; object-fit: cover; }
.member-avatar, .member-avatar.sm {
  width: 40px; height: 40px; border-radius: 50%;
  background: linear-gradient(135deg, #4c6ef5, #7c3aed);
  color: #fff; display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 700;
}
.member-avatar.sm { width: 32px; height: 32px; font-size: 12px; }

.member-info { flex: 1; min-width: 0; }
.member-name { font-size: 14px; font-weight: 500; color: #1a1f36; }
.you-tag { font-size: 12px; color: #a0aec0; margin-left: 4px; }

.role-tag {
  font-size: 11px; font-weight: 600; border-radius: 4px; padding: 2px 8px;
}
.role-admin { background: #fff3cd; color: #856404; }
.role-member { background: #f0f2f5; color: #718096; }

.remove-btn {
  padding: 5px 10px; background: none; border: 1px solid #fed7d7;
  color: #c53030; border-radius: 6px; font-size: 12px; cursor: pointer;
}
.remove-btn:hover { background: #fff5f5; }

/* Danger zone */
.danger-zone { margin-top: 8px; }
.leave-btn {
  width: 100%; padding: 13px;
  background: #fff5f5; color: #c53030;
  border: 1px solid #fed7d7; border-radius: 10px;
  font-size: 15px; font-weight: 600; cursor: pointer;
}
.leave-btn:hover { background: #fed7d7; }
</style>
