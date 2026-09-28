<template>
  <el-tag :type="tagType" size="small" effect="light">
    {{ USER_STATUS_TEXT[value] ?? '未知' }}
  </el-tag>
</template>

<script setup>
import { computed } from 'vue'
import { USER_STATUS_TEXT, USER_STATUS_TYPE } from '@/constants/index.js'

const props = defineProps({
  value: { type: Number, default: 1 }
})
const tagType = computed(() => USER_STATUS_TYPE[props.value] || 'info')
</script>
