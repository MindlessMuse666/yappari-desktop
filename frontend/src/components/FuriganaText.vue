<template>
  <span class="furigana-text" :data-language="Language">
    <span class="main-text">{{ displayText }}</span>
    <span v-if="FuriganaText && FuriganaText !== KanjiText" class="furigana">【{{ FuriganaText }}】</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  KanjiText: string
  FuriganaText?: string | null
  Language?: 'ja' | 'ru'
}>(), {
  FuriganaText: null,
  Language: 'ja',
})

const displayText = computed(() => props.KanjiText)
</script>

<style scoped>
.furigana-text {
  cursor: pointer;
  display: inline;
}

.main-text {
  user-select: text;
  transition: background-position 0.3s ease;
  background-size: 200% 100%;
  background-image: linear-gradient(to right,
      white 50%,
      var(--hover-color) 50%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.furigana-text[data-language="ja"] .main-text {
  --hover-color: #ff0a14;
}

.furigana-text[data-language="ru"] .main-text {
  --hover-color: #004078;
  background-image: linear-gradient(to right,
      #c7cdd8 50%,
      var(--hover-color) 50%);
}

.furigana-text:hover .main-text {
  background-position: -100% 0;
}

.furigana {
  color: #c7cdd8;
  user-select: text;
  transition: background-position 0.3s ease;
  background-size: 200% 100%;
  background-image: linear-gradient(to right,
      #c7cdd8 50%,
      var(--hover-furigana-color) 50%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.furigana-text[data-language="ja"] .furigana {
  --hover-furigana-color: #ff0a14;
}

.furigana-text[data-language="ru"] .furigana {
  --hover-furigana-color: #004078;
}

.furigana-text:hover .furigana {
  background-position: -100% 0;
}
</style>
