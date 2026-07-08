<template>
  <span v-html="rendered"></span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  kanjiText: string
  furiganaText?: string | null
}

const props = defineProps<Props>()

const rendered = computed(() => {
  if (!props.furiganaText || props.furiganaText === props.kanjiText) {
    return props.kanjiText
  }

  const kanji = props.kanjiText
  const furigana = props.furiganaText
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
