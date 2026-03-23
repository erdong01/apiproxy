<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline"
        @keyup.enter="onSubmit">

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable filterable placeholder="请选择状态" style="width: 180px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="Key">
          <el-input v-model="searchInfo.key" clearable placeholder="请输入Key" />
        </el-form-item>

        <el-form-item label="任务ID">
          <el-input v-model="searchInfo.generateTaskId" clearable placeholder="请输入任务ID" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true"
            v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
          @click="onDelete">删除</el-button>

      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="id"
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <!-- <el-table-column align="left" label="id" prop="id" width="120" /> -->
        <el-table-column align="left" label="任务id" prop="generateTaskId" width="220" />
        <el-table-column align="left" label="模型" prop="model" width="200" />
        <el-table-column align="left" label="状态" prop="status" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" effect="light">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="视频花费tokens" prop="completionTokens" width="200" />
        <el-table-column align="left" label="消耗的总tokens" prop="totalTokens" width="200" />
        <el-table-column align="left" label="模型费用" prop="vendorAmount" width="120" />
        <el-table-column align="left" label="销售成本" prop="retailAmount" width="120" />
        <el-table-column align="left" label="失败原因" prop="errorCode" width="120" />
        <el-table-column align="left" label="失败原因" prop="errorMessage" width="120" />
        <el-table-column align="left" label="key字段" prop="key" width="120" />
        <el-table-column align="left" label="含参考视频" width="180">
          <template #default="{ row }">
            {{ row.draft_task_id || row.draftTaskId ? '是' : '否' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>查看</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false"
      :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增' : '编辑' }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="id字段:" prop="id">
          <el-input v-model.number="formData.id" :clearable="true" placeholder="请输入id字段" />
        </el-form-item>
        <el-form-item label="创建时间:" prop="createdAt">
          <el-date-picker v-model="formData.createdAt" type="date" style="width:100%" placeholder="选择日期"
            :clearable="true" />
        </el-form-item>
        <el-form-item label="更新时间:" prop="updatedAt">
          <el-date-picker v-model="formData.updatedAt" type="date" style="width:100%" placeholder="选择日期"
            :clearable="true" />
        </el-form-item>
        <el-form-item label="deletedAt字段:" prop="deletedAt">
          <el-date-picker v-model="formData.deletedAt" type="date" style="width:100%" placeholder="选择日期"
            :clearable="true" />
        </el-form-item>
        <el-form-item label="供应商任务id:" prop="generateTaskId">
          <el-input v-model="formData.generateTaskId" :clearable="true" placeholder="请输入供应商任务id" />
        </el-form-item>
        <el-form-item label="用户id:" prop="userId">
          <el-input v-model.number="formData.userId" :clearable="true" placeholder="请输入用户id" />
        </el-form-item>
        <el-form-item label="ai模型id:" prop="modelId">
          <el-input v-model.number="formData.modelId" :clearable="true" placeholder="请输入ai模型id" />
        </el-form-item>
        <el-form-item label="状态:" prop="status">
          <el-input v-model="formData.status" :clearable="true" placeholder="请输入状态" />
        </el-form-item>
        <el-form-item label="存储原始生成参数 (如 seed, motion_bucket_id 等):" prop="params">
          // 此字段为json结构，可以前端自行控制展示和数据绑定模式 需绑定json的key为 formData.params 后端会按照json的类型进行存取
          {{ formData.params }}
        </el-form-item>
        <el-form-item label="失败原因:" prop="errorMessage">
          <el-input v-model="formData.errorMessage" :clearable="true" placeholder="请输入失败原因" />
        </el-form-item>
        <el-form-item label="输出视频花费的 token 数:" prop="completionTokens">
          <el-input v-model.number="formData.completionTokens" :clearable="true" placeholder="请输入输出视频花费的 token 数" />
        </el-form-item>
        <el-form-item label="本次请求消耗的总 token 数:" prop="totalTokens">
          <el-input v-model.number="formData.totalTokens" :clearable="true" placeholder="请输入本次请求消耗的总 token 数" />
        </el-form-item>
        <el-form-item label="key字段:" prop="key">
          <el-input v-model="formData.key" :clearable="true" placeholder="请输入key字段" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true"
      :before-close="closeDetailShow" title="查看">
      <div class="detail-section-title">基础信息</div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">
          {{ formatDetailValue(detailForm.id) }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTagType(detailForm.status)" effect="light">
            {{ getStatusLabel(detailForm.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="任务ID">
          {{ formatDetailValue(detailForm.generateTaskId) }}
        </el-descriptions-item>
        <el-descriptions-item label="请求ID">
          {{ formatDetailValue(detailForm.requestId) }}
        </el-descriptions-item>
        <!-- <el-descriptions-item label="用户ID">
          {{ formatDetailValue(detailForm.userId) }}
        </el-descriptions-item> -->
        <el-descriptions-item label="模型">
          {{ formatDetailValue(detailForm.model) }}
        </el-descriptions-item>
        <el-descriptions-item label="供应商Key">
          {{ formatDetailValue(detailForm.key) }}
        </el-descriptions-item>
        <el-descriptions-item label="服务等级">
          {{ formatDetailValue(detailForm.serviceTier) }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatDateTime(detailForm.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间">
          {{ formatDateTime(detailForm.updatedAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="任务创建时间">
          {{ formatUnixSeconds(detailForm.taskCreatedAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="任务状态更新时间">
          {{ formatUnixSeconds(detailForm.taskUpdatedAt) }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="detail-section-title">生成参数</div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="分辨率">
          {{ formatDetailValue(detailForm.resolution) }}
        </el-descriptions-item>
        <el-descriptions-item label="宽高比">
          {{ formatDetailValue(detailForm.ratio) }}
        </el-descriptions-item>
        <el-descriptions-item label="时长">
          {{ formatNumberUnit(detailForm.duration, '秒') }}
        </el-descriptions-item>
        <el-descriptions-item label="帧数">
          {{ formatDetailValue(detailForm.frames) }}
        </el-descriptions-item>
        <el-descriptions-item label="帧率">
          {{ formatDetailValue(detailForm.framesPerSecond) }}
        </el-descriptions-item>
        <el-descriptions-item label="Seed">
          {{ formatDetailValue(detailForm.seed) }}
        </el-descriptions-item>
        <el-descriptions-item label="生成音频">
          {{ formatBooleanValue(detailForm.generateAudio) }}
        </el-descriptions-item>
        <el-descriptions-item label="Draft">
          {{ formatBooleanValue(detailForm.draft) }}
        </el-descriptions-item>
        <el-descriptions-item label="Draft任务ID">
          {{ formatDetailValue(detailForm.draftTaskId) }}
        </el-descriptions-item>
        <el-descriptions-item label="超时时间">
          {{ formatNumberUnit(detailForm.executionExpiresAfter, '秒') }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="detail-section-title">计费与错误</div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="视频花费 Tokens">
          {{ formatDetailValue(detailForm.completionTokens) }}
        </el-descriptions-item>
        <el-descriptions-item label="总 Tokens">
          {{ formatDetailValue(detailForm.totalTokens) }}
        </el-descriptions-item>
        <el-descriptions-item label="模型费用">
          {{ formatDetailValue(detailForm.vendorAmount) }}
        </el-descriptions-item>
        <el-descriptions-item label="销售成本">
          {{ formatDetailValue(detailForm.retailAmount) }}
        </el-descriptions-item>
        <el-descriptions-item label="失败原因" :span="2">
          {{ formatDetailValue(detailForm.errorMessage) }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="detail-section-title">原始参数</div>
      <pre class="detail-json-block">{{ formatJsonValue(detailForm.params) }}</pre>

      <div class="detail-section-title">输出内容</div>
      <pre class="detail-json-block">{{ formatJsonValue(detailForm.content) }}</pre>
    </el-drawer>

  </div>
</template>

<script setup>
import {
  createPqAiTask,
  deletePqAiTask,
  deletePqAiTaskByIds,
  updatePqAiTask,
  findPqAiTask,
  getPqAiTaskList
} from '@/api/ai/pqAiTask'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"




defineOptions({
  name: 'PqAiTask'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

const statusOptions = [
  { label: '排队中', value: 'queued' },
  { label: '任务运行中', value: 'running' },
  { label: '取消任务', value: 'cancelled' },
  { label: '任务成功', value: 'succeeded' },
  { label: '任务失败', value: 'failed' },
  { label: '任务超时', value: 'expired' },
]

const statusLabelMap = statusOptions.reduce((map, item) => {
  map[item.value] = item.label
  return map
}, {})

const statusTagTypeMap = {
  queued: 'info',
  running: 'warning',
  cancelled: '',
  succeeded: 'success',
  failed: 'danger',
  expired: 'danger',
}

const getStatusLabel = (status) => statusLabelMap[status] || status || '-'
const getStatusTagType = (status) => statusTagTypeMap[status] || 'info'
const formatDetailValue = (value) => {
  if (value === null || value === undefined || value === '') return '-'
  return value
}

const formatBooleanValue = (value) => {
  if (value === null || value === undefined) return '-'
  return value ? '是' : '否'
}

const formatDateTime = (value) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  return date.toLocaleString('zh-CN', { hour12: false })
}

const formatUnixSeconds = (value) => {
  if (!value) return '-'
  const date = new Date(value * 1000)
  if (Number.isNaN(date.getTime())) return String(value)
  return date.toLocaleString('zh-CN', { hour12: false })
}

const formatNumberUnit = (value, unit) => {
  if (value === null || value === undefined || value === '') return '-'
  return `${value}${unit}`
}

const formatJsonValue = (value) => {
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2)
    } catch (error) {
      return value
    }
  }
  try {
    return JSON.stringify(value, null, 2)
  } catch (error) {
    return String(value)
  }
}

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  id: undefined,
  createdAt: new Date(),
  updatedAt: new Date(),
  deletedAt: new Date(),
  generateTaskId: '',
  userId: undefined,
  modelId: undefined,
  status: '',
  params: {},
  errorMessage: '',
  completionTokens: undefined,
  totalTokens: undefined,
  key: '',
})



// 验证规则
const rule = reactive({
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const initialSearchInfo = () => ({
  status: '',
  key: '',
  generateTaskId: '',
})
const searchInfo = ref(initialSearchInfo())
// 重置
const onReset = () => {
  searchInfo.value = initialSearchInfo()
  page.value = 1
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getPqAiTaskList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deletePqAiTaskFunc(row)
  })
}

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: '请选择要删除的数据'
      })
      return
    }
    multipleSelection.value &&
      multipleSelection.value.map(item => {
        ids.push(item.id)
      })
    const res = await deletePqAiTaskByIds({ ids })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updatePqAiTaskFunc = async (row) => {
  const res = await findPqAiTask({ id: row.id })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}


// 删除行
const deletePqAiTaskFunc = async (row) => {
  const res = await deletePqAiTask({ id: row.id })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    id: undefined,
    createdAt: new Date(),
    updatedAt: new Date(),
    deletedAt: new Date(),
    generateTaskId: '',
    userId: undefined,
    modelId: undefined,
    status: '',
    params: {},
    errorMessage: '',
    completionTokens: undefined,
    totalTokens: undefined,
    key: '',
  }
}
// 弹窗确定
const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    let res
    switch (type.value) {
      case 'create':
        res = await createPqAiTask(formData.value)
        break
      case 'update':
        res = await updatePqAiTask(formData.value)
        break
      default:
        res = await createPqAiTask(formData.value)
        break
    }
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '创建/更改成功'
      })
      closeDialog()
      getTableData()
    }
  })
}

const detailForm = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findPqAiTask({ id: row.id })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}


</script>

<style scoped>
.detail-section-title {
  margin: 0 0 12px;
  font-size: 15px;
  font-weight: 600;
  color: #1f2a37;
}

.detail-section-title:not(:first-child) {
  margin-top: 20px;
}

.detail-json-block {
  margin: 0;
  padding: 12px 14px;
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.6;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #f8fafc;
  color: #334155;
  font-size: 13px;
}
</style>
