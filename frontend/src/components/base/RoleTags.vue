<template>
  <div v-if="roles && roles.length" class="role-tags">
    <el-tag v-for="r in roles" :key="r.id || r.name" size="small" effect="light">
      {{ r.name }}
    </el-tag>
  </div>
  <span v-else class="role-tags__empty">{{ emptyText }}</span>
</template>

<script setup>
// 角色标签组：用户管理列表 / 个人中心 复用（UI 形式 ≥2 次）
defineProps({
  roles: { type: Array, default: () => [] },
  emptyText: { type: String, default: '-' }
})
</script>

<style scoped lang="scss">
.role-tags {
  display: flex;
  flex-wrap: wrap;
  gap: $sp-1;

  &__empty {
    color: $text-disabled;
  }
}
</style>
