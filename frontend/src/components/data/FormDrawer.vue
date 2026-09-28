<template>
  <el-drawer
    v-model="model"
    :title="title"
    :size="420"
    :close-on-click-modal="false"
    @close="onClose"
  >
    <el-form ref="formEl" :model="modelForm" :rules="rules" label-position="top" class="drawer-form">
      <slot :form="modelForm" />
    </el-form>
    <template #footer>
      <div class="drawer-form__footer">
        <el-button @click="onCancel">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="onSubmit">
          保存
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const props = defineProps({
  title: { type: String, default: '' },
  visible: { type: Boolean, default: false },
  submitting: { type: Boolean, default: false },
  // 父组件持有的响应式表单对象（可选；不传则不校验）
  model: { type: Object, default: null },
  // 表单校验规则
  rules: { type: Object, default: null }
})
const emit = defineEmits(['update:visible', 'close', 'submit'])

const model = ref(props.visible)
// 父级 visible 变化时同步打开态（el-drawer 的 model-value 即开关）
watch(() => props.visible, (v) => { model.value = v })
// 暴露给模板的表单对象（校验绑定 + slot 共用）
const modelForm = computed(() => props.model || {})
const formEl = ref(null)

function onClose() {
  emit('close')
}

function onCancel() {
  model.value = false
  emit('update:visible', false)
}

async function onSubmit() {
  if (formEl.value && props.rules) {
    try {
      await formEl.value.validate()
    } catch {
      return
    }
  }
  emit('submit')
}
</script>

<style scoped lang="scss">
.drawer-form {
  display: flex;
  flex-direction: column;
  gap: $sp-4;
  padding-bottom: $sp-5;
  height: 100%;
  overflow-y: auto;

  &__footer {
    display: flex;
    justify-content: flex-end;
    gap: $sp-2;
  }
}
</style>
