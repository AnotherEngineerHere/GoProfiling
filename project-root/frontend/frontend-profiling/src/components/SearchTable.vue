<template>
    <div class="search-table">
      <input 
        v-model="searchQuery" 
        type="text" 
        placeholder="Buscar..."
        class="search-input"
      >
      <table>
        <thead>
          <tr>
            <th>Título</th>
            <th>Fecha</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in filteredItems" 
              :key="item.id"
              @click="$emit('select-item', item)">
            <td>{{ item.title }}</td>
            <td>{{ item.date }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </template>
  
  <script>
  export default {
    name: 'SearchTable',
    data() {
      return {
        searchQuery: ''
      }
    },
    props: {
      items: {
        type: Array,
        required: true
      }
    },
    computed: {
      filteredItems() {
        return this.items.filter(item => {
          const search = this.searchQuery.toLowerCase()
          return Object.values(item).some(value => 
            String(value).toLowerCase().includes(search)
          )
        })
      }
    }
  }
  </script>