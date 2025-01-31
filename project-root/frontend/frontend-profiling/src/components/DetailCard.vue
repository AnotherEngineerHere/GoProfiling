<template>
    <div class="detail-card" v-if="selectedItem">
      <h3>{{ selectedItem.title }}</h3>
      <p>{{ selectedItem.content }}</p>
    </div>
  </template>
  
  <script>
  export default {
    name: 'DetailCard',
    props: {
      selectedItem: {
        type: Object,
        required: true
      }
    }
  }
  </script>