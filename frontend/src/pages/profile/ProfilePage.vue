<template>
  <div class="profile-page">
    <PageHeader title="个人中心" subtitle="查看与修改当前账号资料" />

    <div class="profile-page__grid">
      <PageCard class="profile-page__overview">
        <div class="profile-page__ident">
          <el-avatar class="profile-page__avatar" :size="72">
            {{ initial }}
          </el-avatar>
          <div>
            <h2 class="profile-page__name">{{ user?.nickname || user?.username }}</h2>
            <div class="profile-page__meta">
              <RoleTags :roles="userRoles" empty-text="未分配角色" />
            </div>
          </div>
        </div>

        <el-divider />
        <dl class="profile-page__stats">
          <div class="profile-page__stat">
            <dt>用户 ID</dt>
            <dd>{{ user?.id ?? '-' }}</dd>
          </div>
          <div class="profile-page__stat">
            <dt>账号状态</dt>
            <dd>
              <StatusTag :value="user?.status ?? 1" />
            </dd>
          </div>
          <div class="profile-page__stat">
            <dt>加入时间</dt>
            <dd>{{ formatDateTime(user?.createdAt) }}</dd>
          </div>
        </dl>
      </PageCard>

      <div class="profile-page__forms">
        <PageCard>
          <h3 class="profile-page__card-title">修改资料</h3>
          <el-form :model="profileForm" label-position="top" class="profile-page__form">
            <el-form-item label="昵称">
              <el-input v-model="profileForm.nickname" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="profileForm.email" />
            </el-form-item>
            <el-form-item label="手机号">
              <el-input v-model="profileForm.phone" />
            </el-form-item>
            <el-button type="primary" :loading="savingProfile" @click="saveProfile">
              保存资料
            </el-button>
          </el-form>
        </PageCard>

        <PageCard>
          <h3 class="profile-page__card-title">修改密码</h3>
          <el-form :model="pwdForm" label-position="top" class="profile-page__form">
            <el-form-item label="当前密码">
              <el-input v-model="pwdForm.oldPassword" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="pwdForm.newPassword" type="password" show-password placeholder="至少 6 位" />
            </el-form-item>
            <el-form-item label="确认新密码">
              <el-input v-model="pwdForm.newPassword2" type="password" show-password />
            </el-form-item>
            <el-button type="primary" :loading="savingPwd" @click="changePassword">
              修改密码
            </el-button>
          </el-form>
        </PageCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/base/PageHeader.vue'
import PageCard from '@/components/base/PageCard.vue'
import StatusTag from '@/components/base/StatusTag.vue'
import RoleTags from '@/components/base/RoleTags.vue'
import { useAuth } from '@/composables/useAuth.js'
import { formatDateTime } from '@/utils/format.js'
import { getProfile, updateProfile, changePassword as apiChangePassword } from '@/api/auth.js'

const { user } = useAuth()

const profileForm = reactive({ nickname: '', email: '', phone: '' })
const pwdForm = reactive({ oldPassword: '', newPassword: '', newPassword2: '' })
const savingProfile = ref(false)
const savingPwd = ref(false)

const userRoles = computed(() => user.value?.roles || [])
const initial = computed(() =>
  ((user.value?.nickname || user.value?.username || 'U').slice(0, 1)).toUpperCase()
)

onMounted(async () => {
  try {
    const me = await getProfile()
    user.value = me
    profileForm.nickname = me.nickname || ''
    profileForm.email = me.email || ''
    profileForm.phone = me.phone || ''
  } catch {
    /* 拉取失败不阻断 */
  }
})

async function saveProfile() {
  savingProfile.value = true
  try {
    const updated = await updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
      phone: profileForm.phone
    })
    user.value = { ...user.value, ...updated }
    ElMessage.success('资料已保存')
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    savingProfile.value = false
  }
}

async function changePassword() {
  if (!pwdForm.newPassword) {
    ElMessage.warning('请输入新密码')
    return
  }
  if (pwdForm.newPassword.length < 6) {
    ElMessage.warning('新密码至少 6 位')
    return
  }
  if (pwdForm.newPassword !== pwdForm.newPassword2) {
    ElMessage.warning('两次密码不一致')
    return
  }
  savingPwd.value = true
  try {
    await apiChangePassword({
      oldPassword: pwdForm.oldPassword,
      newPassword: pwdForm.newPassword
    })
    pwdForm.oldPassword = ''
    pwdForm.newPassword = ''
    pwdForm.newPassword2 = ''
    ElMessage.success('密码已修改')
  } catch (e) {
    ElMessage.error(e.message || '修改失败')
  } finally {
    savingPwd.value = false
  }
}
</script>

<style scoped lang="scss">
.profile-page {
  display: flex;
  flex-direction: column;
  gap: $sp-4;

  &__grid {
    display: grid;
    grid-template-columns: $col-form 1fr;
    gap: $sp-4;
  }

  &__overview {
    position: sticky;
    top: 0;
  }

  &__forms {
    display: flex;
    flex-direction: column;
    gap: $sp-4;
  }

  &__ident {
    display: flex;
    align-items: center;
    gap: $sp-4;
  }

  &__avatar {
    background: $avatar-bg;
    color: $white;
  }

  &__name {
    font-size: $fs-xl;
    font-weight: $fw-semibold;
    color: $text-primary;
  }

  &__meta {
    display: flex;
    align-items: center;
    gap: $sp-2;
    margin-top: $sp-2;
  }


  &__stats {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: $sp-3;
  }

  &__stat {
    dt {
      font-size: $fs-xs;
      color: $text-secondary;
      margin-bottom: $sp-1;
    }
    dd {
      margin: 0;
      font-size: $fs-md;
      font-weight: $fw-medium;
      color: $text-regular;
    }
  }

  &__card-title {
    font-size: $fs-lg;
    font-weight: $fw-semibold;
    color: $text-primary;
  }

  &__form {
    display: flex;
    flex-direction: column;
    gap: $sp-4;
    max-width: $drawer-w;
    margin-top: $sp-4;
  }
}

</style>
