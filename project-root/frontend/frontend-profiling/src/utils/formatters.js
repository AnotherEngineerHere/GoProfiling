export const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString()
  }
  
  export const formatContent = (content) => {
    return content.replace(/\r\n/g, '\n').trim()
  }