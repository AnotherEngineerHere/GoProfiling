<template>
  <div class="search-container">
    <div class="search-wrapper">
      <input 
        :value="modelValue"
        @input="$emit('update:modelValue', $event.target.value)"
        type="text" 
        placeholder="Buscar emails..."
        class="search-input"
      >
      <span class="search-icon">🔍</span>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SearchBar',
  props: {
    modelValue: String
  },
  emits: ['update:modelValue']
}
</script>

<style scoped>
.search-container {
  margin-bottom: 1.5rem;
}

.search-wrapper {
  position: relative;
}

.search-input {
  width: 100%;
  padding: 0.5rem 2.5rem 0.5rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 0.375rem;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: #00003e;
}

.search-icon {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #718096;
}
</style>