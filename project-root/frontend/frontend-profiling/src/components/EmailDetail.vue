<template>
  <div class="email-detail">
    <div v-if="email?._source" class="email-content">
      <div class="email-header">
        <h2 class="email-title">{{ email._source.subject }}</h2>
        <div class="email-meta">
          <p class="meta-item">
            <span class="meta-label">From:</span> 
            {{ email._source.sender }}
          </p>
          <p class="meta-item">
            <span class="meta-label">To:</span> 
            {{ email._source.recipient }}
          </p>
          <p class="meta-item">
            <span class="meta-label">Date:</span> 
            {{ formatDate(email._source['@timestamp']) }}
          </p>
        </div>
      </div>
      <div class="email-body">
        <pre>{{ formatContent(email._source.content) }}</pre>
      </div>
    </div>
    <div v-else class="no-selection">
      Selecciona un email para ver los detalles
    </div>
  </div>
</template>

<script>
import { formatDate, formatContent } from '../utils/formatters'

export default {
  name: 'EmailDetail',
  props: {
    email: {
      type: Object,
      default: null
    }
  },
  setup() {
    return {
      formatDate,
      formatContent
    }
  }
}
</script>

<style scoped>
.email-detail {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  height: 100%;
}

.email-content {
  padding: 1.5rem;
}

.email-header {
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 1rem;
  margin-bottom: 1rem;
}

.email-title {
  font-size: 1.25rem;
  font-weight: 500;
  color: #00003e;
  margin-bottom: 1rem;
}

.email-meta {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.meta-item {
  font-size: 0.875rem;
  color: #4a5568;
}

.meta-label {
  font-weight: 500;
  margin-right: 0.5rem;
}

.email-body {
  font-size: 0.875rem;
  color: #4a5568;
  line-height: 1.5;
}

.email-body pre {
  white-space: pre-wrap;
  font-family: inherit;
}

.no-selection {
  padding: 1.5rem;
  text-align: center;
  color: #718096;
}
</style>