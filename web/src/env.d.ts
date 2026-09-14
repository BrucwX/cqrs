/// <reference types="vite/client" />

/**
 * Vite 8 的 vite/client 已不再声明 `*.vue`（Vite 6 起移除），
 * 而本容器没装 Volar 扩展。这个 shim 让 tsserver 能识别 .vue 导入；
 * vue-tsc 自己原生支持 .vue，不受影响。
 */
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
