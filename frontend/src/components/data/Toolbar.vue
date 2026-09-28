<template>
  <div class="toolbar">
    <div class="toolbar__left">
      <el-input
        :model-value="keyword"
        class="toolbar__search"
        :placeholder="placeholder"
        clearable
        @update:model-value="onKeyword"
        @clear="onSearch"
        @keyup.enter="onSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button @click="onSearch">
        <el-icon><Search /></el-icon>查询
      </el-button>
      <slot name="extra-left" />
    </div>
    <div class="toolbar__right">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Search } from '@element-plus/icons-vue'

defineProps({
  placeholder: { type: String, default: '搜索' }
})
const keyword = ref('')
const emit = defineEmits(['search'])

function onKeyword(v) {
  keyword.value = v
  emit('search', v)
}

function onSearch() {
  emit('search', keyword.value)
}
</script>

<style scoped lang="scss">
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;

  &__left,
  &__right {
    display: flex;
    align-items: center;
    gap: $sp-2;
  }

  &__search {
    width: $search-w;
  }
}
</style>
