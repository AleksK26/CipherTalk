<template>
  <div class="not-found-container">
    <h1>404 - Page Not Found</h1>
    <p>The page you are looking for does not exist.</p>
    <router-link to="/" class="btn btn-primary">Go Home</router-link>
  </div>
</template>

<script>
export default {
  name: "NotFoundView"
};
</script>

<style scoped>
.not-found-container {
  text-align: center;
  margin-top: 100px;
}
.not-found-container h1 {
  font-size: 3rem;
  color: #dc3545;
}
.not-found-container p {
  font-size: 1.25rem;
  margin-bottom: 2rem;
}
</style>
