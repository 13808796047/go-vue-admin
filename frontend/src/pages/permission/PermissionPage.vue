<template>
  <div class="perm-page">
    <PageHeader title="权限分配" subtitle="为角色分配菜单权限与接口权限" />

    <div class="perm-page__grid">
      <PageCard class="perm-page__roles">
        <h3 class="perm-page__card-title">选择角色</h3>
        <el-radio-group v-model="selectedRoleId" class="perm-page__role-list" @change="onRoleChange">
          <el-radio
            v-for="r in roles"
            :key="r.id"
            :value="r.id"
            class="perm-page__role"
          >
            <div class="perm-page__role-line">
              <span class="perm-page__role-name">{{ r.name }}</span>
              <CodeChip :code="r.code" />
            </div>
          </el-radio>
        </el-radio-group>
      </PageCard>

      <PageCard class="perm-page__panel">
        <div class="perm-page__panel-head">
          <h3 class="perm-page__card-title">
            {{ selectedRoleName }} 的权限
          </h3>
          <el-button type="primary" size="small" :disabled="!selectedRoleId" :loading="saving" @click="save">
            保存分配
          </el-button>
        </div>

        <div v-if="!selectedRoleId" class="perm-page__hint">请选择左侧角色</div>
        <el-tabs v-else v-model="activeTab" class="perm-page__tabs">
          <el-tab-pane label="菜单权限" name="menu">
            <el-tree
              ref="menuTreeRef"
              :data="menuTree"
              show-checkbox
              node-key="id"
              :props="{ label: 'label', children: 'children' }"
              default-expand-all
              class="perm-page__tree"
            />
          </el-tab-pane>
          <el-tab-pane label="接口权限" name="api">
            <div class="perm-page__api-list">
              <el-checkbox-group v-model="checkedApiIds" class="perm-page__api-group">
                <label
                  v-for="api in allApis"
                  :key="api.id"
                  class="perm-page__api"
                >
                  <el-checkbox :value="api.id" />
                  <span class="perm-page__api-method" :class="`is-${api.method.toLowerCase()}`">
                    {{ api.method }}
                  </span>
                  <code class="perm-page__api-path">{{ api.path }}</code>
                  <span class="perm-page__api-name">{{ api.name }}</span>
                </label>
              </el-checkbox-group>
            </div>
          </el-tab-pane>
        </el-tabs>
      </PageCard>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/base/PageHeader.vue'
import PageCard from '@/components/base/PageCard.vue'
import CodeChip from '@/components/base/CodeChip.vue'
import * as roleApi from '@/api/role.js'
import * as menuApi from '@/api/menu.js'
import * as apiPerm from '@/api/api.js'

const roles = ref([])
const allApis = ref([])
const menuTree = ref([])
const selectedRoleId = ref(null)
const activeTab = ref('menu')
const saving = ref(false)
const checkedApiIds = ref([])
const menuTreeRef = ref(null)

const selectedRoleName = computed(
  () => roles.value.find((r) => r.id === selectedRoleId.value)?.name || ''
)

// 由平铺菜单列表构建树
function buildTree(flat) {
  const map = new Map()
  const roots = []
  for (const m of flat) map.set(m.id, { ...m, label: m.name, children: [] })
  for (const m of flat) {
    const node = map.get(m.id)
    if (m.parentId && map.has(m.parentId)) {
      map.get(m.parentId).children.push(node)
    } else {
      roots.push(node)
    }
  }
  return roots
}

async function loadBase() {
  const [r, menus, apis] = await Promise.all([
    roleApi.allRoles(),
    menuApi.allMenus(),
    apiPerm.allApis()
  ])
  roles.value = r || []
  menuTree.value = buildTree(menus || [])
  allApis.value = apis || []
}

async function onRoleChange() {
  if (!selectedRoleId.value) return
  activeTab.value = 'menu'
  const role = await roleApi.getRole(selectedRoleId.value)
  // 菜单：仅把「树中的叶子节点」设为勾选，祖先自动推导为全选/半选
  const assignedMenuIds = (role.menus || []).map((m) => m.id)
  checkedApiIds.value = (role.apis || []).map((a) => a.id)
  await nextTick()
  setMenuChecked(assignedMenuIds)
}

