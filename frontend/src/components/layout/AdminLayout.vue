<template>
  <div class="admin-layout">
    <Sidebar />
    <div class="admin-layout__main">
      <Header />
      <main class="admin-layout__content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import Header from '@/components/layout/Header.vue'
import { useAuth } from '@/composables/useAuth.js'

const { ensurePermissions } = useAuth()
onMounted(ensurePermissions)
</script>

<style scoped lang="scss">
.admin-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;

  &__main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  &__content {
    flex: 1;
    overflow-y: auto;
    padding: $sp-5;
  }
}
</style>
