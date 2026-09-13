<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, Search } from '@element-plus/icons-vue'

import {
  deleteClassroom,
  listClassrooms,
  saveClassroom,
  type ClassroomForm,
} from '@/api/classrooms'
import {
  CLASSROOM_STATUSES,
  enumLabel,
  enumTag,
  type Classroom,
} from '@/api/types'
import { usePagedList } from '@/composables/usePagedList'
import PagerBar from '@/components/PagerBar.vue'

const paged = usePagedList<Classroom>((page, pageSize) =>
  listClassrooms({ page, page_size: pageSize }),
)

const keyword = ref('')

const rows = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return paged.items.value
  return paged.items.value.filter((row) =>
    [row.id, row.location?.building, row.location?.room]
      .filter((field): field is string => Boolean(field))
      .some((field) => field.toLowerCase().includes(kw)),
  )
})

const dialogVisible = ref(false)
const saving = ref(false)
const editing = ref(false)
const formRef = ref<FormInstance>()

function blankForm(): ClassroomForm {
  // 教室主键是字符串，新增时留空让后端生成，不要自己编一个
  return { id: '', building: '', floor: 1, room: '', capacity: 30, status: 1 }
}

const form = reactive<ClassroomForm>(blankForm())

const rules: FormRules<ClassroomForm> = {
  building: [{ required: true, message: '请填楼栋', trigger: 'blur' }],
  room: [{ required: true, message: '请填房间号', trigger: 'blur' }],
  floor: [{ required: true, message: '请填楼层', trigger: 'blur' }],
  capacity: [{ required: true, message: '请填容量', trigger: 'blur' }],
}

function openCreate(): void {
  Object.assign(form, blankForm())
  editing.value = false
  dialogVisible.value = true
}

function openEdit(row: Classroom): void {
  Object.assign(form, {
    id: row.id,
    building: row.location?.building ?? '',
    floor: row.location?.floor ?? 1,
    room: row.location?.room ?? '',
    capacity: row.capacity ?? 0,
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
    await saveClassroom(form)
    ElMessage.success(editing.value ? '已更新' : '已新增')
    dialogVisible.value = false
    await paged.load()
  } catch (error) {
    ElMessage.error((error as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(row: Classroom): Promise<void> {
  const label = [row.location?.building, row.location?.room].filter(Boolean).join(' ') || row.id
  try {
    await ElMessageBox.confirm(`确定删除教室「${label}」？此操作不可撤销。`, '删除确认', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    await deleteClassroom(row.id)
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

/** 已分配 / 总容量，顺手标一下占用率，超了给红色。 */
function occupancy(row: Classroom): number {
  const capacity = row.capacity ?? 0
  if (capacity <= 0) return 0
  return Math.round(((row.allocated ?? 0) / capacity) * 100)
}

onMounted(paged.load)
</script>

<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h2>教室</h2>
        <p>Classroom 聚合 · 主键是字符串，且没有时间戳字段</p>
      </div>
      <div class="page-head-actions">
        <el-button :icon="Refresh" :loading="paged.loading.value" @click="paged.load()">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">新增教室</el-button>
      </div>
    </div>

    <el-alert type="info" :closable="false" show-icon class="page-note">
      <template #title>
        <code>allocated</code> 为 0 时后端不会返回这个字段（生成的结构体带
        <code>omitempty</code>），所以这里统一按 0 处理。
      </template>
    </el-alert>

    <el-card shadow="never" class="page-card">
      <div class="page-toolbar">
        <el-input
          v-model="keyword"
          :prefix-icon="Search"
          placeholder="按楼栋 / 房间号筛选本页"
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
        <el-table-column label="位置" min-width="200">
          <template #default="{ row }">
            {{ row.location?.building || '—' }} · {{ row.location?.floor ?? '—' }} 层 ·
            {{ row.location?.room || '—' }}
          </template>
        </el-table-column>
        <el-table-column label="容量" width="160">
          <template #default="{ row }">
            {{ row.allocated ?? 0 }} / {{ row.capacity ?? 0 }}
            <el-progress
              :percentage="occupancy(row)"
              :stroke-width="6"
              :status="occupancy(row) >= 100 ? 'exception' : undefined"
              style="margin-top: 4px"
            />
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="enumTag(CLASSROOM_STATUSES, row.status)" effect="light">
              {{ enumLabel(CLASSROOM_STATUSES, row.status) }}
            </el-tag>
          </template>
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
      :title="editing ? '编辑教室' : '新增教室'"
      width="520px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="88px">
        <el-form-item label="楼栋" prop="building">
          <el-input v-model="form.building" placeholder="如：A座" />
        </el-form-item>
        <el-form-item label="楼层" prop="floor">
          <el-input-number v-model="form.floor" :min="1" :max="200" style="width: 100%" />
        </el-form-item>
        <el-form-item label="房间号" prop="room">
          <el-input v-model="form.room" placeholder="如：301" />
        </el-form-item>
        <el-form-item label="容量" prop="capacity">
          <el-input-number v-model="form.capacity" :min="1" :max="2000" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" style="width: 100%">
            <el-option
              v-for="option in CLASSROOM_STATUSES"
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
