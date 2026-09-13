<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, Search, View } from '@element-plus/icons-vue'

import {
  courseAbsences,
  courseMakeups,
  courseSlotChanges,
  courseSlots,
  courseStudents,
  deleteCourse,
  listCourses,
  saveCourse,
  type CourseForm,
} from '@/api/courses'
import {
  ABSENCE_TYPES,
  CHANGE_TYPES,
  MAKEUP_STATUSES,
  STUDENT_STATUSES,
  enumLabel,
  enumTag,
  formatTimeRange,
  formatTimestamp,
  tsToDateStr,
  type AbsenceRecord,
  type Course,
  type CourseSlot,
  type CourseSlotChange,
  type Student,
  type StudentMakeup,
} from '@/api/types'
import { usePagedList } from '@/composables/usePagedList'
import PagerBar from '@/components/PagerBar.vue'

const paged = usePagedList<Course>((page, pageSize) =>
  listCourses({ page, page_size: pageSize }),
)

const keyword = ref('')
const rows = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return paged.items.value
  return paged.items.value.filter((row) =>
    [row.id, row.course_type_id].some((field) => field?.toLowerCase().includes(kw)),
  )
})

// ---------------------------------------------------------------- 表单

const dialogVisible = ref(false)
const saving = ref(false)
const editing = ref(false)
const formRef = ref<FormInstance>()

function blankForm(): CourseForm {
  return {
    course_type_id: '',
    capacity_max: 30,
    capacity_enrolled: 0,
    enroll_start_at: '',
    enroll_end_at: '',
    drop_deadline: '',
    period_start_at: '',
    period_end_at: '',
    total_hours: 48,
    completed_hours: 0,
  }
}

const form = reactive<CourseForm>(blankForm())

const rules: FormRules<CourseForm> = {
  course_type_id: [{ required: true, message: '请填课程类型 ID', trigger: 'blur' }],
  capacity_max: [{ required: true, message: '请填容量', trigger: 'blur' }],
  enroll_start_at: [{ required: true, message: '请选报名开始日', trigger: 'change' }],
  enroll_end_at: [{ required: true, message: '请选报名截止日', trigger: 'change' }],
  drop_deadline: [{ required: true, message: '请选退课截止日', trigger: 'change' }],
  period_start_at: [{ required: true, message: '请选开课日', trigger: 'change' }],
  period_end_at: [{ required: true, message: '请选结课日', trigger: 'change' }],
  total_hours: [{ required: true, message: '请填总课时', trigger: 'blur' }],
}

function openCreate(): void {
  Object.assign(form, blankForm())
  editing.value = false
  dialogVisible.value = true
}

