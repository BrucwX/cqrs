package product

// --- 值对象 (Value Object) ---

// Type 是「商品分类」值对象：**就一个 code 字符串**（offlineCourse / digital / service…）。
//
// 为什么不直接用 string：分类 code 是**跨聚合的约定键**（券的适用范围也按它匹配），
// 单独定义成一个类型之后，编译器会挡住「随手拿别的字符串当分类用」；
// 以后真要给它加东西（校验、展示名）也有地方挂。
//
// ⚠️ 分类是「这类商品是什么」，不是「这件商品怎么卖」—— 后者由具体商品类型表达。
type Type string
