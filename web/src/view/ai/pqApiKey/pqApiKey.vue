<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline"
        @keyup.enter="onSubmit">

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
          <el-form-item label="状态">
            <el-select v-model="searchInfo.status" clearable placeholder="请选择状态" style="width: 160px">
              <el-option label="启用" :value="1" />
              <el-option label="禁用" :value="2" />
            </el-select>
          </el-form-item>
        </template>

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
        <!-- <el-table-column align="left" label="id字段" prop="id" width="120" /> -->
        <!-- 
            <el-table-column align="left" label="updatedAt字段" prop="updatedAt" width="180">
   <template #default="scope">{{ formatDate(scope.row.updatedAt) }}</template>
</el-table-column>
            <el-table-column align="left" label="deletedAt字段" prop="deletedAt" width="180">
   <template #default="scope">{{ formatDate(scope.row.deletedAt) }}</template>
</el-table-column>
            <el-table-column align="left" label="createdAt字段" prop="createdAt" width="180">
   <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
</el-table-column> -->
        <!-- <el-table-column align="left" label="用户id" prop="userId" width="120" /> -->

        <!-- <el-table-column align="left" label="ai模型" prop="aiModelId" width="350">
   <template #default="scope">{{ getModelName(scope.row.aiModelId) }}</template>
</el-table-column> -->
        <el-table-column align="left" label="用户名" prop="UserName" width="200" />
        <el-table-column align="left" label="用户key" prop="UserKey" width="300" />
        <el-table-column align="left" label="调用模型KEY" prop="key" width="300" />
        <el-table-column align="left" label="状态" prop="status" width="150">
          <template #default="scope">
            <el-switch v-model="scope.row.status" :loading="switchLoadingMap[scope.row.id]" :active-value="1"
              :inactive-value="2" active-text="启用" inactive-text="禁用" @change="() => handleStatusSwitch(scope.row)" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="速率" prop="rate" width="120" />
        <!-- <el-table-column align="left" label="拥有tokens数" prop="totalTokens" width="200" /> -->

        <!-- <el-table-column align="left" label="已消耗tokens" prop="useTokens" width="200" /> -->

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>查看</el-button>
            <el-button type="primary" link icon="edit" class="table-button"
              @click="updatePqApiKeyFunc(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
        <!-- <el-form-item label="id字段:" prop="id">
    <el-input v-model.number="formData.id" :clearable="true" placeholder="请输入id字段" />
</el-form-item>
            <el-form-item label="updatedAt字段:" prop="updatedAt">
    <el-date-picker v-model="formData.updatedAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
            <el-form-item label="deletedAt字段:" prop="deletedAt">
    <el-date-picker v-model="formData.deletedAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
            <el-form-item label="createdAt字段:" prop="createdAt">
    <el-date-picker v-model="formData.createdAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item> -->
        <!-- <el-form-item label="用户id:" prop="userId">
    <el-input v-model.number="formData.userId" :clearable="true" placeholder="请输入用户id" />
</el-form-item> -->
        <!-- <el-form-item label="ai模型:" prop="aiModelId">
    <el-select v-model="formData.aiModelId" :clearable="true" placeholder="请选择ai模型" style="width:100%">
      <el-option v-for="item in modelOptions" :key="item.id" :label="item.name + ' (' + item.provider + ' ' + item.version + ')'" :value="item.id" />
    </el-select>
</el-form-item> -->
        <el-form-item label="用户名:" prop="UserName">
          <el-input v-model="formData.UserName" :maxlength="64" :clearable="true" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="用户KEY:" prop="UserKey">
          <el-input v-model="formData.UserKey" :maxlength="64" :clearable="true" placeholder="请输入密钥" />
        </el-form-item>
        <el-form-item label="调用模型KEY:" prop="key">
          <el-input v-model="formData.key" :clearable="true" :maxlength="64" placeholder="请输入密钥" />
        </el-form-item>
        <!-- <el-form-item label="拥有tokens数:" prop="totalTokens">
    <el-input v-model.number="formData.totalTokens" :clearable="true" placeholder="请输入拥有tokens数" />
</el-form-item> -->
        <!-- <el-form-item label="已消耗tokens:" prop="useTokens">
    <el-input v-model="formData.useTokens" :clearable="true" placeholder="请输入已消耗tokens" />
</el-form-item> -->
        <el-form-item label="速率:" prop="rate">
          <el-input v-model.number="formData.rate" :type="number" :clearable="true" placeholder="请输入速率" />
        </el-form-item>
        <el-form-item label="状态:" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="2" active-text="启用"
            inactive-text="禁用" />
        </el-form-item>
      </el-form>

    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true"
      :before-close="closeDetailShow" title="查看">
      <el-descriptions :column="1" border>
        <!-- <el-descriptions-item label="id字段">
    {{ detailForm.id }}
