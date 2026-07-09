<template>
  <div class="deck-manage">
    <div class="header">
      <button @click="goBack" class="icon-btn" title="Назад">
        <span class="icon">←</span>
      </button>
      <InputText 
        v-model="deckName" 
        @keyup.enter="updateDeckName" 
        placeholder="Название колоды" 
        class="deck-name-input"
      />
      <div class="header-actions">
        <button @click="resetDeckProgress" class="secondary-btn" title="Сбросить прогресс">
          <span class="icon">🔄</span>
        </button>
        <button @click="deleteDeck" class="danger-btn" title="Удалить колоду">
          <span class="icon">🗑️</span>
        </button>
      </div>
    </div>

    <div class="add-card">
      <button @click="cardFormModalVisible = true" class="primary-btn">
        <span class="icon">+</span>
        Добавить карточку
      </button>
    </div>

    <div v-if="cards.length === 0" class="empty-state">
      <p>В этой колоде пока нет карточек. Добавь первую!</p>
    </div>

    <DataTable 
      v-else 
      :value="cards" 
      paginator 
      :rows="10" 
      class="custom-table"
      :paginatorTemplate="{'FirstPageLink': '', 'LastPageLink': '', 'PageLinks': 'PageLinks', 'PrevPageLink': 'PrevPageLink', 'NextPageLink': 'NextPageLink', 'RowsPerPageDropdown': '', 'CurrentPageReport': 'CurrentPageReport'}"
    >
      <Column field="KanjiText" header="Японский" class="kanji-col"></Column>
      <Column field="Translation" header="Перевод"></Column>
      <Column header="Действия" style="width: 150px">
        <template #body="slotProps">
          <div class="actions-cell">
            <button @click="editCard(slotProps.data)" class="icon-btn small" title="Редактировать">
              <span class="icon">✏️</span>
            </button>
            <button @click="deleteCard(slotProps.data.ID)" class="icon-btn small danger" title="Удалить">
              <span class="icon">🗑️</span>
            </button>
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog 
      v-model:visible="cardFormModalVisible" 
      :header="editingCard ? 'Редактировать карточку' : 'Новая карточка'" 
      class="custom-dialog"
    >
      <div class="form-content">
        <div class="input-group">
          <label for="kanji-text">Японское слово *</label>
          <InputText 
            id="kanji-text"
            v-model="cardForm.KanjiText" 
            placeholder="Например: 食べる" 
            class="custom-input"
          />
          <div v-if="errors.kanjiText" class="error">{{ errors.kanjiText }}</div>
        </div>
        <div class="input-group">
          <label for="furigana-text">Чтение (фуригана)</label>
          <InputText 
            id="furigana-text"
            v-model="cardForm.FuriganaText" 
            placeholder="Например: たべる" 
            class="custom-input"
          />
        </div>
        <div class="input-group">
          <label for="translation">Перевод *</label>
          <InputText 
            id="translation"
            v-model="cardForm.Translation" 
            placeholder="Например: есть" 
            class="custom-input"
          />
          <div v-if="errors.translation" class="error">{{ errors.translation }}</div>
        </div>
      </div>
      <template #footer>
        <Button label="Отмена" @click="closeCardModal" class="secondary-btn" />
        <Button label="Сохранить" @click="saveCard" class="primary-btn" />
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
import type { Deck, Card, CardInput } from '../types'
import { useWails } from '../composables/useWails'

const router = useRouter()
const route = useRoute()
const { getDecks, getCardsByDeck, updateDeck, resetDeckProgress: resetDeckProgressWails, deleteDeck: deleteDeckWails, createCard: createCardWails, updateCard: updateCardWails, deleteCard: deleteCardWails } = useWails()
const deckId = Number(route.params.id)
const deckName = ref('')
const cards = ref<Card[]>([])
const cardFormModalVisible = ref(false)
const editingCard = ref<Card | null>(null)
const cardForm = ref({
  KanjiText: '',
  FuriganaText: '',
  Translation: ''
})
const errors = ref({ kanjiText: '', translation: '' })

onMounted(async () => {
  await loadDeck()
  await loadCards()
})

const loadDeck = async () => {
  try {
    const decksList = await getDecks()
    const deck = decksList.find((d: Deck) => d.ID === deckId)
    if (deck) {
      deckName.value = deck.Name
    }
  } catch (e) {
    console.error('Ошибка загрузки колоды:', e)
    alert('Не удалось загрузить колоду: ' + e)
  }
}

const loadCards = async () => {
  try {
    cards.value = await getCardsByDeck(deckId)
  } catch (e) {
    console.error('Ошибка загрузки карточек:', e)
    alert('Не удалось загрузить карточки: ' + e)
  }
}

const goBack = () => {
  router.push({ name: 'Home' })
}

const updateDeckName = async () => {
  if (!deckName.value.trim()) return
  try {
    await updateDeck(deckId, deckName.value)
  } catch (e) {
    console.error('Ошибка обновления названия:', e)
    alert('Не удалось обновить название: ' + e)
  }
}

