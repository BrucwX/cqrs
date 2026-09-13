<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Connection, Refresh, Search, User } from '@element-plus/icons-vue'

import {
  assignClassroom,
  assignCourse,
  assignTeacher,
  deleteCourseSlot,
  listCourseSlots,
} from '@/api/courseSlots'
import { listTeachers } from '@/api/teachers'
import { listClassrooms } from '@/api/classrooms'
import { listCourses } from '@/api/courses'
import {
  WEEKDAYS,
  enumLabel,
  formatTimeRange,
  formatTimestamp,
  type CourseSlot,
} from '@/api/types'
import { usePagedList } from '@/composables/usePagedList'
import PagerBar from '@/components/PagerBar.vue'

const paged = usePagedList<CourseSlot>((page, pageSize) =>
  listCourseSlots({ page, page_size: pageSize }),
)

const keyword = ref('')
const rows = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return paged.items.value
  return paged.items.value.filter((row) =>
    [row.id, row.course_id, row.classroom_id].some((field) => field?.toLowerCase().includes(kw)),
  )
})

const selected = ref<CourseSlot[]>([])

function onSelectionChange(selection: CourseSlot[]): void {
  selected.value = selection
}

function slotIds(): string[] {
  return selected.value.map((s) => s.id)
}

function requireSelection(): boolean {
  if (selected.value.length === 0) {
    ElMessage.warning('先勾选要操作的课次')
    return false
  }
  return true
}

// ---------------------------------------------------------------- 安排

type AssignTarget = 'classroom' | 'teacher' | 'course'

const assignVisible = ref(false)
const assignTarget = ref<AssignTarget>('classroom')
const assignSaving = ref(false)
const assignValue = ref('')
const options = ref<{ id: string; label: string }[]>([])

async function openAssign(target: AssignTarget): Promise<void> {
  if (!requireSelection()) return
  assignTarget.value = target
  assignValue.value = ''
  assignVisible.value = true
  try {
    if (target === 'classroom') {
      options.value = (await listClassrooms({ page: 1, page_size: 100 })).map((c) => ({
        id: c.id,
        label: [c.location?.building, c.location?.room].filter(Boolean).join(' ') || c.id,
      }))
    } else if (target === 'teacher') {
      options.value = (await listTeachers({ page: 1, page_size: 100 })).map((t) => ({
        id: t.id,
        label: `${t.name}（${t.id}）`,
      }))
    } else {
      options.value = (await listCourses({ page: 1, page_size: 100 })).map((c) => ({
        id: c.id,
        label: `${c.id}（${c.course_type_id ?? ''}）`,
      }))
    }
  } catch (error) {
    ElMessage.error((error as Error).message)
    assignVisible.value = false
  }
}

async function submitAssign(): Promise<void> {
  if (!assignValue.value) {
    ElMessage.warning('先选择目标')
    return
  }
  assignSaving.value = true
  try {
    const ids = slotIds()
    if (assignTarget.value === 'classroom') await assignClassroom(ids, assignValue.value)
    else if (assignTarget.value === 'teacher') await assignTeacher(ids, assignValue.value)
    else await assignCourse(ids, assignValue.value)
    ElMessage.success(`已应用到 ${ids.length} 个课次`)
    assignVisible.value = false
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    assignSaving.value = false
  }
}

async function remove(row: CourseSlot): Promise<void> {
  try {
    await ElMessageBox.confirm(`确定删除课次「${row.id}」？`, '删除确认', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  try {
    await deleteCourseSlot(row.id)
    ElMessage.success('已删除')
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  }
}

function teacherText(row: CourseSlot): string {
  const tid = row.teacher_id
  if (tid === undefined || tid === null || tid === -1 || tid === '-1') return '未安排'
  return String(tid)
}

onMounted(paged.load)
</script>

<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h2>课次</h2>
        <p>CourseSlot 聚合 · 没有 Save 接口，只有「安排教室 / 讲师 / 课程」三个 Assign 和 Delete</p>
      </div>
      <div class="page-head-actions">
        <el-button :icon="Refresh" :loading="paged.loading.value" @click="paged.load()">
          刷新
        </el-button>
        <el-button :icon="Connection" @click="openAssign('classroom')">安排教室</el-button>
        <el-button :icon="User" @click="openAssign('teacher')">安排讲师</el-button>
        <el-button type="primary" @click="openAssign('course')">安排课程</el-button>
      </div>
    </div>

    <el-alert type="info" :closable="false" show-icon class="page-note">
      <template #title>
        勾选一行或多行后点上面的按钮；Assign 接口是批量的（<code>slot_ids</code> 数组）。讲师 ID 为
        -1、教室为空都表示「未安排」。
      </template>
    </el-alert>

    <el-card shadow="never" class="page-card">
      <div class="page-toolbar">
        <el-input
          v-model="keyword"
          :prefix-icon="Search"
          placeholder="按课程 / 教室筛选本页"
          clearable
          style="width: 320px"
        />
        <span class="page-toolbar-count">已选 {{ selected.length }} · 本页 {{ rows.length }} / {{ paged.items.value.length }} 条</span>
      </div>

      <el-alert
        v-if="paged.error.value"
        type="error"
        :closable="false"
        show-icon
        :title="paged.error.value"
        class="page-note"
      />

      <el-table
        v-loading="paged.loading.value"
        :data="rows"
        stripe
        empty-text="没有数据"
        @selection-change="onSelectionChange"
      >
        <el-table-column type="selection" width="44" />
        <el-table-column prop="course_id" label="课程" width="110" />
        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ enumLabel(WEEKDAYS, row.weekday) }} {{ formatTimeRange(row.time_range) }}
          </template>
        </el-table-column>
        <el-table-column label="讲师" width="120">
          <template #default="{ row }">{{ teacherText(row) }}</template>
        </el-table-column>
        <el-table-column label="教室" width="120">
          <template #default="{ row }">{{ row.classroom_id || '未安排' }}</template>
        </el-table-column>
        <el-table-column label="更新时间" width="160">
          <template #default="{ row }">{{ formatTimestamp(row.updated_at) }}</template>
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

    <el-dialog
      v-model="assignVisible"
      :title="`安排${assignTarget === 'classroom' ? '教室' : assignTarget === 'teacher' ? '讲师' : '课程'}`"
      width="460px"
    >
      <p class="assign-hint">将把下面选择应用到已勾选的 {{ selected.length }} 个课次。</p>
      <el-select v-model="assignValue" filterable placeholder="请选择" style="width: 100%">
        <el-option v-for="option in options" :key="option.id" :label="option.label" :value="option.id" />
      </el-select>
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" :loading="assignSaving" @click="submitAssign">应用</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.assign-hint {
  margin: 0 0 10px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
</style>
