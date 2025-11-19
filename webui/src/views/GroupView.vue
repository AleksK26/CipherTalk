<template>
  <div class="group-container">
    <h2>{{ group.name }}</h2>
    
    <!-- Group Photo -->
    <div class="group-photo-container">
      <img :src="group.photo ? `data:image/jpeg;base64,${group.photo}` : defaultGroupPhoto" 
           alt="Group Photo" 
           class="group-photo" />
      <input type="file" @change="handleGroupPhotoUpload" accept="image/*" />
      <button @click="updateGroupPhoto" :disabled="!newGroupPhoto">Update Group Photo</button>
    </div>

    <!-- Group Name -->
    <div class="group-name-container">
      <input v-model="newGroupName" :placeholder="group.name" />
      <button @click="updateGroupName" :disabled="!newGroupName">Update Group Name</button>
    </div>

    <!-- Members List -->
    <div class="members-list">
      <h3>Members</h3>
      <ul>
        <li v-for="member in group.members" :key="member.id" class="member-item">
           <img 
            :src="member.photo ? `data:image/jpeg;base64,${member.photo}` : defaultUserPhoto" 
            alt="Member Photo" 
            class="member-photo" 
            @error="handleImageError"
          />
          <span>{{ member.name }}</span>
          <button @click="removeMember(member.id)" 
                  v-if="isGroupAdmin" 
                  class="remove-member">Remove</button>
        </li>
      </ul>
    </div>

    <!-- Add Member -->
    <div class="add-member-container" v-if="isGroupAdmin">
      <input v-model="newMemberUsername" placeholder="Enter username to add" />
      <button @click="addMember" :disabled="!newMemberUsername">Add Member</button>
    </div>

    <!-- Leave Group -->
    <button @click="leaveGroup" class="leave-group">Leave Group</button>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";
import defaultUserPhoto from "../assets/default-user.png";
import defaultGroupPhoto from "../assets/default-group.png";

export default {
  name: "GroupView",
  components: { ErrorMsg },
  data() {
    return {
      group: {
        id: "",
        name: "",
        photo: null,
        members: []
      },
      newGroupName: "",
      newGroupPhoto: null,
      newMemberUsername: "",
      errormsg: null,
      isGroupAdmin: false,
      defaultUserPhoto,
      defaultGroupPhoto
    };
  },
  methods: {
    async fetchGroupDetails() {
      try {
        const response = await axios.get(`/groups/${this.$route.params.groupId}`);
        this.group = response.data;

        // DEBUG: Log member photos to verify they're being sent
        console.log("Group members:", this.group.members);
        this.group.members.forEach(m => {
          console.log(`${m.name} photo exists:`, !!m.photo);
        });
      } catch (error) {
        this.errormsg = "Failed to load group details";
      }
    },
    handleImageError(event) {
      console.warn(`Failed to load image for: ${event.target.alt}`);
      event.target.src = this.defaultUserPhoto;
    },
    
    async updateGroupName() {
      try {
        await axios.put(`/groups/${this.group.id}`, {
          name: this.newGroupName
        });
        this.group.name = this.newGroupName;
        this.newGroupName = "";
        this.$toast.success("Group name updated!");
      } catch (error) {
        this.errormsg = "Failed to update group name";
      }
    },
    async addMember() {
      try {
        await axios.post(`/groups/${this.group.id}/members`, {
          username: this.newMemberUsername
        });
        this.newMemberUsername = "";
        await this.fetchGroupDetails();
        this.$toast.success("Member added successfully!");
      } catch (error) {
        this.errormsg = error.response?.data?.message || "Failed to add member";
      }
    },
    async leaveGroup() {
      if (confirm("Are you sure you want to leave this group?")) {
        try {
          await axios.delete(`/groups/${this.group.id}/members/me`);
          this.$router.push("/dashboard");
        } catch (error) {
          this.errormsg = "Failed to leave group";
        }
      }
    }
  },
  mounted() {
    this.fetchGroupDetails();
  }
  
};
</script>

<style scoped>
.group-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.group-photo-container {
  margin: 20px 0;
  text-align: center;
}

.group-photo {
  width: 150px;
  height: 150px;
  border-radius: 75px;
  object-fit: cover;
}

.member-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.member-photo {
  width: 40px;
  height: 40px;
  border-radius: 20px;
  margin-right: 10px;
}

.leave-group {
  margin-top: 20px;
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
}
</style>