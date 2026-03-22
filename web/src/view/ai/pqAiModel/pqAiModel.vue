<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button v-if="!showAllQuery" link type="primary" icon="arrow-down" @click="showAllQuery = true">
            展开
          </el-button>
          <el-button v-else link type="primary" icon="arrow-up" @click="showAllQuery = false">
            收起
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
        <el-button
          icon="delete"
          style="margin-left: 10px;"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >
          删除
        </el-button>
      </div>

      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <!-- <el-table-column align="left" label="id字段" prop="id" width="120" /> -->
        <el-table-column align="left" label="模型ID" prop="name" width="380" />
        <el-table-column align="left" label="供应商" prop="provider" width="180" />
        <!-- <el-table-column align="left" label="版本" prop="version" width="180" /> -->
        <el-table-column align="left" label="价格条数" width="100">
          <template #default="scope">
            {{ scope.row.PqAiModelPrice?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看
            </el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updatePqAiModelFunc(scope.row)">
              编辑
            </el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="dialogFormVisible"
      destroy-on-close
      :size="formDrawerSize"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增' : '编辑' }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="80px">
        <el-form-item label="模型ID:" prop="name">
          <el-input v-model="formData.name" :clearable="true" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="供应商:" prop="provider">
          <el-input v-model="formData.provider" :clearable="true" placeholder="请输入供应商" />
        </el-form-item>
        <!-- <el-form-item label="版本:" prop="version">
          <el-input v-model="formData.version" :clearable="true" placeholder="请输入版本" />
        </el-form-item> -->

        <el-form-item label="模型价格">
          <div class="price-editor">
            <div class="price-editor__header">
              <span>支持按行添加和删除价格</span>
              <el-button type="primary" link icon="plus" @click="addPriceRow">新增价格</el-button>
            </div>

            <el-table
              v-if="formData.PqAiModelPrice.length"
              :data="formData.PqAiModelPrice"
              border
              class="price-table"
            >
              <el-table-column type="index" label="#" width="60" />
              <el-table-column label="分辨率" min-width="180">
                <template #default="scope">
                  <el-select
                    v-model="scope.row.Resolution"
                    clearable
                    filterable
                    allow-create
                    default-first-option
                    placeholder="请选择或输入分辨率"
                    style="width: 100%;"
                  >
                    <el-option
                      v-for="item in resolutionOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="生成模式" min-width="180">
                <template #default="scope">
                  <el-select
                    v-model="scope.row.GenerationModes"
                    clearable
                    filterable
                    allow-create
                    default-first-option
                    placeholder="请选择或输入生成模式"
                    style="width: 100%;"
                  >
                    <el-option
                      v-for="item in generationModeOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="供应商价格" min-width="140">
                <template #default="scope">
                  <el-input-number
                    v-model="scope.row.VendorPrice"
                    :min="0"
                    :precision="4"
                    controls-position="right"
                    placeholder="请输入"
                    style="width: 100%;"
                  />
                </template>
              </el-table-column>
              <el-table-column label="供应商单位" min-width="150">
                <template #default="scope">
                  <el-select
                    v-model="scope.row.VendorUnit"
                    clearable
                    placeholder="请选择供应商单位"
                    style="width: 100%;"
                  >
                    <el-option
                      v-for="item in vendorUnitOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="零售价格" min-width="140">
                <template #default="scope">
                  <el-input-number
                    v-model="scope.row.RetailPrice"
                    :min="0"
                    :precision="4"
                    controls-position="right"
                    placeholder="请输入"
                    style="width: 100%;"
                  />
                </template>
              </el-table-column>
              <el-table-column label="零售单位" min-width="140">
                <template #default="scope">
                  <el-select
                    v-model="scope.row.RetailUnit"
                    clearable
                    placeholder="请选择零售单位"
                    style="width: 100%;"
                  >
                    <el-option
                      v-for="item in retailUnitOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="90" fixed="right">
                <template #default="scope">
                  <el-button type="primary" link icon="delete" @click="removePriceRow(scope.$index)">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <el-empty v-else description="暂无价格，请新增一行" />
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-model="detailShow"
      destroy-on-close
      :size="appStore.drawerSize"
      :show-close="true"
      :before-close="closeDetailShow"
      title="查看"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="创建时间">
          {{ detailForm.createdAt || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间">
          {{ detailForm.updatedAt || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="模型ID">
          {{ detailForm.name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="供应商">
          {{ detailForm.provider || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="版本">
          {{ detailForm.version || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="detail-price-box">
        <div class="detail-price-box__title">模型价格</div>
        <el-table v-if="detailForm.PqAiModelPrice?.length" :data="detailForm.PqAiModelPrice" border>
          <el-table-column type="index" label="#" width="60" />
          <el-table-column prop="Resolution" label="分辨率" min-width="140" />
          <el-table-column label="生成模式" min-width="140">
            <template #default="scope">
              {{ getGenerationModeLabel(scope.row.GenerationModes) }}
            </template>
          </el-table-column>
          <el-table-column prop="VendorPrice" label="供应商价格" min-width="120" />
          <el-table-column label="供应商单位" min-width="120">
            <template #default="scope">
              {{ getVendorUnitLabel(scope.row.VendorUnit) }}
            </template>
          </el-table-column>
          <el-table-column prop="RetailPrice" label="零售价格" min-width="120" />
          <el-table-column label="零售单位" min-width="120">
            <template #default="scope">
              {{ getRetailUnitLabel(scope.row.RetailUnit) }}
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-else description="暂无价格数据" />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createPqAiModel,
  deletePqAiModel,
  deletePqAiModelByIds,
  updatePqAiModel,
  findPqAiModel,
  getPqAiModelList
} from '@/api/ai/pqAiModel'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, reactive, ref } from 'vue'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'PqAiModel'
})

const resolutionOptions = [
  { label: '1080P', value: '1080p' },
  { label: '720P', value: '720p' },
  { label: '480P', value: '480p' }
]

const generationModeOptions = [
  { label: '无输入视频', value: 'draft_false' },
  { label: '有输入视频', value: 'draft_true' },
  { label: '有声视频', value: 'generate_audio_true' },
  { label: '无声视频', value: 'generate_audio_false' }
]

const vendorUnitOptions = [
  { label: '千tokens', value: '1' },
  { label: '张', value: '2' }
]

const retailUnitOptions = [
  { label: '单秒', value: '1' },
  { label: '单次', value: '2' }
]

const createDefaultPriceItem = () => ({
  Id: undefined,
  PqAiModelId: undefined,
  Resolution: '',
  GenerationModes: '',
  VendorPrice: undefined,
  VendorUnit: '1',
  RetailPrice: undefined,
  RetailUnit: ''
})

const createDefaultFormData = () => ({
  id: undefined,
  name: '',
  provider: '',
  version: '',
  PqAiModelPrice: [createDefaultPriceItem()]
})

const normalizePriceItem = (price = {}) => ({
  ...createDefaultPriceItem(),
  ...price,
  Resolution: price.Resolution || '',
  GenerationModes: price.GenerationModes || '',
  VendorUnit: price.VendorUnit === undefined || price.VendorUnit === null || price.VendorUnit === '' ? '1' : String(price.VendorUnit),
  RetailUnit: price.RetailUnit || '',
  VendorPrice: price.VendorPrice ?? undefined,
  RetailPrice: price.RetailPrice ?? undefined
})

const normalizeFormData = (data = {}) => ({
  ...createDefaultFormData(),
  ...data,
  name: data.name || '',
  provider: data.provider || '',
  version: data.version || '',
  PqAiModelPrice: Array.isArray(data.PqAiModelPrice) && data.PqAiModelPrice.length
    ? data.PqAiModelPrice.map(item => normalizePriceItem(item))
    : [createDefaultPriceItem()]
})

const normalizeDetailData = (data = {}) => ({
  ...data,
  PqAiModelPrice: Array.isArray(data.PqAiModelPrice) ? data.PqAiModelPrice.map(item => normalizePriceItem(item)) : []
})

const hasPriceContent = (price = {}) => {
  const resolution = String(price.Resolution || '').trim()
  const generationModes = String(price.GenerationModes || '').trim()
  const retailUnit = String(price.RetailUnit || '').trim()
  return generationModes !== '' || resolution !== '' || retailUnit !== '' || price.VendorPrice !== undefined || price.RetailPrice !== undefined
}

const buildSubmitData = () => ({
  ...formData.value,
  PqAiModelPrice: (formData.value.PqAiModelPrice || [])
    .map(item => ({
      ...item,
      Resolution: String(item.Resolution || '').trim(),
      GenerationModes: String(item.GenerationModes || '').trim(),
      VendorUnit: item.VendorUnit === undefined || item.VendorUnit === null || item.VendorUnit === '' ? '' : String(item.VendorUnit),
      RetailUnit: String(item.RetailUnit || '').trim()
    }))
    .filter(item => hasPriceContent(item))
})

const getGenerationModeLabel = (value) => {
  const target = generationModeOptions.find(item => item.value === String(value))
  return target ? target.label : value || '-'
}

const getVendorUnitLabel = (value) => {
  const target = vendorUnitOptions.find(item => item.value === String(value))
  return target ? target.label : value || '-'
}

const getRetailUnitLabel = (value) => {
  const target = retailUnitOptions.find(item => item.value === String(value))
  return target ? target.label : value || '-'
}

const btnLoading = ref(false)
const appStore = useAppStore()
const formDrawerSize = computed(() => (appStore.drawerSize === '100%' ? '100%' : '1200px'))
const showAllQuery = ref(false)
const multipleTable = ref()
const formData = ref(createDefaultFormData())
const rule = reactive({})
const elFormRef = ref()
const elSearchFormRef = ref()

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const getTableData = async () => {
  const table = await getPqAiModelList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const onSubmit = () => {
  elSearchFormRef.value?.validate(async valid => {
    if (!valid) {
      return
    }
    page.value = 1
    getTableData()
  })
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const multipleSelection = ref([])

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deletePqAiModelFunc(row)
  })
}

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
    multipleSelection.value.forEach(item => {
      ids.push(item.id)
    })
    const res = await deletePqAiModelByIds({ ids })
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

const type = ref('')

const addPriceRow = () => {
  formData.value.PqAiModelPrice.push(createDefaultPriceItem())
}

const removePriceRow = (index) => {
  formData.value.PqAiModelPrice.splice(index, 1)
}

const updatePqAiModelFunc = async (row) => {
  const res = await findPqAiModel({ id: row.id })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = normalizeFormData(res.data)
    dialogFormVisible.value = true
  }
}

const deletePqAiModelFunc = async (row) => {
  const res = await deletePqAiModel({ id: row.id })
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

const dialogFormVisible = ref(false)

const openDialog = () => {
  type.value = 'create'
  formData.value = createDefaultFormData()
  dialogFormVisible.value = true
}

const closeDialog = (done) => {
  dialogFormVisible.value = false
  formData.value = createDefaultFormData()
  elFormRef.value?.clearValidate()
  done?.()
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) {
      btnLoading.value = false
      return
    }
    let res
    const submitData = buildSubmitData()
    switch (type.value) {
      case 'create':
        res = await createPqAiModel(submitData)
        break
      case 'update':
        res = await updatePqAiModel(submitData)
        break
      default:
        res = await createPqAiModel(submitData)
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

const detailForm = ref({
  PqAiModelPrice: []
})

const detailShow = ref(false)

const openDetailShow = () => {
  detailShow.value = true
}

const getDetails = async (row) => {
  const res = await findPqAiModel({ id: row.id })
  if (res.code === 0) {
    detailForm.value = normalizeDetailData(res.data)
    openDetailShow()
  }
}

const closeDetailShow = (done) => {
  detailShow.value = false
  detailForm.value = {
    PqAiModelPrice: []
  }
  done?.()
}
</script>

<style scoped>
.price-editor {
  width: 100%;
}

.price-editor__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  color: var(--el-text-color-secondary);
}

.price-table {
  width: 100%;
}

.detail-price-box {
  margin-top: 20px;
}

.detail-price-box__title {
  margin-bottom: 12px;
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}
</style>
