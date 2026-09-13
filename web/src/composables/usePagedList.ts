import type { Ref } from 'vue'
import { ref, shallowRef } from 'vue'

/** 往上翻页时用来报错的对象，界面只需要 message。 */
function messageOf(e: unknown): string {
    return e instanceof Error ? e.message : String(e)
}

export interface PagedList<T> {
    items: Ref<T[]>
    page: Ref<number>
    pageSize: Ref<number>
    loading: Ref<boolean>
    /** 当前页是否可能还有下一页。 */
    hasMore: Ref<boolean>
    error: Ref<string>
    load: () => Promise<void>
    goTo: (page: number) => Promise<void>
    setPageSize: (size: number) => Promise<void>
}

/**
 * 后端分页的一个现实约束：`XxxSet` 里只有 `repeated Xxx xxx = 1`，
 * **既没有 total 也没有 next_page_token**。所以：
 *
 *   - 算不出总页数，只能「本页拿满了 pageSize 条 ⇒ 大概率还有下一页」；
 *   - 因此翻页控件不显示总条数，也不显示页码总数（显示了就是编数字）。
 */
export function usePagedList<T>(
    fetcher: (page: number, pageSize: number) => Promise<T[]>,
    initialPageSize = 20,
): PagedList<T> {
    const items = shallowRef<T[]>([])
    const page = ref(1)
    const pageSize = ref(initialPageSize)
    const loading = ref(false)
    const hasMore = ref(false)
    const error = ref('')

    async function load(): Promise<void> {
        loading.value = true
        error.value = ''
        try {
            const rows = await fetcher(page.value, pageSize.value)
            items.value = rows
            hasMore.value = rows.length >= pageSize.value
        } catch (e) {
            error.value = messageOf(e)
            items.value = []
            hasMore.value = false
        } finally {
            loading.value = false
        }
    }

    async function goTo(target: number): Promise<void> {
        if (target < 1 || target === page.value) return
        page.value = target
        await load()
    }

    async function setPageSize(size: number): Promise<void> {
        pageSize.value = size
        page.value = 1
        await load()
    }

    return { items, page, pageSize, loading, hasMore, error, load, goTo, setPageSize }
}
