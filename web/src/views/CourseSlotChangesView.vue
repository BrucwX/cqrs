<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'

import {
  changeCourseSlot,
  deleteCourseSlotChange,
  listCourseSlotChanges,
  type ChangeForm,
} from '@/api/courseSlotChanges'
import {
  CHANGE_TYPES,
  enumLabel,
  enumTag,
  formatTimestamp,
  type CourseSlotChange,
} from '@/api/types'
import { usePagedList } from '@/composables/usePagedList'
import PagerBar from '@/components/PagerBar.vue'

const paged = usePagedList<CourseSlotChange>((page, pageSize) =>
  listCourseSlotChanges({ page, page_size: pageSize }),
)

const keyword = ref('')
const rows = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return paged.items.value
  return paged.items.value.filter((row) =>
    [row.id, row.course_id, row.reason].some((field) =>
      String(field ?? '').toLowerCase().includes(kw),
    ),
  )
})

// ---------------------------------------------------------------- 表单

const dialogVisible = ref(false)
const saving = ref(false)
const formRef = ref<FormInstance>()

function blankForm(): ChangeForm {
  return {
    course_id: '',
    applicant_id: '',
    change_type: 1,
    reason: '',
    original_slot_id: '',
    original_date: '',
    original_teacher_id: '',
    original_classroom_id: '',
    original_start_time: '',
    original_end_time: '',
    target_start_at: '',
    target_end_at: '',
    target_teacher_id: '',
    target_classroom_id: '',
  }
}

const form = reactive<ChangeForm>(blankForm())

const rules: FormRules<ChangeForm> = {
  course_id: [{ required: true, message: '请填课程 ID', trigger: 'blur' }],
  applicant_id: [{ required: true, message: '请填申请人 ID', trigger: 'blur' }],
  reason: [{ required: true, message: '请填变更事由', trigger: 'blur' }],
  target_teacher_id: [{ required: true, message: '请填目标讲师 ID', trigger: 'blur' }],
  target_classroom_id: [{ required: true, message: '请填目标教室 ID', trigger: 'blur' }],
  target_start_at: [{ required: true, message: '请选目标开始时间', trigger: 'change' }],
  target_end_at: [{ required: true, message: '请选目标结束时间', trigger: 'change' }],
}

function openCreate(): void {
  Object.assign(form, blankForm())
  dialogVisible.value = true
}

async function submit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    await changeCourseSlot(form)
    ElMessage.success('已登记变更')
    dialogVisible.value = false
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(row: CourseSlotChange): Promise<void> {
  try {
    await ElMessageBox.confirm(`确定删除调课记录「${row.id}」？`, '删除确认', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  try {
    await deleteCourseSlotChange(row.id)
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
        <h2>调课记录</h2>
        <p>CourseSlotChange 聚合 · 登记即生效，没有审批环节</p>
      </div>
      <div class="page-head-actions">
        <el-button :icon="Refresh" :loading="paged.loading.value" @click="paged.load()">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">发起调课</el-button>
      </div>
    </div>

    <el-card shadow="never" class="page-card">
      <div class="page-toolbar">
        <el-input
          v-model="keyword"
          :prefix-icon="Search"
          placeholder="按课程 / 事由筛选本页"
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
        <el-table-column prop="course_id" label="课程" width="110" />
        <el-table-column label="类型" width="110">
          <template #default="{ row }">
            <el-tag :type="enumTag(CHANGE_TYPES, row.change_type)" effect="light">
              {{ enumLabel(CHANGE_TYPES, row.change_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="原计划" min-width="200">
          <template #default="{ row }">
            {{ row.original_plan?.start_time || '—' }}–{{ row.original_plan?.end_time || '—' }}
            <span v-if="row.original_plan?.classroom_id">@{{ row.original_plan.classroom_id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="目标计划" min-width="200">
          <template #default="{ row }">
            {{ formatTimestamp(row.target_plan?.target_start_at) }}
            <br />
            <span class="dim">讲师 {{ row.target_plan?.teacher_id ?? '—' }} · 教室 {{ row.target_plan?.classroom_id || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="事由" min-width="160" show-overflow-tooltip />
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

    <el-dialog v-model="dialogVisible" title="发起调课" width="620px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="课程 ID" prop="course_id">
          <el-input v-model="form.course_id" placeholder="如 C001" />
        </el-form-item>
        <el-form-item label="申请人 ID" prop="applicant_id">
          <el-input v-model="form.applicant_id" placeholder="原讲师或教务的 ID" />
        </el-form-item>
        <el-form-item label="变更类型" prop="change_type">
          <el-select v-model="form.change_type" style="width: 100%">
            <el-option v-for="option in CHANGE_TYPES" :key="option.value" :label="option.label" :value="option.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="事由" prop="reason">
          <el-input v-model="form.reason" type="textarea" :rows="2" placeholder="为什么要调课" />
        </el-form-item>

        <el-divider content-position="left">原计划（可选，用于留档）</el-divider>
        <el-form-item label="原课次 ID">
          <el-input v-model="form.original_slot_id" placeholder="被调动的槽位 ID" />
        </el-form-item>
        <el-form-item label="原日期">
          <el-date-picker v-model="form.original_date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="原时间">
          <div style="display: flex; gap: 8px; width: 100%">
            <el-time-picker v-model="form.original_start_time" format="HH:mm" value-format="HH:mm" placeholder="开始" style="flex: 1" />
            <el-time-picker v-model="form.original_end_time" format="HH:mm" value-format="HH:mm" placeholder="结束" style="flex: 1" />
          </div>
        </el-form-item>
        <el-form-item label="原讲师 ID">
          <el-input v-model="form.original_teacher_id" placeholder="选填" />
        </el-form-item>
        <el-form-item label="原教室 ID">
          <el-input v-model="form.original_classroom_id" placeholder="选填" />
        </el-form-item>

        <el-divider content-position="left">目标计划（必填）</el-divider>
        <el-form-item label="开始" prop="target_start_at">
          <el-date-picker v-model="form.target_start_at" type="datetime" value-format="YYYY-MM-DDTHH:mm" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结束" prop="target_end_at">
          <el-date-picker v-model="form.target_end_at" type="datetime" value-format="YYYY-MM-DDTHH:mm" style="width: 100%" />
        </el-form-item>
        <el-form-item label="目标讲师 ID" prop="target_teacher_id">
          <el-input v-model="form.target_teacher_id" placeholder="必须 &gt; 0" />
        </el-form-item>
        <el-form-item label="目标教室 ID" prop="target_classroom_id">
          <el-input v-model="form.target_classroom_id" placeholder="如 R101" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">登记</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.dim {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
