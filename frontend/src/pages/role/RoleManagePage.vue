<template>
  <div class="role-page">
    <PageHeader title="角色管理" subtitle="管理角色、状态与关联权限">
      <el-button v-if="hasPermission('sys:role:add')" type="primary" @click="openCreate">
        <el-icon><Plus /></el-icon>新建角色
      </el-button>
    </PageHeader>

    <PageCard>
      <Toolbar
        placeholder="搜索角色名 / 标识"
        @search="onSearch"
      >
        <template #actions>
          <el-button @click="reload"><el-icon><Refresh /></el-icon>刷新</el-button>
        </template>
      </Toolbar>

      <div class="table-wrap">
        <el-table v-loading="loading" :data="rows" style="width: 100%">
          <el-table-column prop="name" label="角色名" min-width="140" />
          <el-table-column label="标识" min-width="140">
            <template #default="{ row }">
              <CodeChip :code="row.code" />
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="180">
            <template #default="{ row }">
              {{ row.remark || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="菜单数" width="90" align="center">
            <template #default="{ row }">
              <el-tag size="small" effect="light" type="primary">
                {{ (row.menus || []).length }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="接口数" width="90" align="center">
            <template #default="{ row }">
              <el-tag size="small" effect="light" type="primary">
                {{ (row.apis || []).length }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small" effect="light">
                {{ row.status === 1 ? '启用' : '停用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="150">
            <template #default="{ row }">
              {{ formatDateTime(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <RowActions>
                <el-button
                  v-if="hasPermission('sys:role:add')"
                  link
                  type="primary"
                  @click="openEdit(row)"
                >
                  编辑
                </el-button>
                <el-button
                  v-if="hasPermission('sys:role:delete')"
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

      <div class="pager">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="onPage"
          @size-change="onPageSize"
        />
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
      <el-form-item label="角色名" prop="name">
        <el-input v-model="form.name" placeholder="如：编辑" />
      </el-form-item>
      <el-form-item v-if="!form.id" label="标识" prop="code">
        <el-input v-model="form.code" placeholder="如：editor" />
      </el-form-item>
      <el-form-item label="排序">
        <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
      </el-form-item>
      <el-form-item label="备注">
        <el-input
          v-model="form.remark"
          type="textarea"
          :rows="2"
          placeholder="角色说明"
        />
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
import { reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import PageHeader from '@/components/base/PageHeader.vue'
import PageCard from '@/components/base/PageCard.vue'
import Toolbar from '@/components/data/Toolbar.vue'
import EmptyState from '@/components/base/EmptyState.vue'
import RowActions from '@/components/base/RowActions.vue'
import CodeChip from '@/components/base/CodeChip.vue'
import FormDrawer from '@/components/data/FormDrawer.vue'
import { useTable } from '@/composables/useTable.js'
import { useAuth } from '@/composables/useAuth.js'
import { formatDateTime } from '@/utils/format.js'
import * as roleApi from '@/api/role.js'

const { hasPermission } = useAuth()
const { data: rows, total, loading, pagination, reload, onSearch, onPage, onPageSize, confirmRemove } =
  useTable((params) => roleApi.listRoles(params))

onMounted(reload)

const drawer = reactive({ visible: false, title: '新建角色', submitting: false })
const form = reactive({
  id: null,
  name: '',
  code: '',
  sort: 0,
  remark: '',
  status: 1
})

const formRules = {
  name: [{ required: true, message: '请输入角色名', trigger: 'blur' }],
  code: [
    { required: true, message: '请输入标识', trigger: 'blur' },
    { pattern: /^[a-z][a-z0-9_:]*$/, message: '标识需以小写字母开头（字母/数字/下划线/冒号）', trigger: 'blur' }
  ]
}

function resetForm() {
  form.id = null
  form.name = ''
  form.code = ''
  form.sort = 0
  form.remark = ''
  form.status = 1
}

function openCreate() {
  resetForm()
  drawer.title = '新建角色'
  drawer.visible = true
}

function openEdit(row) {
  resetForm()
  Object.assign(form, {
    id: row.id,
    name: row.name,
    code: row.code,
    sort: row.sort ?? 0,
    remark: row.remark || '',
    status: row.status
  })
  drawer.title = '编辑角色'
  drawer.visible = true
}

async function onSubmit() {
  drawer.submitting = true
  try {
    if (form.id) {
      const { id, code, ...payload } = form
      await roleApi.updateRole(id, payload)
      ElMessage.success('已更新')
    } else {
      const { id, ...payload } = form
      await roleApi.createRole(payload)
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
    await roleApi.removeRole(row.id)
    ElMessage.success('已删除')
    await reload()
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  }
}
</script>

<style scoped lang="scss">
.role-page {
  display: flex;
  flex-direction: column;
  gap: $sp-4;
}

.table-wrap {
  background: $bg-card;
}

.pager {
  display: flex;
  justify-content: flex-end;
  padding-top: $sp-4;
}

</style>
