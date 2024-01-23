package bptree

// B+Tree 多路平衡搜索树
// 1. 每个节点最多有 m 个子节点
// 2. 每个节点最多有 m-1 个关键字
// 3. 根节点至少有两个子节点
// 4. 除根节点外，所有非叶子节点至少有 m/2 个子节点
// 5. 所有叶子节点都在同一层
// 6. 每个节点的关键字都按照升序排列
// 7. 每个节点的子节点指针数与关键字数相同
// 8. 非叶子节点的关键字 k1,k2,k3,...,ki 代表子节点指针 p1,p2,p3,...,pi-1,pi 的范围
//    其中 p1 指向关键字小于 k1 的子树，pi 指向关键字大于 ki 的子树
//    p1,p2,p3,...,pi-1,pi 的关键字范围为 k1,k2,k3,...,ki-1,ki
// 9. 叶子节点的关键字 k1,k2,k3,...,ki 代表数据指针 p1,p2,p3,...,pi 的范围
//    其中 p1 指向关键字等于 k1 的数据，pi 指向关键字等于 ki 的数据
//    p1,p2,p3,...,pi 的关键字范围为 k1,k2,k3,...,ki-1,ki
