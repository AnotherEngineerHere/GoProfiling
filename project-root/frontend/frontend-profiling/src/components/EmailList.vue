<template>
  <div class="email-list">
    <div class="email-items">
      <div v-for="email in emails" 
           :key="email._id" 
           @click="$emit('select', email)"
           class="email-item">
        <div class="email-item-content">
          <div class="email-item-header">
            <p class="email-subject">{{ email._source.subject }}</p>
            <span class="email-date">{{ formatDate(email._source['@timestamp']) }}</span>
          </div>
          <p class="email-sender">From: {{ email._source.sender }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { formatDate } from '../utils/formatters'

export default {
  name: 'EmailList',
  props: {
    emails: {
      type: Array,
      required: true
    }
  },
  emits: ['select'],
  setup() {
    return {
      formatDate
    }
  }
}
</script>

<style scoped>
.email-list {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.email-item {
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
  cursor: pointer;
  transition: background-color 0.2s;
}

.email-item:hover {
  background-color: #f8fafc;
}

.email-item-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.5rem;
}

.email-subject {
  font-size: 0.875rem;
  font-weight: 500;
  color: #00003e;
}

.email-date {
  font-size: 0.75rem;
  color: #718096;
}

.email-sender {
  font-size: 0.75rem;
  color: #718096;
}
</style>