<template>
  <div class="chat-container">
    <div class="chat-header">
      <div class="chat-photo" v-if="conversationPhoto">
        <img :src="'data:image/jpeg;base64,' + conversationPhoto" alt="Chat Thumbnail" />
      </div>
      <h3>{{ convName }}</h3>
      <button @click="toggleGroupInfo" v-if="conversationType === 'group'" class="info-button">
        ℹ️ Members
      </button>
    </div>

    <!-- Main chat area with optional sidebar -->
    <div style="display: flex; flex: 1; overflow: hidden;">
      <!-- Messages section -->
      <div style="flex: 1; display: flex; flex-direction: column;">
        <div class="chat-messages" ref="chatMessages">
          <p v-if="messages.length === 0">No messages yet...</p>
          <div
            v-for="message in messages"
            :key="message.id"
            class="message"
            :class="message.senderId === userToken ? 'self' : 'other'"
            :style="message.senderId !== userToken && conversationType === 'group' ? { paddingLeft: '45px' } : {}"
          >
            <div v-if="conversationType === 'group' && message.senderId !== userToken" class="sender-thumbnail">
              <img :src="'data:image/jpeg;base64,' + message.senderPhoto" alt="Sender Photo" />
            </div>
            <div class="message-content">
              <div v-if="message.replyTo" class="reply-preview">
                <small>Replying to {{ message.replySenderName || 'Unknown' }}: {{ message.replyContent }}</small>
                <img
                  v-if="message.replyAttachment"
                  :src="'data:image/jpeg;base64,' + message.replyAttachment"
                  alt="Reply Attachment"
                  class="reply-attachment"
                />
              </div>
              <p>
                <strong>
                  {{ message.senderId === userToken ? 'You' : (message.senderName || 'Unknown Sender') }}:
                </strong>
                {{ message.content }}
              </p>
              <div v-if="message.attachment" class="attachment-container">
                <img :src="'data:image/jpeg;base64,' + message.attachment" alt="Attachment" class="attachment-image" />
              </div>
              <small>{{ formatTimestamp(message.timestamp) }}</small>
              <div v-if="message.reactionCount > 0" class="reaction-count">
                ❤️ × {{ message.reactionCount }}
                <div class="reactors-list">
                  <ul>
                    <li v-for="(reactor, idx) in message.reactingUserNames" :key="reactor">
                      {{ idx + 1 }}. {{ reactor }}
                    </li>
                  </ul>
                </div>
              </div>
              <div class="action-buttons">
                <button v-if="message.senderId !== userToken" class="action-button reply-button" @click.stop="setReply(message)">
                  ↩
                </button>
                <button
                  v-if="message.senderId !== userToken"
                  class="action-button heart-button"
                  :class="{ 'has-reacted': (message.reactingUserNames || []).includes(userName) }"
                  :disabled="message.reactionLoading"
                  @click.stop="toggleReaction(message)"
                >
                  ❤️
                </button>
                <button class="action-button forward-button" @click.stop="showForwardOptions(message.id)">
                  →
                </button>
                <button v-if="message.senderId === userToken" class="action-button delete-button" @click.stop="deleteMessage(message)">
                  ✖
                </button>
              </div>
              <div v-if="messageOptions[message.id]?.showForwardMenu" class="forward-options" @click.stop>
                <select id="forward-select" class="forward-select" v-model="messageOptions[message.id].selectedConversationId">
                  <option value="" disabled>Select conversation</option>
                  <option v-for="conv in messageOptions[message.id].forwardConversations" :key="conv.id" :value="conv.id">
                    {{ conv.name }}
                  </option>
                  <option value="new">New contact</option>
                </select>
                <div v-if="messageOptions[message.id].selectedConversationId === 'new'" class="contact-search">
                  <input type="text" v-model="messageOptions[message.id].contactQuery" placeholder="Enter contact name" @input="searchContact(message.id)" />
                  <ul v-if="messageOptions[message.id].contactResults.length > 0" class="contact-results">
                    <li v-for="contact in messageOptions[message.id].contactResults" :key="contact.id" @click="selectContact(contact, message.id)" class="contact-result">
                      {{ contact.name }}
                    </li>
                  </ul>
                </div>
                <div class="forward-buttons-container">
                  <button class="button-style" v-if="messageOptions[message.id].selectedConversationId !== 'new'" :disabled="!messageOptions[message.id].selectedConversationId" @click.stop="forwardMessage(messageOptions[message.id].selectedConversationId, message.id)">
                    Send
                  </button>
                  <button class="button-style" v-if="messageOptions[message.id].selectedConversationId === 'new'" :disabled="!messageOptions[message.id].selectedContactId" @click.stop="forwardToContact(messageOptions[message.id].selectedContactId, message.id)">
                    Send
                  </button>
                  <button class="button-style" @click.stop="closeForwardMenu(message.id)">Cancel</button>
                </div>
                <div v-if="messageOptions[message.id].forwardConversations.length === 0">
                  No conversation found.
                </div>
              </div>
            </div>

            <div v-if="message.senderId === userToken" class="message-status">
              <span v-if="message.isRead" class="checkmark read" title="Read">✓✓</span>
              <span v-else-if="message.isDelivered" class="checkmark delivered" title="Delivered">✓✓</span>
              <span v-else class="checkmark" title="Sent">✓</span>
            </div>
          </div>
        </div>

        <div v-if="replyToMessage" class="reply-preview-box">
          <div class="reply-info">
            <strong>Replying to {{ replyToMessage.senderName || 'Unknown' }}:</strong>
            <span class="reply-text">{{ replyToMessage.content }}</span>
            <img v-if="replyToMessage.attachment" :src="'data:image/jpeg;base64,' + replyToMessage.attachment" alt="Reply Attachment" class="reply-attachment-preview" />
          </div>
          <button class="cancel-reply-button" @click="cancelReply">✖</button>
        </div>

        <div class="chat-input">
          <input type="file" ref="fileInput" style="display: none" accept="image/*, .gif" @change="handleFileSelect" />
          <button class="attach-button" @click="triggerFileInput">
            Attach Image or GIF
            <span v-if="selectedFile" class="file-icon">🖼️</span>
          </button>
          <input v-model="message" class="message-input" type="text" placeholder="Type a message..." @input="toggleSendButton" />
          <button v-if="message.trim() || selectedFile" class="send-button" @click="sendMessage" :disabled="sending">
            {{ sending ? 'Sending...' : 'Send' }}
          </button>
        </div>
      </div>

      <!-- Group Info Sidebar -->
      <div v-if="showGroupInfo && conversationType === 'group'" class="group-info-sidebar">
        <div class="sidebar-header">
          <h3>Members <span class="member-count">{{ groupMembers.length }}</span></h3>
          <button class="sidebar-close" @click="showGroupInfo = false">✕</button>
        </div>
        <ul class="members-list">
          <li v-for="member in groupMembers" :key="member.id" class="member-item">
            <div class="member-avatar-wrap">
              <img
                v-if="member.photo"
                :src="`data:image/jpeg;base64,${member.photo}`"
                alt="Photo"
                class="member-photo"
              />
              <div v-else class="member-initials">{{ member.name.slice(0,2).toUpperCase() }}</div>
            </div>
            <div class="member-info">
              <span class="member-name">{{ member.name }}</span>
              <span v-if="member.role === 'admin'" class="admin-badge">Admin</span>
            </div>
            <button
              v-if="isCurrentUserAdmin && member.id !== userToken"
              class="remove-member-btn"
              @click="removeMember(member)"
              title="Remove member"
            >✕</button>
          </li>
        </ul>

        <div class="sidebar-actions">
          <button @click="leaveGroup" class="leave-group-btn">Leave Group</button>
        </div>
        <ErrorMsg v-if="errormsg" :msg="errormsg" />
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "ChatView",
  components: { ErrorMsg },
  data() {
    return {
      message: "",
      messages: [],
      userToken: localStorage.getItem("token"),
      convName: localStorage.getItem("conversationName") || "Unknown User",
      conversationPhoto: null,
      conversationType: null,
      conversationId: this.$route.params.uuid,
      messageOptions: {},
      selectedFile: null,
      pollIntervalId: null,
      firstLoad: true,
      replyToMessage: null,
      showGroupInfo: false,
      groupMembers: [],
      errormsg: null,
      sending: false,
      defaultUserPhoto: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR4nGMAAQAABQABDQottAAAAABJRU5ErkJggg=="
    };
  },
  computed: {
    userName() {
      return localStorage.getItem("name");
    },
    userId() {
      return localStorage.getItem("userId");
    },
    isCurrentUserAdmin() {
      return this.groupMembers.some(m => m.id === this.userToken && m.role === "admin");
    },
  },
  methods: {
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    handleFileSelect(event) {
      this.selectedFile = event.target.files[0];
    },
    async sendMessage() {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
        return;
      }
      if (this.sending) return;
      this.sending = true;
      try {
        const formData = new FormData();
        formData.append("content", this.message);
        if (this.replyToMessage) {
          formData.append("replyTo", this.replyToMessage.id);
        }
        if (this.selectedFile) {
          formData.append("attachment", this.selectedFile);
        }
        await axios.post(`/conversations/${this.conversationId}/message`, formData, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.message = "";
        this.selectedFile = null;
        this.$refs.fileInput.value = "";
        this.replyToMessage = null;
        await this.fetchMessages();
        this.$nextTick(() => {
          this.forceScrollToBottom();
        });
      } catch (err) {
        console.error("Failed to send message:", err);
        this.errormsg = "Failed to send message. Please try again.";
      } finally {
        this.sending = false;
      }
    },
    async fetchMessages() {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
        return;
      }
      const response = await axios.get(`/conversations/${this.conversationId}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      this.messages = (response.data.messages || []).map(msg => ({
        ...msg,
        reactingUserNames: msg.reactingUserNames || [],
        showReactedList: false
      }));
      if (response.data.name) {
        this.convName = response.data.name;
      }
      if (response.data.conversationPhoto && response.data.conversationPhoto.String) {
        this.conversationPhoto = response.data.conversationPhoto.String;
      } else {
        this.conversationPhoto = null;
      }
      this.conversationType = response.data.type || "direct";
      
      // Fetch group members if it's a group
      if (this.conversationType === "group") {
        await this.fetchGroupMembers();
      }
      
      this.$nextTick(() => {
        if (this.firstLoad) {
          this.forceScrollToBottom();
          this.firstLoad = false;
        }
      });
    },
    async fetchGroupMembers() {
      try {
        const token = localStorage.getItem("token");
        const response = await axios.get(`/conversations/${this.conversationId}/members`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.groupMembers = response.data || [];
        console.log("Group members:", this.groupMembers);
      } catch (error) {
        console.error("Failed to load group members:", error);
      }
    },
    toggleGroupInfo() {
      this.showGroupInfo = !this.showGroupInfo;
    },
    async leaveGroup() {
      if (confirm("Are you sure you want to leave this group?")) {
        try {
          const token = localStorage.getItem("token");
          await axios.delete(`/conversations/${this.conversationId}/members/me`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.$router.push("/home");
        } catch (error) {
          this.errormsg = error.response?.data?.message || "Failed to leave group";
        }
      }
    },
    forceScrollToBottom() {
      const chat = this.$refs.chatMessages;
      if (chat) {
        chat.scrollTop = chat.scrollHeight;
      }
    },
    async toggleReaction(message) {
      const token = localStorage.getItem("token");
      if (!token || message.senderId === this.userToken) return;
      const hasReacted = (message.reactingUserNames || []).includes(this.userName);
      try {
        if (hasReacted) {
          await axios.delete(`/conversations/${this.conversationId}/message/${message.id}/comment`, {
            headers: { Authorization: `Bearer ${token}` }
          });
        } else {
          await axios.post(`/conversations/${this.conversationId}/message/${message.id}/comment`, {},
            { headers: { Authorization: `Bearer ${token}` } }
          );
        }
      } catch (err) {
        console.error("Error toggling reaction", err);
      } finally {
        await this.fetchMessages();
      }
    },
    async deleteMessage(message) {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
        return;
      }
      if (!confirm("Delete this message?")) return;
      try {
        await axios.delete(`/conversations/${this.conversationId}/message/${message.id}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.messages = this.messages.filter(m => m.id !== message.id);
      } catch (err) {
        console.error("Failed to delete message:", err);
        this.errormsg = "Failed to delete message.";
      }
    },
    formatTimestamp(timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleString();
    },
    showForwardOptions(messageId) {
      this.closeAllMenus();
      if (!this.messageOptions[messageId]) {
        this.messageOptions[messageId] = {
          showForwardMenu: true,
          forwardConversations: [],
          selectedConversationId: "",
          contactQuery: "",
          contactResults: [],
          selectedContactId: ""
        };
        this.fetchForwardConversations(messageId);
      } else {
        this.messageOptions[messageId].showForwardMenu = !this.messageOptions[messageId].showForwardMenu;
      }
    },
    closeForwardMenu(messageId) {
      if (this.messageOptions[messageId]) {
        this.messageOptions[messageId].showForwardMenu = false;
      }
    },
    closeAllMenus() {
      for (const id in this.messageOptions) {
        this.messageOptions[id].showForwardMenu = false;
      }
    },
    handleOutsideClick(event) {
      const messageContent = this.$el.querySelector('.message-content');
      if (messageContent && !messageContent.contains(event.target)) {
        this.closeAllMenus();
      }
    },
    async fetchForwardConversations(messageId) {
      const token = localStorage.getItem("token");
      const response = await axios.get('/conversations', {
        headers: { Authorization: `Bearer ${token}` }
      });
      const conversations = response.data.filter(conv => conv.id !== this.conversationId);
      this.messageOptions[messageId].forwardConversations = conversations;
    },
    async searchContact(messageId) {
      const query = this.messageOptions[messageId].contactQuery;
      if (!query.trim()) {
        this.messageOptions[messageId].contactResults = [];
        return;
      }
      const token = localStorage.getItem("token");
      const response = await axios.get('/search', {
        params: { username: query },
        headers: { Authorization: `Bearer ${token}` }
      });
      this.messageOptions[messageId].contactResults = response.data;
    },
    selectContact(contact, messageId) {
      this.messageOptions[messageId].selectedContactId = contact.id;
      this.messageOptions[messageId].contactQuery = contact.name;
      this.messageOptions[messageId].contactResults = [];
    },
    async forwardToContact(selectedContactId, messageId) {
      const token = localStorage.getItem("token");
      if (!token) {
        this.$router.push({ path: "/" });
        return;
      }
      const conversationResponse = await axios.post(
        `/conversations`,
        { senderId: token, recipientId: selectedContactId },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      const targetConversationId = conversationResponse.data.conversationId;
      const forwarderName = localStorage.getItem("name") || "Unknown";
      await axios.post(
        `/conversations/${this.conversationId}/message/${messageId}/forward`,
        { sourceMessageId: messageId, targetConversationId: targetConversationId, forwarderName: forwarderName },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert("Message forwarded successfully!");
      this.closeForwardMenu(messageId);
    },
    async forwardMessage(targetConversationId, messageId) {
      const message = this.messages.find(m => m.id === messageId);
      if (!message) return;
      const token = localStorage.getItem("token");
      const forwarderName = localStorage.getItem("name") || "Unknown";
      await axios.post(
        `/conversations/${this.conversationId}/message/${messageId}/forward`,
        { targetConversationId: targetConversationId, forwarderName: forwarderName },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      alert("Message forwarded successfully!");
      this.closeForwardMenu(messageId);
    },
    setReply(message) {
      this.replyToMessage = message;
    },
    cancelReply() {
      this.replyToMessage = null;
    },
    async removeMember(member) {
      if (!confirm(`Remove ${member.name} from the group?`)) return;
      const token = localStorage.getItem("token");
      try {
        await axios.delete(`/groups/${this.conversationId}/members/${member.id}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.groupMembers = this.groupMembers.filter(m => m.id !== member.id);
      } catch (err) {
        this.errormsg = err.response?.data || "Failed to remove member.";
      }
    }
  },
  mounted() {
    this.fetchMessages();
    this.pollIntervalId = setInterval(() => {
      this.fetchMessages();
    }, 5000);
    document.addEventListener("click", this.handleOutsideClick);
  },
  beforeUnmount() {
    document.removeEventListener("click", this.handleOutsideClick);
    clearInterval(this.pollIntervalId);
  }
};
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 92vh;
  overflow: hidden;
}
.chat-header {
  display: flex;
  align-items: center;
  padding: 14px 20px;
  background: #fff;
  border-bottom: 1px solid #e2e8f0;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}
.chat-photo {
  width: 40px;
  height: 40px;
  margin-right: 10px;
}
.chat-photo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}
.chat-messages {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #f7f8fc;
  border-top: 1px solid #e2e8f0;
  border-bottom: 1px solid #e2e8f0;
}
.message {
  position: relative;
  max-width: 70%;
  margin-bottom: 10px;
  border-radius: 14px;
  padding: 10px 14px;
  background: #f0f2f5;
  border-bottom-left-radius: 4px;
}
.message.self {
  margin-left: auto;
  background: #e8f0fe;
  border-bottom-left-radius: 14px;
  border-bottom-right-radius: 4px;
}
.sender-thumbnail {
  position: absolute;
  left: 10px;
  top: 10px;
  width: 30px;
  height: 30px;
}
.sender-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}
.message-content {
  position: relative;
  min-height: 40px;
}
.message p {
  margin: 0;
  color: #333;
  word-wrap: break-word;
  white-space: pre-wrap;
}
.message small {
  margin-top: 5px;
  color: #666;
  font-size: 0.8em;
}
.attachment-container {
  margin-top: 8px;
  width: 300px;
  height: 300px;
  overflow: hidden;
  border: 1px solid #ddd;
  border-radius: 8px;
}
.attachment-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.action-buttons {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
  justify-content: flex-end;
}
.message.self .action-buttons {
  justify-content: flex-start;
}
.action-button {
  position: static;
  width: 17px;
  height: 17px;
  border: 1px solid #aaa;
  border-radius: 50%;
  background-color: white;
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
  opacity: 0.9;
  transition: opacity 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  padding: 0;
}
.action-button:hover {
  opacity: 1;
}
.reply-button {
  font-size: 10px;
  margin-right: 5px;
}
.forward-options {
  position: absolute;
  top: 30px;
  right: 0;
  background-color: #ffffff;
  border-radius: 5px;
  padding: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  z-index: 100;
  width: 250px;
}
.forward-select {
  width: 100%;
  padding: 8px;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  margin-bottom: 8px;
  font-size: 14px;
}
.forward-buttons-container {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-top: 10px;
}
.button-style {
  background-color: #128c7e;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}
.button-style:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
.contact-search input {
  width: 100%;
  padding: 6px;
  margin-bottom: 4px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.contact-results {
  list-style: none;
  padding: 0;
  margin: 0 0 6px 0;
  max-height: 100px;
  overflow-y: auto;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.contact-result {
  padding: 4px;
  cursor: pointer;
  border-bottom: 1px solid #eee;
}
.contact-result:hover {
  background-color: #f0f0f0;
}
.file-icon {
  font-size: 18px;
  margin-left: 5px;
}
.message-status {
  position: absolute;
  bottom: 5px;
  right: 10px;
  font-size: 12px;
  color: #555;
}
.reactors-list ul {
  margin: 0;
  padding: 0;
  list-style: none;
  font-size: 0.8em;
  color: #444;
}
.reactors-list li {
  margin: 2px 0;
}
.reply-preview-box {
  background-color: #f0f0f0;
  border-left: 4px solid #128c7e;
  padding: 8px;
  margin: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.reply-info {
  font-size: 0.9em;
  color: #444;
}
.reply-attachment,
.reply-attachment-preview {
  width: 21px;
  height: 21px;
  object-fit: cover;
  margin-left: 10px;
  border-radius: 4px;
}
.cancel-reply-button {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #888;
}
.chat-input {
  display: flex;
  align-items: flex-start;
  flex-wrap: wrap;
  padding: 10px;
  gap: 8px;
}
.attach-button {
  background: #f0f2f5;
  color: #4a5568;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  padding: 10px 15px;
  transition: background 0.2s;
}
.attach-button:hover { background: #e2e8f0; }
.message-input {
  flex: 1;
  min-width: 200px;
  padding: 12px;
  border: 1px solid #dee2e6;
  border-radius: 20px;
  font-size: 14px;
  outline: none;
}
.send-button {
  background: #4c6ef5;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 20px;
  margin-left: 10px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: background 0.2s;
}
.send-button:hover { background: #3b5bdb; }
.checkmark {
  color: #2c8b47; /* green for sent */
  font-size: 1.2em;
  margin-left: 4px;
}
.checkmark.delivered {
  color: #2c8b47; /* green for delivered */
  font-weight: bold;
}
.checkmark.read {
  color: #2196f3; /* blue for read */
  font-weight: bold;
}

.info-button {
  background: #e8f0fe;
  color: #4c6ef5;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  margin-left: auto;
  font-size: 13px;
  font-weight: 600;
  transition: background 0.2s;
}
.info-button:hover { background: #c5d5fc; }

.group-info-sidebar {
  width: 280px;
  border-left: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  background: #fff;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 16px 12px;
  border-bottom: 1px solid #f0f2f5;
}
.sidebar-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 700;
  color: #1a1f36;
  display: flex;
  align-items: center;
  gap: 8px;
}
.member-count {
  background: #e8f0fe;
  color: #4c6ef5;
  font-size: 12px;
  border-radius: 12px;
  padding: 2px 8px;
  font-weight: 600;
}
.sidebar-close {
  background: none;
  border: none;
  font-size: 16px;
  color: #a0aec0;
  cursor: pointer;
  padding: 4px;
}
.sidebar-close:hover { color: #4a5568; }

.members-list {
  list-style: none;
  padding: 8px 0;
  margin: 0;
  flex: 1;
  overflow-y: auto;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
}
.member-item:hover { background: #f7f8fc; }

.member-avatar-wrap { flex-shrink: 0; }
.member-photo {
  width: 38px; height: 38px;
  border-radius: 50%; object-fit: cover;
}
.member-initials {
  width: 38px; height: 38px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4c6ef5, #7c3aed);
  color: #fff; display: flex; align-items: center; justify-content: center;
  font-size: 13px; font-weight: 700;
}

.member-info { flex: 1; min-width: 0; }
.member-name { font-size: 14px; font-weight: 500; color: #1a1f36; display: block; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.admin-badge {
  display: inline-block;
  font-size: 10px;
  background: #fff3cd;
  color: #856404;
  border-radius: 4px;
  padding: 1px 6px;
  font-weight: 600;
  margin-top: 2px;
}

.remove-member-btn {
  background: none; border: none;
  color: #cbd5e0; font-size: 14px; cursor: pointer; padding: 4px;
}
.remove-member-btn:hover { color: #e53e3e; }

.sidebar-actions {
  padding: 12px 16px;
  border-top: 1px solid #f0f2f5;
}

.leave-group-btn {
  width: 100%;
  background: #fff5f5;
  color: #c53030;
  border: 1px solid #fed7d7;
  padding: 10px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.leave-group-btn:hover { background: #fed7d7; }

@media (max-width: 767px) {
  .chat-container { height: calc(100vh - 68px); }

  .group-info-sidebar {
    position: fixed;
    top: 0; bottom: 68px;
    right: 0;
    width: 85%;
    max-width: 320px;
    box-shadow: -4px 0 20px rgba(0,0,0,0.15);
    z-index: 50;
  }

  .message { max-width: 88%; }

  .attachment-container {
    width: 100%;
    max-width: 220px;
    height: auto;
    aspect-ratio: 1;
  }

  .message-input {
    min-width: 80px;
  }

  .attach-button {
    font-size: 12px;
    padding: 8px 10px;
  }

  .chat-input {
    padding: 8px;
    gap: 6px;
  }

  .send-button {
    padding: 10px 16px;
    font-size: 13px;
    margin-left: 0;
  }
}
</style>
