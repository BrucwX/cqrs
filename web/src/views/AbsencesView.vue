<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'

import { deleteAbsence, listAbsences, recordAbsence, type AbsenceForm } from '@/api/absences'
import { listCourses } from '@/api/courses'
import { listStudents } from '@/api/students'
import {
  ABSENCE_TYPES,
  enumLabel,
  enumTag,
  formatTimestamp,
  type AbsenceRecord,
} from '@/api/types'
import { usePagedList } from '@/composables/usePagedList'
import PagerBar from '@/components/PagerBar.vue'

const paged = usePagedList<AbsenceRecord>((page, pageSize) =>
  listAbsences({ page, page_size: pageSize }),
)

const keyword = ref('')
const rows = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return paged.items.value
  return paged.items.value.filter((row) =>
    [row.id, row.student_id, row.course_id, row.reason].some((field) =>
      String(field ?? '').toLowerCase().includes(kw),
    ),
  )
})

// ---------------------------------------------------------------- 表单

const dialogVisible = ref(false)
const saving = ref(false)
const formRef = ref<FormInstance>()

function blankForm(): AbsenceForm {
  return {
    student_id: '',
    course_id: '',
    course_slot_id: '',
    schedule_date: '',
    missed_hours: 2,
    absence_type: 1,
    reason: '',
  }
}

const form = reactive<AbsenceForm>(blankForm())
const studentOptions = ref<{ id: string; label: string }[]>([])
const courseOptions = ref<{ id: string; label: string }[]>([])

const rules: FormRules<AbsenceForm> = {
  student_id: [{ required: true, message: '请选学员', trigger: 'change' }],
  course_id: [{ required: true, message: '请选课程', trigger: 'change' }],
  course_slot_id: [{ required: true, message: '请填课次 ID', trigger: 'blur' }],
  schedule_date: [{ required: true, message: '请选缺勤日期', trigger: 'change' }],
  missed_hours: [{ required: true, message: '请填缺勤课时', trigger: 'blur' }],
}

async function openCreate(): Promise<void> {
  Object.assign(form, blankForm())
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
    await recordAbsence(form)
    ElMessage.success('已记录缺勤')
    dialogVisible.value = false
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(row: AbsenceRecord): Promise<void> {
  try {
    await ElMessageBox.confirm(`确定删除缺勤记录「${row.id}」？`, '删除确认', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  try {
    await deleteAbsence(row.id)
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
        <h2>缺勤</h2>
        <p>AbsenceRecord 聚合 · 只记录事实，不参与状态流转</p>
      </div>
      <div class="page-head-actions">
        <el-button :icon="Refresh" :loading="paged.loading.value" @click="paged.load()">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">记录缺勤</el-button>
      </div>
    </div>

    <el-card shadow="never" class="page-card">
      <div class="page-toolbar">
        <el-input
          v-model="keyword"
          :prefix-icon="Search"
          placeholder="按学员 / 课程 / 事由筛选本页"
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
        <el-table-column prop="student_id" label="学员" width="150" />
        <el-table-column prop="course_id" label="课程" width="120" />
        <el-table-column label="缺勤日期" width="160">
          <template #default="{ row }">{{ formatTimestamp(row.schedule_date) }}</template>
        </el-table-column>
        <el-table-column prop="missed_hours" label="课时" width="80" />
        <el-table-column label="类型" width="90">
          <template #default="{ row }">
            <el-tag :type="enumTag(ABSENCE_TYPES, row.absence_type)" effect="light">
              {{ enumLabel(ABSENCE_TYPES, row.absence_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="事由" min-width="140" show-overflow-tooltip />
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

    <el-dialog v-model="dialogVisible" title="记录缺勤" width="520px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
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
        <el-form-item label="课次 ID" prop="course_slot_id">
          <el-input v-model="form.course_slot_id" placeholder="被缺的槽位 UUID" />
        </el-form-item>
        <el-form-item label="缺勤日期" prop="schedule_date">
          <el-date-picker v-model="form.schedule_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="缺勤课时" prop="missed_hours">
          <el-input-number v-model="form.missed_hours" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="类型" prop="absence_type">
          <el-select v-model="form.absence_type" style="width: 100%">
            <el-option v-for="option in ABSENCE_TYPES" :key="option.value" :label="option.label" :value="option.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="事由" prop="reason">
          <el-input v-model="form.reason" type="textarea" :rows="2" placeholder="选填" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">记录</el-button>
      </template>
    </el-dialog>
  </div>
</template>
