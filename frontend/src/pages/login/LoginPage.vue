<template>
  <div class="login">
    <div class="login__card">
      <div class="login__brand">
        <div class="login__logo">
          <el-icon :size="22"><Platform /></el-icon>
        </div>
        <h1 class="login__title">Go Vue Admin</h1>
        <p class="login__sub">后台管理系统</p>
      </div>

      <el-form class="login__form" @submit.prevent>
        <el-form-item>
          <el-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
            @keyup.enter="submit"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            show-password
            size="large"
            @keyup.enter="submit"
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-button
          type="primary"
          size="large"
          class="login__submit"
          :loading="loading"
          @click="submit"
        >
          登录
        </el-button>
      </el-form>
      <p class="login__hint">演示账号：admin / 123456</p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Platform, User, Lock } from '@element-plus/icons-vue'
import { useAuth } from '@/composables/useAuth.js'

const router = useRouter()
const { login } = useAuth()

const form = reactive({ username: 'admin', password: '123456' })
const loading = ref(false)

async function submit() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    await login({ username: form.username, password: form.password })
    ElMessage.success('登录成功')
    router.push('/')
  } catch (e) {
    ElMessage.error(e.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, $indigo-50 0%, $bg-page 55%, $bg-login 100%);

  &__card {
    width: $login-card-w;
    background: $bg-card;
    border-radius: $radius-lg;
    box-shadow: $shadow-lg;
    padding: $sp-7 $sp-6;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: $sp-5;
  }

  &__brand {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: $sp-3;
  }

  &__logo {
    width: $logo-lg;
    height: $logo-lg;
    border-radius: $radius-lg;
    background: $indigo-500;
    color: $white;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: $shadow-lg;
  }

  &__title {
    font-size: $fs-2xl;
    font-weight: $fw-bold;
    color: $text-primary;
    margin-top: $sp-2;
  }

  &__sub {
    font-size: $fs-sm;
    color: $text-secondary;
  }

  &__form {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: $sp-4;
    margin: 0;
  }

  &__submit {
    width: 100%;
  }

  &__hint {
    font-size: $fs-xs;
    color: $text-disabled;
  }
}
</style>
