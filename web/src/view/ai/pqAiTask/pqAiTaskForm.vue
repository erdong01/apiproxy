
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="id字段:" prop="id">
    <el-input v-model.number="formData.id" :clearable="true" placeholder="请输入id字段" />
</el-form-item>
        <el-form-item label="创建时间:" prop="createdAt">
    <el-date-picker v-model="formData.createdAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
        <el-form-item label="更新时间:" prop="updatedAt">
    <el-date-picker v-model="formData.updatedAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
        <el-form-item label="deletedAt字段:" prop="deletedAt">
    <el-date-picker v-model="formData.deletedAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
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
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createPqAiTask,
  updatePqAiTask,
  findPqAiTask
} from '@/api/ai/pqAiTask'

defineOptions({
    name: 'PqAiTaskForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
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

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findPqAiTask({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
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
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
