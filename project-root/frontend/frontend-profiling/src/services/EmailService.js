import api from './api'

export default {
  // Obtener todos los emails
  async getAll() {
    try {
      const response = await api.get('/emails')
      return response.data.hits.hits // Extraer los emails de la estructura correcta
    } catch (error) {
      console.error('Error fetching emails:', error)
      throw error
    }
  },
  // Buscar emails
  async search(query) {
    try {
      const response = await api.get(`/search?q=${encodeURIComponent(query)}`)
      return response.data.hits.hits // Extraer los emails de la estructura correcta
    } catch (error) {
      console.error('Error searching emails:', error)
      throw error
    }
  },

  // Obtener un email por ID
  async getById(id) {
    try {
      const response = await api.get(`/emails/${id}`)
      return response.data
    } catch (error) {
      console.error('Error fetching email:', error)
      throw error
    }
  },

  // Crear un nuevo email
  async create(emailData) {
    try {
      const response = await api.post('/emails', emailData)
      return response.data
    } catch (error) {
      console.error('Error creating email:', error)
      throw error
    }
  },

  // Actualizar un email
  async update(id, emailData) {
    try {
      const response = await api.put(`/emails/${id}`, emailData)
      return response.data
    } catch (error) {
      console.error('Error updating email:', error)
      throw error
    }
  },

  // Eliminar un email
  async delete(id) {
    try {
      const response = await api.delete(`/emails/${id}`)
      return response.data
    } catch (error) {
      console.error('Error deleting email:', error)
      throw error
    }
  }
}