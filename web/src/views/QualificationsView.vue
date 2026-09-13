<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, Search, View } from '@element-plus/icons-vue'

import {
  courseTypesByTeacher,
  deleteQualification,
  listQualifications,
  qualify,
  type QualifyForm,
} from '@/api/qualifications'
import { listTeachers } from '@/api/teachers'
import {
  QUALIFICATION_STATUSES,
  enumLabel,
  enumTag,
  formatTimestamp,
  type CourseType,
  type Qualification,
} from '@/api/types'
import { usePagedList } from '@/composables/usePagedList'
import PagerBar from '@/components/PagerBar.vue'

const paged = usePagedList<Qualification>((page, pageSize) =>
  listQualifications({ page, page_size: pageSize }),
)

const keyword = ref('')
const rows = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return paged.items.value
  return paged.items.value.filter((row) =>
    [row.id, row.teacher_id, row.course_type_id].some((field) =>
      String(field ?? '').toLowerCase().includes(kw),
    ),
  )
})

// ---------------------------------------------------------------- 表单

const dialogVisible = ref(false)
const saving = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<QualifyForm>({ teacher_id: '', course_type_id: '' })
const teacherOptions = ref<{ id: string; label: string }[]>([])

const rules: FormRules<QualifyForm> = {
  teacher_id: [{ required: true, message: '请选讲师', trigger: 'change' }],
  course_type_id: [{ required: true, message: '请填课程类型 ID', trigger: 'blur' }],
}

async function openCreate(): Promise<void> {
  Object.assign(form, { teacher_id: '', course_type_id: '' })
  dialogVisible.value = true
  try {
    teacherOptions.value = (await listTeachers({ page: 1, page_size: 100 })).map((t) => ({
      id: t.id,
      label: `${t.name}（${t.id}）`,
    }))
  } catch (error) {
    ElMessage.error((error as Error).message)
    dialogVisible.value = false
  }
}

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    await qualify(form)
    ElMessage.success('已授予资质')
    dialogVisible.value = false
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(row: Qualification): Promise<void> {
  try {
    await ElMessageBox.confirm(`确定删除资质记录「${row.id}」？`, '删除确认', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  try {
    await deleteQualification(row.id)
    ElMessage.success('已删除')
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  }
}

// ---------------------------------------------------------------- 讲师课程类型

const lookupVisible = ref(false)
const lookupLoading = ref(false)
const lookupTypes = ref<CourseType[]>([])
const lookupTeacher = ref('')

async function openLookup(row: Qualification): Promise<void> {
  const tid = String(row.teacher_id ?? '')
  lookupTeacher.value = tid
  lookupVisible.value = true
  lookupLoading.value = true
  try {
    lookupTypes.value = await courseTypesByTeacher(tid)
  } catch (error) {
    ElMessage.error((error as Error).message)
    lookupTypes.value = []
  } finally {
    lookupLoading.value = false
  }
}

onMounted(paged.load)
</script>

<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h2>授课资质</h2>
        <p>Qualification 聚合 · 只读入讲师 ID + 课程类型 ID，授予时间由领域层决定</p>
      </div>
      <div class="page-head-actions">
        <el-button :icon="Refresh" :loading="paged.loading.value" @click="paged.load()">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">授予资质</el-button>
      </div>
    </div>

    <el-card shadow="never" class="page-card">
      <div class="page-toolbar">
        <el-input
          v-model="keyword"
          :prefix-icon="Search"
          placeholder="按讲师 / 课程类型筛选本页"
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
        <el-table-column prop="teacher_id" label="讲师" width="150" />
        <el-table-column prop="course_type_id" label="课程类型" width="140" />
        <el-table-column label="授予时间" width="160">
          <template #default="{ row }">{{ formatTimestamp(row.certified_at) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="enumTag(QUALIFICATION_STATUSES, row.status)" effect="light">
              {{ enumLabel(QUALIFICATION_STATUSES, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" :icon="View" @click="openLookup(row)">课程类型</el-button>
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

    <el-dialog v-model="dialogVisible" title="授予资质" width="460px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="88px">
        <el-form-item label="讲师" prop="teacher_id">
          <el-select v-model="form.teacher_id" filterable placeholder="选择讲师" style="width: 100%">
            <el-option v-for="option in teacherOptions" :key="option.id" :label="option.label" :value="option.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="课程类型 ID" prop="course_type_id">
          <el-input v-model="form.course_type_id" placeholder="如 ct-0001（少儿编程）" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">授予</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="lookupVisible" :title="`讲师 ${lookupTeacher} 有资质的课程类型`" width="560px">
      <el-table v-loading="lookupLoading" :data="lookupTypes" size="small" stripe empty-text="没有数据">
        <el-table-column prop="id" label="ID" width="140" />
        <el-table-column prop="name" label="名称" width="140" />
        <el-table-column prop="description" label="描述" />
      </el-table>
    </el-dialog>
  </div>
</template>
