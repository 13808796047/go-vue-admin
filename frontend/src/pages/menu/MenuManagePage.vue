<template>
  <div class="menu-page">
    <PageHeader title="菜单管理" subtitle="管理目录、菜单与按钮权限">
      <el-button v-if="hasPermission('sys:menu:add')" type="primary" @click="openCreate()">
        <el-icon><Plus /></el-icon>新增
      </el-button>
    </PageHeader>

    <PageCard>
      <div class="table-wrap">
        <el-table
          v-loading="loading"
          :data="tree"
          row-key="id"
          default-expand-all
          :tree-props="{ children: 'children' }"
        >
          <el-table-column prop="name" label="名称" min-width="200">
            <template #default="{ row }">
              <div class="menu-page__name">
                <el-icon class="menu-page__icon"><component :is="iconOf(row)" /></el-icon>
                <span>{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="类型" width="110">
            <template #default="{ row }">
              <el-tag :type="typeOf(row).type" size="small" effect="light">
                {{ typeOf(row).label }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="路由" min-width="150">
            <template #default="{ row }">
              <CodeChip v-if="row.url" :code="row.url" />
              <span v-else class="menu-page__muted">-</span>
            </template>
          </el-table-column>
          <el-table-column label="权限标识" min-width="160">
            <template #default="{ row }">
              <CodeChip v-if="row.permKey" :code="row.permKey" />
              <span v-else class="menu-page__muted">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="70" align="center" />
          <el-table-column label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small" effect="light">
                {{ row.status === 1 ? '启用' : '停用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="{ row }">
              <RowActions>
                <el-button
                  v-if="row.type !== 3 && hasPermission('sys:menu:add')"
                  link
                  type="primary"
                  @click="openCreate(row)"
                >
                  新增子级
                </el-button>
                <el-button
                  v-if="hasPermission('sys:menu:add')"
                  link
                  type="primary"
                  @click="openEdit(row)"
                >
                  编辑
                </el-button>
                <el-button
                  v-if="hasPermission('sys:menu:delete')"
                  link
                  type="danger"
                  @click="remove(row)"
                >
                  删除
                </el-button>
              </RowActions>
            </template>
          </el-table-column>
          <template #empty>
            <EmptyState />
          </template>
        </el-table>
      </div>
    </PageCard>

    <FormDrawer
      v-model:visible="drawer.visible"
      :title="drawer.title"
      :submitting="drawer.submitting"
      :model="form"
      :rules="formRules"
      @close="drawer.visible = false"
      @submit="onSubmit"
    >
      <el-form-item v-if="!drawer.parentLocked" label="上级菜单">
        <el-select v-model="form.parentId" placeholder="根节点" style="width: 100%">
          <el-option :value="0" label="根节点" />
          <el-option
            v-for="m in parentOptions"
            :key="m.id"
            :label="m.name"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="菜单名称" />
      </el-form-item>
      <el-form-item label="类型">
        <el-radio-group v-model="form.type">
          <el-radio :value="1">目录</el-radio>
          <el-radio :value="2">菜单</el-radio>
          <el-radio :value="3">按钮/权限</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item v-if="form.type !== 3" label="路由">
        <el-input v-model="form.url" placeholder="/xxx" />
      </el-form-item>
      <el-form-item v-if="form.type !== 3" label="图标">
        <el-input v-model="form.icon" placeholder="图标标识（可选）" />
      </el-form-item>
      <el-form-item v-if="form.type === 2" label="组件">
        <el-input v-model="form.component" placeholder="组件路径（可选）" />
      </el-form-item>
      <el-form-item v-if="form.type !== 1" label="权限标识">
        <el-input v-model="form.permKey" placeholder="如 sys:user:add" />
      </el-form-item>
      <el-form-item label="排序">
        <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
      </el-form-item>
      <el-form-item label="状态">
        <el-radio-group v-model="form.status">
          <el-radio :value="1">启用</el-radio>
          <el-radio :value="0">停用</el-radio>
        </el-radio-group>
      </el-form-item>
    </FormDrawer>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Menu, Grid, Pointer } from '@element-plus/icons-vue'
import PageHeader from '@/components/base/PageHeader.vue'
import PageCard from '@/components/base/PageCard.vue'
import EmptyState from '@/components/base/EmptyState.vue'
import RowActions from '@/components/base/RowActions.vue'
import CodeChip from '@/components/base/CodeChip.vue'
import FormDrawer from '@/components/data/FormDrawer.vue'
import { useTable } from '@/composables/useTable.js'
import { useAuth } from '@/composables/useAuth.js'
import * as menuApi from '@/api/menu.js'

const { hasPermission } = useAuth()
const { data: tree, loading, reload, confirmRemove } = useTable(async () => {
  const list = await menuApi.menuTree()
  return { list, total: 0 }
})

onMounted(reload)

// 目录/菜单 可作为上级（type 1/2）
const parentOptions = computed(() => flatten(tree.value).filter((m) => m.type !== 3))

function flatten(nodes, out = []) {
  for (const n of nodes || []) {
    out.push(n)
    if (n.children && n.children.length) flatten(n.children, out)
  }
  return out
}

function iconOf(row) {
  return row.type === 1 ? Grid : row.type === 3 ? Pointer : Menu
}

function typeOf(row) {
  const map = {
    1: { label: '目录', type: 'info' },
    2: { label: '菜单', type: 'primary' },
    3: { label: '按钮/权限', type: 'warning' }
  }
  return map[row.type] || { label: String(row.type), type: 'info' }
}

const drawer = reactive({ visible: false, title: '新增菜单', submitting: false, parentId: 0, parentLocked: false })
const form = reactive({
  id: null,
  parentId: 0,
  name: '',
  type: 2,
  url: '',
  icon: '',
  component: '',
  permKey: '',
  sort: 0,
  status: 1
})

const formRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
}

function resetForm() {
  form.id = null
  form.parentId = 0
  form.name = ''
  form.type = 2
  form.url = ''
  form.icon = ''
  form.component = ''
  form.permKey = ''
  form.sort = 0
  form.status = 1
}

function openCreate(parent) {
  resetForm()
  drawer.title = '新增菜单'
  drawer.parentId = parent?.id || 0
  drawer.parentLocked = !!parent
  form.parentId = parent?.id || 0
  if (parent) {
    form.type = 3 // 默认按钮/权限
  }
  drawer.visible = true
}

function openEdit(row) {
  resetForm()
  Object.assign(form, {
    id: row.id,
    parentId: row.parentId || 0,
    name: row.name,
    type: row.type,
    url: row.url || '',
    icon: row.icon || '',
    component: row.component || '',
    permKey: row.permKey || '',
    sort: row.sort ?? 0,
    status: row.status ?? 1
  })
  drawer.title = '编辑菜单'
  drawer.parentLocked = false
  drawer.visible = true
}

async function onSubmit() {
  drawer.submitting = true
  try {
    const { id, ...payload } = form
    if (form.id) {
      await menuApi.updateMenu(id, payload)
      ElMessage.success('已更新')
    } else {
      await menuApi.createMenu(payload)
      ElMessage.success('已创建')
    }
    drawer.visible = false
    await reload()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    drawer.submitting = false
  }
}

async function remove(row) {
  if (!(await confirmRemove(row.name))) return
  try {
    await menuApi.removeMenu(row.id)
    ElMessage.success('已删除')
    await reload()
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  }
}
</script>

<style scoped lang="scss">
.menu-page {
  display: flex;
  flex-direction: column;
  gap: $sp-4;

  &__name {
    display: flex;
    align-items: center;
    gap: $sp-2;
  }

  &__icon {
    color: $indigo-500;
  }

  &__muted {
    color: $text-disabled;
  }
}

.table-wrap {
  background: $bg-card;
}

</style>