</el-descriptions-item> -->
        <!-- <el-descriptions-item label="updatedAt字段">
    {{ detailForm.updatedAt }}
</el-descriptions-item> -->
        <!-- <el-descriptions-item label="deletedAt字段">
    {{ detailForm.deletedAt }}
</el-descriptions-item> -->
        <!-- <el-descriptions-item label="createdAt字段">
    {{ detailForm.createdAt }}
</el-descriptions-item> -->
        <!-- <el-descriptions-item label="用户id">
    {{ detailForm.userId }}
</el-descriptions-item> -->
        <!-- <el-descriptions-item label="ai模型">
    {{ getModelName(detailForm.aiModelId) }}
</el-descriptions-item> -->
        <el-descriptions-item label="用户名">
          {{ detailForm.UserName }}
        </el-descriptions-item>
        <el-descriptions-item label="用户KEY">
          {{ detailForm.UserKey }}
        </el-descriptions-item>
        <el-descriptions-item label="调用模型KEY">
          {{ detailForm.key }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-switch :model-value="detailForm.status" :active-value="1" :inactive-value="2" active-text="启用"
            inactive-text="禁用" disabled />
        </el-descriptions-item>
        <el-descriptions-item label="速率">
          {{ detailForm.rate }}
        </el-descriptions-item>
        <!-- <el-descriptions-item label="拥有tokens数">
          {{ detailForm.totalTokens }}
        </el-descriptions-item>
        <el-descriptions-item label="已消耗tokens">
          {{ detailForm.useTokens }}
        </el-descriptions-item> -->
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import {
  createPqApiKey,
  deletePqApiKey,
  deletePqApiKeyByIds,
  updatePqApiKey,
  findPqApiKey,
  getPqApiKeyList
} from '@/api/ai/pqApiKey'
import { getPqAiModelList } from '@/api/ai/pqAiModel'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"




defineOptions({
  name: 'PqApiKey'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  id: undefined,
  updatedAt: new Date(),
  deletedAt: new Date(),
  createdAt: new Date(),
  userId: undefined,
  aiModelId: undefined,
  key: '',
  totalTokens: undefined,
  useTokens: '',
  UserName: '',
  status: 1,
  rate:0
})



// 验证规则
const rule = reactive({
  UserName: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  UserKey: [
    { required: true, message: '请输入用户KEY', trigger: 'blur' }
  ],
  key: [
    { required: true, message: '请输入调用模型KEY', trigger: 'blur' }
  ],
  rate: [
    { required: true, message: '请输入速率', trigger: 'blur' }
  ]
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const switchLoadingMap = ref({})
const initialSearchInfo = () => ({
  status: undefined,
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
  const table = await getPqApiKeyList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// ai模型下拉选项列表
const modelOptions = ref([])

// 获取ai模型列表
const getModelOptions = async () => {
  const res = await getPqAiModelList({ page: 1, pageSize: 999 })
  if (res.code === 0) {
    modelOptions.value = res.data.list
  }
}

// 根据模型id获取模型名称
const getModelName = (id) => {
  const model = modelOptions.value.find(item => item.id === id)
  return model ? `${model.name} (${model.provider} )` : id
}

// 初始化加载模型选项
getModelOptions()


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
    deletePqApiKeyFunc(row)
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
    const res = await deletePqApiKeyByIds({ ids })
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
const updatePqApiKeyFunc = async (row) => {
  const res = await findPqApiKey({ id: row.id })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}

const handleStatusSwitch = async (row) => {
  const previousStatus = row.status === 1 ? 2 : 1
  switchLoadingMap.value[row.id] = true

  try {
    const res = await updatePqApiKey({ ...row, status: row.status })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: `${row.status === 1 ? '启用' : '禁用'}成功`
      })
      return
    }

    row.status = previousStatus
    ElMessage({
      type: 'error',
      message: res.msg || '状态切换失败'
    })
  } catch (error) {
    row.status = previousStatus
    ElMessage({
      type: 'error',
      message: '状态切换失败'
    })
  } finally {
    switchLoadingMap.value[row.id] = false
  }
}


// 删除行
const deletePqApiKeyFunc = async (row) => {
  const res = await deletePqApiKey({ id: row.id })
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
    updatedAt: new Date(),
    deletedAt: new Date(),
    createdAt: new Date(),
    userId: undefined,
    aiModelId: undefined,
    key: '',
    totalTokens: undefined,
    useTokens: '',
    UserName: '',
    status: 1,
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
        res = await createPqApiKey(formData.value)
        break
      case 'update':
        res = await updatePqApiKey(formData.value)
        break
      default:
        res = await createPqApiKey(formData.value)
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
  const res = await findPqApiKey({ id: row.id })
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

<style></style>