function leafIdsInTree() {
  const out = new Set()
  const walk = (nodes) => {
    for (const n of nodes) {
      if (n.children && n.children.length) walk(n.children)
      else out.add(n.id)
    }
  }
  walk(menuTree.value)
  return out
}

function setMenuChecked(assignedIds) {
  if (!menuTreeRef.value) return
  const leaves = leafIdsInTree()
  const toSet = assignedIds.filter((id) => leaves.has(id))
  menuTreeRef.value.setCheckedKeys(toSet, true)
}

async function save() {
  if (!selectedRoleId.value) return
  saving.value = true
  try {
    const checked = menuTreeRef.value.getCheckedKeys()
    const half = menuTreeRef.value.getHalfCheckedKeys()
    const menuIds = [...checked, ...half]
    await Promise.all([
      roleApi.assignRoleMenus(selectedRoleId.value, menuIds),
      roleApi.assignRoleApis(selectedRoleId.value, checkedApiIds.value)
    ])
    ElMessage.success('权限已保存')
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await loadBase()
  if (roles.value.length) {
    selectedRoleId.value = roles.value[0].id
    await onRoleChange()
  }
})
</script>

<style scoped lang="scss">
.perm-page {
  display: flex;
  flex-direction: column;
  gap: $sp-4;

  &__grid {
    display: grid;
    grid-template-columns: $col-narrow 1fr;
    gap: $sp-4;
  }

  &__roles {
    position: sticky;
    top: 0;
  }

  &__panel {
    display: flex;
    flex-direction: column;
  }

  &__card-title {
    font-size: $fs-lg;
    font-weight: $fw-semibold;
    color: $text-primary;
  }

  &__role-list {
    display: flex;
    flex-direction: column;
    gap: $sp-1;
    width: 100%;
  }

  &__role {
    height: auto;
    padding: $sp-3;
    border: 1px solid $border;
    border-radius: $radius;
    margin: 0;
    transition: border-color 0.15s ease, background 0.15s ease;

    &:hover {
      background: $bg-hover;
    }

    &.is-checked {
      border-color: $indigo-500;
      background: $indigo-50;
    }
  }

  &__role-line {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    gap: $sp-2;
  }

  &__role-name {
    font-size: $fs-md;
    color: $text-regular;
  }

  &__panel-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: $sp-4;
  }

  &__tabs {
    flex: 1;
  }

  &__hint {
    color: $text-disabled;
    font-size: $fs-sm;
    padding: $sp-5 0;
    text-align: center;
  }

  &__tree {
    background: transparent;
  }

  &__api-list {
    max-height: $api-list-h;
    overflow-y: auto;
    padding: $sp-2;
  }

  &__api-group {
    display: flex;
    flex-direction: column;
    gap: $sp-1;
  }

  &__api {
    display: flex;
    align-items: center;
    gap: $sp-2;
    padding: $sp-2;
    border-radius: $radius-sm;
    cursor: pointer;
    transition: background 0.15s ease;

    &:hover {
      background: $bg-hover;
    }

    .el-checkbox {
      margin-right: 0;
      flex: none;
    }
  }

  &__api-method {
    flex: none;
    width: $method-chip-w;
    text-align: center;
    font-size: $fs-xs;
    font-family: $code-font;
    font-weight: $fw-semibold;
    padding: $chip-pad 0;
    border-radius: $radius-sm;

    &.is-get {
      color: $success;
      background: rgba(16, 185, 129, 0.12);
    }
    &.is-post {
      color: $indigo-600;
      background: $indigo-50;
    }
    &.is-put {
      color: $warning;
      background: rgba(245, 158, 11, 0.12);
    }
    &.is-delete {
      color: $danger;
      background: rgba(239, 68, 68, 0.1);
    }
  }

  &__api-path {
    flex: 1;
    font-family: $code-font;
    font-size: $fs-sm;
    color: $text-regular;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__api-name {
    flex: none;
    font-size: $fs-xs;
    color: $text-secondary;
  }
}

</style>
