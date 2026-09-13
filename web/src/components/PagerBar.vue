<script setup lang="ts">
import { computed } from 'vue'
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'

const props = defineProps<{
  page: number
  pageSize: number
  hasMore: boolean
  loading: boolean
  /** 本页实际拿到的条数 */
  count: number
}>()

const emit = defineEmits<{
  'update:page': [value: number]
  'update:pageSize': [value: number]
}>()

const canPrev = computed(() => props.page > 1 && !props.loading)
const canNext = computed(() => props.hasMore && !props.loading)

const sizeOptions = [10, 20, 50, 100]
</script>

<template>
  <div class="pager-bar">
    <span class="pager-hint">
      本页 {{ count }} 条
      <el-tooltip
        placement="top"
        content="后端 XxxSet 里只有 repeated 字段，没有 total 也没有 next_page_token，算不出总条数和总页数，所以这里不显示。"
      >
        <span class="pager-why">为什么没有总数？</span>
      </el-tooltip>
    </span>

    <div class="pager-actions">
      <span class="pager-size">
        每页
        <el-select
          :model-value="pageSize"
          size="small"
          style="width: 84px; margin: 0 4px"
          @change="(value: number) => emit('update:pageSize', value)"
        >
          <el-option v-for="size in sizeOptions" :key="size" :label="String(size)" :value="size" />
        </el-select>
        条
      </span>

      <el-button :disabled="!canPrev" :icon="ArrowLeft" @click="emit('update:page', page - 1)">
        上一页
      </el-button>
      <span class="pager-page">第 {{ page }} 页</span>
      <el-button :disabled="!canNext" @click="emit('update:page', page + 1)">
        下一页
        <el-icon class="el-icon--right"><ArrowRight /></el-icon>
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.pager-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 14px;
}

.pager-hint {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.pager-why {
  margin-left: 6px;
  color: var(--el-color-primary);
  cursor: help;
  text-decoration: underline dotted;
}

.pager-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pager-size {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.pager-page {
  min-width: 72px;
  font-size: 13px;
  text-align: center;
  color: var(--el-text-color-regular);
}
</style>
