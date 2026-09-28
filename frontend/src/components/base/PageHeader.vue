<template>
  <div class="page-header">
    <div class="page-header__title">
      <h1 class="page-header__h1">{{ title }}</h1>
      <p v-if="subtitle" class="page-header__sub">{{ subtitle }}</p>
    </div>
    <div class="page-header__actions">
      <slot />
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' }
})
</script>

<style scoped lang="scss">
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: $sp-4;

  &__h1 {
    font-size: $fs-2xl;
    font-weight: $fw-semibold;
    line-height: $lh-tight;
    color: $text-primary;
  }

  &__sub {
    margin-top: $sp-1;
    font-size: $fs-sm;
    color: $text-secondary;
  }

  &__actions {
    display: flex;
    align-items: center;
    gap: $sp-2;
    flex-shrink: 0;
  }
}
</style>
