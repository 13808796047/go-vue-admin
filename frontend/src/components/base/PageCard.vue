<template>
  <div class="page-card">
    <slot name="header" />
    <div class="page-card__body">
      <slot />
    </div>
    <div v-if="$slots.footer" class="page-card__footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup>
// 卡片壳：统一 12px 圆角 + 柔和阴影 + 24px 内边距（《设计规范》卡片）
</script>

<style scoped lang="scss">
.page-card {
  background: $bg-card;
  border-radius: $radius;
  box-shadow: $shadow;
  padding: $sp-5;
  display: flex;
  flex-direction: column;
  gap: $sp-4;

  &__body {
    flex: 1;
  }

  &__footer {
    padding-top: $sp-4;
  }
}
</style>
