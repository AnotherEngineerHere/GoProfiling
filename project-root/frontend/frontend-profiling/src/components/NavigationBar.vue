<template>
  <nav class="navbar">
    <div class="navbar-container">
      <h1 class="navbar-title">Email Dashboard</h1>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'NavigationBar'
}
</script>

<style scoped>
.navbar {
  background-color: #00003e;
  color: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.navbar-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
  height: 64px;
  display: flex;
  align-items: center;
}

.navbar-title {
  font-size: 1.25rem;
  font-weight: 600;
}
</style>