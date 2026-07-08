<template>
  <div class="home">
    <h1>Yappari</h1>
    <div class="deck-list">
      <div v-for="deck in decks" :key="deck.id" class="deck-item">
        <input type="checkbox" :id="`deck-${deck.id}`" v-model="selectedDeckIds" :value="deck.id" />
        <label :for="`deck-${deck.id}`">{{ deck.name }}</label>
        <button @click="goToDeck(deck.id)">Управлять</button>
      </div>
    </div>
    <div class="actions">
      <button @click="createDeckModalVisible = true">Новая колода</button>
    </div>
    <div class="training-buttons">
      <button @click="startTraining('interval')" :disabled="selectedDeckIds.length === 0">
        Интервальное повторение
      </button>
      <button @click="startTraining('normal')" :disabled="selectedDeckIds.length === 0">
        Обычная зубрёжка
      </button>
      <button @click="startTraining('lazy')" :disabled="selectedDeckIds.length === 0">
        Ленивое заучивание
      </button>
    </div>

    <Dialog v-model:visible="createDeckModalVisible" header="Новая колода">
      <InputText v-model="newDeckName" placeholder="Название колоды" />
      <template #footer>
        <Button label="Отмена" @click="createDeckModalVisible = false" />
        <Button label="Создать" @click="createDeck" />
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

interface Deck {
  id: number
  name: string
  created_at: string
}

const router = useRouter()
const decks = ref<Deck[]>([])
const selectedDeckIds = ref<number[]>([])
const createDeckModalVisible = ref(false)
const newDeckName = ref('')

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
    // @ts-ignore
    decks.value = await window.go.main.App.GetDecks()
  } catch (e) {
    console.error(e)
  }
}

const createDeck = async () => {
  if (!newDeckName.value.trim()) return
  try {
    // @ts-ignore
    await window.go.main.App.CreateDeck(newDeckName.value)
    createDeckModalVisible.value = false
    newDeckName.value = ''
    await loadDecks()
  } catch (e) {
    console.error(e)
  }
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
}

.deck-list {
  margin: 2rem 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.deck-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem;
  background: #1a1a1a;
  border-radius: 0.5rem;
}

.deck-item label {
  flex: 1;
}

.actions {
  margin: 1rem 0;
}

.training-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}
</style>
