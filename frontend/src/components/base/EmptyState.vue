<template>
  <div class="empty-state">
    <slot name="icon">
      <el-icon :size="48" class="empty-state__icon"><Document /></el-icon>
    </slot>
    <p class="empty-state__text">
      <slot>{{ defaultText }}</slot>
    </p>
  </div>
</template>

<script setup>
import { Document } from '@element-plus/icons-vue'
defineProps({
  defaultText: { type: String, default: '暂无数据' }
})
</script>

<style scoped lang="scss">
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $sp-3;
  padding: $sp-7 0;
  color: $text-secondary;

  &__icon {
    color: $text-disabled;
  }
}
</style>
