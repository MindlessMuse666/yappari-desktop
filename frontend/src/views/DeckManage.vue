<template>
  <div class="deck-manage">
    <div class="header">
      <button @click="goBack">← Назад</button>
      <InputText v-model="deckName" @keyup.enter="updateDeckName" />
      <button @click="resetDeckProgress">Сбросить прогресс</button>
      <button @click="deleteDeck" class="danger">Удалить колоду</button>
    </div>
    <div class="add-card">
      <Button label="Добавить карточку" @click="cardFormModalVisible = true" />
    </div>
    <DataTable :value="cards" paginator :rows="10">
      <Column field="kanjiText" header="Японский" />
      <Column field="translation" header="Перевод" />
      <Column header="Действия">
        <template #body="slotProps">
          <Button label="Редактировать" @click="editCard(slotProps.data)" />
          <Button label="Удалить" @click="deleteCard(slotProps.data.id)" />
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="cardFormModalVisible" :header="editingCard ? 'Редактировать карточку' : 'Новая карточка'">
      <div class="form">
        <label>Японское слово *</label>
        <InputText v-model="cardForm.kanjiText" />
        <label>Чтение</label>
        <InputText v-model="cardForm.furiganaText" />
        <label>Перевод *</label>
        <InputText v-model="cardForm.translation" />
      </div>
      <template #footer>
        <Button label="Отмена" @click="cardFormModalVisible = false" />
        <Button label="Сохранить" @click="saveCard" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

interface Card {
  id: number
  deck_id: number
  kanjiText: string
  furiganaText?: string | null
  translation: string
  easeFactor: number
  interval: number
  repetitions: number
  nextReview: string
  lastReview?: string | null
  createdAt: string
  updatedAt: string
}

interface Deck {
  id: number
  name: string
  created_at: string
}

const router = useRouter()
const route = useRoute()
const deckId = Number(route.params.id)
const deckName = ref('')
const cards = ref<Card[]>([])
const cardFormModalVisible = ref(false)
const editingCard = ref<Card | null>(null)
const cardForm = ref({
  kanjiText: '',
  furiganaText: '',
  translation: ''
})

onMounted(async () => {
  await loadDeck()
  await loadCards()
})

const loadDeck = async () => {
  try {
    // @ts-ignore
    const decks = await window.go.main.App.GetDecks()
    const deck = decks.find((d: Deck) => d.id === deckId)
    if (deck) {
      deckName.value = deck.name
    }
  } catch (e) {
    console.error(e)
  }
}

const loadCards = async () => {
  try {
    // @ts-ignore
    cards.value = await window.go.main.App.GetCardsByDeck(deckId)
  } catch (e) {
    console.error(e)
  }
}

const goBack = () => {
  router.push({ name: 'Home' })
}

const updateDeckName = async () => {
  try {
    // @ts-ignore
    await window.go.main.App.UpdateDeck(deckId, deckName.value)
  } catch (e) {
    console.error(e)
  }
}

const resetDeckProgress = async () => {
  try {
    // @ts-ignore
    await window.go.main.App.ResetDeckProgress(deckId)
  } catch (e) {
    console.error(e)
  }
}

const deleteDeck = async () => {
  if (!confirm('Удалить колоду и все карточки?')) return
  try {
    // @ts-ignore
    await window.go.main.App.DeleteDeck(deckId)
    router.push({ name: 'Home' })
  } catch (e) {
    console.error(e)
  }
}

const editCard = (card: Card) => {
  editingCard.value = card
  cardForm.value = {
    kanjiText: card.kanjiText,
    furiganaText: card.furiganaText || '',
    translation: card.translation
  }
  cardFormModalVisible.value = true
}

const deleteCard = async (id: number) => {
  if (!confirm('Удалить карточку?')) return
  try {
    // @ts-ignore
    await window.go.main.App.DeleteCard(id)
    await loadCards()
  } catch (e) {
    console.error(e)
  }
}

const saveCard = async () => {
  if (!cardForm.value.kanjiText.trim() || !cardForm.value.translation.trim()) return
  try {
    if (editingCard.value) {
      // @ts-ignore
      await window.go.main.App.UpdateCard(editingCard.value.id, {
        deck_id: deckId,
        kanji_text: cardForm.value.kanjiText,
        furigana_text: cardForm.value.furiganaText || null,
        translation: cardForm.value.translation
      })
    } else {
      // @ts-ignore
      await window.go.main.App.CreateCard({
        deck_id: deckId,
        kanji_text: cardForm.value.kanjiText,
        furigana_text: cardForm.value.furiganaText || null,
        translation: cardForm.value.translation
      })
    }
    cardFormModalVisible.value = false
    editingCard.value = null
    cardForm.value = { kanjiText: '', furiganaText: '', translation: '' }
    await loadCards()
  } catch (e) {
    console.error(e)
  }
}
</script>

<style scoped>
.deck-manage {
  padding: 2rem;
}

.header {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-bottom: 2rem;
}

.add-card {
  margin-bottom: 2rem;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.danger {
  background: #ff0a14;
  border-color: #ff0a14;
}
</style>
