<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">Profile</h1>
      <p class="page-sub">Manage your account</p>
    </div>

    <div class="profile-card">
      <!-- Avatar + name -->
      <div class="identity-row">
        <div class="avatar-wrap">
          <img v-if="userPhoto" :src="userPhoto" alt="Photo" class="avatar-img" />
          <div v-else class="avatar-initials">{{ initials }}</div>
        </div>
        <div class="identity-info">
          <h2 class="display-name">{{ userName }}</h2>
          <span class="online-dot">● Online</span>
        </div>
      </div>

      <ErrorMsg v-if="errormsg" :msg="errormsg" />
      <div v-if="successMsg" class="success-banner">{{ successMsg }}</div>

      <hr class="divider" />

      <!-- Update username -->
      <div class="section">
        <label class="field-label">Change Username</label>
        <div class="input-row">
          <input
            v-model="newUserName"
            class="field-input"
            type="text"
            placeholder="New username (3–16 chars)"
            maxlength="16"
            @keyup.enter="updateUsername"
          />
          <button
            class="save-btn"
            @click="updateUsername"
            :disabled="!newUserName || newUserName === userName || savingName"
          >
            {{ savingName ? '…' : 'Save' }}
          </button>
        </div>
      </div>

      <!-- Update photo -->
      <div class="section">
        <label class="field-label">Profile Photo</label>
        <div class="photo-row">
          <label class="upload-btn">
            Choose Photo
            <input type="file" @change="handlePhotoUpload" accept="image/*" style="display:none" />
          </label>
          <span v-if="newPhoto" class="file-selected">{{ newPhoto.name }}</span>
          <button
            class="save-btn"
            @click="updatePhoto"
            :disabled="!newPhoto || savingPhoto"
          >
            {{ savingPhoto ? '…' : 'Upload' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "ProfileView",
  components: { ErrorMsg },
  data() {
    return {
      userName: localStorage.getItem("name") || "",
      userPhoto: null,
      newUserName: "",
      newPhoto: null,
      errormsg: "",
      successMsg: "",
      savingName: false,
      savingPhoto: false,
    };
  },
  computed: {
    initials() {
      return (this.userName || "U").slice(0, 2).toUpperCase();
    },
  },
  methods: {
    async fetchUserProfile() {
      const token = localStorage.getItem("token");
      if (!token) { this.$router.push("/"); return; }
      try {
        const response = await axios.get("/users/photo", {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.userName = localStorage.getItem("name");
        const { photo } = response.data;
        this.userPhoto = photo ? `data:image/jpeg;base64,${photo}` : null;
      } catch {
        this.errormsg = "Failed to load profile.";
      }
    },
    handlePhotoUpload(event) {
      this.newPhoto = event.target.files[0] || null;
    },
    async updatePhoto() {
      if (!this.newPhoto) return;
      this.savingPhoto = true;
      this.errormsg = "";
      this.successMsg = "";
      try {
        const token = localStorage.getItem("token");
        const formData = new FormData();
        formData.append("photo", this.newPhoto);
        await axios.put("/users/photo", formData, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.successMsg = "Photo updated!";
        this.newPhoto = null;
        await this.fetchUserProfile();
      } catch {
        this.errormsg = "Failed to update photo.";
      } finally {
        this.savingPhoto = false;
      }
    },
    async updateUsername() {
      if (!this.newUserName || this.newUserName === this.userName) return;
      this.savingName = true;
      this.errormsg = "";
      this.successMsg = "";
      try {
        const token = localStorage.getItem("token");
        const response = await axios.put("/users/name", { name: this.newUserName }, {
          headers: { Authorization: `Bearer ${token}` },
        });
        localStorage.setItem("name", response.data.name);
        this.userName = response.data.name;
        this.newUserName = "";
        this.successMsg = "Username updated!";
      } catch (error) {
        if (error.response?.status === 409) this.errormsg = "Username already taken.";
        else this.errormsg = error.response?.data || "Failed to update username.";
      } finally {
        this.savingName = false;
      }
    },
  },
  mounted() {
    this.fetchUserProfile();
  },
};
</script>

<style scoped>
.page {
  max-width: 600px;
  margin: 0 auto;
  padding: 40px 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
.page-header { margin-bottom: 28px; }
.page-title { font-size: 28px; font-weight: 700; color: #1a1f36; margin: 0 0 6px; }
.page-sub { font-size: 15px; color: #718096; margin: 0; }

.profile-card {
  background: #fff;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.07);
}

.identity-row {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 8px;
}
.avatar-wrap { flex-shrink: 0; }
.avatar-img {
  width: 80px; height: 80px;
  border-radius: 50%; object-fit: cover;
}
.avatar-initials {
  width: 80px; height: 80px; border-radius: 50%;
  background: linear-gradient(135deg, #4c6ef5, #7c3aed);
  color: #fff;
  display: flex; align-items: center; justify-content: center;
  font-size: 28px; font-weight: 700;
}
.identity-info { flex: 1; }
.display-name { font-size: 22px; font-weight: 700; color: #1a1f36; margin: 0 0 4px; }
.online-dot { font-size: 12px; color: #68d391; font-weight: 500; }

.success-banner {
  background: #f0fff4; border: 1px solid #9ae6b4; color: #276749;
  border-radius: 8px; padding: 10px 14px; font-size: 13px; margin-top: 12px;
}

.divider { border: none; border-top: 1px solid #f0f2f5; margin: 24px 0; }

.section { margin-bottom: 24px; }
.section:last-child { margin-bottom: 0; }

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #4a5568;
  margin-bottom: 8px;
}
.input-row { display: flex; gap: 10px; }
.field-input {
  flex: 1;
  padding: 11px 14px;
  border: 1.5px solid #e2e8f0;
  border-radius: 10px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}
.field-input:focus { border-color: #4c6ef5; }

.save-btn {
  padding: 11px 20px;
  background: #4c6ef5; color: #fff;
  border: none; border-radius: 10px;
  font-size: 14px; font-weight: 600; cursor: pointer;
  white-space: nowrap; transition: background 0.2s;
}
.save-btn:hover:not(:disabled) { background: #3b5bdb; }
.save-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.photo-row { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.upload-btn {
  padding: 11px 18px;
  background: #f0f2f5; color: #4a5568;
  border-radius: 10px; font-size: 14px; font-weight: 500; cursor: pointer;
}
.upload-btn:hover { background: #e2e8f0; }
.file-selected { font-size: 13px; color: #718096; flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

@media (max-width: 480px) {
  .page { padding: 20px 12px; }
  .profile-card { padding: 20px; }
  .input-row { flex-direction: column; }
  .save-btn { width: 100%; }
  .identity-row { flex-direction: column; align-items: flex-start; gap: 12px; }
  .page-title { font-size: 22px; }
}
</style>
