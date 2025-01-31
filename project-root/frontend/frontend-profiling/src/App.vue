<template>
  <div class="app">
    <div class="container">
      <SearchTable 
        :items="items"
        @select-item="selectItem"
      />
      <DetailCard 
        :selectedItem="selectedItem"
      />
    </div>
  </div>
</template>

<script>
import SearchTable from './components/SearchTable.vue'
import DetailCard from './components/DetailCard.vue'

export default {
  name: 'App',
  components: {
    SearchTable,
    DetailCard
  },
  data() {
    return {
      items: [
        // Aquí van tus datos
        {id: 1, title: 'Item 1', content: 'Contenido 1', date: '2023-01-01'},
        {id: 2, title: 'Item 2', content: 'Contenido 2', date: '2023-01-02'},
      ],
      selectedItem: null
    }
  },
  methods: {
    selectItem(item) {
      this.selectedItem = item
    }
  }
}
</script>

<style>
.container {
  display: flex;
  gap: 20px;
  padding: 20px;
}

.search-table {
  flex: 1;
}

.detail-card {
  flex: 1;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.search-input {
  width: 100%;
  padding: 8px;
  margin-bottom: 20px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 8px;
  border: 1px solid #ddd;
  text-align: left;
}

tr:hover {
  background-color: #f5f5f5;
  cursor: pointer;
}
</style>