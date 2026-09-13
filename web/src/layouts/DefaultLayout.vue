<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Calendar, Clock, Medal, Notebook, School, Switch, Tickets, Timer, User, UserFilled, Warning } from '@element-plus/icons-vue'

const route = useRoute()

const navs = [
  { path: '/students', label: '学员', icon: User },
  { path: '/teachers', label: '讲师', icon: UserFilled },
  { path: '/classrooms', label: '教室', icon: School },
  { path: '/courses', label: '课程', icon: Notebook },
  { path: '/course-slots', label: '课次', icon: Timer },
  { path: '/course-slot-changes', label: '调课', icon: Switch },
  { path: '/enrollments', label: '选课', icon: Tickets },
  { path: '/absences', label: '缺勤', icon: Warning },
  { path: '/makeups', label: '补课', icon: Clock },
  { path: '/qualifications', label: '资质', icon: Medal },
]

const active = computed(() => route.path)
</script>

<template>
  <el-container class="app-shell">
    <el-aside width="210px" class="app-aside">
      <div class="app-brand">
        <el-icon :size="20"><Calendar /></el-icon>
        <span>课程排课后台</span>
      </div>
      <el-menu :default-active="active" router class="app-menu">
        <el-menu-item v-for="nav in navs" :key="nav.path" :index="nav.path">
          <el-icon><component :is="nav.icon" /></el-icon>
          <span>{{ nav.label }}</span>
        </el-menu-item>
      </el-menu>
      <div class="app-footnote">
        course_scheduling
        <br />
        <span>读写分离 · CQRS</span>
      </div>
    </el-aside>

    <el-main class="app-main">
      <router-view />
    </el-main>
  </el-container>
</template>

<style scoped>
.app-shell {
  height: 100vh;
}

.app-aside {
  display: flex;
  flex-direction: column;
  background: #1f2d3d;
  color: #fff;
}

.app-brand {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 56px;
  padding: 0 18px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.5px;
  border-bottom: 1px solid rgb(255 255 255 / 10%);
}

.app-menu {
  flex: 1;
  border-right: none;
  background: transparent;
  --el-menu-bg-color: transparent;
  --el-menu-text-color: rgb(255 255 255 / 70%);
  --el-menu-hover-bg-color: rgb(255 255 255 / 8%);
  --el-menu-active-color: #fff;
}

.app-menu :deep(.el-menu-item.is-active) {
  background: var(--el-color-primary);
  color: #fff;
}

.app-footnote {
  padding: 14px 18px;
  font-size: 12px;
  line-height: 1.7;
  color: rgb(255 255 255 / 40%);
  border-top: 1px solid rgb(255 255 255 / 10%);
}

.app-main {
  padding: 20px 24px;
  background: #f5f7fa;
  overflow-y: auto;
}
</style>