function openEdit(row: Course): void {
  Object.assign(form, {
    id: row.id,
    course_type_id: row.course_type_id ?? '',
    capacity_max: row.capacity?.max ?? 0,
    capacity_enrolled: row.capacity?.enrolled ?? 0,
    enroll_start_at: tsToDateStr(row.enrollment?.start_at),
    enroll_end_at: tsToDateStr(row.enrollment?.end_at),
    drop_deadline: tsToDateStr(row.enrollment?.drop_deadline),
    period_start_at: tsToDateStr(row.period?.start_at),
    period_end_at: tsToDateStr(row.period?.end_at),
    total_hours: row.period?.total_hours ?? 0,
    completed_hours: row.period?.completed_hours ?? 0,
  })
  editing.value = true
  dialogVisible.value = true
}

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    await saveCourse(form)
    ElMessage.success(editing.value ? '已更新' : '已新增')
    dialogVisible.value = false
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(row: Course): Promise<void> {
  try {
    await ElMessageBox.confirm(`确定删除课程「${row.id}」？`, '删除确认', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  try {
    await deleteCourse(row.id)
    ElMessage.success('已删除')
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  }
}

// ---------------------------------------------------------------- 详情抽屉

const detailVisible = ref(false)
const detailLoading = ref(false)
const detailCourse = ref<Course | null>(null)
const students = ref<Student[]>([])
const slots = ref<CourseSlot[]>([])
const absences = ref<AbsenceRecord[]>([])
const makeups = ref<StudentMakeup[]>([])
const changes = ref<CourseSlotChange[]>([])

async function openDetail(row: Course): Promise<void> {
  detailCourse.value = row
  detailVisible.value = true
  detailLoading.value = true
  try {
    const [a, b, c, d, e] = await Promise.all([
      courseStudents(row.id),
      courseSlots(row.id),
      courseAbsences(row.id),
      courseMakeups(row.id),
      courseSlotChanges(row.id),
    ])
    students.value = a
    slots.value = b
    absences.value = c
    makeups.value = d
    changes.value = e
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    detailLoading.value = false
  }
}

onMounted(paged.load)
</script>

<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h2>课程</h2>
        <p>Course 聚合 · 主键是字符串，无时间戳；容量/报名窗口/课时都是值对象</p>
      </div>
      <div class="page-head-actions">
        <el-button :icon="Refresh" :loading="paged.loading.value" @click="paged.load()">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">新增课程</el-button>
      </div>
    </div>

    <el-card shadow="never" class="page-card">
      <div class="page-toolbar">
        <el-input
          v-model="keyword"
          :prefix-icon="Search"
          placeholder="按课程类型筛选本页"
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
        <el-table-column prop="course_type_id" label="课程类型" width="140" />
        <el-table-column label="容量（已报/上限）" width="150">
          <template #default="{ row }">
            {{ row.capacity?.enrolled ?? 0 }} / {{ row.capacity?.max ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column label="报名窗口" min-width="200">
          <template #default="{ row }">
            {{ formatTimestamp(row.enrollment?.start_at) }}
            <br />
            {{ formatTimestamp(row.enrollment?.end_at) }}
          </template>
        </el-table-column>
        <el-table-column label="课时（已上/总）" width="140">
          <template #default="{ row }">
            {{ row.period?.completed_hours ?? 0 }} / {{ row.period?.total_hours ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" :icon="View" @click="openDetail(row)">详情</el-button>
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
      :title="editing ? '编辑课程' : '新增课程'"
      width="560px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="96px">
        <el-form-item label="课程类型 ID" prop="course_type_id">
          <el-input v-model="form.course_type_id" placeholder="如 ct-0001（少儿编程）" />
        </el-form-item>
        <el-divider content-position="left">容量</el-divider>
        <el-form-item label="容量" prop="capacity_max">
          <el-input-number v-model="form.capacity_max" :min="1" :max="2000" style="width: 100%" />
        </el-form-item>
        <el-form-item label="已报名" prop="capacity_enrolled">
          <el-input-number v-model="form.capacity_enrolled" :min="0" :max="2000" style="width: 100%" />
        </el-form-item>
        <el-divider content-position="left">报名窗口</el-divider>
        <el-form-item label="开始" prop="enroll_start_at">
          <el-date-picker v-model="form.enroll_start_at" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="截止" prop="enroll_end_at">
          <el-date-picker v-model="form.enroll_end_at" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="退课截止" prop="drop_deadline">
          <el-date-picker v-model="form.drop_deadline" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-divider content-position="left">教学周期</el-divider>
        <el-form-item label="开课" prop="period_start_at">
          <el-date-picker v-model="form.period_start_at" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结课" prop="period_end_at">
          <el-date-picker v-model="form.period_end_at" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="总课时" prop="total_hours">
          <el-input-number v-model="form.total_hours" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="已上" prop="completed_hours">
          <el-input-number v-model="form.completed_hours" :min="0" style="width: 100%" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">
          {{ editing ? '保存修改' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="detailVisible" :title="`课程 ${detailCourse?.id ?? ''} 关联数据`" size="720px">
      <div v-loading="detailLoading">
        <el-tabs>
          <el-tab-pane :label="`学员 ${students.length}`">
            <el-table :data="students" size="small" stripe>
              <el-table-column prop="id" label="ID" width="140" />
              <el-table-column prop="name" label="姓名" />
              <el-table-column label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="enumTag(STUDENT_STATUSES, row.status)" effect="light">
                    {{ enumLabel(STUDENT_STATUSES, row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          <el-tab-pane :label="`课次 ${slots.length}`">
            <el-table :data="slots" size="small" stripe>
              <el-table-column prop="id" label="ID" width="260" show-overflow-tooltip />
              <el-table-column label="时间">
                <template #default="{ row }">
                  周{{ ['', '日', '一', '二', '三', '四', '五', '六'][row.weekday ?? 0] }}
                  {{ formatTimeRange(row.time_range) }}
                </template>
              </el-table-column>
              <el-table-column label="讲师" prop="teacher_id" width="100" />
              <el-table-column label="教室" prop="classroom_id" width="100" />
            </el-table>
          </el-tab-pane>
          <el-tab-pane :label="`缺勤 ${absences.length}`">
            <el-table :data="absences" size="small" stripe>
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="student_id" label="学员" width="90" />
              <el-table-column label="类型" width="90">
                <template #default="{ row }">{{ enumLabel(ABSENCE_TYPES, row.absence_type) }}</template>
              </el-table-column>
              <el-table-column prop="reason" label="事由" />
            </el-table>
          </el-tab-pane>
          <el-tab-pane :label="`补课 ${makeups.length}`">
            <el-table :data="makeups" size="small" stripe>
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="student_id" label="学员" width="90" />
              <el-table-column label="状态" width="90">
                <template #default="{ row }">{{ enumLabel(MAKEUP_STATUSES, row.status) }}</template>
              </el-table-column>
              <el-table-column prop="target_date" label="补课时间">
                <template #default="{ row }">{{ formatTimestamp(row.target_date) }}</template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          <el-tab-pane :label="`调课 ${changes.length}`">
            <el-table :data="changes" size="small" stripe>
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column label="类型" width="90">
                <template #default="{ row }">{{ enumLabel(CHANGE_TYPES, row.change_type) }}</template>
              </el-table-column>
              <el-table-column prop="reason" label="事由" />
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-drawer>
  </div>
</template>
