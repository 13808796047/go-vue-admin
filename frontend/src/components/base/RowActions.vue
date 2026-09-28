<template>
  <div class="row-actions"><slot /></div>
</template>

<script setup>
// 表格行操作区：统一横向排布 link 型按钮（间距 $sp-3）
</script>

<style scoped lang="scss">
.row-actions {
  display: flex;
  align-items: center;
  gap: $sp-3;
}
</style>
