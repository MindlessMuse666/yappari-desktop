<template>
  <span v-html="rendered" class="furigana-text"></span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  KanjiText: string
  FuriganaText?: string | null
}

const props = defineProps<Props>()

const rendered = computed(() => {
  if (!props.FuriganaText || props.FuriganaText === props.KanjiText) {
    return props.KanjiText
  }

  const kanji = props.KanjiText
  const furigana = props.FuriganaText
  let result = ''
  let kIndex = 0
  let fIndex = 0

  while (kIndex < kanji.length) {
    const char = kanji[kIndex]
    if (/[\u4e00-\u9faf]/.test(char)) {
      let kanjiPart = char
      kIndex++
      while (kIndex < kanji.length && /[\u4e00-\u9faf]/.test(kanji[kIndex])) {
        kanjiPart += kanji[kIndex]
        kIndex++
      }
      let furiganaPart = ''
      while (fIndex < furigana.length && !/[\u4e00-\u9faf]/.test(furigana[fIndex])) {
        furiganaPart += furigana[fIndex]
        fIndex++
      }
      result += `<ruby>${kanjiPart}<rt>${furiganaPart}</rt></ruby>`
    } else {
      result += char
      kIndex++
      fIndex++
    }
  }

  return result
})
</script>

<style scoped>
.furigana-text {
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  cursor: default;
}

.furigana-text :deep(ruby) {
  display: inline-block;
  text-align: center;
  line-height: 1.2;
}

.furigana-text :deep(rt) {
  display: block;
  font-size: 0.5em;
  color: #c7cdd8;
  line-height: 1;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}
</style>
