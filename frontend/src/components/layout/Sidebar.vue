<template>
  <aside class="sidebar">
    <div class="sidebar__brand">
      <div class="sidebar__logo">
        <el-icon :size="20"><Platform /></el-icon>
      </div>
      <span class="sidebar__name">Go Vue Admin</span>
    </div>

    <nav class="sidebar__nav">
      <router-link
        v-for="item in nav"
        :key="item.path"
        :to="item.path"
        class="sidebar__link"
        :class="{ 'sidebar__link--active': $route.path === item.path }"
      >
        <el-icon><component :is="item.icon" /></el-icon>
        <span>{{ item.label }}</span>
      </router-link>
    </nav>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
import { usePermission } from '@/composables/usePermission.js'
import { Platform } from '@element-plus/icons-vue'

const { visibleMenus } = usePermission()
const nav = computed(() =>
  visibleMenus.value.map((m) => ({
    path: m.path,
    label: m.title,
    icon: m.icon
  }))
)
</script>

<style scoped lang="scss">
.sidebar {
  width: $sidebar-width;
  height: 100vh;
  background: $bg-card;
  border-right: 1px solid $border-light;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;

  &__brand {
    display: flex;
    align-items: center;
    gap: $sp-3;
    height: $header-height;
    padding: 0 $sp-5;
  }

  &__logo {
    width: $logo-sm;
    height: $logo-sm;
    border-radius: $radius-sm;
    background: $indigo-500;
    color: $white;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: $shadow-sm;
  }

  &__name {
    font-size: $fs-lg;
    font-weight: $fw-semibold;
    color: $text-primary;
  }

  &__nav {
    padding: $sp-4 $sp-3;
    display: flex;
    flex-direction: column;
    gap: $sp-1;
  }

  &__link {
    display: flex;
    align-items: center;
    gap: $sp-3;
    padding: $sp-3;
    border-radius: $radius;
    color: $text-regular;
    font-size: $fs-md;
    text-decoration: none;
    transition: background 0.15s ease, color 0.15s ease;

    .el-icon {
      font-size: $icon-base;
      color: $text-secondary;
    }

    &:hover {
      background: $bg-hover;
    }

    &--active {
      background: $indigo-50;
      color: $indigo-600;
      font-weight: $fw-medium;

      .el-icon {
        color: $indigo-500;
      }
    }
  }
}
</style>