const resetDeckProgress = async () => {
  if (!confirm('Сбросить прогресс всей колоды?')) return
  try {
    await resetDeckProgressWails(deckId)
    alert('Прогресс успешно сброшен!')
  } catch (e) {
    console.error('Ошибка сброса прогресса:', e)
    alert('Не удалось сбросить прогресс: ' + e)
  }
}

const deleteDeck = async () => {
  if (!confirm('Удалить колоду и все карточки?')) return
  try {
    await deleteDeckWails(deckId)
    router.push({ name: 'Home' })
  } catch (e) {
    console.error('Ошибка удаления колоды:', e)
    alert('Не удалось удалить колоду: ' + e)
  }
}

const validateCardForm = (): boolean => {
  errors.value = { kanjiText: '', translation: '' }
  let valid = true
  
  if (!cardForm.value.KanjiText.trim()) {
    errors.value.kanjiText = 'Японское слово обязательно'
    valid = false
  }
  
  if (!cardForm.value.Translation.trim()) {
    errors.value.translation = 'Перевод обязателен'
    valid = false
  }
  
  return valid
}

const editCard = (card: Card) => {
  editingCard.value = card
  cardForm.value = {
    KanjiText: card.KanjiText,
    FuriganaText: card.FuriganaText || '',
    Translation: card.Translation
  }
  cardFormModalVisible.value = true
}

const deleteCard = async (id: number) => {
  if (!confirm('Удалить карточку?')) return
  try {
    await deleteCardWails(id)
    await loadCards()
  } catch (e) {
    console.error('Ошибка удаления карточки:', e)
    alert('Не удалось удалить карточку: ' + e)
  }
}

const saveCard = async () => {
  if (!validateCardForm()) return

  try {
    const input: CardInput = {
      DeckID: deckId,
      KanjiText: cardForm.value.KanjiText,
      FuriganaText: cardForm.value.FuriganaText.trim() || null,
      Translation: cardForm.value.Translation
    }

    if (editingCard.value) {
      await updateCardWails(editingCard.value.ID, input)
    } else {
      await createCardWails(input)
    }

    closeCardModal()
    await loadCards()
  } catch (e) {
    console.error('Ошибка сохранения карточки:', e)
    alert('Не удалось сохранить карточку: ' + e)
  }
}

const closeCardModal = () => {
  cardFormModalVisible.value = false
  editingCard.value = null
  cardForm.value = { KanjiText: '', FuriganaText: '', Translation: '' }
  errors.value = { kanjiText: '', translation: '' }
}
</script>

<style scoped>
.deck-manage {
  padding: 2rem;
  max-width: 1000px;
  margin: 0 auto;
}

.header {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.deck-name-input {
  flex: 1;
  min-width: 200px;
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  border: 1px solid #333333;
  background: #111111;
  color: white;
  font-family: inherit;
  font-size: 1.1rem;
  font-weight: 600;
  outline: none;
  transition: all 0.2s;
}

.deck-name-input:focus {
  border-color: #ff0a14;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.add-card {
  margin-bottom: 2rem;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #c7cdd8;
  font-size: 1.1rem;
}

.custom-table {
  border-radius: 0.75rem;
  overflow: hidden;
  border: 1px solid #222222;
}

:deep(.custom-table .p-datatable-thead > tr > th) {
  background: #1a1a1a !important;
  color: white !important;
  border: 1px solid #222222 !important;
  font-weight: 600;
}

:deep(.custom-table .p-datatable-tbody > tr > td) {
  background: #111111 !important;
  color: white !important;
  border: 1px solid #222222 !important;
}

:deep(.custom-table .p-paginator) {
  background: #111111 !important;
  border: none !important;
  padding: 1rem;
}

:deep(.custom-table .p-paginator .p-paginator-page),
:deep(.custom-table .p-paginator .p-paginator-prev),
:deep(.custom-table .p-paginator .p-paginator-next) {
  background: transparent !important;
  border: none !important;
  color: white !important;
}

:deep(.custom-table .p-paginator .p-paginator-page.p-highlight) {
  background: #ff0a14 !important;
  border-radius: 0.5rem !important;
}

:deep(.custom-table .p-paginator .p-paginator-page:hover),
:deep(.custom-table .p-paginator .p-paginator-prev:hover),
:deep(.custom-table .p-paginator .p-paginator-next:hover) {
  background: #222222 !important;
}

.kanji-col {
  font-family: 'Noto Sans JP', sans-serif;
  font-size: 1.1rem;
}

.actions-cell {
  display: flex;
  gap: 0.5rem;
}

.form-content {
  padding: 0.5rem 0;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
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

.icon-btn.small {
  padding: 0.4rem 0.6rem;
}

.icon-btn.small .icon {
  font-size: 1rem;
}

.icon-btn.danger:hover {
  border-color: #ff4444;
  color: #ff4444;
}

.icon {
  font-size: 1.25rem;
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

.primary-btn:hover {
  background-color: #e00912;
  transform: translateY(-1px);
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

.danger-btn {
  background-color: #ff4444;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  cursor: pointer;
  font-family: inherit;
  font-size: 1rem;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.danger-btn:hover {
  background-color: #e03a3a;
  transform: translateY(-1px);
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
