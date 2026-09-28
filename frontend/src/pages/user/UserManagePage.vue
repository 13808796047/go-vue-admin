<template>
  <div class="user-page">
    <PageHeader title="用户管理" subtitle="管理登录用户、角色与状态">
      <el-button v-if="hasPermission('sys:user:add')" type="primary" @click="openCreate">
        <el-icon><Plus /></el-icon>新建用户
      </el-button>
    </PageHeader>

    <PageCard>
      <Toolbar
        placeholder="搜索用户名 / 昵称 / 邮箱"
        @search="onSearch"
      >
        <template #actions>
          <el-button @click="reload"><el-icon><Refresh /></el-icon>刷新</el-button>
        </template>
      </Toolbar>

      <div class="table-wrap">
        <el-table
          v-loading="loading"
          :data="rows"
          style="width: 100%"
          row-key="id"
        >
          <el-table-column prop="username" label="用户名" min-width="120" />
          <el-table-column prop="nickname" label="昵称" min-width="120" />
          <el-table-column prop="email" label="邮箱" min-width="180">
            <template #default="{ row }">
              {{ row.email || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="角色" min-width="140">
 <template #default="{ row }">
   <RoleTags :roles="row.roles" />
 </template>
          </el-table-column>
          <el-table-column label="创建时间" width="150">
            <template #default="{ row }">
              {{ formatDateTime(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <RowActions>
                <el-button
                  v-if="hasPermission('sys:user:edit')"
                  link
                  type="primary"
                  @click="openEdit(row)"
                >
                  编辑
                </el-button>
                <el-button
                  v-if="hasPermission('sys:user:edit')"
                  link
                  type="primary"
                  @click="toggleStatus(row)"
                >
                  {{ row.status === 1 ? '停用' : '启用' }}
                </el-button>
                <el-button
                  v-if="hasPermission('sys:user:delete')"
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
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="form.username"
          placeholder="登录名（3-32 位）"
          :disabled="!!form.id"
        />
      </el-form-item>
      <el-form-item label="昵称" prop="nickname">
        <el-input v-model="form.nickname" placeholder="昵称" />
      </el-form-item>
      <el-form-item v-if="!form.id" label="初始密码" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          show-password
          placeholder="至少 6 位"
        />
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" placeholder="邮箱" />
      </el-form-item>
      <el-form-item label="手机号" prop="phone">
        <el-input v-model="form.phone" placeholder="手机号" />
      </el-form-item>
      <el-form-item label="角色">
        <el-select v-model="form.roleIds" multiple placeholder="选择角色" style="width: 100%">
          <el-option
            v-for="r in roles"
            :key="r.id"
            :label="r.name"
            :value="r.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-radio-group v-model="form.status">
          <el-radio :value="1">正常</el-radio>
          <el-radio :value="0">停用</el-radio>
        </el-radio-group>
      </el-form-item>
    </FormDrawer>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import PageHeader from '@/components/base/PageHeader.vue'
import PageCard from '@/components/base/PageCard.vue'
import Toolbar from '@/components/data/Toolbar.vue'
import RowActions from '@/components/base/RowActions.vue'
import StatusTag from '@/components/base/StatusTag.vue'
import EmptyState from '@/components/base/EmptyState.vue'
import RoleTags from '@/components/base/RoleTags.vue'
import FormDrawer from '@/components/data/FormDrawer.vue'
import { useTable } from '@/composables/useTable.js'
import { useAuth } from '@/composables/useAuth.js'
import { formatDateTime } from '@/utils/format.js'
import * as userApi from '@/api/user.js'
import * as roleApi from '@/api/role.js'

const { hasPermission } = useAuth()
const { data: rows, total, loading, pagination, reload, onSearch, onPage, onPageSize, confirmRemove } =
  useTable((params) => userApi.listUsers(params))

const roles = ref([])
onMounted(async () => {
  await Promise.all([
    reload(),
    roleApi.allRoles().then((list) => (roles.value = list || []))
  ])
})

// —— 表单抽屉 ——
const drawer = reactive({
  visible: false,
  title: '新建用户',
  submitting: false
})
const form = reactive({
  id: null,
  username: '',
  nickname: '',
  password: '',
  email: '',
  phone: '',
  roleIds: [],
  status: 1
})

const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '用户名 3-32 位', trigger: 'blur' }
  ],
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入初始密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 位', trigger: 'blur' }
  ],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }]
}

function resetForm() {
  form.id = null
  form.username = ''
  form.nickname = ''
  form.password = ''
  form.email = ''
  form.phone = ''
  form.roleIds = []
  form.status = 1
}

function openCreate() {
  resetForm()
  drawer.title = '新建用户'
  drawer.visible = true
}

function openEdit(row) {
  resetForm()
  Object.assign(form, {
    id: row.id,
    username: row.username,
    nickname: row.nickname,
    email: row.email || '',
    phone: row.phone || '',
    status: row.status,
    roleIds: (row.roles || []).map((r) => r.id)
  })
  drawer.title = '编辑用户'
  drawer.visible = true
}

async function onSubmit() {
  drawer.submitting = true
  try {
    if (form.id) {
      const { id, username, password, ...payload } = form
      await userApi.updateUser(id, payload)
      ElMessage.success('已更新')
    } else {
      const { id, ...payload } = form
      await userApi.createUser(payload)
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

async function toggleStatus(row) {
  try {
    await userApi.setUserStatus(row.id, row.status === 1 ? 0 : 1)
    ElMessage.success(row.status === 1 ? '已停用' : '已启用')
    await reload()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function remove(row) {
  if (!(await confirmRemove(row.nickname || row.username))) return
  try {
    await userApi.removeUser(row.id)
    ElMessage.success('已删除')
    await reload()
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  }
}
</script>

<style scoped lang="scss">
.user-page {
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
