<template>
  <div class="home">
    <div class="header">
      <img src="/yappari_logo.png" alt="Yappari Logo" class="logo" draggable="false" />
      <h1>Yappari</h1>
    </div>

    <div v-if="decks.length === 0" class="empty-state">
      <p>У тебя пока нет колод. Создай первую!</p>
    </div>

    <div v-else class="deck-list">
      <div v-for="deck in decks" :key="deck.ID" class="deck-item">
        <input type="checkbox" :id="`deck-${deck.ID}`" v-model="selectedDeckIds" :value="deck.ID" />
        <label :for="`deck-${deck.ID}`">{{ deck.Name }}</label>
        <button @click="goToDeck(deck.ID)" class="icon-btn" title="Управлять колодой">
          <span class="icon">⚙️</span>
        </button>
      </div>
    </div>

    <div class="actions">
      <button @click="createDeckModalVisible = true" class="primary-btn">
        <span class="icon">+</span>
        Новая колода
      </button>
    </div>

    <div v-if="decks.length > 0" class="training-buttons">
      <button @click="startTraining('interval')" :disabled="selectedDeckIds.length === 0" class="training-btn">
        <span class="icon">🔄</span>
        Повторение
      </button>
      <button @click="startTraining('free')" :disabled="selectedDeckIds.length === 0" class="training-btn">
        <span class="icon">📖</span>
        Свободный режим
      </button>
    </div>

    <Dialog v-model:visible="createDeckModalVisible" header="Создать колоду" class="custom-dialog">
      <div class="form-content">
        <div class="input-group">
          <label for="deck-name">Название колоды</label>
          <InputText 
            id="deck-name"
            v-model="newDeckName" 
            placeholder="Например: Н5 слова" 
            class="custom-input"
            @keyup.enter="createDeck"
          />
          <div v-if="errors.deckName" class="error">{{ errors.deckName }}</div>
        </div>
      </div>
      <template #footer>
        <Button label="Отмена" @click="closeCreateModal" class="secondary-btn" />
        <Button label="Создать" @click="createDeck" class="primary-btn" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import type { Deck } from '../types'
import { useWails } from '../composables/useWails'

const router = useRouter()
const { getDecks, createDeck: createDeckWails } = useWails()
const decks = ref<Deck[]>([])
const selectedDeckIds = ref<number[]>([])
const createDeckModalVisible = ref(false)
const newDeckName = ref('')
const isLoading = ref(false)
const errors = ref({ deckName: '' })

onMounted(async () => {
  const saved = localStorage.getItem('selectedDeckIds')
  if (saved) {
    selectedDeckIds.value = JSON.parse(saved)
  }
  await loadDecks()
})

watch(selectedDeckIds, (newVal) => {
  localStorage.setItem('selectedDeckIds', JSON.stringify(newVal))
})

const loadDecks = async () => {
  try {
    isLoading.value = true
    decks.value = await getDecks()
  } catch (e) {
    console.error('Ошибка загрузки колод:', e)
    alert('Не удалось загрузить колоды: ' + e)
  } finally {
    isLoading.value = false
  }
}

const validateForm = (): boolean => {
  errors.value = { deckName: '' }
  if (!newDeckName.value.trim()) {
    errors.value.deckName = 'Название колоды не может быть пустым'
    return false
  }
  if (newDeckName.value.trim().length < 2) {
    errors.value.deckName = 'Название колоды должно быть не менее 2 символов'
    return false
  }
  return true
}

const createDeck = async () => {
  if (!validateForm()) return
  
  try {
    isLoading.value = true
    await createDeckWails(newDeckName.value.trim())
    closeCreateModal()
    await loadDecks()
  } catch (e) {
    console.error('Ошибка создания колоды:', e)
    alert('Не удалось создать колоду: ' + e)
  } finally {
    isLoading.value = false
  }
}

const closeCreateModal = () => {
  createDeckModalVisible.value = false
  newDeckName.value = ''
  errors.value = { deckName: '' }
}

const goToDeck = (id: number) => {
  router.push({ name: 'DeckManage', params: { id } })
}

const startTraining = (mode: string) => {
  router.push({ name: 'Training', query: { mode, deckIds: selectedDeckIds.value.join(',') } })
}
</script>

<style scoped>
.home {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.logo {
  width: 64px;
  height: 64px;
  user-select: none;
  -webkit-user-drag: none;
  pointer-events: none;
}

.header h1 {
  margin: 0;
  font-size: 2.5rem;
  font-weight: 700;
  user-select: none;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #c7cdd8;
  font-size: 1.1rem;
}

.deck-list {
  margin: 2rem 0;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.deck-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background-color: #111111;
  border: 1px solid #222222;
  border-radius: 0.75rem;
  transition: all 0.2s;
}

.deck-item:hover {
  border-color: #333333;
  background-color: #151515;
}

.deck-item label {
  flex: 1;
  cursor: pointer;
  font-size: 1.1rem;
  user-select: none;
}

.deck-item input[type="checkbox"] {
  width: auto;
  transform: scale(1.3);
  cursor: pointer;
}

.icon-btn {
  background: transparent;
  border: 1px solid #333333;
  color: #c7cdd8;
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-btn:hover {
  background: #1a1a1a;
  border-color: #ff0a14;
  color: #ffffff;
}

.icon {
  font-size: 1.25rem;
}

.actions {
  margin: 1.5rem 0;
}

.primary-btn {
  background-color: #ff0a14;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  cursor: pointer;
  font-family: inherit;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.primary-btn:hover:not(:disabled) {
  background-color: #e00912;
  transform: translateY(-1px);
}

.primary-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.secondary-btn {
  background-color: #222222;
  color: white;
  border: 1px solid #333333;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  cursor: pointer;
  font-family: inherit;
  font-size: 1rem;
  transition: all 0.2s;
}

.secondary-btn:hover {
  background-color: #333333;
}

.training-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.training-btn {
  flex: 1;
  min-width: 200px;
  padding: 1rem 1.5rem;
  font-size: 1.05rem;
}

.form-content {
  padding: 1rem 0;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.input-group label {
  font-size: 0.9rem;
  color: #c7cdd8;
  font-weight: 500;
}

.custom-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  border: 1px solid #333333;
  background: #111111;
  color: white;
  font-family: inherit;
  font-size: 1rem;
  outline: none;
  transition: all 0.2s;
}

.custom-input:focus {
  border-color: #ff0a14;
}

.error {
  color: #ff4444;
  font-size: 0.85rem;
  margin-top: 0.25rem;
}
</style>

<style>
.custom-dialog .p-dialog {
  background: #111111 !important;
  border: 1px solid #222222 !important;
  border-radius: 1rem !important;
}

.custom-dialog .p-dialog-header {
  background: #111111 !important;
  border-bottom: 1px solid #222222 !important;
  color: white !important;
  border-radius: 1rem 1rem 0 0 !important;
  padding: 1.5rem !important;
}

.custom-dialog .p-dialog-content {
  background: #111111 !important;
  color: white !important;
  padding: 1.5rem !important;
}

.custom-dialog .p-dialog-footer {
  background: #111111 !important;
  border-top: 1px solid #222222 !important;
  border-radius: 0 0 1rem 1rem !important;
  padding: 1.5rem !important;
}

.custom-dialog .p-dialog-title {
  font-size: 1.25rem !important;
  font-weight: 600 !important;
}
</style>
