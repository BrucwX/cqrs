<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'

import { deleteTeacher, listTeachers, saveTeacher, type TeacherForm } from '@/api/teachers'
import {
  PHONE_PATTERN,
  TEACHER_STATUSES,
  enumLabel,
  enumTag,
  formatTimestamp,
  type Teacher,
} from '@/api/types'
import { usePagedList } from '@/composables/usePagedList'
import PagerBar from '@/components/PagerBar.vue'

const paged = usePagedList<Teacher>((page, pageSize) =>
  listTeachers({ page, page_size: pageSize }),
)

const keyword = ref('')

const rows = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return paged.items.value
  return paged.items.value.filter((row) =>
    [row.id, row.name, row.title, row.contact?.phone, row.contact?.email]
      .filter((field): field is string => Boolean(field))
      .some((field) => field.toLowerCase().includes(kw)),
  )
})

const dialogVisible = ref(false)
const saving = ref(false)
const editing = ref(false)
const formRef = ref<FormInstance>()

function blankForm(): TeacherForm {
  return { student_id: '', name: '', title: '', phone: '', email: '', status: 1 }
}

const form = reactive<TeacherForm>(blankForm())

const rules: FormRules<TeacherForm> = {
  name: [{ required: true, message: '请填姓名', trigger: 'blur' }],
  title: [{ required: true, message: '请填职称', trigger: 'blur' }],
  phone: [
    { required: true, message: '请填手机号', trigger: 'blur' },
    { pattern: PHONE_PATTERN, message: '手机号要 11 位、1 开头', trigger: 'blur' },
  ],
}

function openCreate(): void {
  Object.assign(form, blankForm())
  editing.value = false
  dialogVisible.value = true
}

function openEdit(row: Teacher): void {
  Object.assign(form, {
    id: row.id,
    student_id: row.student_id ? String(row.student_id) : '',
    name: row.name ?? '',
    title: row.title ?? '',
    phone: row.contact?.phone ?? '',
    email: row.contact?.email ?? '',
    status: row.status ?? 1,
  })
  editing.value = true
  dialogVisible.value = true
}

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    await saveTeacher(form)
    ElMessage.success(editing.value ? '已更新' : '已新增')
    dialogVisible.value = false
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(row: Teacher): Promise<void> {
  try {
    await ElMessageBox.confirm(`确定删除讲师「${row.name}」？此操作不可撤销。`, '删除确认', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    await deleteTeacher(row.id)
    ElMessage.success('已删除')
    if (paged.items.value.length === 1 && paged.page.value > 1) {
      await paged.goTo(paged.page.value - 1)
    } else {
      await paged.load()
    }
  } catch (error) {
    ElMessage.error((error as Error).message)
  }
}

onMounted(paged.load)
</script>

<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h2>讲师</h2>
        <p>Teacher 聚合 · 可关联学员 ID（讲师本人也是学员时）</p>
      </div>
      <div class="page-head-actions">
        <el-button :icon="Refresh" :loading="paged.loading.value" @click="paged.load()">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">新增讲师</el-button>
      </div>
    </div>

    <el-card shadow="never" class="page-card">
      <div class="page-toolbar">
        <el-input
          v-model="keyword"
          :prefix-icon="Search"
          placeholder="按姓名 / 职称 / 手机号筛选本页"
          clearable
          style="width: 320px"
        />
        <span class="page-toolbar-count">本页 {{ rows.length }} / {{ paged.items.value.length }} 条</span>
      </div>

      <el-alert
        v-if="paged.error.value"
        type="error"
        :closable="false"
        show-icon
        :title="paged.error.value"
        class="page-note"
      />

      <el-table v-loading="paged.loading.value" :data="rows" stripe empty-text="没有数据">
        <el-table-column prop="name" label="姓名" width="140" />
        <el-table-column prop="title" label="职称" width="150" />
        <el-table-column label="关联学员 ID" width="130">
          <template #default="{ row }">
            {{ row.student_id ? row.student_id : '—' }}
          </template>
        </el-table-column>
        <el-table-column label="手机号" width="140">
          <template #default="{ row }">{{ row.contact?.phone || '—' }}</template>
        </el-table-column>
        <el-table-column label="邮箱" min-width="190" show-overflow-tooltip>
          <template #default="{ row }">{{ row.contact?.email || '—' }}</template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="enumTag(TEACHER_STATUSES, row.status)" effect="light">
              {{ enumLabel(TEACHER_STATUSES, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">{{ formatTimestamp(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="130" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <PagerBar
        :page="paged.page.value"
        :page-size="paged.pageSize.value"
        :has-more="paged.hasMore.value"
        :loading="paged.loading.value"
        :count="paged.items.value.length"
        @update:page="paged.goTo($event)"
        @update:page-size="paged.setPageSize($event)"
      />
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="editing ? '编辑讲师' : '新增讲师'"
      width="520px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="讲师姓名" />
        </el-form-item>
        <el-form-item label="职称" prop="title">
          <el-input v-model="form.title" placeholder="如：金牌讲师" />
        </el-form-item>
        <el-form-item label="关联学员 ID" prop="student_id">
          <el-input v-model="form.student_id" placeholder="选填，纯数字" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="11 位手机号" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="选填" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" style="width: 100%">
            <el-option
              v-for="option in TEACHER_STATUSES"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">
          {{ editing ? '保存修改' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>
