<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{ userName }}, here is your profile</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">Refresh</button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="logOut">Log Out</button>
        </div>
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="newGroup">New group</button>
        </div>
      </div>
    </div>
    
    <div class="profile-container">
      <div class="profile-header">
        <div class="photo-container">
          <img v-if="userPhoto" :src="userPhoto" alt="User Photo" class="profile-photo" />
          <p v-else class="no-photo-placeholder">No Photo</p>
        </div>
        <div class="username-container">
          <h1 class="username">{{ userName }}</h1>
          <div class="update-username-section">
            <input
              v-model="newUserName"
              placeholder="Enter new username"
              maxlength="16"
              minlength="3"
            />
            <button
              class="custom-button"
              @click="updateUsername"
              :disabled="!newUserName || newUserName === userName"
            >
              Update Username
            </button>
          </div>
          <div class="update-photo-section">
            <input type="file" @change="handlePhotoUpload" accept="image/*" />
            <button class="custom-button" @click="updatePhoto" :disabled="!newPhoto">
              Update Photo
            </button>
          </div>
        </div>
      </div>
      <ErrorMsg v-if="errormsg" :msg="errormsg" />
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "ProfileView",
  components: {
    ErrorMsg,
  },
  data() {
    return {
      userName: localStorage.getItem("name") || "",
      newUserName: "",
      errormsg: "",
      userId: localStorage.getItem("userId"),
      token: localStorage.getItem("token")
    };
  },
  methods: {
    async fetchUserProfile() {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          this.$router.push({ path: "/" });
          return;
        }
        const response = await axios.get("/users/photo", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const { photo } = response.data;
        this.userName = localStorage.getItem("name");
        this.userPhoto = photo ? `data:image/jpeg;base64,${photo}` : null;
      } catch (error) {
        console.error("Failed to fetch user profile:", error);
        this.errormsg = "Failed to load user profile. Please try again later.";
      }
    },
    handlePhotoUpload(event) {
      const file = event.target.files[0];
      if (file) {
        this.newPhoto = file;
      }
    },
    async updatePhoto() {
      if (!this.newPhoto) return;
      try {
        const token = localStorage.getItem("token");
        const formData = new FormData();
        formData.append("photo", this.newPhoto);
        await axios.put("/users/photo", formData, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        alert("Photo updated successfully!");
        this.newPhoto = null;
        this.fetchUserProfile(); 
      } catch (error) {
        console.error("Failed to update photo:", error);
        this.errormsg = "Failed to update photo. Please try again.";
      }
    },
    async updateUsername() {
      if (!this.newUserName || this.newUserName === this.userName) {
        this.errormsg = "Please enter a different username";
        return;
      }

      try {
        const response = await axios.put(
          "/users/name",
          { name: this.newUserName },
          {
            headers: {
              Authorization: `Bearer ${this.token}`
            }
          }
        );

        // Check for specific error message
        if (!response.ok) {
          const error = await response.json();
          this.errormsg = error.message; // will show "username already exists" msg
          return;
        }

        // Update local storage with new username
        localStorage.setItem("name", response.data.name);
        this.userName = response.data.name;
        this.newUserName = "";
        this.errormsg = "";

        // Emit event to update header/navbar display
        this.$root.$emit("usernameUpdated", response.data.name);
        
        // Show success message
        this.$toast.success("Username updated successfully!");
        
        // Emit event to update parent/header components
        this.$root.$emit("usernameUpdated", response.data.name);
      } catch (error) {
        console.error("Error updating username:", error);
        
        if (error.response?.data?.message) {
          this.errormsg = error.response.data.message;
        } else if (error.response?.status === 409) {
          this.errormsg = "Username already exists. Please choose another.";
        } else {
          this.errormsg = "Failed to update username. Please try again.";
        }
      }
        
    },
    refresh() {
      this.fetchUserProfile();
    },
    logOut() {
      localStorage.clear();
      this.$router.push({ path: "/" });
    },
    newGroup() {
      this.$router.push({ path: "/new-group" });
    }
  },
  mounted() {
    this.fetchUserProfile();
  },
};
</script>

<style scoped>
.profile-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.profile-header {
  display: flex;
  align-items: flex-start;
  gap: 20px;
  width: 100%;
  max-width: 800px;
}

.photo-container {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  border: 1px solid #ccc;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
}

.profile-photo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-photo-placeholder {
  color: #aaa;
  font-size: 14px;
}

.username-container {
  flex: 1;
}

.username {
  margin: 0;
  font-size: 24px;
  font-weight: bold;
}

.update-username-section,
.update-photo-section {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
}

input {
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  flex: 1;
  max-width: 300px;
}

.custom-button {
  padding: 8px 16px;
  background-color: transparent;
  border: 1px solid #007bff;
  color: #007bff;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.custom-button:hover:not(:disabled) {
  background-color: #007bff;
  color: white;
}

.custom-button:disabled {
  border-color: #cccccc;
  color: #cccccc;
  cursor: not-allowed;
  background-color: transparent;
}
</style>