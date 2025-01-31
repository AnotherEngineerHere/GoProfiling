<template>
  <div class="app">
    <NavigationBar />
    
    <main class="main-content">
      <SearchBar 
        v-model="searchQuery"
        @update:modelValue="handleSearch"
      />

      <div v-if="loading" class="loading">
        <div class="loader"></div>
      </div>

      <div v-else-if="error" class="error">
        {{ error }}
      </div>

      <div v-else class="email-container">
        <EmailList 
          :emails="emails"
          @select="selectEmail"
        />
        <EmailDetail 
          :email="selectedEmail"
        />
      </div>
    </main>
  </div>
</template>

<script>
import { ref } from 'vue'
import NavigationBar from './components/NavigationBar.vue'
import SearchBar from './components/SearchBar.vue'
import EmailList from './components/EmailList.vue'
import EmailDetail from './components/EmailDetail.vue'
import emailService from './services/EmailService'

export default {
  name: 'App',
  components: {
    NavigationBar,
    SearchBar,
    EmailList,
    EmailDetail
  },
  setup() {
    const emails = ref([])
    const selectedEmail = ref(null)
    const searchQuery = ref('')
    const loading = ref(false)
    const error = ref(null)

    const loadEmails = async () => {
      try {
        loading.value = true
        const response = await emailService.getAll()
        emails.value = response
      } catch (err) {
        error.value = 'Error al cargar los emails'
        console.error(err)
      } finally {
        loading.value = false
      }
    }

    const selectEmail = (email) => {
      selectedEmail.value = email
    }

    const handleSearch = async () => {
      try {
        loading.value = true
        if (searchQuery.value.trim()) {
          const response = await emailService.search(searchQuery.value)
          emails.value = response
        } else {
          await loadEmails()
        }
      } catch (err) {
        error.value = 'Error en la búsqueda'
        console.error(err)
      } finally {
        loading.value = false
      }
    }

    // Cargar emails al iniciar
    loadEmails()

    return {
      emails,
      selectedEmail,
      searchQuery,
      loading,
      error,
      selectEmail,
      handleSearch
    }
  }
}
</script>

<style>
/* Estilos globales */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: Arial, sans-serif;
  background-color: #f7f8fa;
  color: #2d3748;
}

.app {
  min-height: 100vh;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.email-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.loading {
  text-align: center;
  padding: 2rem;
}

.loader {
  display: inline-block;
  width: 2rem;
  height: 2rem;
  border: 2px solid #e2e8f0;
  border-radius: 50%;
  border-top-color: #00003e;
  animation: spin 1s linear infinite;
}

.error {
  background-color: #fff5f5;
  border-left: 4px solid #f56565;
  padding: 1rem;
  margin-bottom: 1.5rem;
  color: #c53030;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .email-container {
    grid-template-columns: 1fr;
  }
}
</style>