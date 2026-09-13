<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'

import { deleteEnrollment, enrollStudent, listEnrollments, type EnrollForm } from '@/api/enrollments'
import { listCourses } from '@/api/courses'
import { listStudents } from '@/api/students'
import {
  ENROLLMENT_STATUSES,
  enumLabel,
  enumTag,
  formatTimestamp,
  type CourseEnrollment,
} from '@/api/types'
import { usePagedList } from '@/composables/usePagedList'
import PagerBar from '@/components/PagerBar.vue'

const paged = usePagedList<CourseEnrollment>((page, pageSize) =>
  listEnrollments({ page, page_size: pageSize }),
)

const keyword = ref('')
const rows = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return paged.items.value
  return paged.items.value.filter((row) =>
    [row.id, row.student_id, row.course_id].some((field) =>
      String(field ?? '').toLowerCase().includes(kw),
    ),
  )
})

// ---------------------------------------------------------------- 表单

const dialogVisible = ref(false)
const saving = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<EnrollForm>({ student_id: '', course_id: '' })
const studentOptions = ref<{ id: string; label: string }[]>([])
const courseOptions = ref<{ id: string; label: string }[]>([])

const rules: FormRules<EnrollForm> = {
  student_id: [{ required: true, message: '请选学员', trigger: 'change' }],
  course_id: [{ required: true, message: '请选课程', trigger: 'change' }],
}

async function openCreate(): Promise<void> {
  Object.assign(form, { student_id: '', course_id: '' })
  dialogVisible.value = true
  try {
    ;[studentOptions.value, courseOptions.value] = await Promise.all([
      listStudents({ page: 1, page_size: 100 }).then((list) =>
        list.map((s) => ({ id: s.id, label: `${s.name}（${s.id}）` })),
      ),
      listCourses({ page: 1, page_size: 100 }).then((list) =>
        list.map((c) => ({ id: c.id, label: `${c.id}（${c.course_type_id ?? ''}）` })),
      ),
    ])
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
    await enrollStudent(form)
    ElMessage.success('已报名')
    dialogVisible.value = false
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(row: CourseEnrollment): Promise<void> {
  try {
    await ElMessageBox.confirm(`确定删除选课记录「${row.id}」？`, '删除确认', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  try {
    await deleteEnrollment(row.id)
    ElMessage.success('已删除')
    await paged.load()
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
        <h2>选课</h2>
        <p>Enrollment 聚合 · 只读入学员 ID + 课程 ID，状态与时间由领域层决定</p>
      </div>
      <div class="page-head-actions">
        <el-button :icon="Refresh" :loading="paged.loading.value" @click="paged.load()">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">新增选课</el-button>
      </div>
    </div>

    <el-card shadow="never" class="page-card">
      <div class="page-toolbar">
        <el-input
          v-model="keyword"
          :prefix-icon="Search"
          placeholder="按学员 / 课程筛选本页"
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
        <el-table-column label="学员" width="150">
          <template #default="{ row }">{{ row.student_id ?? '—' }}</template>
        </el-table-column>
        <el-table-column prop="course_id" label="课程" width="120" />
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="enumTag(ENROLLMENT_STATUSES, row.status)" effect="light">
              {{ enumLabel(ENROLLMENT_STATUSES, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="报名时间" width="160">
          <template #default="{ row }">{{ formatTimestamp(row.enrolled_at) }}</template>
        </el-table-column>
        <el-table-column label="结业/退课" min-width="160">
          <template #default="{ row }">
            {{ formatTimestamp(row.completed_at) }} {{ formatTimestamp(row.dropped_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
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

    <el-dialog v-model="dialogVisible" title="新增选课" width="460px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="学员" prop="student_id">
          <el-select v-model="form.student_id" filterable placeholder="选择学员" style="width: 100%">
            <el-option v-for="option in studentOptions" :key="option.id" :label="option.label" :value="option.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="课程" prop="course_id">
          <el-select v-model="form.course_id" filterable placeholder="选择课程" style="width: 100%">
            <el-option v-for="option in courseOptions" :key="option.id" :label="option.label" :value="option.id" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">报名</el-button>
      </template>
    </el-dialog>
  </div>
</template>
