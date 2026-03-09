<script>
export default {
  data() {
    localStorage.removeItem("token");
    localStorage.removeItem("name");
    localStorage.removeItem("conversationName");
    localStorage.removeItem("recipientId");
    return {
      mode: "signin",
      name: "",
      password: "",
      confirmPassword: "",
      errormsg: null,
      loading: false,
    };
  },
  methods: {
    async loadDefaultPhoto() {
      try {
        const response = await fetch('/nopfp.jpg');
        const blob = await response.blob();
        return new Promise((resolve) => {
          const reader = new FileReader();
          reader.onload = () => resolve(reader.result.toString().split(',')[1]);
          reader.readAsDataURL(blob);
        });
      } catch {
        return null;
      }
    },
    switchMode(m) {
      this.mode = m;
      this.errormsg = null;
      this.password = "";
      this.confirmPassword = "";
    },
    async submit() {
      this.errormsg = null;
      if (!this.name.trim()) { this.errormsg = "Username is required."; return; }
      if (this.name.length < 3 || this.name.length > 16) { this.errormsg = "Username must be 3–16 characters."; return; }
      if (!this.password) { this.errormsg = "Password is required."; return; }
      if (this.mode === "signup" && this.password !== this.confirmPassword) {
        this.errormsg = "Passwords do not match."; return;
      }
      this.loading = true;
      try {
        const photoData = this.mode === "signup" ? await this.loadDefaultPhoto() : null;
        const response = await this.$axios.post("/session", {
          name: this.name,
          password: this.password,
          photo: photoData,
          mode: this.mode,
        });
        localStorage.setItem("token", response.data.identifier);
        localStorage.setItem("name", this.name);
        this.$router.push({ path: "/home" });
      } catch (e) {
        const status = e.response?.status;
        if (status === 404) this.errormsg = "User not found. Did you mean to Sign Up?";
        else if (status === 401) this.errormsg = "Incorrect password.";
        else if (status === 409) this.errormsg = "Username already taken.";
        else if (status === 400) this.errormsg = e.response?.data || "Invalid input.";
        else this.errormsg = "Something went wrong. Please try again.";
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<template>
  <div class="auth-root">
    <div class="auth-left">
      <div class="brand">
        <div class="brand-icon">💬</div>
        <h1 class="brand-name">WASAText</h1>
        <p class="brand-tagline">Simple, fast, and secure messaging.</p>
      </div>
    </div>
    <div class="auth-right">
      <div class="auth-card">
        <h2 class="auth-heading">{{ mode === 'signin' ? 'Welcome back' : 'Create account' }}</h2>
        <p class="auth-sub">{{ mode === 'signin' ? 'Sign in to continue' : 'Join WASAText today' }}</p>

        <div class="tab-row">
          <button :class="['tab-btn', mode === 'signin' ? 'active' : '']" @click="switchMode('signin')">Sign In</button>
          <button :class="['tab-btn', mode === 'signup' ? 'active' : '']" @click="switchMode('signup')">Sign Up</button>
        </div>

        <div v-if="errormsg" class="err-banner">{{ errormsg }}</div>

        <form @submit.prevent="submit" class="auth-form">
          <label class="field-label">Username</label>
          <input v-model="name" class="field-input" type="text" placeholder="3–16 characters" autocomplete="username" />

          <label class="field-label">Password</label>
          <input v-model="password" class="field-input" type="password" placeholder="Enter password" autocomplete="current-password" />

          <template v-if="mode === 'signup'">
            <label class="field-label">Confirm Password</label>
            <input v-model="confirmPassword" class="field-input" type="password" placeholder="Repeat password" autocomplete="new-password" />
          </template>

          <button class="submit-btn" type="submit" :disabled="loading">
            <span v-if="loading">Please wait…</span>
            <span v-else>{{ mode === 'signin' ? 'Sign In' : 'Create Account' }}</span>
          </button>
        </form>

        <p class="switch-hint">
          {{ mode === 'signin' ? "Don't have an account?" : "Already have an account?" }}
          <a class="switch-link" @click="switchMode(mode === 'signin' ? 'signup' : 'signin')">
            {{ mode === 'signin' ? 'Sign Up' : 'Sign In' }}
          </a>
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-root {
  display: flex;
  min-height: 100vh;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

/* Left branding panel */
.auth-left {
  flex: 1;
  background: linear-gradient(135deg, #1a1f36 0%, #2d3461 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}
.brand { text-align: center; color: #fff; }
.brand-icon { font-size: 64px; margin-bottom: 16px; }
.brand-name { font-size: 42px; font-weight: 700; letter-spacing: -1px; margin: 0 0 12px; }
.brand-tagline { font-size: 18px; color: #a0aec0; margin: 0; }

/* Right form panel */
.auth-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7f8fc;
  padding: 40px 20px;
}
.auth-card {
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 16px;
  padding: 40px 36px;
  box-shadow: 0 4px 32px rgba(0,0,0,0.08);
}
.auth-heading {
  font-size: 26px;
  font-weight: 700;
  color: #1a1f36;
  margin: 0 0 4px;
}
.auth-sub {
  font-size: 14px;
  color: #718096;
  margin: 0 0 24px;
}

/* Tabs */
.tab-row {
  display: flex;
  background: #f0f2f5;
  border-radius: 10px;
  padding: 4px;
  margin-bottom: 24px;
}
.tab-btn {
  flex: 1;
  padding: 8px;
  border: none;
  background: transparent;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #718096;
  cursor: pointer;
  transition: all 0.2s;
}
.tab-btn.active {
  background: #fff;
  color: #1a1f36;
  box-shadow: 0 1px 4px rgba(0,0,0,0.12);
}

/* Error */
.err-banner {
  background: #fff5f5;
  border: 1px solid #fed7d7;
  color: #c53030;
  border-radius: 8px;
  padding: 10px 14px;
  font-size: 13px;
  margin-bottom: 16px;
}

/* Form */
.auth-form { display: flex; flex-direction: column; gap: 6px; }
.field-label {
  font-size: 13px;
  font-weight: 600;
  color: #4a5568;
  margin-top: 10px;
}
.field-input {
  padding: 11px 14px;
  border: 1.5px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  color: #1a1f36;
  outline: none;
  transition: border-color 0.2s;
}
.field-input:focus { border-color: #4c6ef5; }

.submit-btn {
  margin-top: 20px;
  padding: 13px;
  background: #4c6ef5;
  color: #fff;
  border: none;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.submit-btn:hover:not(:disabled) { background: #3b5bdb; }
.submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* Switch hint */
.switch-hint {
  text-align: center;
  margin-top: 20px;
  font-size: 13px;
  color: #718096;
}
.switch-link {
  color: #4c6ef5;
  font-weight: 600;
  cursor: pointer;
  margin-left: 4px;
}
.switch-link:hover { text-decoration: underline; }

/* Responsive: stack on mobile */
@media (max-width: 700px) {
  .auth-root { flex-direction: column; }
  .auth-left { padding: 40px 20px; }
  .brand-name { font-size: 28px; }
}
</style>
