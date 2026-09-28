<template>
  <code class="code-chip">{{ code }}</code>
</template>

<script setup>
defineProps({
  code: { type: String, default: '' }
})
</script>

<style scoped lang="scss">
.code-chip {
  background: $indigo-50;
  color: $indigo-600;
  font-family: $code-font;
  border-radius: $radius-sm;
  font-size: $fs-xs;
  font-family: 'SFMono-Regular', Consolas, monospace;
}
</style>
