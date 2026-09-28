<template>
  <header class="header">
    <div class="header__crumb">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>{{ $route.meta.title }}</el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="header__user">
      <el-dropdown trigger="click" @command="onCommand">
        <div class="header__avatar">
          <el-avatar :size="36">
            {{ initial }}
          </el-avatar>
          <span class="header__nickname">{{ nickname }}</span>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">个人中心</el-dropdown-item>
            <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { useAuth } from '@/composables/useAuth.js'

const router = useRouter()
const { user, logout } = useAuth()

const nickname = computed(() => user.value?.nickname || user.value?.username || '用户')
const initial = computed(() => (nickname.value || 'U').slice(0, 1).toUpperCase())

function onCommand(cmd) {
  if (cmd === 'profile') router.push('/profile')
  if (cmd === 'logout') {
    logout().then(() => {
      ElMessage.success('已退出登录')
      router.push('/login')
    })
  }
}
</script>

<style scoped lang="scss">
.header {
  height: $header-height;
  background: $bg-card;
  border-bottom: 1px solid $border-light;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 $sp-5;
  position: sticky;
  top: 0;
  z-index: 10;

  &__crumb {
    font-size: $fs-md;
    color: $text-secondary;
  }

  &__user {
    display: flex;
    align-items: center;
  }

  &__avatar {
    display: flex;
    align-items: center;
    gap: $sp-2;
    cursor: pointer;
    padding: $sp-1 $sp-2;
    border-radius: $radius;
    transition: background 0.15s ease;

    &:hover {
      background: $bg-hover;
    }

    .el-avatar {
      background: $avatar-bg;
      color: $white;
    }
  }

  &__nickname {
    font-size: $fs-md;
    color: $text-regular;
    font-weight: $fw-medium;
  }
}
</style>
